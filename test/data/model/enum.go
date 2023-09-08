// Code generated . DO NOT EDIT.

package types

import (
	"fmt"
)

const (
	_ = iota
	COMPUTERROOMTYPE_DATACENTER
	COMPUTERROOMTYPE_WORKPLACECOMPUTERROOM
)

const (
	COMPUTERROOMTYPE_DATACENTERALIAS            = "dataCenter"
	COMPUTERROOMTYPE_WORKPLACECOMPUTERROOMALIAS = "workplaceComputerRoom"
)

var _ComputerRoomTypeValue = map[int]ComputerRoomType{
	1: COMPUTERROOMTYPE_DATACENTER,
	2: COMPUTERROOMTYPE_WORKPLACECOMPUTERROOM,
}

func ParseComputerRoomType(id int) (ComputerRoomType, error) {
	if x, ok := _ComputerRoomTypeValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e ComputerRoomType) String() string {
	switch e {
	case 1:
		return COMPUTERROOMTYPE_DATACENTERALIAS
	case 2:
		return COMPUTERROOMTYPE_WORKPLACECOMPUTERROOMALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}

const (
	_ = iota
	REPORTWAY_BEE
	REPORTWAY_WECHAT
	REPORTWAY_EMAIL
	REPORTWAY_PHONE
	REPORTWAY_OTHER
	REPORTWAY_INSPECTION
)

const (
	REPORTWAY_BEEALIAS        = "Bee"
	REPORTWAY_WECHATALIAS     = "Wechat"
	REPORTWAY_EMAILALIAS      = "Email"
	REPORTWAY_PHONEALIAS      = "Phone"
	REPORTWAY_OTHERALIAS      = "Other"
	REPORTWAY_INSPECTIONALIAS = "Inspection"
)

var _ReportWayValue = map[int]ReportWay{
	1: REPORTWAY_BEE,
	2: REPORTWAY_WECHAT,
	3: REPORTWAY_EMAIL,
	4: REPORTWAY_PHONE,
	5: REPORTWAY_OTHER,
	6: REPORTWAY_INSPECTION,
}

func ParseReportWay(id int) (ReportWay, error) {
	if x, ok := _ReportWayValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e ReportWay) String() string {
	switch e {
	case 1:
		return REPORTWAY_BEEALIAS
	case 2:
		return REPORTWAY_WECHATALIAS
	case 3:
		return REPORTWAY_EMAILALIAS
	case 4:
		return REPORTWAY_PHONEALIAS
	case 5:
		return REPORTWAY_OTHERALIAS
	case 6:
		return REPORTWAY_INSPECTIONALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}

const (
	_ = iota
	RESOURCETYPE_SERVER
	RESOURCETYPE_NETWORKDEVICE
	RESOURCETYPE_INTERNETCONNECTION
	RESOURCETYPE_DEDICATEDLINE
	RESOURCETYPE_IPADDRESS
	RESOURCETYPE_IDC
)

const (
	RESOURCETYPE_SERVERALIAS             = "Server"
	RESOURCETYPE_NETWORKDEVICEALIAS      = "NetworkDevice"
	RESOURCETYPE_INTERNETCONNECTIONALIAS = "InternetConnection"
	RESOURCETYPE_DEDICATEDLINEALIAS      = "DedicatedLine"
	RESOURCETYPE_IPADDRESSALIAS          = "IPAddress"
	RESOURCETYPE_IDCALIAS                = "IDC"
)

var _ResourceTypeValue = map[int]ResourceType{
	1: RESOURCETYPE_SERVER,
	2: RESOURCETYPE_NETWORKDEVICE,
	3: RESOURCETYPE_INTERNETCONNECTION,
	4: RESOURCETYPE_DEDICATEDLINE,
	5: RESOURCETYPE_IPADDRESS,
	6: RESOURCETYPE_IDC,
}

func ParseResourceType(id int) (ResourceType, error) {
	if x, ok := _ResourceTypeValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e ResourceType) String() string {
	switch e {
	case 1:
		return RESOURCETYPE_SERVERALIAS
	case 2:
		return RESOURCETYPE_NETWORKDEVICEALIAS
	case 3:
		return RESOURCETYPE_INTERNETCONNECTIONALIAS
	case 4:
		return RESOURCETYPE_DEDICATEDLINEALIAS
	case 5:
		return RESOURCETYPE_IPADDRESSALIAS
	case 6:
		return RESOURCETYPE_IDCALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}

const (
	_ = iota
	FAULTSTATUS_FAULTINPROGRESS
	FAULTSTATUS_FAULTRESOLVED
)

const (
	FAULTSTATUS_FAULTINPROGRESSALIAS = "FaultInProgress"
	FAULTSTATUS_FAULTRESOLVEDALIAS   = "FaultResolved"
)

var _FaultStatusValue = map[int]FaultStatus{
	1: FAULTSTATUS_FAULTINPROGRESS,
	2: FAULTSTATUS_FAULTRESOLVED,
}

func ParseFaultStatus(id int) (FaultStatus, error) {
	if x, ok := _FaultStatusValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e FaultStatus) String() string {
	switch e {
	case 1:
		return FAULTSTATUS_FAULTINPROGRESSALIAS
	case 2:
		return FAULTSTATUS_FAULTRESOLVEDALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}

const (
	_ = iota
	FAULTTYPE_HARDDISK
	FAULTTYPE_MOTHERBOARD
	FAULTTYPE_MEMORY
	FAULTTYPE_NETWORKCARD
	FAULTTYPE_FIRMWAREVERSION
	FAULTTYPE_OPERATINGSYSTEM
)

const (
	FAULTTYPE_HARDDISKALIAS        = "HardDisk"
	FAULTTYPE_MOTHERBOARDALIAS     = "Motherboard"
	FAULTTYPE_MEMORYALIAS          = "Memory"
	FAULTTYPE_NETWORKCARDALIAS     = "NetworkCard"
	FAULTTYPE_FIRMWAREVERSIONALIAS = "FirmwareVersion"
	FAULTTYPE_OPERATINGSYSTEMALIAS = "OperatingSystem"
)

var _FaultTypeValue = map[int]FaultType{
	1: FAULTTYPE_HARDDISK,
	2: FAULTTYPE_MOTHERBOARD,
	3: FAULTTYPE_MEMORY,
	4: FAULTTYPE_NETWORKCARD,
	5: FAULTTYPE_FIRMWAREVERSION,
	6: FAULTTYPE_OPERATINGSYSTEM,
}

func ParseFaultType(id int) (FaultType, error) {
	if x, ok := _FaultTypeValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e FaultType) String() string {
	switch e {
	case 1:
		return FAULTTYPE_HARDDISKALIAS
	case 2:
		return FAULTTYPE_MOTHERBOARDALIAS
	case 3:
		return FAULTTYPE_MEMORYALIAS
	case 4:
		return FAULTTYPE_NETWORKCARDALIAS
	case 5:
		return FAULTTYPE_FIRMWAREVERSIONALIAS
	case 6:
		return FAULTTYPE_OPERATINGSYSTEMALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}

const (
	_ = iota
	ADDRESSBLOCKTYPE_BUSINESS
	ADDRESSBLOCKTYPE_MANAGEMENT
	ADDRESSBLOCKTYPE_STORAGE
)

const (
	ADDRESSBLOCKTYPE_BUSINESSALIAS   = "Business"
	ADDRESSBLOCKTYPE_MANAGEMENTALIAS = "Management"
	ADDRESSBLOCKTYPE_STORAGEALIAS    = "Storage"
)

var _AddressBlockTypeValue = map[int]AddressBlockType{
	1: ADDRESSBLOCKTYPE_BUSINESS,
	2: ADDRESSBLOCKTYPE_MANAGEMENT,
	3: ADDRESSBLOCKTYPE_STORAGE,
}

func ParseAddressBlockType(id int) (AddressBlockType, error) {
	if x, ok := _AddressBlockTypeValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e AddressBlockType) String() string {
	switch e {
	case 1:
		return ADDRESSBLOCKTYPE_BUSINESSALIAS
	case 2:
		return ADDRESSBLOCKTYPE_MANAGEMENTALIAS
	case 3:
		return ADDRESSBLOCKTYPE_STORAGEALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}

const (
	_ = iota
	NETWORKSEGMENTENVIRONMENT_WORKPLACE
	NETWORKSEGMENTENVIRONMENT_PRODUCTION
	NETWORKSEGMENTENVIRONMENT_TEST
)

const (
	NETWORKSEGMENTENVIRONMENT_WORKPLACEALIAS  = "Workplace"
	NETWORKSEGMENTENVIRONMENT_PRODUCTIONALIAS = "Production"
	NETWORKSEGMENTENVIRONMENT_TESTALIAS       = "Test"
)

var _NetworkSegmentEnvironmentValue = map[int]NetworkSegmentEnvironment{
	1: NETWORKSEGMENTENVIRONMENT_WORKPLACE,
	2: NETWORKSEGMENTENVIRONMENT_PRODUCTION,
	3: NETWORKSEGMENTENVIRONMENT_TEST,
}

func ParseNetworkSegmentEnvironment(id int) (NetworkSegmentEnvironment, error) {
	if x, ok := _NetworkSegmentEnvironmentValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e NetworkSegmentEnvironment) String() string {
	switch e {
	case 1:
		return NETWORKSEGMENTENVIRONMENT_WORKPLACEALIAS
	case 2:
		return NETWORKSEGMENTENVIRONMENT_PRODUCTIONALIAS
	case 3:
		return NETWORKSEGMENTENVIRONMENT_TESTALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}

const (
	_ = iota
	NETWORKSEGMENNETWORKTYPE_INTRANET
	NETWORKSEGMENNETWORKTYPE_EXTRANET
)

const (
	NETWORKSEGMENNETWORKTYPE_INTRANETALIAS = "Intranet"
	NETWORKSEGMENNETWORKTYPE_EXTRANETALIAS = "Extranet"
)

var _NetworkSegmenNetworkTypeValue = map[int]NetworkSegmenNetworkType{
	1: NETWORKSEGMENNETWORKTYPE_INTRANET,
	2: NETWORKSEGMENNETWORKTYPE_EXTRANET,
}

func ParseNetworkSegmenNetworkType(id int) (NetworkSegmenNetworkType, error) {
	if x, ok := _NetworkSegmenNetworkTypeValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e NetworkSegmenNetworkType) String() string {
	switch e {
	case 1:
		return NETWORKSEGMENNETWORKTYPE_INTRANETALIAS
	case 2:
		return NETWORKSEGMENNETWORKTYPE_EXTRANETALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}
