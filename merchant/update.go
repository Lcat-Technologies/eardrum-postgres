package merchant

import (
	pgerror "errors"

	"github.com/Lcat-Technologies/eardrum-interfaces/errors" // Replace with your actual custom errors package path
	"github.com/Lcat-Technologies/eardrum-interfaces/merchant"
	"gorm.io/gorm"
)

// UpdatePassword updates the merchant's password using their username.
func UpdatePassword(Db *gorm.DB, encryptedpassword string, userName string) (merchant.Merchant, error) {
	var m *Merchant

	// 1. Fetch the record by username
	if err := Db.Where("user_name = ?", userName).First(&m).Error; err != nil {
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			err1 := errors.New(errors.EARMerchantNotFoundByUsername, err)
			err1.Log()
			return nil, err1 // 404 Not Found
		}
		err1 := errors.New(errors.EARMerchantLookupFailedByUsername, err)
		err1.Log()
		return nil, err1 // 500 Internal Error
	}

	// 2. Update the password field
	if err := Db.Model(&m).Update("password", encryptedpassword).Error; err != nil {
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1 // 500 Internal Server Error
	}

	return m, nil
}

// ---

// UpdatePinCode updates the merchant's PIN code using their username.
func UpdatePinCode(Db *gorm.DB, encryptedpincode string, userName string) (merchant.Merchant, error) {
	var m *Merchant

	// 1. Fetch the record by username
	if err := Db.Where("user_name = ?", userName).First(&m).Error; err != nil {
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			err1 := errors.New(errors.EARMerchantNotFoundByUsername, err)
			err1.Log()
			return nil, err1
		}
		err1 := errors.New(errors.EARMerchantLookupFailedByUsername, err)
		err1.Log()
		return nil, err1 // 500 Internal Error
	}

	// 2. Update the pin_code field
	if err := Db.Model(&m).Update("pin_code", encryptedpincode).Error; err != nil {
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1 // 500 Internal Server Error
	}

	return m, nil
}
