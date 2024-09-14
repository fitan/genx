package model

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

	NamespaceService TblServicetree `gorm:"->;foreignKey:NamespaceID;references:Pri" json:"namespaceService"`

	NameID int `gorm:"column:name_id;notnull;comment:'服务名称'" json:"nameId"`

	NameService TblServicetree `gorm:"->;foreignKey:NameID;references:Pri" json:"nameService"`

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

	// 配件关联
	Accessories []Accessories `gorm:"polymorphic:Owner;foreignKey:UUID" json:"accessories"`
	// 关联库房
	//Warehouse Warehouse `gorm:"polymorphic:Warehouse" json:"warehouse"`

	// 维护信息
	//Maintenances []Maintenance `gorm:"foreignKey:SourceID;references:UUID" json:"maintain"`
	// 采集日志
	Collects []Collect `gorm:"foreignKey:SourceID;references:UUID" json:"collect"`

	//PhysicalMachineAllocation *PhysicalMachineAllocation `gorm:"foreignKey:PhysicalMachineUUID;references:UUID" json:"physicalMachineAllocation"`
}

// @crud id
// @crud-update-method(name="updateStatus", fields=["Status","IpmiIp","Env"])
// @crud-update(name="updateOsBrandUUID", fields="OsBrandUUID")
type PmChangeLog struct {
	gorm.Model
	// 物理机UUID
	// @crud-list-op = >=
	UUID string `gorm:"column:uuid;notnull;comment:'物理机UUID'" json:"physicalMachineUUID"`
	// sn
	Sn string `gorm:"column:sn;comment:'sn'" json:"sn"`
	// 业务ip
	// @crud-where
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

type Brand struct {
	gorm.Model
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 品牌
	Brand string `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`

	// 产品类型 1:服务器 2:交换机 3:路由器 4:防火墙 5:负载均衡器 6:存储设备 7: cpu 8:内存 9:硬盘 10:网卡 10: 系统 11: gpu
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

type ProductParam struct {
	Cpu *CpuParam `json:"cpu,omitempty"`

	Memory *MemoryParam `json:"memory,omitempty"`

	HardDisk *HardDiskParam `json:"hardDisk,omitempty"`

	Gpu *GpuParam `json:"gpu,omitempty"`
}

type GpuParam struct {
	// serial name memory.total
	Capacity string `json:"capacity"`
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

func (Device) TableName() string {
	return "assets_device"
}

type DevicePort struct {
	gorm.Model
	UUID           string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	DeviceID       string `gorm:"column:device_id;notnull;comment:'设备id'" json:"deviceId"`
	DevicePortName string `gorm:"column:device_port_name;notnull;comment:'设备端口名称'" json:"devicePortName"`

	ToDeviceID     string      `gorm:"column:to_device_id;notnull;comment:'连接至设备id'" json:"toDeviceId"`
	ToDevicePortID string      `gorm:"column:to_device_port_id;notnull;comment:'对端设备端口id'" json:"toDevicePortId"`
	ToDevicePort   *DevicePort `gorm:"foreignKey:ToDevicePortID;references:UUID" json:"toDevicePort"`

	FromDevicePort *DevicePort `gorm:"foreignKey:UUID;references:ToDevicePortID" json:"fromDevicePort"`

	Device Device `gorm:"foreignKey:DeviceID;references:DeviceID" json:"sourceDevice"`
}

func (DevicePort) TableName() string {
	return "assets_device_port"
}

func (devicePort *DevicePort) BeforeCreate(tx *gorm.DB) (err error) {
	if devicePort.UUID == "" {
		devicePort.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

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

type SparePart struct {
	gorm.Model

	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	WarehouseUUID string    `gorm:"column:warehouse_uuid;notnull;comment:'仓库uuid'" json:"warehouseUuid"`
	Warehouse     Warehouse `gorm:"foreignKey:WarehouseUUID;references:UUID" json:"warehouse"`

	OwnerID   string
	OwnerType string

	BrandUUID    string
	SupplierUUID string

	// 物理机
	PhysicalMachine *PhysicalMachine `gorm:"foreignKey:UUID;references:OwnerID" json:"physicalMachine"`

	// 网络设备
	NetworkDevice *NetworkDevice `gorm:"foreignKey:UUID;references:OwnerID" json:"networkDevice"`

	// 配件
	Accessories *Accessories `gorm:"foreignKey:UUID;references:OwnerID" json:"accessories"`
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

type Accessories struct {
	gorm.Model

	// 配置UUID
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	// 采购批次
	Batch string `gorm:"column:batch;notnull;comment:'采购批次'" json:"batch"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	// sn号
	Sn string `gorm:"column:sn;notnull;comment:'sn号'" json:"sn"`

	OwnerID   string
	OwnerType string

	// 物理机
	PhysicalMachine *PhysicalMachine `gorm:"foreignKey:OwnerID;references:UUID" json:"physicalMachine"`

	// 网络设备
	NetworkDevice *NetworkDevice `gorm:"foreignKey:OwnerID;references:UUID" json:"networkDevice"`

	// 备件关联
	SparePart SparePart `gorm:"polymorphic:Owner;foreignKey:UUID" json:"sparePart"`

	WarehouseUUID string    `gorm:"column:warehouse_uuid;notnull;comment:'仓库uuid'" json:"warehouseUUID"`
	Warehouse     Warehouse `gorm:"foreignKey:WarehouseUUID;references:UUID" json:"warehouse"`

	BrandUUID string `gorm:"column:brand_uuid;notnull;comment:'品牌UUID'" json:"brandUuid"`
	Brand     Brand  `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`

	SupplierUUID string   `gorm:"column:supplier_uuid;notnull;comment:'供应商UUID'" json:"supplierUuid"`
	Supplier     Supplier `gorm:"foreignKey:SupplierUUID;references:UUID" json:"supplier"`

	// 数据来源
	// 1: 自动采集 2: 手动录入 3: 导入
	DataSource int `gorm:"column:data_source;notnull;comment:'数据来源'" json:"dataSource"`

	// 维保起始时间
	MaintenanceStartTime string `gorm:"column:maintenance_start_time;comment:'维保起始时间'" json:"maintenanceStartTime"`

	// 维保到期时间
	MaintenanceExpireTime string `gorm:"column:expire_time;comment:'维保到期时间'" json:"maintenanceExpireTime"`
}

func (Accessories) TableName() string {
	return "assets_accessories"
}

func (a *Accessories) BeforeCreate(tx *gorm.DB) (err error) {
	if a.UUID == "" {
		a.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

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

// 网络设备
type NetworkDevice struct {
	gorm.Model
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	// 网络设备名字
	NetworkDeviceName string `gorm:"column:network_device_name;notnull;comment:'网络设备名字'" json:"networkDeviceName"`

	BrandUUID string `gorm:"column:brand_uuid;notnull;comment:'品牌uuid'" json:"brandUUID"`
	Brand     Brand  `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`

	// 供应商UUID
	SupplierUUID string   `gorm:"column:supplier_uuid;notnull;comment:'供应商uuid'" json:"supplierUUID"`
	Supplier     Supplier `gorm:"foreignKey:SupplierUUID;references:UUID" json:"supplier"`

	ManagementIP string `gorm:"column:management_ip;notnull;comment:'管理ip'" json:"managementIp"`
	// sn
	Sn string `gorm:"column:sn;notnull;comment:'sn'" json:"sn"`
	// 端口数量
	PortCount int `gorm:"column:port_count;notnull;default:0;comment:'端口数量'" json:"portCount"`
	// 环境  1:生产 2:测试 3:开发
	Env int `gorm:"column:env;notnull;default:1;comment:'环境'" json:"env"`
	// 状态  1:使用中 2：故障 3：关机 4: 进入库房
	Status int `gorm:"column:status;notnull;default:1;comment:'状态'" json:"status"`

	// 维护状态 1:正常 2:维护中
	MaintainStatus int `gorm:"column:maintain_status;notnull;default:1;comment:'维护状态'" json:"maintainStatus"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	// 关联的设备
	Device Device `gorm:"polymorphic:Device;foreignKey:UUID" json:"device"`

	// 备件关联
	SparePart SparePart `gorm:"polymorphic:Owner;foreignKey:UUID" json:"sparePart"`

	// 配件关联
	Accessories []Accessories `gorm:"polymorphic:Owner;foreignKey:UUID" json:"accessories"`

	// 维保起始时间
	MaintenanceStartTime string `gorm:"column:maintenance_start_time;notnull;comment:'维保起始时间'" json:"maintenanceStartTime"`

	// 维保到期时间
	MaintenanceExpireTime string `gorm:"column:expire_time;notnull;comment:'维保到期时间'" json:"maintenanceExpireTime"`

	NamespaceID int `gorm:"column:namespace_id;notnull;comment:'项目ID'" json:"namespaceId"`

	NamespaceService TblServicetree `gorm:"->;foreignKey:NamespaceID;references:Pri" json:"namespaceService"`

	NameID int `gorm:"column:name_id;notnull;comment:'服务名称'" json:"nameId"`

	NameService TblServicetree `gorm:"->;foreignKey:NameID;references:Pri" json:"nameService"`

	// 关联库房
	//Warehouse Warehouse `gorm:"polymorphic:Warehouse" json:"warehouse"`
}

func (networkDevice *NetworkDevice) BeforeCreate(tx *gorm.DB) (err error) {
	if networkDevice.UUID == "" {
		networkDevice.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

func (NetworkDevice) TableName() string {
	return "assets_network_device"
}

type TblServicetree struct {
	Pri      int              `gorm:"primary_key;column:pri" json:"pri"`
	Id       int              `gorm:"column:id;primary_key;notnull;comment:'id'" json:"id"`
	Name     string           `gorm:"column:name;notnull;comment:'名称'" json:"name"`
	Aname    string           `gorm:"column:aname;notnull;comment:'别名'" json:"aname"`
	Pnode    int              `gorm:"column:pnode;notnull;comment:'父节点'" json:"pnode"`
	Type     int              `gorm:"column:type;notnull;comment:'类型'" json:"type"`
	Key      string           `gorm:"column:key;notnull;comment:'key'" json:"key"`
	Origin   string           `gorm:"column:origin;notnull;comment:'来源'" json:"origin"`
	Services []TblServicetree `json:"services" gorm:"FOREIGNKEY:pnode;ASSOCIATION_FOREIGNKEY:id"`
}

func (t TblServicetree) TableName() string {
	return "tbl_servicetree"
}
