package device

import (
	"github.com/Lcat-Technologies/eardrum-interfaces/device"
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"gorm.io/gorm"
	pgerror"errors"
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