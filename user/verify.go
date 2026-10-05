package user

import (
	pgerror "errors"

	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"github.com/Lcat-Technologies/eardrum-interfaces/user"
	"gorm.io/gorm"
)

// Transforms an unverified user record to a verified user record
func VerifyUser(phoneNumber string, Db *gorm.DB) (verifiedUser user.User, finalErr error) {

	// Start a new transaction
	tx := Db.Begin()
	if tx.Error != nil {
		// If the transaction can't even start (e.g., connection issue)
		err1 := errors.New(errors.EARInternalError, tx.Error)
		err1.Log()
		return nil, err1
	}

	// Defer a rollback that will only execute if an error occurred in the transaction body.
	defer func() {
		if finalErr != nil {
			tx.Rollback()
		}
	}()

	var unverifieduser *UnverifiedUser

	// 1. Find the unverified user within the transaction.
	if err := tx.Where("phone_number = ?", phoneNumber).First(&unverifieduser).Error; err != nil {
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			// User not found (expected business failure) -> 404 Not Found
			err1 := errors.New(errors.EARUserNotFoundByPhone, err)
			err1.Log()
			finalErr = err1
			return
		}
		// Other lookup failure (e.g., connection issue) -> 500 Internal
		err1 := errors.New(errors.EARUserLookupFailedByPhone, err)
		err1.Log()
		finalErr = err1
		return
	}

	// transform the unverified user model into a user model
	verifiedUser = &User{
		// ... field assignments
		UserName:              unverifieduser.UserName,
		PhoneNumber:           unverifieduser.PhoneNumber,
		Password:              unverifieduser.Password,
		AccountBalanceInCents: unverifieduser.AccountBalanceInCents,
		OverdraftBalanceInCents: unverifieduser.OverdraftBalanceInCents,
		PinCode:               unverifieduser.PinCode,
		QrCode:                unverifieduser.QrCode,
	}

	// 2. Create the verified user in the transaction.
	if err := tx.Create(verifiedUser).Error; err != nil {
		// If creation fails, check for a unique constraint violation (e.g., QR Code clash)
		// Generic create failure -> 500 Internal
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		finalErr = err1
		return
	}

	// 3. Delete the unverified user in the transaction.
	if err := tx.Delete(unverifieduser).Error; err != nil {
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

	// All successful, return the verified user (finalErr is nil)
	return
}
