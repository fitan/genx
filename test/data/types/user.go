package types

import (
	"time"

	"gorm.io/gorm"
)

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

type ListResponseDTO struct {
	UUID string `json:"uuid"`
	// 设备UUID
	DeviceUUID string `json:"deviceUuid"`
	// 品牌信息
	Brand Brand `json:"brand"`

	// cpu品牌信息
	CpuBrand Brand `json:"cpuBrand"`

	// 内存品牌信息
	MemoryBrand Brand `json:"memoryBrand"`
	// 内存数量
	MemoryCount int `json:"memoryCount"`

	// 硬盘品牌信息
	DiskBrand Brand `json:"diskBrand"`
	// 硬盘数量
	DiskCount int `json:"diskCount"`

	// 网卡品牌信息
	NetworkCardBrand Brand `json:"networkCardBrand"`
	// 网卡数量
	NetworkCardCount int `json:"networkCardCount"`

	// 系统品牌信息
	OsBrand Brand `json:"osBrand"`

	// 供应商信息
	Supplier Supplier `json:"supplier"`

	// 业务ip
	BusinessIp string `json:"businessIp"`
	// 业务子网掩码
	BusinessSubnetMask string `json:"businessSubnetMask"`
	// nasIp
	NasIp string `json:"nasIp"`
	// nas子网掩码
	NasSubnetMask string `json:"nasSubnetMask"`
	// IPMIIP
	IpmiIp string `json:"ipmiIp"`
	// ipmi子网掩码
	IpmiSubnetMask string `json:"ipmiSubnetMask"`
	// sn
	Sn string `json:"sn"`

	// 跳板机登录端口
	JumpPort int `json:"jumpPort"`
	// 环境  1:生产 2:测试 3:开发
	Env int `json:"env"`
	// 状态  1:使用中 2：故障 3：关机
	Status int `json:"status"`

	// 维护状态 1:正常 2:维护中
	MaintainStatus int `json:"maintainStatus"`

	// 备注
	Remark string `json:"remark"`

	// namespace 项目名称
	Namespace string `json:"namespace"`
	// namespaceAlias 项目别名
	NamespaceAlias string `json:"namespaceAlias"`
	// name 服务名称
	Name string `json:"name"`
	// nameAlias 服务别名
	NameAlias string `json:"nameAlias"`

	// 维保起始时间
	MaintenanceStartTime string `json:"maintenanceStartTime"`

	// 维保到期时间
	MaintenanceExpireTime string `json:"maintenanceExpireTime"`

	// 负责人
	Owners []string `json:"owners"`

	// 机房名字
	ComputerRoomName string `json:"computerRoomName"`
	ComputerRoomUUID string `json:"computerRoomUuid"`
	// 机柜名字
	CabinetName string `json:"cabinetName"`

	// 机柜UUID
	CabinetUUID string `json:"cabinetUuid" valid:"required"`
	// 起始层数
	StartLayer int `json:"startLayer" valid:"required"`
	// 终止层数
	EndLayer int `json:"endLayer" valid:"required"`

	PhysicalMachineName string `json:"physicalMachineName"`

	CpuCount int64 `json:"cpuCount"`

	MemorySize int64 `json:"memorySize"`

	DiskSize int64 `json:"diskSize"`

	// 去连接设备
	//ToDevices []Device `json:"toDevice"`
	//连接我的设备
	//FromDevices []Device `json:"fromDevice"`
}

type Supplier struct {
	// uuid
	UUID string `json:"uuid"`
	// 名字
	Name string `json:"name"`
	// 联系方式
	Contact string `json:"contact"`
	// 地址
	Address string `json:"address"`
	// 备注
	Remark string `json:"remark"`
}
