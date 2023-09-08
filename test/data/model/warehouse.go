package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type Warehouse struct {
	gorm.Model
	// 库房UUID
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	// 库房类型
	Type string `gorm:"column:type;notnull;comment:'类型'" json:"type"`

	// 库房名称
	Name string `gorm:"column:name;notnull;comment:'名称'" json:"name"`

	// 库房地址
	Address string `gorm:"column:address;notnull;comment:'地址'" json:"address"`

	// 维护人邮箱
	ManagerEmail []string `gorm:"column:manager_email;serializer:json;notnull;comment:'维护人邮箱'" json:"managerEmail"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	// 备件
	SparePart []SparePart `gorm:"gorm:foreignKey:WarehouseUUID;references:UUID" json:"sparePart"`

	// 设备类型
	//WarehouseID   uint   `gorm:"column:warehouse_id;notnull;comment:'设备ID'" json:"warehouseId"`
	//WarehouseType string `gorm:"column:warehouse_type;notnull;comment:'设备类型'" json:"warehouseType"`
	// 机房id
	//ComputerRoomID uint `gorm:"column:computer_room_id;notnull;comment:'机房ID'" json:"computerRoomId"`
	// 机房
	//ComputerRoom ComputerRoom `gorm:"foreignKey:ComputerRoomID" json:"computerRoom"`
}

func (w *Warehouse) BeforeCreate(tx *gorm.DB) (err error) {
	if w.UUID == "" {
		w.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

func (Warehouse) TableName() string {
	return "assets_warehouse"
}

type SparePart struct {
	gorm.Model

	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	WarehouseUUID string `gorm:"column:warehouse_uuid;notnull;comment:'仓库uuid'" json:"warehouseUuid"`

	OwnerID   string
	OwnerType string

	BrandUUID    string
	SupplierUUID string
}

func (s *SparePart) BeforeCreate(tx *gorm.DB) (err error) {
	if s.UUID == "" {
		s.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

func (SparePart) TableName() string {
	return "assets_spare_part"
}
