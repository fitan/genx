package types

// 关联用户
type RelationUser struct {
	// 创建用户
	CreateUser string `gorm:"column:create_user;notnull;comment:'创建用户'" json:"createUserId"`
	// 更新用户
	UpdateUser string `gorm:"column:update_user;notnull;comment:'更新用户'" json:"updateUserId"`
	// 维护用户
	Maintainer []string `gorm:"column:maintainer;notnull;comment:'维护用户'" json:"maintainerId"`
}
