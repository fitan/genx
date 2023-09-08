package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type DevicePort struct {
	gorm.Model
	UUID           string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	DeviceID       string `gorm:"column:device_id;notnull;comment:'设备id'" json:"deviceId"`
	DevicePortName string `gorm:"column:device_port_name;notnull;comment:'设备端口名称'" json:"devicePortName"`

	ToDeviceID     string      `gorm:"column:to_device_id;notnull;comment:'连接至设备id'" json:"toDeviceId"`
	ToDevicePortID string      `gorm:"column:to_device_port_id;notnull;comment:'对端设备端口id'" json:"toDevicePortId"`
	ToDevicePort   *DevicePort `gorm:"foreignKey:ToDevicePortID;references:UUID" json:"toDevicePort"`

	FromDevicePort *DevicePort `gorm:"foreignKey:UUID;references:ToDevicePortID" json:"fromDevicePort"`

	Device Device `gorm:"foreignKey:DeviceID;references:DeviceID" json:"sourceDevice"`
}

func (DevicePort) TableName() string {
	return "assets_device_port"
}

func (devicePort *DevicePort) BeforeCreate(tx *gorm.DB) (err error) {
	if devicePort.UUID == "" {
		devicePort.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}
