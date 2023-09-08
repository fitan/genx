package types

import "gorm.io/gorm"

type PhysicalMachineAllocation struct {
	gorm.Model
	// 项目id
	NamespaceID int `gorm:"column:namespace_id;comment:'项目id'" json:"namespaceId"`
	// 服务id
	NameID int `gorm:"column:name_id;comment:'服务id'" json:"nameId"`
	// 物理机UUID
	PhysicalMachineUUID string `gorm:"column:physical_machine_uuid;comment:'物理机UUID'" json:"physicalMachineUUID"`
	// 物理机关联
	PhysicalMachine PhysicalMachine `gorm:"references:UUID;foreignKey:PhysicalMachineUUID" json:"physicalMachine"`

	NamespaceService TblServicetree `gorm:"->;foreignKey:NamespaceID;references:Pri" json:"namespaceService"`

	NameService TblServicetree `gorm:"->;foreignKey:NameID;references:Pri" json:"nameService"`
}

func (PhysicalMachineAllocation) TableName() string {
	return "assets_physical_machine_allocation"
}
