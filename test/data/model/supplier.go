package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

// 供应商
type Supplier struct {
	gorm.Model
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 名字
	Name string `gorm:"column:name;notnull;comment:'名字'" json:"name"`
	// 联系方式
	Contact string `gorm:"column:contact;notnull;comment:'联系方式'" json:"contact"`
	// 地址
	Address string `gorm:"column:address;notnull;comment:'地址'" json:"address"`
	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`
}

func (s *Supplier) BeforeCreate(tx *gorm.DB) (err error) {
	if s.UUID == "" {
		s.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

func (Supplier) TableName() string {
	return "assets_supplier"
}
