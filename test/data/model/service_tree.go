package types

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
