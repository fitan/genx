package types

// grom User表
type User struct {
	// ID
	ID int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	// 用户名
	Username string `json:"username" gorm:"column:username;type:varchar(32);not null;default:'';comment:用户名"`
	// 密码
	Password string `json:"password" gorm:"column:password;type:varchar(32);not null;default:'';comment:密码"`
	// 昵称
	Nickname string `json:"nickname" gorm:"column:nickname;type:varchar(32);not null;default:'';comment:昵称"`
	// 头像
	Avatar string `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像"`
	// 朋友
	Friends []Friend `json:"friends" gorm:"foreignKey:UserID;references:ID;comment:朋友"`
}
	

// gorm 朋友表
type Friend struct {
	// ID
	ID int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	// 名字
	Name string `json:"name" gorm:"column:name;type:varchar(32);not null;default:'';comment:名字"`
	// 用户ID
	UserID int64 `json:"user_id" gorm:"column:user_id;type:bigint(20);not null;default:0;comment:用户ID"`
	// 朋友ID
	FriendID int64 `json:"friend_id" gorm:"column:friend_id;type:bigint(20);not null;default:0;comment:朋友ID"`
	// 备注
	Remark string `json:"remark" gorm:"column:remark;type:varchar(32);not null;default:'';comment:备注"`
	// 分组
	Group string `json:"group" gorm:"column:group;type:varchar(32);not null;default:'';comment:分组"`
	// 创建时间
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;type:bigint(20);not null;default:0;comment:创建时间"`
	// 更新时间
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;type:bigint(20);not null;default:0;comment:更新时间"`
}