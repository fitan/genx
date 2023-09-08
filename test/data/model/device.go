package types

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	// 设备UUID
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 机柜ID
	CabinetID uint `gorm:"column:cabinet_id;notnull;comment:'机柜id'" json:"cabinetId"`
	// 设备id
	DeviceID string `gorm:"column:device_id;notnull;comment:'设备id'" json:"deviceId"`
	// 设备类型
	DeviceType string `gorm:"column:device_type;notnull;comment:'设备类型'" json:"deviceType"`
	// 起始层
	StartLayer int `gorm:"column:start_layer;notnull;comment:'起始层'" json:"startLayer"`
	// 终止层
	EndLayer int `gorm:"column:end_layer;notnull;comment:'终止层'" json:"endLayer"`
	// 网线连接设备
	NetworkToDevices []DevicePort `gorm:"foreignKey:DeviceID;references:DeviceID" json:"networkToDevices"`

	// 别的设备连接至此设备
	NetworkFromDevices []DevicePort `gorm:"foreignKey:ToDeviceID;references:DeviceID" json:"networkFromDevices"`
	// 关联的机柜
	Cabinet Cabinet `gorm:"foreignKey:CabinetID" json:"cabinet"`
}
