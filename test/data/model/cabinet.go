package types

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

// 机柜
type Cabinet struct {
	gorm.Model
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 机柜名称
	Name string `gorm:"column:name;notnull;comment:'名称'" json:"name"`
	// 机柜编号
	Code string `gorm:"column:code;notnull;comment:'编号'" json:"code"`
	// 机柜类型
	Type int `gorm:"column:type;notnull;default:0;comment:'类型'" json:"type"`
	// 机柜大小 多少USIZE
	Size int `gorm:"column:size;notnull;default:0;comment:'大小'" json:"size"`
	// 第几行
	Row int `gorm:"column:row;notnull;default:0;comment:'第几行'" json:"row"`
	// 第几列
	Column int `gorm:"column:column;notnull;default:0;comment:'第几列'" json:"column"`
	// 机柜状态 1:使用中 2：关闭
	Status int `gorm:"column:status;notnull;default:1;comment:'状态'" json:"status"`
	// 管理人邮箱
	ManagerEmail []string `gorm:"column:manager_email;serializer:json;notnull;comment:'管理人邮箱'" json:"managerEmail"`
	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	// 机房id
	ComputerRoomID uint `gorm:"column:computer_room_id;notnull;comment:'机房ID'" json:"computerRoomId"`

	// 机房
	ComputerRoom ComputerRoom `gorm:"foreignKey:ComputerRoomID" json:"computerRoom"`
	// 设备
	Devices []Device `gorm:"foreignKey:CabinetID" json:"device"`
}

func (cabinet *Cabinet) BeforeCreate(tx *gorm.DB) (err error) {
	err = tx.Model(&Cabinet{}).
		Where(map[string]interface{}{
			"computer_room_id": cabinet.ComputerRoomID,
			"row":              cabinet.Row,
			"column":           cabinet.Column,
		}).Where("uuid != ?", cabinet.UUID).
		First(&Cabinet{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if cabinet.UUID == "" {
			cabinet.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
		}
		return nil
	}

	if err == nil {
		err = errors.Errorf("机柜已存在 row:%d column:%d", cabinet.Row, cabinet.Column)
		return
	}

	return
}

func (Cabinet) TableName() string {
	return "assets_cabinet"
}
