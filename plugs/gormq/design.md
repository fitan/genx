```go
package main

import "gorm.io/gorm"

// @gq model.User
type User struct {
	// @gq-op = 
	ID uint `gorm:"primary_key"`
	// @gq-op !=
	Name string `json:"name"`
	// @gq-op in
	Email []string `json:"email"`
	
	FriendID uint `json:"friendID"`
	
	
	
	
	// 会递归解析
	Friend 
	
	// 也会递归解析
	FriendNest Friend
	
	// 会变成 where((id, name, email) in ?, [][]interface{}{})
	FriendIn []Friend
	
	// 如果是数组会变成 where("name in ?", names)
	// @gq-column name
	Names []string
	
}

// @gq model.Friend
type Friend struct {
	ID uint `gorm:"primary_key"`
	// @gq-op =
	Name string `json:"name"`
	// @gq-op in
	Email []string `json:"email"`
}



func (u User) GenQuery(db *gorm.DB) *gorm.DB {
	return db.
		Where("name = ?", u.Name).
		Where("email in ?", u.Email).
		Where("id > ?", u.ID).
		Where("friend_id in (?)", db.Model(&model.Friend{}).Select("id").Where("name = ?", u.Friend.Name).Where("email in ?", u.Friend.Email).Where("id > ?", u.Friend.ID))
}
```

### 用法
#### @gq
使用@gq标记model: @gq model.User
#### @gq-column
使用@gq-column标记字段: @gq-column uuid
可以使用连接多个 @gq-column uuid name age
这样 会生成Where("uuid = ?", xx).Or("name = ?", xx).Or("age = ?", xx)
#### @gq-sub
使用@gq-sub标记子查询: @gq-sub brand_uuid uuid
会生成Where("brand_uuid in (?)", db.Model(&model.Brand{}).Select("uuid").Where("uuid = ?", xx))
#### @gq-group
使用@gq-group标记分组: @gq-group
#### @gq-clause
使用@gq-clause标记子句: @gq-clause or 或者 @gq-clause where
#### @gq-op
使用@gq-op标记操作符: @gq-op = 
#### @gq-struct
使用@gq-struct标记结构体: @gq-struct
会生成Where(&user)

执行的操作符有:
* `=`
* `!=`
* `>`
* `>=`
* `<`
* `<=`
* `><`
* `!><`
* `like`
* `in`
* `!in`
* `null`
* `!null`


