package merchant

import (
	pgerror "errors"

	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"github.com/Lcat-Technologies/eardrum-interfaces/merchant"
	"gorm.io/gorm"
	// Import your custom error helpers
)

// GetMerchantWithPhoneNumber finds a merchant by phone number.
func GetMerchantWithPhoneNumber(Db *gorm.DB, PhoneNumber string) (merchant.Merchant, error) {
	var m *Merchant

	// Find the first merchant that matches the input phonenumber
	if err := Db.Where("phone_number = ?", PhoneNumber).First(&m).Error; err != nil {

		// 1. Check for the known "Not Found" condition
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			// Returns 404 Not Found
			err1 := errors.New(errors.EARMerchantNotFoundByPhone, err)
			err1.Log()
			return nil, err1
		}

		// 2. All other errors (connection, query syntax, etc.) -> 500 Internal
		err1 := errors.New(errors.EARMerchantLookupFailedByPhone, err)
		err1.Log()
		return nil, err1
	}

	return m, nil
}

// GetMerchantWithUserName finds a merchant by username.
func GetMerchantWithUserName(Db *gorm.DB, userName string) (merchant.Merchant, error) {
	var m *Merchant

	// Find the first merchant that matches the input username
	if err := Db.Where("user_name = ?", userName).First(&m).Error; err != nil {

		// 1. Check for the known "Not Found" condition
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			// Returns 404 Not Found
			err1 := errors.New(errors.EARMerchantNotFoundByUsername, err)
			err1.Log()
			return nil, err1
		}

		// 2. All other errors (connection, query syntax, etc.) -> 500 Internal
		err1 := errors.New(errors.EARMerchantLookupFailedByUsername, err)
		err1.Log()
		return nil, err1
	}

	return m, nil
}

// GetMerchants retrieves a paginated list of merchants registered in the database.
// - limit: The maximum number of records to return (e.g., 5)
// - offset: The number of records to skip before starting to return (e.g., 0 for the first page)
func GetMerchants(Db *gorm.DB, limit int, offset int) ([]merchant.Merchant, error) {
	var merchants []*Merchant

	// Find all records with limit and offset applied
	if err := Db.Limit(limit).Offset(offset).Find(&merchants).Error; err != nil {
		// Db.Find only returns an error on connection or query issue, not if the table is empty.
		err1 := errors.New(errors.EARMerchantListRetrievalFailed, err)
		err1.Log()
		return nil, err1
	}

	// Transform [](*Merchant) to []merchant.Merchant
	var merchantlist []merchant.Merchant
	for _, m := range merchants {
		merchantlist = append(merchantlist, m)
	}

	return merchantlist, nil
}
