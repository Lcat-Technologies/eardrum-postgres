package device

import (
	pgerror "errors"

	"github.com/Lcat-Technologies/eardrum-interfaces/device"
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"gorm.io/gorm"
)

// GetDeviceByID retrieves a device from the database by its device ID.
// Returns an empty Device struct if not found, or an error for internal failures.
func GetDeviceByID(deviceID string, Db *gorm.DB) (device.Device, error) {
	var d Device

	// 1. Query the device record by device ID
	if err := Db.Where("device_id = ?", deviceID).First(&d).Error; err != nil {

		// 2. If the record is not found, return an empty struct without an error
		if pgerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		// 3. All other persistence failures -> Internal Server Error
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}

	return &d, nil
}

// GetDevicesByIDs retrieves multiple devices from the database by a slice of device IDs.
func GetDevicesByIDs(deviceIDs []string, Db *gorm.DB) ([]device.Device, error) {
	if len(deviceIDs) == 0 {
		return nil, nil
	}

	var dbDevices []*Device

	// Use `IN ?` to fetch all matching devices in a single database round-trip
	if err := Db.Where("device_id IN ?", deviceIDs).Find(&dbDevices).Error; err != nil {
		// All other persistence failures -> Internal Server Error
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}

	// Convert concrete slice to interface slice if your interface requires it
	result := make([]device.Device, len(dbDevices))
	for i, d := range dbDevices {
		result[i] = d
	}

	return result, nil
}
