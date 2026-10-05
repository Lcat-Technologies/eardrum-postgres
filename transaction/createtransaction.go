package transaction

import (
	pgerror "errors"
	"os"
	"strconv"

	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"github.com/Lcat-Technologies/eardrum-interfaces/transaction"
	"github.com/Lcat-Technologies/eardrum-postgres/merchant"
	"github.com/Lcat-Technologies/eardrum-postgres/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProcessTransaction handles a complete, atomic financial transaction for a single amount.
// It validates credentials, processes the amount, updates accounts, and creates the
// necessary database records securely within a single transaction.
// checkPIN is a function passed as an argument to decouple logic.
func ProcessTransaction(db *gorm.DB, merchantUsername string, newTx transaction.NewTransaction) (transaction.Transaction, error) {
	// We use GORM's built-in transaction helper to ensure all operations are
	// either fully completed or fully rolled back if an error occurs.
	var newTransaction *Transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// =========================================================================
		// 1. USER VALIDATION & ROW LOCKING
		// =========================================================================
		var u user.User
		// Find the user by their UUID (QrCode) and immediately lock the row
		// for the duration of this transaction to prevent race conditions.
		var code uuid.UUID
		var err1 *errors.PublicError
		var err2 error

		if len(newTx.GetUUID()) < 36 {
			code, err1 = PartialToFullUuid(db, newTx.GetUUID(), newTx.GetFacialEmbedding())
			if err1 != nil {
				return err1
			}
		} else if len(newTx.GetUUID()) == 36 {
			code, err2 = uuid.Parse(newTx.GetUUID())
			if err2 != nil {
				err3 := errors.New(errors.EARInternalError, err2)
				err3.Log()
				return err3
			}

		} else {
			err3 := errors.New(errors.EARTxUserAccountNotFound, pgerror.New("Invalid QR code"))
			err3.Log()
			return err3
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("qr_code = ?", code).
			First(&u).Error; err != nil {
			// User not found (QR Code invalid) -> 404 NotFound
			if pgerror.Is(err, gorm.ErrRecordNotFound) {
				err3 := errors.New(errors.EARTxUserAccountNotFound, err)
				err3.Log()
				return err3
			}
			// Other DB error during lookup -> 500 Internal Server Error
			err3 := errors.New(errors.EARInternalError, err)
			err3.Log()

			return err3
		}

		// perform facial match
		if !u.MatchFace(newTx.GetFacialEmbedding(), FacialMatchThreshold) {
			err3 := errors.New(errors.EARTxInvalidAuthentication, pgerror.New("Facial mismatch"))
			err3.Log()
			return err3

		}

		// =========================================================================
		// 2. TRANSACTION COST & FINAL BALANCE CHECK
		// =========================================================================
		// We use the provided total amount directly, no need to calculate from products.
		totalAmount := newTx.GetTotalAmountInCents()

		// Basic validation: Amount must be positive.
		if totalAmount <= 0 {
			// Invalid amount -> 400 Bad Request
			err3 := errors.New(errors.EARTxAmountMustBeGreaterThanZero, pgerror.New("Transaction amount must be greater than zero."))
			err3.Log()
			return err3
		}

		// Retrieve the transaction cost percentage from an environment variable.
		feePercentageStr := os.Getenv("TRANSACTION_FEE_PERCENT")
		feePercentage, err := strconv.ParseUint(feePercentageStr, 10, 64)
		if err != nil {
			// Fail the transaction if the environment variable is not set or invalid. -> 500 Internal Server Error
			err3 := errors.New(errors.EARInternalError, err)
			err3.Log()
			return err3
		}

		// Calculate the transaction cost using integer math to prevent floating point inaccuracies.
		transactionCost := (totalAmount * uint(feePercentage)) / 100
		finalDeduction := totalAmount + transactionCost

		// Check if the user has a sufficient balance to cover the total amount.
		if u.AccountBalanceInCents < finalDeduction {
			// Insufficient Balance -> 402 Payment Required
			err3 := errors.New(errors.EARTxInsufficientBalance, pgerror.New("Insufficient balance to complete transaction."))
			err3.Log()
			return err3
		}

		// =========================================================================
		// 3. ACCOUNT UPDATES (LOCKING BOTH USER AND MERCHANT)
		// =========================================================================
		// Find the merchant and lock its row within this same transaction to prevent race conditions.
		var s merchant.Merchant
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&s, "user_name = ?", merchantUsername).Error; err != nil {
			// Merchant not found -> 403 Forbidden (prevents enumeration) or 500
			if pgerror.Is(err, gorm.ErrRecordNotFound) {
				err3 := errors.New(errors.EARTxMerchantAccountNotFound, err)
				err3.Log()
				return err3
			}
			err3 := errors.New(errors.EARMerchantLookupFailedByUsername, err)
			err3.Log()
			return err3
		}

		// Deduct from the user's account and add to the shop's.
		u.AccountBalanceInCents -= finalDeduction
		s.AccountBalanceInCents += totalAmount

		// Save the updated balances to the database within the transaction.
		if err := tx.Save(&u).Error; err != nil {
			err3 := errors.New(errors.EARInternalError, err)
			err3.Log()
			return err3 // 500 Internal (Save error)
		}
		if err := tx.Save(&s).Error; err != nil {
			err3 := errors.New(errors.EARInternalError, err)
			err3.Log()
			return err3 // 500 Internal (Save error)
		}

		// =========================================================================
		// 4. CREATE TRANSACTION RECORD
		// =========================================================================
		// Create the main transaction record with all calculated final amounts.
		newTransaction = &Transaction{
			UserUserName:           u.UserName,
			MerchantUserName:       merchantUsername,
			TotalAmountInCents:     totalAmount,
			TransactionCostInCents: transactionCost,
		}
		if err := tx.Create(newTransaction).Error; err != nil {
			err3 := errors.New(errors.EARInternalError, err)
			err3.Log()
			return err3 // 500 Internal (Create error)
		}

		// This type of transaction does not have purchase records.
		// The `db.Transaction` helper will automatically commit.
		return nil
	})

	// If the transaction failed, the structured error will be returned here.
	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}



