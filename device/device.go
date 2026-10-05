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

// Returns the unique ID of the device
func (s Device) GetDeviceID() string {
	return s.DeviceID
}

// Returns the model of the device
func (s Device) GetDeviceModel() *string {
	m := s.DeviceModel
	return &m
}