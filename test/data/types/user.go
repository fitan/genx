package types

import (
	"gorm.io/gorm"
	"time"
)

// grom User表
type User struct {
	// ID
	ID int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	// 用户名
	Username string `json:"username" gorm:"column:username;type:varchar(32);not null;default:'';comment:用户名"`
	// 密码
	Password string `json:"password" gorm:"column:password;type:varchar(32);not null;default:'';comment:密码"`
	// 昵称
	Nickname string `json:"nickname" gorm:"column:nickname;type:varchar(32);not null;default:'';comment:昵称"`
	// 头像
	Avatar string `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像"`
	// 朋友
	Friends []Friend `json:"friends" gorm:"foreignKey:UserID;references:ID;comment:朋友"`
}

// gorm 朋友表
type Friend struct {
	// ID
	ID int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	// 名字
	Name string `json:"name" gorm:"column:name;type:varchar(32);not null;default:'';comment:名字"`
	// 用户ID
	UserID int64 `json:"user_id" gorm:"column:user_id;type:bigint(20);not null;default:0;comment:用户ID"`
	// 朋友ID
	FriendID int64 `json:"friend_id" gorm:"column:friend_id;type:bigint(20);not null;default:0;comment:朋友ID"`
	// 备注
	Remark string `json:"remark" gorm:"column:remark;type:varchar(32);not null;default:'';comment:备注"`
	// 分组
	Group string `json:"group" gorm:"column:group;type:varchar(32);not null;default:'';comment:分组"`
	// 创建时间
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;type:bigint(20);not null;default:0;comment:创建时间"`
	// 更新时间
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;type:bigint(20);not null;default:0;comment:更新时间"`
}

type PhysicalMachine struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	//品牌信息
	//Brand string `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`
	//型号
	//BrandModel string `gorm:"column:brand_model;notnull;comment:'型号'" json:"brandModel"`

	// 是否采集完成

	BrandUUID string `gorm:"column:brand_uuid;notnull;comment:'品牌uuid'" json:"brandUUID"`
	Brand     Brand  `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`
}

func (PhysicalMachine) TableName() string {
	return "assets_physical_machine"
}

type Brand struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 品牌
	Brand string `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`

	// 产品类型 1:服务器 2:交换机 3:路由器 4:防火墙 5:负载均衡器 6:存储设备 7: cpu 8:内存 9:硬盘 10:网卡 10: 系统
	ProductType string `gorm:"column:product_type;notnull;comment:'产品类型'" json:"productType"`

	// 产品型号
	ProductModel string `gorm:"column:product_model;notnull;comment:'产品型号'" json:"productModel"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	Users []User `gorm:"foreignKey:UUID;references:UUID" json:"users"`
}

type User struct {
	gorm.Model
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	Name string `gorm:"column:name;notnull;comment:'姓名'" json:"name"`
}

func (b *Brand) TableName() string {
	return "assets_brand"
}
