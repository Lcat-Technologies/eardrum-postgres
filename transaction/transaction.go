package transaction

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"crypto/rand"

	"github.com/Lcat-Technologies/eardrum-postgres/merchant"
	"github.com/Lcat-Technologies/eardrum-postgres/user"
	"gorm.io/gorm"
)

// Transaction represents a financial transaction.
// It links a User and a Merchant via their unique UserNames.
type Transaction struct {
	gorm.Model

    // TransactionID is a 12-character reference code (e.g., 260627R8K4WX)
	TransactionID string `gorm:"uniqueIndex;not null;type:varchar(12)"`

	// OfflineTransactionID combines IMEI and Millisecond Unix Timestamp (e.g., 864209041234567_1773130225123)
	OfflineTransactionID string `gorm:"uniqueIndex;type:varchar(36)"`

	// TotalAmountInCents is the transaction amount, stored in cents.
	TotalAmountInCents uint `gorm:"not null"`

	// TransactionCostInCents is the cost of the transaction.
	TransactionCostInCents uint `gorm:"not null"`

	// ScanLog stores the base64-encoded image of the scan that authorized the transaction.
	ScanLog string `gorm:"type:text"`

	// UserUserName acts as the foreign key referencing the User's UserName.
	UserUserName string `gorm:"not null"`

	// User is the GORM association to the User model.
	// We specify both foreignKey (on this struct) and references (on the User struct).
	User user.User `gorm:"foreignKey:UserUserName;references:UserName"`

	// MerchantUserName acts as the foreign key referencing the Merchant's UserName.
	MerchantUserName string `gorm:"not null"`

	// Merchant is the GORM association to the Merchant model.
	Merchant merchant.Merchant `gorm:"foreignKey:MerchantUserName;references:UserName"`
}

// Returns the unique identifier of the transaction
func (t Transaction) GetTransactionID() string {
	return t.TransactionID
}

// Returns the creation timestamp of the transaction
func (t Transaction) GetCreatedAt() time.Time {
	return t.CreatedAt.UTC()
}

// Returns the update timestamp of the transaction
func (t Transaction) GetUpdatedAt() time.Time {
	return t.UpdatedAt.UTC()
}

// Returns the deletion timestamp of the transaction
func (t Transaction) GetDeletedAt() time.Time {
	return t.DeletedAt.Time.UTC()
}

// Returns the total amount in cents spent in the transaction
func (t Transaction) GetTotalAmountInCents() uint {
	return t.TotalAmountInCents
}

// Returns the transaction cost in cents spent in the transaction
func (t Transaction) GetTransactionCostInCents() uint {
	return t.TransactionCostInCents
}

//Returns image of the scan that authorized transaction
func(t Transaction) GetScanLog() string{
	return t.ScanLog
}

// Returns the unique username of the merchant the transaction was made to. 🏪
func (t Transaction) GetMerchantName() string {
	return t.MerchantUserName
}


// Returns the unique username of the user that made the transtion
func (t Transaction) GetUserName() string {
	return t.UserUserName
}

// Omitted confusing characters 'I', 'O', '0', '1' to keep strings clear
const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" 

// BeforeCreate automatically runs right before the database insert operation.
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	// 1. Get the current date (6 characters)
	datePrefix := time.Now().Format("060102")

	// 2. Generate a high-entropy 6-character suffix safely
	randomSuffix, err := GenerateSecureSuffix(6)
	if err != nil {
		return err
	}

	// 3. Assemble the final 12-character collision-free token
	t.TransactionID = datePrefix + randomSuffix
	return nil
}

// GenerateSecureSuffix uses crypto/rand for secure distributed generation at scale.
func GenerateSecureSuffix(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		// Map the random byte cleanly into our 32-character alphabet
		bytes[i] = charset[bytes[i]%uintbyte(len(charset))]
	}
	return string(bytes), nil
}

// Helper block for type conversion safety
func uintbyte(n int) byte {
	return byte(n)
}

// GetTransactionUnixTimeStamp extracts the millisecond timestamp from the OfflineTransactionID
// and returns both the raw int64 millisecond value and a UTC time.Time object.
func (t Transaction) GetTransactionUnixTimeStamp() (int64, time.Time, error) {
	if t.OfflineTransactionID == "" {
		return 0, time.Time{}, errors.New("offline transaction ID is empty")
	}

	parts := strings.Split(t.OfflineTransactionID, "_")
	if len(parts) < 2 {
		return 0, time.Time{}, errors.New("invalid offline transaction ID format")
	}

	millis, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, time.Time{}, errors.New("failed to parse millisecond timestamp: " + err.Error())
	}

	// Convert milliseconds to time.Time and explicitly enforce UTC location
	txTime := time.Unix(millis/1000, (millis%1000)*1000000).UTC()

	return millis, txTime, nil
}


// GetTransactionDeviceID extracts the device ID from the OfflineTransactionID.
func (t Transaction) GetTransactionDeviceID() (string, error) {
	if t.OfflineTransactionID == "" {
		return "", errors.New("offline transaction ID is empty")
	}

	parts := strings.Split(t.OfflineTransactionID, "_")
	if len(parts) < 2 || parts[0] == "" {
		return "", errors.New("invalid offline transaction ID format")
	}

	return parts[0], nil
}

//GetTransactionDeviceModel returns the model of the device that processed transaction
func(t Transaction) GetTransactionDeviceModel() *string{
	return nil
}