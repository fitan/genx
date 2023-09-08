package types

import "gorm.io/gorm"

type Maintenance struct {
	gorm.Model
	// 来源id
	SourceID string `gorm:"column:source_id;notnull;comment:'来源id'" json:"sourceId"`

	// 维护类型 1:日常维护 2:故障维护
	MaintenanceType string `gorm:"column:maintenance_type;notnull;comment:'维护类型'" json:"maintenanceType"`

	// 维护标题
	Title string `gorm:"column:title;notnull;comment:'维护标题'" json:"title"`
	// 维护信息
	MaintenanceInfo string `gorm:"column:maintenance_info;notnull;comment:'维护信息'" json:"maintenanceInfo"`
	// 维护人员
	MaintenanceUser string `gorm:"column:maintenance_user;notnull;comment:'维护人员'" json:"maintenanceUser"`
}

func (Maintenance) TableName() string {
	return "assets_maintenance"
}
