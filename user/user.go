package user

import (
	"encoding/base64"
	"encoding/binary"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

// User represents the customer model for the system.
type User struct {
	gorm.Model
	UserName                string         `gorm:"uniqueIndex;not null"`  // Unique Username of the user
	PhoneNumber             string         `gorm:"uniqueIndex;not null"`  // Unique Phone number of the user
	Password                string         `gorm:"not null"`              // The user's password. It should be stored as a secure hash.
	AccountBalanceInCents   uint           `gorm:"not null;default:0"`    // The user's account balance in cents
	OverdraftBalanceInCents uint           `gorm:"not null;default:0"`    // The user's overdraft account balance in cents
	PinCode                 *string        `gorm:""`                      // The user's security PIN code (can be NULL). Must be stored as a secure hash.
	QrCode                  uuid.UUID      `gorm:"uniqueIndex;type:uuid"` // The user's unique Qr Code
	FacialEmbeddings        pq.StringArray `gorm:"type:text[]"`           // Stores array of base64 user facial embeddings
	FacialImages            pq.StringArray `gorm:"type:text[]"`           // Stores array of the user's facial image URLs
	Passport                *string        `gorm:""`                      //Stores URL of the user's passport photo
}

// Returns the creation timestamp of the user
func (u User) GetCreatedAt() time.Time {
	return u.CreatedAt.UTC()
}

// Returns the update timestamp of the user
func (u User) GetUpdatedAt() time.Time {
	return u.UpdatedAt.UTC()
}

// Returns the deletion timestamp of the user
func (u User) GetDeletedAt() time.Time {
	return u.DeletedAt.Time.UTC()
}

// Returns the username of the user
func (u User) GetUserName() string {
	return u.UserName
}

// Returns the phone number of the user
func (u User) GetPhoneNumber() string {
	return u.PhoneNumber
}

// Returns the account balance in cents of the user
func (u User) GetAccountBalanceInCents() int64 {
	return int64(u.AccountBalanceInCents)
}

// Returns the user's overdraft account balance in cents
func (u User) GetOverdraftBalanceInCents() int64 {
	return int64(u.OverdraftBalanceInCents)
}

// Returns the security password of the user
func (u User) GetPassword() string {
	return u.Password
}

// Returns the security PIN code of the user
func (u User) GetPinCode() *string {
	return u.PinCode
}

// GetQrCodeUUID returns the UUID stored in the QrCode field.
func (u User) GetUUID() string {
	return u.QrCode.String()
}

// GetQrCodeBase64 returns the Base64 string of the QR code image.
// It generates a QR code image from the user's UUID and encodes it.
func (u User) GetQrCodeBase64() string {
	// Convert the UUID to a string
	uuidStr := u.QrCode.String()

	// Generate the QR code image as a PNG byte slice
	// The size is set to 256 for a clear image.
	png, err := qrcode.Encode(uuidStr, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error generating QR code for UUID %s: %v", uuidStr, err)
		return "" // Return an empty string or handle the error as needed
	}

	// Encode the PNG byte slice to a Base64 string
	base64Str := base64.StdEncoding.EncodeToString(png)
	return base64Str
}

// GetFacialEmbeddings returns a pointer to a list of the user's facial embeddings
func (u User) GetFacialEmbeddings() *[]string {
	if u.FacialEmbeddings == nil {
		return nil
	}
	// Cast pq.StringArray back to standard []string
	embeddingsSlice := []string(u.FacialEmbeddings)
	return &embeddingsSlice
}

// GetFacialImages returns a pointer to a list of the user's facial images
func (u User) GetFacialImages() *[]string {
	if u.FacialImages == nil {
		return nil
	}
	// Cast pq.StringArray back to standard []string
	imagesSlice := []string(u.FacialImages)
	return &imagesSlice
}

// GetPassport returns the pointer to the passport string
func (u User) GetPassport() *string {
	return u.Passport
}

//returns true if the user has facial embeddings, otherwise false.
func (u User) GetFaceEnrollmentStatus()bool{
return len(u.FacialEmbeddings) > 0
}

// GetPinStatus returns true if the PIN code is populated.
func (u User) GetPinStatus() bool {
    if u.PinCode == nil {
        return false
    }
    return *u.PinCode != ""
}

// MatchFace takes an incoming base64 embedding and matches it against the user's saved embeddings.
// It returns true immediately if any saved embedding meets or exceeds the required threshold.
func (u User) MatchFace(incomingB64 string, threshold float32) bool {
	if len(u.FacialEmbeddings) == 0 {
		return false
	}

	// 1. Convert the incoming base64 embedding into a float32 slice
	incomingVector, err := base64ToFloat32Slice(incomingB64)
	if err != nil {
		return false // Or handle error gracefully based on your logging strategy
	}

	// 2. Loop through the user's stored facial embeddings
	for _, storedB64 := range u.FacialEmbeddings {
		storedVector, err := base64ToFloat32Slice(storedB64)
		if err != nil {
			continue // Skip corrupted or invalid embeddings if they exist
		}

		// 3. Calculate cosine similarity between the two vectors
		similarity := cosineSimilarity(incomingVector, storedVector)
        log.Println("similarity score: ", similarity)
		// 4. Early Exit: If it clears the threshold, we have a match!
		if similarity >= threshold {
			return true
		}
	}

	// None of the embeddings met the threshold
	return false
}

// --- Math and Decoding Helpers ---

func base64ToFloat32Slice(base64Str string) ([]float32, error) {
	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	floats := make([]float32, len(bytes)/4)
	for i := range floats {
		bits := binary.LittleEndian.Uint32(bytes[i*4 : (i+1)*4])
		floats[i] = math.Float32frombits(bits)
	}
	return floats, nil
}

func cosineSimilarity(vecA, vecB []float32) float32 {
	if len(vecA) != len(vecB) || len(vecA) == 0 {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := 0; i < len(vecA); i++ {
		dotProduct += vecA[i] * vecB[i]
		normA += vecA[i] * vecA[i]
		normB += vecB[i] * vecB[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// UnverifiedUser represents the unverified customer model for the system.
type UnverifiedUser struct {
	gorm.Model
	UserName                string         `gorm:"uniqueIndex;not null"`  // Unique Username of the unverified user
	PhoneNumber             string         `gorm:"uniqueIndex;not null"`  // Unique Phone number of the unverified user
	Password                string         `gorm:"not null"`              // The user's password. It should be stored as a secure hash.
	AccountBalanceInCents   uint           `gorm:"not null;default:0"`    // The user's account balance in cents
	OverdraftBalanceInCents uint           `gorm:"not null;default:0"`    // The user's overdraft account balance in cents
	PinCode                 *string        `gorm:""`                      // The user's security PIN code (can be NULL). Must be stored as a secure hash.
	QrCode                  uuid.UUID      `gorm:"uniqueIndex;type:uuid"` // The user's unique Qr Code
	FacialEmbeddings        pq.StringArray `gorm:"type:text[]"`           // Stores array of base64 user facial embeddings
	FacialImages            pq.StringArray `gorm:"type:text[]"`           // Stores array of the user's facial image URLs
	Passport                *string        `gorm:""`                      // Stores URL of the user's passport photo
}

// Returns the creation timestamp of the unverified user
func (u UnverifiedUser) GetCreatedAt() time.Time {
	return u.CreatedAt.UTC()
}

// Returns the update timestamp of the unverified user
func (u UnverifiedUser) GetUpdatedAt() time.Time {
	return u.UpdatedAt.UTC()
}

// Returns the deletion timestamp of the unverified user
func (u UnverifiedUser) GetDeletedAt() time.Time {
	return u.DeletedAt.Time.UTC()
}

// Returns the username of the unverified user
func (u UnverifiedUser) GetUserName() string {
	return u.UserName
}

// Returns the phone number of the unverified user
func (u UnverifiedUser) GetPhoneNumber() string {
	return u.PhoneNumber
}

// Returns the account balance in cents of the unverified user
func (u UnverifiedUser) GetAccountBalanceInCents() int64 {
	return int64(u.AccountBalanceInCents)
}

// Returns the unverified user's overdraft account balance in cents
func (u UnverifiedUser) GetOverdraftBalanceInCents() int64 {
	return int64(u.OverdraftBalanceInCents)
}

// Returns the security password of the unverified user
func (u UnverifiedUser) GetPassword() string {
	return u.Password
}

// Returns the security PIN code of the unverified user
func (u UnverifiedUser) GetPinCode() *string {
	return u.PinCode
}

// GetQrCodeUUID returns the UUID stored in the QrCode field.
func (u UnverifiedUser) GetUUID() string {
	return u.QrCode.String()
}

// GetQrCodeBase64 returns the Base64 string of the QR code image.
// It generates a QR code image from the user's UUID and encodes it.
func (u UnverifiedUser) GetQrCodeBase64() string {
	uuidStr := u.QrCode.String()

	png, err := qrcode.Encode(uuidStr, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error generating QR code for UUID %s: %v", uuidStr, err)
		return ""
	}

	base64Str := base64.StdEncoding.EncodeToString(png)
	return base64Str
}

// GetPinStatus returns true if the PIN code is populated.
func (u UnverifiedUser) GetPinStatus() bool {
    if u.PinCode == nil {
        return false
    }
    return *u.PinCode != ""
}

// GetFacialEmbeddings returns a pointer to a list of the unverified user's facial embeddings
func (u UnverifiedUser) GetFacialEmbeddings() *[]string {
	if u.FacialEmbeddings == nil {
		return nil
	}
	embeddingsSlice := []string(u.FacialEmbeddings)
	return &embeddingsSlice
}

// GetFacialImages returns a pointer to a list of the unverified user's facial images
func (u UnverifiedUser) GetFacialImages() *[]string {
	if u.FacialImages == nil {
		return nil
	}
	imagesSlice := []string(u.FacialImages)
	return &imagesSlice
}

// GetPassport returns the pointer to the passport string
func (u UnverifiedUser) GetPassport() *string {
	return u.Passport
}

// MatchFace takes an incoming base64 embedding and matches it against the unverified user's saved embeddings.
// It returns true immediately if any saved embedding meets or exceeds the required threshold.
func (u UnverifiedUser) MatchFace(incomingB64 string, threshold float32) bool {
	if len(u.FacialEmbeddings) == 0 {
		return false
	}

	incomingVector, err := base64ToFloat32Slice(incomingB64)
	if err != nil {
		return false
	}

	for _, storedB64 := range u.FacialEmbeddings {
		storedVector, err := base64ToFloat32Slice(storedB64)
		if err != nil {
			continue
		}

		similarity := cosineSimilarity(incomingVector, storedVector)

		if similarity >= threshold {
			return true
		}
	}

	return false
}

//returns true if the user has facial embeddings, otherwise false.
func (u UnverifiedUser) GetFaceEnrollmentStatus()bool{
return false
}

// UniquenessCheck represents the availability status of a unique identifier (like phone or username),
// indicating whether it is already taken and if the existing account is verified.
type UniquenessCheck struct {
	Exists     bool `json:"exists"`      // True if the identifier is found in either verified or unverified tables.
	IsVerified bool `json:"is_verified"` // True if the existing record belongs to a fully verified user.
}
