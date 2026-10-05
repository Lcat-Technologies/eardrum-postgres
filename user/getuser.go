package user

import (
	pgerror "errors" // Standard Go errors package

	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"github.com/Lcat-Technologies/eardrum-interfaces/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetUserWithPhoneNumber finds a user by phone number.
func GetUserWithPhoneNumber(Db *gorm.DB, PhoneNumber string) (user.User, error) {
	var u *User

	// Find the first user that matches the input phonenumber
	if err := Db.Where("phone_number = ?", PhoneNumber).First(&u).Error; err != nil {
		
			// 1. Check for the known "Not Found" condition
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			// Returns 404 Not Found
			err1 := errors.New(errors.EARUserNotFoundByPhone, err)
			err1.Log()
			return nil, err1
		}

		// 2. All other errors (connection, query syntax, etc.) -> 500 Internal
		err1 := errors.New(errors.EARUserLookupFailedByPhone, err)
		err1.Log()
		return nil, err1
	}

	return u, nil
}

// GetUserWithUsername finds a user by username.
func GetUserWithUsername(Db *gorm.DB, Username string) (user.User, error) {
    var u *User

    // Find the first user that matches the input username
    if err := Db.Where("user_name = ?", Username).First(&u).Error; err != nil {
        
        // 1. Check for the known "Not Found" condition
        if pgerror.Is(err, gorm.ErrRecordNotFound) {
            // Returns 404 Not Found (Assuming EARUserNotFoundByUsername exists in your interfaces)
            err1 := errors.New(errors.EARUserNotFoundByUsername, err)
            err1.Log()
            return nil, err1
        }

        // 2. All other errors (connection, query syntax, etc.) -> 500 Internal
        err1 := errors.New(errors.EARUserLookupFailedByUsername, err)
        err1.Log()
        return nil, err1
    }

    return u, nil
}
// GetUsers retrieves a paginated list of users registered in the database.
// - limit: The maximum number of records to return (e.g., 5)
// - offset: The number of records to skip before starting to return (e.g., 0 for the first page)
func GetUsers(Db *gorm.DB, limit int, offset int) ([]user.User, error) {
    var users []*User

    // Find all records with limit and offset applied
    if err := Db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
        // Db.Find only returns an error on connection or query issue, not if the table is empty.
        err1 := errors.New(errors.EARUserListRetrievalFailed, err)
        err1.Log()
        return nil, err1
    }

    // Transform [](*User) to []user.User
    var userslist []user.User
    for _, u := range users {
        userslist = append(userslist, u)
    }

    return userslist, nil
}


// GetUsersByUUIDs retrieves all users whose QrCode UUID matches any of the given UUIDs.
// Returns a nil slice if no matching users are found or if the input slice is empty.
func GetUsersByUUIDs(Db *gorm.DB, uuids []uuid.UUID) ([]user.User, error) {
	// 1. Guard against empty input to prevent generating invalid SQL queries like `WHERE qr_code IN ()`
	if len(uuids) == 0 {
		return nil, nil
	}

	var users []*User

	// 2. Query matching records using Postgres IN clause
	if err := Db.Where("qr_code IN (?)", uuids).Find(&users).Error; err != nil {
		err1 := errors.New(errors.EARUserListRetrievalFailed, err)
		err1.Log()
		return nil, err1
	}

	// 3. Pre-allocate capacity (0 length, len(users) capacity) so append starts at index 0 without pre-filling nil values
	usersList := make([]user.User, 0, len(users))
	for _, u := range users {
		if u != nil {
			usersList = append(usersList, u)
		}
	}

	// 4. Return explicit nil if no matching users exist in the database
	if len(usersList) == 0 {
		return nil, nil
	}

	return usersList, nil
}


