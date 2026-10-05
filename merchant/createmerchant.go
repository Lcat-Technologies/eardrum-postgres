package merchant

import (
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"github.com/Lcat-Technologies/eardrum-interfaces/merchant"
	"gorm.io/gorm"
)

// CreateMerchant creates an unverified merchant record after ensuring the phone number and username are unique.
func CreateMerchant(s merchant.NewMerchant, Db *gorm.DB) (merchant.Merchant, error) {
	// 1. Check if the phone number already exists
	phoneCheck, err := CheckMerchantPhoneNumber(Db, s.GetPhoneNumber())
	if err != nil {
		// Database failure during the lookup/count operation -> 500 Internal Server Error
		return nil, err
	}

	if phoneCheck.Exists {
		if !phoneCheck.IsVerified {
			err1 := errors.New(errors.EARMerchantPhoneExistsUnverified, err)
			err1.Log()
			return nil, err1
		}
		// Business logic conflict: User fully exists -> 409 Conflict
		err1 := errors.New(errors.EARMerchantPhoneExistsVerified, err)
		err1.Log()
		return nil, err1
	}

	// 2. Check if the username already exists
	usernameCheck, err := CheckMerchantUserName(Db, s.GetUserName())
	if err != nil {
		// Database failure during the lookup/count operation -> 500 Internal Server Error
		return nil, err
	}

	if usernameCheck.Exists {
		if !usernameCheck.IsVerified {
			// Business logic conflict: Username exists in unverified state -> 409 Conflict
			err1 := errors.New(errors.EARMerchantUsernameExistsUnverified, err)
			err1.Log()
			return nil, err1
		}
		// Business logic conflict: Username fully exists -> 409 Conflict
		err1 := errors.New(errors.EARMerchantUsernameExistsVerified, err)
		err1.Log()
		return nil, err1
	}

	// 3. Create unverified merchant data
	unverifiedmerchant := &UnverifiedMerchant{
		UserName:              s.GetUserName(),
		PhoneNumber:           s.GetPhoneNumber(),
		Password:              s.GetPassword(),
		AccountBalanceInCents: 0,
	}

	// 4. Create an unverified merchant record in the database
	if err := Db.Create(unverifiedmerchant).Error; err != nil {

		// All other persistence failures -> 500 Internal Server Error
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}

	return unverifiedmerchant, nil
}
