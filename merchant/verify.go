package merchant

import (
	pgerror "errors"

	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"github.com/Lcat-Technologies/eardrum-interfaces/merchant"
	"gorm.io/gorm"
)

// Transforms an unverified shop record to a verified merchant record
// This function uses named return variables 'verifiedMerchant' and 'err'
func VerifyMerchant(phoneNumber string, Db *gorm.DB) (verifiedMerchant merchant.Merchant, finalErr error) {

	// Start a new transaction
	tx := Db.Begin()
	if tx.Error != nil {

		// If the transaction can't even start (e.g., connection issue)
		// All other persistence failures -> 500 Internal Server Error
		err1 := errors.New(errors.EARInternalError, tx.Error)
		err1.Log()
		return nil, err1
	}

	// Defer a rollback that will only execute if an error occurred.
	// The named return variable 'finalErr' is used for structured errors.
	defer func() {
		if finalErr != nil {
			tx.Rollback()
		}
	}()

	var unverifiedmerchant *UnverifiedMerchant

	// 1. Find the unverified merchant within the transaction.
	if err := tx.Where("phone_number = ?", phoneNumber).First(&unverifiedmerchant).Error; err != nil {
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			// Merchant not found (expected business failure) -> 404 Not Found
			err1 := errors.New(errors.EARMerchantNotFoundByPhone, err)
			err1.Log()
			finalErr = err1
			return
		}
		// Other lookup failure (e.g., connection issue) -> 500 Internal
		err1 := errors.New(errors.EARMerchantLookupFailedByPhone, err)
		err1.Log()
		finalErr = err1
		return
	}

	// transform the unverified merchant model into a merchant model
	verifiedMerchant = &Merchant{
		// ... field assignments
		UserName:              unverifiedmerchant.UserName,
		PhoneNumber:           unverifiedmerchant.PhoneNumber,
		Password:              unverifiedmerchant.Password,
		AccountBalanceInCents: unverifiedmerchant.AccountBalanceInCents,
		PinCode:               unverifiedmerchant.PinCode,
	}

	// 2. Create the verified merchant in the transaction.
	if err := tx.Create(verifiedMerchant).Error; err != nil {
		// Generic create failure -> 500 Internal
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		finalErr = err1
		return
	}

	// 3. Delete the unverified merchant in the transaction.
	if err := tx.Delete(unverifiedmerchant).Error; err != nil {
		// Delete failure -> 500 Internal
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		finalErr = err1
		return
	}

	// 4. Commit the transaction.
	if err := tx.Commit().Error; err != nil {
		// Commit failure -> 500 Internal
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		finalErr = err1
		return
	}

	// The function returns here with finalErr == nil
	return
}
