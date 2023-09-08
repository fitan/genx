package types

import "gorm.io/gorm"

type MonitorPageConfig struct {
	gorm.Model
	Name   string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Config Config `gorm:"serializer:json;not null"`
}

func (MonitorPageConfig) TableName() string {
	return "assets_monitor_page_config"
}

type Config struct {
	// 监控页面名称
	Name string `yaml:"name" json:"name"`
	// 监控页面描述
	Description string `yaml:"description" json:"description"`
	// 监控页面配置
	Groups []Group `yaml:"groups" json:"groups"`
}

type Group struct {
	// 名字
	Name string `yaml:"name" json:"name"`
	// 描述
	Description string `yaml:"description" json:"description"`
	// 图标
	Charts []Chart `yaml:"charts" json:"charts"`
}

type Chart struct {
	// 名字
	Name string `yaml:"name" json:"name"`
	// 描述
	Description string `yaml:"description" json:"description"`
	// 单位
	Unit string `yaml:"unit" json:"unit"`
	// 宽度
	Width int `yaml:"width" json:"width"`
	// 查询语句模板
	Queries []Query `yaml:"queries" json:"queries"`
}

type Query struct {
	// 模板
	Metrics string `yaml:"metrics" json:"metrics"`
	// 标签
	Tag string `yaml:"tag" json:"tag"`
}

type PromqlTemplate struct {
	gorm.Model
	// promql 名字
	Name string `gorm:"column:name;type:varchar(255)" json:"name"`
	// promql 模板
	Template string `gorm:"column:template;type:varchar(255)" json:"template"`
	// tag
	TagTemplate string `gorm:"column:tag_template;type:varchar(255)" json:"tagTemplate"`
}

func (PromqlTemplate) TableName() string {
	return "assets_promql_template"
}

type PromQuery struct {
	gorm.Model
	// query 名字
	Name string `gorm:"column:name;type:varchar(255)" json:"name"`
	// 图标title
	Title string `gorm:"column:title;type:varchar(255)" json:"title"`

	// 关联的promql模板 多对多
	PromqlTemplates []PromqlTemplate `gorm:"many2many:prom_query_promql_template;" json:"promqlTemplates"`

	// 指标单位
	Unit string `gorm:"column:unit;type:varchar(255)" json:"unit"`

	// 顺序
	Order int `gorm:"column:order;type:int" json:"order"`

	PromPageId uint `gorm:"column:prom_page_id;type:int" json:"promPageId"`
	// 宽度
	Width int `gorm:"column:width;type:int" json:"width"`
}

func (PromQuery) TableName() string {
	return "assets_prom_query"
}

type PromPage struct {
	gorm.Model
	// 页面名字
	Name string `gorm:"column:name;type:varchar(255)" json:"name"`
	// 页面描述
	Desc string `gorm:"column:desc;type:varchar(255)" json:"desc"`
	// 指标类型
	Tag string `gorm:"column:tag;type:varchar(255)" json:"tag"`
	// 关联PromQuery 一对多
	PromQueries []PromQuery `gorm:"foreignkey:prom_page_id;references:id;" json:"promQueries"`
	// 顺序
	Order int `gorm:"column:order;type:int" json:"order"`
}

func (PromPage) TableName() string {
	return "assets_prom_page"
}
