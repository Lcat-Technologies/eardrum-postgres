package device

import (
	"gorm.io/gorm"
)

// Device represents a device once logged into the system
type Device struct {
	gorm.Model
	DeviceID    string `gorm:"uniqueIndex;not null"` // Unique identifier of the device
	DeviceModel string `gorm:"not null"`             // The model of the device
}

// GetDeviceID returns the unique ID of the device.
// Returns an empty string if receiver is nil.
func (s *Device) GetDeviceID() string {
	if s == nil {
		return ""
	}
	return s.DeviceID
}

// GetDeviceModel returns a pointer to the model of the device.
// Returns nil if receiver is nil.
func (s *Device) GetDeviceModel() *string {
	if s == nil {
		return nil
	}
	// Copy string value to a new variable to safely return its address
	m := s.DeviceModel
	return &m
}