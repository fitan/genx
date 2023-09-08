package types

import "gorm.io/gorm"

type PubDnsAudit struct {
	gorm.Model
	Event         string `gorm:"column:event;type:varchar(50)" json:"event"`                              // 事件
	OperateType   string `gorm:"column:operate_type;type:varchar(20);default:domain" json:"operate_type"` // 操作类型：域名、记录
	DomainName    string `gorm:"column:domain_name;type:varchar(255)" json:"domain_name"`                 // 域名
	RecordName    string `gorm:"column:record_name;type:varchar(255)" json:"record_name"`                 // 记录名
	Remark        string `gorm:"column:remark;type:varchar(255)" json:"remark"`                           // 备注
	ExtData       []byte `gorm:"column:ext_data;type:varbinary(20480)" json:"ext_data"`
	OperatorEmail string `gorm:"column:operator_email;type:varchar(255)" json:"operator_email"` // 操作人邮箱
	OperatorName  string `gorm:"column:operator_name;type:varchar(255)" json:"operator_name"`   // 操作人姓名
}

func (m *PubDnsAudit) TableName() string {
	return "pub_dns_audit"
}
