package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

var NetworkCardPortNames = []string{
	"nasPort1",
	"nasPort2",
	"businessPort1",
	"businessPort2",
	"managementPort1",
}

type PhysicalMachine struct {
	gorm.Model
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	//品牌信息
	//Brand string `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`
	//型号
	//BrandModel string `gorm:"column:brand_model;notnull;comment:'型号'" json:"brandModel"`

	// 是否采集完成
	IsCollectFinish bool `gorm:"DEFAULT false; comment:'是否采集完成';column:is_collect_finish" json:"isCollectFinish"`
	// 采集状态 1: 完成过第一次采集
	CollectStatus int `gorm:"DEFAULT 0; comment:'采集状态';column:collect_status" json:"collectStatus"`

	BrandUUID string `gorm:"column:brand_uuid;notnull;comment:'品牌uuid'" json:"brandUUID"`
	Brand     Brand  `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`

	// 物理机名字
	PhysicalMachineName string `gorm:"column:physical_machine_name;notnull;comment:'物理机名字'" json:"physicalMachineName"`

	// 业务ip
	BusinessIp string `gorm:"column:business_ip;notnull;comment:'业务ip'" json:"businessIp"`
	// 业务子网掩码
	BusinessSubnetMask string `gorm:"column:business_subnet_mask;notnull;comment:'业务子网掩码'" json:"businessSubnetMask"`
	// nasIp
	NasIp string `gorm:"column:nas_ip;notnull;comment:'nasIp'" json:"nasIp"`
	// nas子网掩码
	NasSubnetMask string `gorm:"column:nas_subnet_mask;notnull;comment:'nas子网掩码'" json:"nasSubnetMask"`
	// IPMIIP
	IpmiIp string `gorm:"column:ipmi_ip;notnull;comment:'ipmiIp'" json:"ipmiIp"`
	// ipmi子网掩码
	IpmiSubnetMask string `gorm:"column:ipmi_subnet_mask;notnull;comment:'ipmi子网掩码'" json:"ipmiSubnetMask"`
	// sn
	Sn string `gorm:"column:sn;notnull;comment:'sn'" json:"sn"`

	// cpu品牌uuid
	//CpuBrandUUID string `gorm:"column:cpu_brand_uuid;notnull;comment:'cpu品牌uuid'" json:"cpuBrandUUID"`
	//CpuBrand     Brand  `gorm:"foreignKey:CpuBrandUUID;references:UUID" json:"cpuBrand"`
	CpuCount int64 `gorm:"column:cpu_count;notnull;comment:'cpu个数'" json:"cpuCount"`

	MemorySize int64 `gorm:"column:memory_size;comment:'内存大小'" json:"memorySize"`

	DiskSize int64 `gorm:"column:disk_size;comment:'硬盘大小'" json:"diskSize"`

	// 网卡品牌UUID
	//NetworkCardBrandUUID string `gorm:"column:network_card_brand_uuid;notnull;comment:'网卡品牌uuid'" json:"networkCardBrandUUID"`
	//NetworkCardBrand     Brand  `gorm:"foreignKey:NetworkCardBrandUUID;references:UUID" json:"networkCardBrand"`
	//NetworkCardCount     int    `gorm:"column:network_card_count;notnull;comment:'网卡数量'" json:"networkCardCount"`

	//MemorySize int `gorm:"column:memory_size;notnull;comment:'内存大小'" json:"memorySize"`
	// 内存数量
	// 硬盘大小
	//DiskSize int `gorm:"column:disk_size;notnull;comment:'硬盘大小'" json:"diskSize"`
	// 硬盘数量
	//DiskCount int `gorm:"column:disk_count;notnull;comment:'硬盘数量'" json:"diskCount"`
	// 网卡数量
	//NetworkCardCount int `gorm:"column:network_card_count;notnull;comment:'网卡数量'" json:"networkCardCount"`
	// 操作系统品牌UUID
	OsBrandUUID string `gorm:"column:os_brand_uuid;notnull;comment:'操作系统品牌uuid'" json:"osBrandUUID"`
	OsBrand     Brand  `gorm:"foreignKey:OsBrandUUID;references:UUID" json:"osBrand"`

	// 供应商UUID
	SupplierUUID string   `gorm:"column:supplier_uuid;notnull;comment:'供应商uuid'" json:"supplierUUID"`
	Supplier     Supplier `gorm:"foreignKey:SupplierUUID;references:UUID" json:"supplier"`

	// 跳板机登录端口
	JumpPort int `gorm:"column:jump_port;notnull;default:0;comment:'跳板机登录端口'" json:"jumpPort"`
	// 环境  1:生产 2:测试 3:开发
	Env int `gorm:"column:env;notnull;default:1;comment:'环境'" json:"env"`
	// 状态  1:开机 2：故障 3：关机 4: 进入库房
	Status int `gorm:"column:status;notnull;default:1;comment:'状态'" json:"status"`

	// 维护状态 1:正常 2:维护中
	MaintainStatus int `gorm:"column:maintain_status;notnull;default:1;comment:'维护状态'" json:"maintainStatus"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	NamespaceID int `gorm:"column:namespace_id;notnull;comment:'项目ID'" json:"namespaceId"`

	NameID int `gorm:"column:name_id;notnull;comment:'服务名称'" json:"nameId"`

	// 维保起始时间
	MaintenanceStartTime string `gorm:"column:maintenance_start_time;notnull;comment:'维保起始时间'" json:"maintenanceStartTime"`

	// 维保到期时间
	MaintenanceExpireTime string `gorm:"column:expire_time;notnull;comment:'维保到期时间'" json:"maintenanceExpireTime"`

	// 负责人
	Owners []string `gorm:"column:owners;serializer:json;notnull;comment:'负责人'" json:"owners"`

	// 物理核心数
	PhysicalCore int `gorm:"column:physical_core;comment:'物理核心数'" json:"physicalCore"`

	// 内存槽位数
	MemorySlots int `gorm:"column:memory_slots;comment:'内存槽位数'" json:"memorySlots"`

	// 几U
	Height int `gorm:"column:height;comment:'几U'" json:"height"`

	//盘位信息
	//DiskSlots map[string]int `gorm:"column:disk_slots;serializer:json;comment:'盘位信息'" json:"diskSlots"`

	// 2.5寸盘位个数
	DiskSlots25 int `json:"diskSlots25" csv:"diskSlots25"`
	// 3.5寸盘位个数
	DiskSlots35 int `json:"diskSlots35" csv:"diskSlots35"`

	// 关联的设备
	Device Device `gorm:"polymorphic:Device;foreignKey:UUID" json:"device"`

	// 备件关联
	SparePart SparePart `gorm:"polymorphic:Owner;foreignKey:UUID" json:"sparePart"`

	// 维护信息
	Maintenances []Maintenance `gorm:"foreignKey:SourceID;references:UUID" json:"maintain"`
	// 采集日志
	Collects []Collect `gorm:"foreignKey:SourceID;references:UUID" json:"collect"`
}

type PmChangeLog struct {
	gorm.Model
	// 物理机UUID
	UUID string `gorm:"column:uuid;notnull;comment:'物理机UUID'" json:"physicalMachineUUID"`
	// sn
	Sn string `gorm:"column:sn;comment:'sn'" json:"sn"`
	// 业务ip
	BusinessIp string `gorm:"column:business_ip;comment:'业务ip'" json:"businessIp"`
	// 业务子网掩码
	BusinessSubnetMask string `gorm:"column:business_subnet_mask;comment:'业务子网掩码'" json:"businessSubnetMask"`
	// nasIp
	NasIp string `gorm:"column:nas_ip;comment:'nasIp'" json:"nasIp"`
	// nas子网掩码
	NasSubnetMask string `gorm:"column:nas_subnet_mask;comment:'nas子网掩码'" json:"nasSubnetMask"`
	// IPMIIP
	IpmiIp string `gorm:"column:ipmi_ip;comment:'ipmiIp'" json:"ipmiIp"`
	// ipmi子网掩码
	IpmiSubnetMask string `gorm:"column:ipmi_subnet_mask;comment:'ipmi子网掩码'" json:"ipmiSubnetMask"`
	// 跳板机登录端口
	JumpPort int `gorm:"column:jump_port;default:0;comment:'跳板机登录端口'" json:"jumpPort"`
	// 环境  1:生产 2:测试 3:开发
	Env int `gorm:"column:env;default:1;comment:'环境'" json:"env"`
	// 状态  1:使用中 2：故障 3：关机 4: 进入库房
	Status int `gorm:"column:status;default:1;comment:'状态'" json:"status"`
	// 机房名字
	RoomName string `gorm:"column:room_name;comment:'机房名字'" json:"roomName"`
	// 机柜名字
	CabinetName string `gorm:"column:cabinet_name;comment:'机柜名字'" json:"cabinetName"`
	// 机柜起始层
	CabinetStart int `gorm:"column:cabinet_start;comment:'机柜起始层'" json:"cabinetStart"`
	// 机柜结束层
	CabinetEnd int `gorm:"column:cabinet_end;comment:'机柜结束层'" json:"cabinetEnd"`
	// 操作人
	Operator string `gorm:"column:operator;comment:'操作人'" json:"operator"`
}

func (PmChangeLog) TableName() string {
	return "assets_pm_change_log"
}

type Collect struct {
	gorm.Model

	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	// 采集信息
	Msg string `gorm:"column:msg;comment:'采集信息'" json:"msg"`
	// 采集状态 1:成功 2:失败
	Status int `gorm:"column:status;comment:'采集状态'" json:"status"`
	// 采集来源
	SourceID string `gorm:"column:source_id;comment:'采集来源'" json:"sourceId"`
}

func (Collect) TableName() string {
	return "assets_collect"
}

func (collect *Collect) BeforeCreate(tx *gorm.DB) (err error) {
	if collect.UUID == "" {
		collect.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

func (physicalMachine *PhysicalMachine) BeforeCreate(tx *gorm.DB) (err error) {
	if physicalMachine.UUID == "" {
		physicalMachine.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

//func (PhysicalMachine *PhysicalMachine) BeforeUpdate(tx *gorm.DB) (err error) {
//	if PhysicalMachine.Status == 4 {
//		tx.
//	}
//}

func (PhysicalMachine) TableName() string {
	return "assets_physical_machine"
}
