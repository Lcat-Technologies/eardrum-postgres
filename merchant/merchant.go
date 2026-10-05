package merchant

import (
	"time"

	"gorm.io/gorm"
)

// Merchant represents the merchant model for the system.
type Merchant struct {
	gorm.Model
	UserName              string  `gorm:"uniqueIndex;not null"` // UserName of the merchant (must be unique)
	PhoneNumber           string  `gorm:"uniqueIndex;not null"` // Phone number of the merchant
	Password              string  `gorm:"not null"`             // The merchant's password. It should be stored as a secure hash.
	AccountBalanceInCents uint    `gorm:"not null;default:0"`   // The merchant's account balance in cents
	PinCode               *string `gorm:""`                     // The merchant's security PIN code. It should be stored as a secure hash (NULLABLE).
}

// Returns the creation timestamp of the merchant
func (s Merchant) GetCreatedAt() time.Time {
	return s.CreatedAt.UTC()
}

// Returns the update timestamp of the merchant
func (s Merchant) GetUpdatedAt() time.Time {
	return s.UpdatedAt.UTC()
}

// Returns the deletion timestamp of the merchant
func (s Merchant) GetDeletedAt() time.Time {
	return s.DeletedAt.Time.UTC()
}

// Returns the unique username of the merchant
func (s Merchant) GetUserName() string {
	return s.UserName
}

// Returns the phone number of the merchant
func (s Merchant) GetPhoneNumber() string {
	return s.PhoneNumber
}

// Returns the account balance in cents of the merchant
func (s Merchant) GetAccountBalanceInCents() int64 {
	return int64(s.AccountBalanceInCents)
}

// Returns the security password of the shop
func (s Merchant) GetPassword() string {
	return s.Password
}

// Returns the security PIN code of the shop
func (s Merchant) GetPinCode() *string {
	return s.PinCode
}

// GetPinStatus returns true if the PIN code is populated.
func (s Merchant) GetPinStatus() bool {
    if s.PinCode == nil {
        return false
    }
    return *s.PinCode != ""
}

// UnverifiedMerchant represents the unverified merchant model for the system.
type UnverifiedMerchant struct {
	gorm.Model
	UserName              string  `gorm:"uniqueIndex;not null"` // Name of the merchant (must be unique)
	PhoneNumber           string  `gorm:"uniqueIndex;not null"` // Phone number of the merchant
	Password              string  `gorm:"not null"`             // The merchant's password. It should be stored as a secure hash.
	AccountBalanceInCents uint    `gorm:"not null;default:0"`   // The merchant's account balance in cents
	PinCode               *string `gorm:""`                     // The merchant's security PIN code. It should be stored as a secure hash.
}

// Returns the creation timestamp of the unverified merchant
func (s UnverifiedMerchant) GetCreatedAt() time.Time {
	return s.CreatedAt.UTC()
}

// Returns the update timestamp of the unverified merchant
func (s UnverifiedMerchant) GetUpdatedAt() time.Time {
	return s.UpdatedAt.UTC()
}

// Returns the deletion timestamp of the unverified merchant
func (s UnverifiedMerchant) GetDeletedAt() time.Time {
	return s.DeletedAt.Time.UTC()
}

// Returns the name of the unverified merchant
func (s UnverifiedMerchant) GetUserName() string {
	return s.UserName
}

// Returns the phone number of the unverified merchant
func (s UnverifiedMerchant) GetPhoneNumber() string {
	return s.PhoneNumber
}

// Returns the account balance in cents of the unverified merchant
func (s UnverifiedMerchant) GetAccountBalanceInCents() int64 {
	return int64(s.AccountBalanceInCents)
}

// Returns the security password of the unverified merchant
func (s UnverifiedMerchant) GetPassword() string {
	return s.Password
}

// Returns the security PIN code of the unverified merchant
func (s UnverifiedMerchant) GetPinCode() *string {
	return s.PinCode
}

// GetPinStatus returns true if the PIN code is populated.
func (s UnverifiedMerchant) GetPinStatus() bool {
    if s.PinCode == nil {
        return false
    }
    return *s.PinCode != ""
}

// UniquenessCheck represents the availability status of a unique identifier (like phone or username),
// indicating whether it is already taken and if the existing account is verified.
type UniquenessCheck struct {
	Exists     bool `json:"exists"`      // True if the identifier is found in either verified or unverified tables.
	IsVerified bool `json:"is_verified"` // True if the existing record belongs to a fully verified merchant.
}
