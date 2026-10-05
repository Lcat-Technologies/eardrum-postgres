package user

import (
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"github.com/Lcat-Technologies/eardrum-interfaces/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// creates an unverified user record
func CreateUser(s user.NewUser, Db *gorm.DB) (user.User, error) {
	// 1. Check if the phone number already exists
	phoneCheck, err := CheckUserPhoneNumber(Db, s.GetPhoneNumber())

	if err != nil {
		// Database connection/query failure on lookup -> 500 Internal
		return nil, err
	}

	if phoneCheck.Exists {
		if !phoneCheck.IsVerified {
			err1 := errors.New(errors.EARUserPhoneExistsUnverified, err)
			err1.Log()
			return nil, err1
		}
		// Business logic conflict: User fully exists -> 409 Conflict
		err1 := errors.New(errors.EARUserPhoneExistsVerified, err)
		err1.Log()
		return nil, err1
	}

	// 2. Check if the username already exists
	usernameCheck, err := CheckUserUserName(Db, s.GetUserName())
	if err != nil {
		// Database failure during the lookup/count operation -> 500 Internal Server Error
		return nil, err
	}

	if usernameCheck.Exists {
		if !usernameCheck.IsVerified {
			// Business logic conflict: Username exists in unverified state -> 409 Conflict
			err1 := errors.New(errors.EARUserUsernameExistsUnverified, err)
			err1.Log()
			return nil, err1
		}
		// Business logic conflict: Username fully exists -> 409 Conflict
		err1 := errors.New(errors.EARUserUsernameExistsVerified, err)
		err1.Log()
		return nil, err1
	}

	// 2. Create unverified user data
	// (Omitted fields for brevity, assuming UnverifiedUser struct exists)
	unverifieduser := &UnverifiedUser{
		UserName:    s.GetUserName(),
		PhoneNumber: s.GetPhoneNumber(),
		Password:    s.GetPassword(),
	}

	// Generate a new UUID and assign it to the QrCode field
	unverifieduser.QrCode = uuid.New()

	// 3. Create record in the database
	if err := Db.Create(unverifieduser).Error; err != nil {


		// B. All other unexpected DB errors -> 500 Internal
		// The helper wraps the raw 'err' to preserve it for logging.
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}

	return unverifieduser, nil
}
