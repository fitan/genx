package types

import (
	"fmt"
	"testing"

	"github.com/fitan/genx/plugs/gormq/test/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestPm(t *testing.T) {
	dsn := "spider_dev:spider_dev123@tcp(10.170.34.22:3307)/spider_dev?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	q := &PmList{
		UUID: "fsdafasdfsf",
		Sn:   "fsadfas",
		BrandSub: PmListBrand{
			ProductType: "fsadfadf",
			Brand:       "fadsfas",
		},
		OsBrandSub: PmListOsBrand{
			ProductType: "sdafas",
			Brand:       "fasdfads",
			ListIn:      nil,
		},
		MaintainStatus: 100,
		Sn1:            nil,
		Keyword: PmListKeyword{
			Keyword: "123",
		},
		ListIn: []PmListBrand{
			{
				ProductType: "adsfasdf",
				Brand:       "fsdafadsf",
			},
		},
		PmListKeyword: PmListKeyword{
			Keyword: "456",
		},
		PmListNest: PmListNest{
			OsBrandSub: PmListOsBrand{
				ProductType: "1",
				Brand:       "1",
				ListIn:      nil,
			},
		},
		PointPmListNest: nil,
	}

	qDb := db.Model(&model.PhysicalMachine{})
	var res []model.PhysicalMachine
	sql := qDb.Scopes(q.Scope).Find(&res).Statement.SQL.String()

	db1 := db.Model(&model.PhysicalMachine{})
	sql1 := db1.Where("uuid = ?", q.UUID).Or(db.Where("id = ?", 100).Where("id = ?", 100)).Find(&res).Statement.SQL.String()
	fmt.Println(sql)
	fmt.Println(sql1)

	db2 := db.Model(&model.PhysicalMachine{})
	sql2 := q.Scope(db2).Find(&res).Statement.SQL.String()
	fmt.Println(sql2)
}
