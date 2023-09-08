// Code generated . DO NOT EDIT.

package gormq

import (
	"fmt"
)

const (
	_ = iota
	SQLWHEREOP_EQ
	SQLWHEREOP_NEQ
	SQLWHEREOP_IN
	SQLWHEREOP_NOTIN
	SQLWHEREOP_LT
	SQLWHEREOP_LTE
	SQLWHEREOP_GT
	SQLWHEREOP_GTE
	SQLWHEREOP_OR
	SQLWHEREOP_AND
)

const (
	SQLWHEREOP_EQALIAS    = "Eq"
	SQLWHEREOP_NEQALIAS   = "Neq"
	SQLWHEREOP_INALIAS    = "In"
	SQLWHEREOP_NOTINALIAS = "NotIn"
	SQLWHEREOP_LTALIAS    = "Lt"
	SQLWHEREOP_LTEALIAS   = "Lte"
	SQLWHEREOP_GTALIAS    = "Gt"
	SQLWHEREOP_GTEALIAS   = "Gte"
	SQLWHEREOP_ORALIAS    = "OR"
	SQLWHEREOP_ANDALIAS   = "AND"
)

var _SqlWhereOpValue = map[int]SqlWhereOp{
	1:  SQLWHEREOP_EQ,
	10: SQLWHEREOP_AND,
	2:  SQLWHEREOP_NEQ,
	3:  SQLWHEREOP_IN,
	4:  SQLWHEREOP_NOTIN,
	5:  SQLWHEREOP_LT,
	6:  SQLWHEREOP_LTE,
	7:  SQLWHEREOP_GT,
	8:  SQLWHEREOP_GTE,
	9:  SQLWHEREOP_OR,
}

func ParseSqlWhereOp(id int) (SqlWhereOp, error) {
	if x, ok := _SqlWhereOpValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e SqlWhereOp) String() string {
	switch e {
	case 1:
		return SQLWHEREOP_EQALIAS
	case 2:
		return SQLWHEREOP_NEQALIAS
	case 3:
		return SQLWHEREOP_INALIAS
	case 4:
		return SQLWHEREOP_NOTINALIAS
	case 5:
		return SQLWHEREOP_LTALIAS
	case 6:
		return SQLWHEREOP_LTEALIAS
	case 7:
		return SQLWHEREOP_GTALIAS
	case 8:
		return SQLWHEREOP_GTEALIAS
	case 9:
		return SQLWHEREOP_ORALIAS
	case 10:
		return SQLWHEREOP_ANDALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}

const (
	_ = iota
	FIELDTYPE_CUSTOMIZE
	FIELDTYPE_COLUMN
	FIELDTYPE_ASSOCIATE
)

const (
	FIELDTYPE_CUSTOMIZEALIAS = "customize"
	FIELDTYPE_COLUMNALIAS    = "column"
	FIELDTYPE_ASSOCIATEALIAS = "associate"
)

var _FieldTypeValue = map[int]FieldType{
	1: FIELDTYPE_CUSTOMIZE,
	2: FIELDTYPE_COLUMN,
	3: FIELDTYPE_ASSOCIATE,
}

func ParseFieldType(id int) (FieldType, error) {
	if x, ok := _FieldTypeValue[id]; ok {
		return x, nil
	}
	return 0, fmt.Errorf("unknown enum value: %s", id)
}

func (e FieldType) String() string {
	switch e {
	case 1:
		return FIELDTYPE_CUSTOMIZEALIAS
	case 2:
		return FIELDTYPE_COLUMNALIAS
	case 3:
		return FIELDTYPE_ASSOCIATEALIAS
	}
	return fmt.Sprintf("unknown %d", e)
}
