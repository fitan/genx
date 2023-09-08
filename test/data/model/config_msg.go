package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type Brand struct {
	gorm.Model
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 品牌
	Brand string `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`

	// 产品类型 1:服务器 2:交换机 3:路由器 4:防火墙 5:负载均衡器 6:存储设备 7: cpu 8:内存 9:硬盘 10:网卡 10: 系统
	ProductType string `gorm:"column:product_type;notnull;comment:'产品类型'" json:"productType"`

	// 产品型号
	ProductModel string `gorm:"column:product_model;notnull;comment:'产品型号'" json:"productModel"`

	// 产品参数
	ProductParam ProductParam `gorm:"column:product_param;serializer:json;default:'{}';comment:'产品参数'" json:"productParam"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`
}

func (b *Brand) BeforeCreate(tx *gorm.DB) (err error) {
	if b.UUID == "" {
		b.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

func (b *Brand) TableName() string {
	return "assets_brand"
}

type ProductParam struct {
	Cpu *CpuParam `json:"cpu,omitempty"`

	Memory *MemoryParam `json:"memory,omitempty"`

	HardDisk *HardDiskParam `json:"hardDisk,omitempty"`
}

type CpuParam struct {
	// 核心数
	CoreCount string `json:"coreCount"`
	// 核心频率
	CoreFrequency int64 `json:"coreFrequency"`
}

type MemoryParam struct {
	// 频率
	Frequency int64 `json:"frequency"`
	// 容量
	Capacity int64 `json:"capacity"`
}

type HardDiskParam struct {
	// 容量
	Capacity int64 `json:"capacity"`
	// 类型 1 ssd 2 hdd
	Type string `json:"type"`
}
