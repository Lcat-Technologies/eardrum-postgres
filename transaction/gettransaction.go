package transaction

import (
	pgerror "errors"
	"time"

	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"github.com/Lcat-Technologies/eardrum-interfaces/transaction"
	"gorm.io/gorm"
)

// GetTransactionByReference gets a transaction by its 12-character TransactionID string.
func GetTransactionByReference(Db *gorm.DB, transactionID string) (transaction.Transaction, error) {
	var tx Transaction
	if err := Db.Where("transaction_id = ?", transactionID).First(&tx).Error; err != nil {
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			err1 := errors.New(errors.EARFileNotFound, err)
			err1.Log()
			return nil, err
		}
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}
	return tx, nil
}

// GetTransactions retrieves all transactions with optional pagination and timestamp filters.
func GetTransactions(Db *gorm.DB, limit *int, offset *int, startTime *time.Time, endTime *time.Time) ([]transaction.Transaction, error) {
	var txs []Transaction
	query := Db.Order("created_at desc")

	// Apply optional timestamp filters
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Apply optional pagination
	if limit != nil && *limit > 0 {
		query = query.Limit(*limit)
	} else {
		query = query.Limit(10) // Safe default
	}

	if offset != nil && *offset > 0 {
		query = query.Offset(*offset)
	}

	if err := query.Find(&txs).Error; err != nil {
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}

	var transactionslist []transaction.Transaction
	for _, t := range txs {
		transactionslist = append(transactionslist, t)
	}

	return transactionslist, nil
}

// GetTransactionsForUser retrieves all transactions for a specific user with optional pagination and timestamp filters.
func GetTransactionsForUser(db *gorm.DB, username string, limit *int, offset *int, startTime *time.Time, endTime *time.Time) ([]transaction.Transaction, error) {
	var txs []Transaction
	query := db.Where("user_user_name = ?", username).Order("created_at desc")

	// Apply optional timestamp filters
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Apply optional pagination
	if limit != nil && *limit > 0 {
		query = query.Limit(*limit)
	} else {
		query = query.Limit(10)
	}

	if offset != nil && *offset > 0 {
		query = query.Offset(*offset)
	}

	if err := query.Find(&txs).Error; err != nil {
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}

	var transactionslist []transaction.Transaction
	for _, t := range txs {
		transactionslist = append(transactionslist, t)
	}

	return transactionslist, nil
}

// GetTransactionsForMerchant retrieves all transactions for a specific merchant with optional pagination and timestamp filters.
func GetTransactionsForMerchant(db *gorm.DB, merchantname string, limit *int, offset *int, startTime *time.Time, endTime *time.Time) ([]transaction.Transaction, error) {
	var txs []Transaction
	query := db.Where("merchant_user_name = ?", merchantname).Order("created_at desc")

	// Apply optional timestamp filters
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// Apply optional pagination
	if limit != nil && *limit > 0 {
		query = query.Limit(*limit)
	} else {
		query = query.Limit(10)
	}

	if offset != nil && *offset > 0 {
		query = query.Offset(*offset)
	}

	if err := query.Find(&txs).Error; err != nil {
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}

	var transactionslist []transaction.Transaction
	for _, t := range txs {
		transactionslist = append(transactionslist, t)
	}

	return transactionslist, nil
}
