package user

import (
    pgerror "errors"

    "github.com/Lcat-Technologies/eardrum-interfaces/errors"
    "github.com/Lcat-Technologies/eardrum-interfaces/user"
    "gorm.io/gorm"
    "github.com/lib/pq"
    "github.com/GigaDesk/eardrum-prefix/validate"
)

// UpdatePassword updates the user's password using their username.
func UpdatePassword(Db *gorm.DB, encryptedpassword string, username string) (user.User, error) {
    var u *User

    // 1. Fetch the record by username
    if err := Db.Where("user_name = ?", username).First(&u).Error; err != nil {
        if pgerror.Is(err, gorm.ErrRecordNotFound) {
            err1 := errors.New(errors.EARUserNotFoundByUsername, err)
            err1.Log()
            return nil, err1 // 404 Not Found
        }
        err1 := errors.New(errors.EARUserLookupFailedByUsername, err)
        err1.Log()
        return nil, err1 // 500 Internal Error
    }

    // 2. Update the password field
    if err := Db.Model(&u).Update("password", encryptedpassword).Error; err != nil {
        err1 := errors.New(errors.EARInternalError, err)
        err1.Log()
        return nil, err1 // 500 Internal Server Error
    }

    return u, nil
}

// ---

// UpdatePinCode updates the user's PIN code using their username.
func UpdatePinCode(Db *gorm.DB, encryptedpincode string, username string) (user.User, error) {
    var u *User

    // 1. Fetch the record by username
    if err := Db.Where("user_name = ?", username).First(&u).Error; err != nil {
        if pgerror.Is(err, gorm.ErrRecordNotFound) {
            err1 := errors.New(errors.EARUserNotFoundByUsername, err)
            err1.Log()
            return nil, err1 // 404 Not Found
        }
        err1 := errors.New(errors.EARUserLookupFailedByUsername, err)
        err1.Log()
        return nil, err1 // 500 Internal Error
    }

    // 2. Update the pin_code field
    if err := Db.Model(&u).Update("pin_code", encryptedpincode).Error; err != nil {
        err1 := errors.New(errors.EARInternalError, err)
        err1.Log()
        return nil, err1 // 500 Internal Server Error
    }

    return u, nil
}

// UpdateFacialEmbeddings updates the user's facial embeddings using their username.
func UpdateFacialEmbeddings(Db *gorm.DB, embeddings []string, username string) (user.User, error) {
    
    //Validate embeddings
    for _, embedding:=range embeddings{
        _, err:=validate.ValidateMobileFaceNetEmbeddingB64(embedding, false)
        if err!=nil{
            return nil, err
        }
    }
   
    // 1. Fetch the user using the existing helper function
    u, err := GetUserWithUsername(Db, username)
    if err != nil {
        return nil, err // Returns 404 or 500 error handled by GetUserWithUsername
    }

    // 2. Cast []string to pq.StringArray and update the facial_embeddings column
    if err := Db.Model(u).Update("facial_embeddings", pq.StringArray(embeddings)).Error; err != nil {
        err1 := errors.New(errors.EARInternalError, err)
        err1.Log()
        return nil, err1 // 500 Internal Server Error
    }

    return u, nil
}