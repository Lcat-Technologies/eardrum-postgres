package device

import (
	"github.com/Lcat-Technologies/eardrum-interfaces/device"
	"github.com/Lcat-Technologies/eardrum-interfaces/errors"
	"gorm.io/gorm"
)

// CreateDevice creates a new device record
func CreateDevice(s device.NewDevice, Db *gorm.DB) (device.Device, error) {

	// 3. Create device data
	d := &Device{
		DeviceID:              s.GetDeviceID(),
		DeviceModel:           s.GetDeviceModel(),
	}

	// 4. Create a device record in the database
	if err := Db.Create(d).Error; err != nil {

		// All other persistence failures -> 500 Internal Server Error
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
		return nil, err1
	}

	return d, nil
}