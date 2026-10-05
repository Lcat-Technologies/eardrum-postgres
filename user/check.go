package user

import (
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"gorm.io/gorm"
	// Assuming the error helpers are available in this package
)

// checks if phone number exists in both the unverified user table and the user table.
// Takes database instance and phone number as arguments
func CheckUserPhoneNumber(Db *gorm.DB, phoneNumber string) (*UniquenessCheck, error) {
	var verifiedcount int64
	var unverifiedcount int64

	// 1. Check Verified User Count
	if err := Db.Model(&User{}).Where("phone_number = ?", phoneNumber).Count(&verifiedcount).Error; err != nil {
		// If an error occurs here, it's a database failure (connection, query syntax, etc.)
		err1 := errors.New(errors.EARUserLookupFailedByPhone, err)
		err1.Log()
		return nil, err1
	}

	// 2. Check Unverified User Count
	if err := Db.Model(&UnverifiedUser{}).Where("phone_number = ?", phoneNumber).Count(&unverifiedcount).Error; err != nil {
		// If an error occurs here, it's a database failure (connection, query syntax, etc.)
		err1 := errors.New(errors.EARUserLookupFailedByPhone, err)
		err1.Log()
		return nil, err1
	}

	// Success: Return the result struct
	return &UniquenessCheck{
		Exists:     (verifiedcount > 0) || (unverifiedcount > 0),
		IsVerified: verifiedcount > 0,
	}, nil
}

// CheckUserUserName checks if a username exists in both the unverified user table and the user table.
// Takes database instance and username as arguments.
func CheckUserUserName(Db *gorm.DB, userName string) (*UniquenessCheck, error) {
	var verifiedcount int64
	var unverifiedcount int64

	// 1. Check Verified User Count
	if err := Db.Model(&User{}).Where("user_name = ?", userName).Count(&verifiedcount).Error; err != nil {
		// If an error occurs here, it's a database failure (connection, query syntax, etc.)
		err1 := errors.New(errors.EARUserLookupFailedByUsername, err)
		err1.Log()
		return nil, err1
	}

	// 2. Check Unverified User Count
	if err := Db.Model(&UnverifiedUser{}).Where("user_name = ?", userName).Count(&unverifiedcount).Error; err != nil {
		// If an error occurs here, it's a database failure (connection, query syntax, etc.)
		err1 := errors.New(errors.EARUserLookupFailedByUsername, err)
		err1.Log()
		return nil, err1
	}

	// Success: Return the result struct
	return &UniquenessCheck{
		Exists:     (verifiedcount > 0) || (unverifiedcount > 0),
		IsVerified: verifiedcount > 0,
	}, nil
}
