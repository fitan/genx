package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

// 1 数据中心 2 职场机房
// @enum("dataCenter:dataCenter", "workplaceComputerRoom:workplaceComputerRoom")
type ComputerRoomType int

// 机房
type ComputerRoom struct {
	gorm.Model
	// 唯一定位符
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 机房类型
	Type ComputerRoomType `gorm:"column:type;notnull;comment:'机房类型'" json:"type"`
	// 机房名称
	Name string `gorm:"column:name;notnull;comment:'名称'" json:"name"`
	// 机房状态 1:使用中 2：关闭
	Status int `gorm:"column:status;notnull;default:1;comment:'状态'" json:"status"`
	// 城市
	City string `gorm:"column:city;notnull;comment:'城市'" json:"city"`
	// 地点
	Location string `gorm:"column:location;notnull;comment:'地点'" json:"location"`
	// 管理人邮箱
	ManagerEmail []string `gorm:"column:manager_email;serializer:json;notnull;comment:'管理人邮箱'" json:"managerEmail"`
	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	// 包含的机柜
	Cabinets []Cabinet `gorm:"foreignKey:ComputerRoomID;references:ID" json:"cabinet"`
}

func (computerRoom *ComputerRoom) BeforeCreate(tx *gorm.DB) (err error) {
	if computerRoom.UUID == "" {
		computerRoom.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

func (ComputerRoom) TableName() string {
	return "assets_computer_room"
}
