package user

import (
	pgerror "errors" // Standard Go errors package

	"github.com/Lcat-Technologies/eardrum-interfaces/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
)

// RegenerateQrCode generates a new UUID for the user's QR code by looking them up via username.
func RegenerateQrCode(Db *gorm.DB, username string) (user.User, error) {
    var verifiedUser *User

    err := Db.Transaction(func(tx *gorm.DB) error {
        // 1. Find the user with a lock using username
        // Use a local variable for the error to keep the main 'err' clean for the outer block.
        if findErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_name = ?", username).First(&verifiedUser).Error; findErr != nil {

            // Critical: If the user is not found during the lock/find step,
            // return the error so the outer block can handle it.
            return findErr 
        }

        // 2. Generate new UUID
        newQrCode := uuid.New()

        // 3. Update the user's record with the new UUID.
        if updateErr := tx.Model(&verifiedUser).Update("qr_code", newQrCode).Error; updateErr != nil {

            // Database failed to execute the update query.
            return updateErr
        }

        // Update the local struct to reflect the new UUID
        verifiedUser.QrCode = newQrCode

        // Return nil to commit the transaction
        return nil
    })

    // --- Error Handling Block ---
    if err != nil {

        // 1. Check for the known 'Not Found' error from the find/lock step
        if pgerror.Is(err, gorm.ErrRecordNotFound) {
            // Returns 404 Not Found
            err1 := errors.New(errors.EARUserNotFoundByUsername, err)
            err1.Log()
            return nil, err1
        }

        // 2. All other errors are internal database or concurrency issues
        // Returns 500 Internal Server Error
        err1 := errors.New(errors.EARInternalError, err)
        err1.Log()
        return nil, err1
    }

    // The transaction was successful.
    return verifiedUser, nil
}