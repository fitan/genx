package data

type CheckIpIn map[string]CheckIpInItem

type CheckIpInItem struct {
	Ip   string
	UUID string
}

// @gormM(types.User)
// @gormP(Friends)
type ListRequest struct {
	// Page 页码
	// @gormPage
	Page int `param:"query,page" json:"page"`
	// PageSize 每页数量
	// @gormPageSize
	PageSize int `param:"query,pageSize" json:"pageSize"`
	// Order 排序
	// @gormq(Username)
	Order string `param:"query,order" json:"order"`
	// @gormq(Keyword, like)
	// Keyword 关键字
	Keyword string `param:"query,keyword" json:"keyword"`
	// 品牌
	Brand string `param:"query,brand" json:"brand"`
	// 产品型号
	BrandUuid string `param:"query,brandUuid" json:"brandUuid"`

	// 朋友名字
	// @gormq(Friends.Name, "like")
	FriendName string `param:"query,friendName" json:"friendName"`

	// ComputerRoomUUID 机房UUID
	ComputerRoomUUID string `param:"query,computerRoomUuid" json:"computerRoomUuid"`
}

type Brand struct {
	// uuid
	UUID string `json:"uuid"`
	// 品牌
	Brand string `json:"brand"`

	// 产品类型 1:服务器 2:交换机 3:路由器 4:防火墙 5:负载均衡器 6:存储设备 7: cpu 8:内存 9:硬盘 10:网卡 10: 系统
	ProductType string `json:"productType"`

	// 产品型号
	ProductModel string `json:"productModel"`

	// 备注
	Remark string `json:"remark"`
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

type ListResponse struct {
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

type Device struct {
	// device uuid
	UUID string `json:"uuid"`
	// 设备类型
	DeviceType string `json:"deviceType"`
	// 设备ID
	DeviceUUID string `json:"deviceUuid"`

	// 机柜名字
	CabinetName string `json:"cabinetName"`
	// 机柜UUID
	CabinetUUID string `json:"cabinetUuid"`

	// 起始层
	StartLayer int `json:"startLayer"`
	// 终止层
	EndLayer int `json:"endLayer"`
}

type GetByUUIDRequest struct {
	UUID string `param:"path,uuid" json:"uuid"`
}

type GetByUUIDResponse struct {
}

type CreateRequest struct {
	// brand UUID
	BrandUUID string `json:"brandUuid"`
	// 业务ip
	BusinessIp string `json:"businessIp"`
	// 业务子网掩码
	BusinessSubnetMask string `json:"businessSubnetMask"`
	// nasIp
	NasIp string `json:"nasIp"`
	// nas子网掩码
	NasSubnetMask string `json:"nasSubnetMask"`

	// 物理机名字
	PhysicalMachineName string `json:"physicalMachineName"`

	// IPMIIP
	IpmiIp string `json:"ipmiIp"`
	// ipmi子网掩码
	IpmiSubnetMask string `json:"ipmiSubnetMask"`
	// sn
	Sn string `json:"sn"`

	// cpu brand uuid
	//CpuBrandUUID string `json:"cpuBrandUuid"`
	// cpu数量
	//CpuCount int `json:"cpuCount"`

	// mem brand uuid
	//MemoryBrandUUID string `json:"memoryBrandUuid"`
	// 内存数量
	//MemoryCount int `json:"memoryCount"`

	// disk brand uuid
	//DiskBrandUUID string `json:"diskBrandUuid"`
	// 硬盘数量
	//DiskCount int `json:"diskCount"`

	// NetworkCardBrandUUID
	//NetworkCardBrandUUID string `json:"networkCardBrandUuid"`
	// 网卡数量
	//NetworkCardCount int `json:"networkCardCount"`

	// os brand uuid
	OsBrandUUID string `json:"osBrandUuid"`
	// 供应商UUID
	SupplierUUID string `json:"supplierUuid"`

	// 维保起始时间
	MaintenanceStartTime string `json:"maintenanceStartTime"`

	// 维保到期时间
	MaintenanceExpireTime string `json:"maintenanceExpireTime"`

	// 负责人
	Owners []string `json:"owners"`

	// 跳板机登录端口
	JumpPort int `json:"jumpPort"`
	// 环境  1:生产 2:测试 3:开发
	Env int `json:"env"`
	// 状态  1:开机 2：故障 3：关机
	Status int `json:"status"`
	// 维护状态 1:正常 2:维护中
	MaintainStatus int `json:"maintainStatus"`

	// 备注
	Remark string `json:"remark"`

	// namespace 项目名称
	Namespace string `json:"namespace"`
	// name 服务名称
	Name string `json:"name"`

	// 机柜UUID
	CabinetUUID string `json:"cabinetUuid" valid:"required"`
	// 起始层数
	StartLayer int `json:"startLayer" valid:"required"`
	// 终止层数
	EndLayer int `json:"endLayer" valid:"required"`

	//网络关联设备的UUID
	//NetworkToDevices []string `json:"networkToDevices"`
}

type CreateResponse struct {
}

type UpdateRequest struct {
	UUID            string        `param:"path,uuid" json:"uuid" valid:"required"`
	PhysicalMachine CreateRequest `param:"body,physicalMachine" json:"computerRoom"`
}

type UpdateResponse struct {
}

type DeleteRequest struct {
	UUID string `param:"path,uuid" json:"uuid"`
}

type DeleteResponse struct {
}

type PutInStorageRequest struct {
	UUID string                  `param:"path,uuid" json:"uuid"`
	Body PutInStorageRequestBody `param:"body,body"`
}

type PutInStorageRequestBody struct {
	WarehouseUUID string `json:"warehouseUuid"`
}

type PutInStorageResponse struct {
}

type IncludeRequest struct {
	//Body *multipart.FileHeader `param:"file,file" json:"body"`
	Body []byte `json:"body" param:"body,body"`
}

type Include struct {
	PhysicalMachineName string `json:"physicalMachineName" csv:"physicalMachineName"`
	Env                 int    `json:"env" csv:"env"`
	StartLayer          int    `json:"startLayer" csv:"startLayer"`
	EndLayer            int    `json:"endLayer" csv:"endLayer"`
	OsBrand             string `json:"osBrand" csv:"osBrand"`
	OsBrandProductModel string `json:"osBrandProductModel" csv:"osBrandProductModel"`
	BusinessIp          string `json:"businessIp" csv:"businessIp"`
	BusinessSubnetMask  string `json:"businessSubnetMask" csv:"businessSubnetMask"`
	NasIp               string `json:"nasIp" csv:"nasIp"`
	NasSubnetMask       string `json:"nasSubnetMask" csv:"nasSubnetMask"`
	IpmiIp              string `json:"ipmiIp" csv:"ipmiIp"`
	IpmiSubnetMask      string `json:"ipmiSubnetMask" csv:"ipmiSubnetMask"`
	JumpPort            int    `json:"jumpPort" csv:"jumpPort"`
	Remark              string `json:"remark" csv:"remark"`
	Owners              string `json:"owners" csv:"owners"`
	CabinetName         string `json:"cabinetName" csv:"cabinetName"`
	ComputerRoomName    string `json:"computerRoomName" csv:"computerRoomName"`
	Name                string `json:"name" csv:"name"`
	Namespace           string `json:"namespace" csv:"namespace"`
	SupplierName        string `json:"supplierName" csv:"supplierName"`
}

//type Include struct {
//	PhysicalMachineName string `json:"physicalMachineName" csv:"物理机名字"`
//	Env                 int    `json:"env" csv:"环境"`
//	StartLayer          int    `json:"startLayer" csv:"起始u"`
//	EndLayer            int    `json:"endLayer" csv:"结束u"`
//	OsBrand             string `json:"osBrand" csv:"系统型号"`
//	OsBrandProductModel string `json:"osBrandProductModel" csv:"系统版本"`
//	BusinessIp          string `json:"businessIp" csv:"业务ip"`
//	BusinessSubnetMask  string `json:"businessSubnetMask" csv:"业务ip子网掩码"`
//	NasIp               string `json:"nasIp" csv:"nasIp"`
//	NasSubnetMask       string `json:"nasSubnetMask" csv:"nas子网掩码"`
//	IpmiIp              string `json:"ipmiIp" csv:"管理ip"`
//	IpmiSubnetMask      string `json:"ipmiSubnetMask" csv:"管理子网掩码"`
//	JumpPort            int    `json:"jumpPort" csv:"跳板机端口"`
//	Remark              string `json:"remark" csv:"备注"`
//	Owners              string `json:"owners" csv:"负责人邮箱"`
//	CabinetName         string `json:"cabinetName" csv:"机柜名字"`
//	ComputerRoomName    string `json:"computerRoomName" csv:"机房名字"`
//	Name                string `json:"name" csv:"服务名字"`
//	Namespace           string `json:"namespace" csv:"项目名字"`
//	SupplierName        string `json:"supplierName" csv:"供应商名字"`
//}

type IncludeResponse []IncludeError

type IncludeError struct {
	Include
	Error string `json:"error" csv:"error"`
}

type ProjectServiceRequest struct {
	Body ProjectServiceRequestBody `param:"body,body" json:"body"`
}

type ProjectServiceRequestBody struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"service"`
	Ips       []string `json:"ips"`
}

type ProjectServiceAvailableRequest struct {
	Page     int    `param:"query,page" json:"page"`
	PageSize int    `param:"query,pageSize" json:"pageSize"`
	Ip       string `param:"query,ip" json:"ip"`
}

type RemoveProjectServiceRequest struct {
	UUID string `param:"path,uuid" json:"uuid"`
}

type ImpiRequest struct {
	UUID string `json:"uuid" param:"path,uuid"`
	// chassisControlPowerUp 电源开启
	// chassisControlSoftShutdown 电源软关机
	// chassisControlPowerCycle 电源重启
	// chassisControlPowerDown 电源关闭
	Action string `json:"action" param:"path,action"`
}
type ImpiUserPassword struct {
	DELL struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	HuaWei struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	IBM struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	HP struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	SuperMicro struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
}

var impiUserPassword = ImpiUserPassword{
	DELL: struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		"root",
		"calvin",
	},
	HuaWei: struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		"root",
		"Huawei12#$",
	},
	IBM: struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		"USERID",
		"PASSW0RD",
	},
	HP: struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		"root",
		"Huawei12#$",
	},
	SuperMicro: struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		"ADMIN",
		"ADMIN",
	},
}
