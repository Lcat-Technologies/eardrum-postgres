package device

import (
	"github.com/Lcat-Technologies/eardrum-interfaces/device"
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"gorm.io/gorm"
	pgerror"errors"
)

// CreateDevice creates a new device record
func CreateDevice(s device.NewDevice, Db *gorm.DB) (device.Device, error) {

    // 1. Prepare device data
    d := Device{
        DeviceID:    s.GetDeviceID(),
        DeviceModel: s.GetDeviceModel(),
    }

    // 2. Create device record in the database
    if err := Db.Create(&d).Error; err != nil {
        // Check if the failure is due to a duplicate key (device already exists)
        if pgerror.Is(err, gorm.ErrDuplicatedKey) {
            return nil, nil
        }

        // All other persistence failures -> Internal Server Error
        err1 := errors.New(errors.EARInternalError, err)
        err1.Log()
        return nil, err1
    }

    return &d, nil
}