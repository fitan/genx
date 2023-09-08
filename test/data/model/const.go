package types

// 多态映射表名
var PolymorphicMap = map[string]string{
	"physicalMachine": PhysicalMachine{}.TableName(),
	"networkDevice":   NetworkDevice{}.TableName(),
}
