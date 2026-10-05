package merchant

import (
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"gorm.io/gorm"
)

// CheckMerchantPhoneNumber checks if a phone number exists in either the merchant or unverified merchant table.
func CheckMerchantPhoneNumber(Db *gorm.DB, phoneNumber string) (*UniquenessCheck, error) {
	var verifiedCount int64
	var unverifiedCount int64

	// 1. Check Verified Merchant Count
	if err := Db.Model(&Merchant{}).Where("phone_number = ?", phoneNumber).Count(&verifiedCount).Error; err != nil {
		err1 := errors.New(errors.EARMerchantLookupFailedByPhone, err)
		err1.Log()
		return nil, err1
	}

	// 2. Check Unverified Merchant Count
	if err := Db.Model(&UnverifiedMerchant{}).Where("phone_number = ?", phoneNumber).Count(&unverifiedCount).Error; err != nil {
		err1 := errors.New(errors.EARMerchantLookupFailedByPhone, err)
		err1.Log()
		return nil, err1
	}

	return &UniquenessCheck{
		Exists:     (verifiedCount > 0) || (unverifiedCount > 0),
		IsVerified: verifiedCount > 0,
	}, nil
}

// CheckMerchantUserName checks if a username exists in either the merchant or unverified merchant table.
func CheckMerchantUserName(Db *gorm.DB, userName string) (*UniquenessCheck, error) {
	var verifiedCount int64
	var unverifiedCount int64

	// 1. Check Verified Merchant Count
	if err := Db.Model(&Merchant{}).Where("user_name = ?", userName).Count(&verifiedCount).Error; err != nil {
		err1 := errors.New(errors.EARMerchantLookupFailedByUsername, err)
		err1.Log()
		return nil, err1
	}

	// 2. Check Unverified Merchant Count
	if err := Db.Model(&UnverifiedMerchant{}).Where("user_name = ?", userName).Count(&unverifiedCount).Error; err != nil {
		err1 := errors.New(errors.EARMerchantLookupFailedByUsername, err)
		err1.Log()
		return nil, err1
	}

	return &UniquenessCheck{
		Exists:     (verifiedCount > 0) || (unverifiedCount > 0),
		IsVerified: verifiedCount > 0,
	}, nil
}
