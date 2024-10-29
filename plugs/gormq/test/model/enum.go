// Code generated . DO NOT EDIT.

package model

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

const (
	COMPUTERROOMTYPE_DATACENTERREMARK            = "dataCenter"
	COMPUTERROOMTYPE_WORKPLACECOMPUTERROOMREMARK = "workplaceComputerRoom"
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

func (e ComputerRoomType) Remark() string {
	switch e {
	case 1:
		return COMPUTERROOMTYPE_DATACENTERREMARK
	case 2:
		return COMPUTERROOMTYPE_WORKPLACECOMPUTERROOMREMARK
	}
	return fmt.Sprintf("unknown %d", e)
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
