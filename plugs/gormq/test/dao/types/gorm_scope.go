package types

import (
	"github.com/fitan/genx/plugs/gormq/test/model"
	"gorm.io/gorm"
)

func (q *PmList) Scope(db *gorm.DB) *gorm.DB {
	db = db.Model(&model.PhysicalMachine{})

	db = db.Where("uuid = ?", q.UUID)

	db = db.Where("sn = ?", q.Sn)

	db = db.Where("maintain_status = ?", q.MaintainStatus)

	db = db.Where("name in ?", q.Name)
	if q.Sn1 != nil {

		db = db.Where("sn = ?", q.Sn1)
	}

	db = db.Where("sn = ?", q.Keyword.Keyword).Or("uuid = ?", q.Keyword.Keyword)

	ListInValue := make([][]interface{}, 0, 0)
	for _, v := range q.ListIn {
		ListInValue = append(ListInValue, []interface{}{v.ProductType, v.Brand})
	}
	db = db.Where("(product_type, product_model) IN ?", ListInValue)
	if q.PmListNest.MaintainStatus != nil {

		db = db.Where("maintain_status = ?", q.PmListNest.MaintainStatus)
	}
	db = db.Where(&q.Match, "MaintainStatus")
	db = db.Where("brand_uuid in (?)", q.BrandSub.Scope(db.Session(&gorm.Session{NewDB: true}).Select("uuid")))
	db = db.Where("os_brand_uuid in (?)", q.OsBrandSub.Scope(db.Session(&gorm.Session{NewDB: true}).Select("uuid")))
	db = db.Where("os_brand_uuid in (?)", q.PmListNest.OsBrandSub.Scope(db.Session(&gorm.Session{NewDB: true}).Select("uuid")))
	db = db.Or(func(db *gorm.DB) *gorm.DB {

		db = db.Where("sn = ?", q.PmListKeyword.Keyword).Or("uuid = ?", q.PmListKeyword.Keyword)
		return db
	}(db.Session(&gorm.Session{NewDB: true})))
	if q.PointPmListNest != nil {
		if q.PointPmListNest.MaintainStatus != nil {

			db = db.Where("maintain_status = ?", q.PointPmListNest.MaintainStatus)
		}
		db = db.Where("os_brand_uuid in (?)", q.PointPmListNest.OsBrandSub.Scope(db.Session(&gorm.Session{NewDB: true}).Select("uuid")))
	}
	return db
}
func (q *PmListBrand) Scope(db *gorm.DB) *gorm.DB {
	db = db.Model(&model.Brand{})

	db = db.Where("product_type like ?", "%"+q.ProductType+"%")

	db = db.Where("product_model like ?", "%"+q.Brand+"%")
	return db
}
func (q *PmListOsBrand) Scope(db *gorm.DB) *gorm.DB {
	db = db.Model(&model.Brand{})

	db = db.Where("product_type like ?", "%"+q.ProductType+"%")

	db = db.Where("product_model like ?", "%"+q.Brand+"%")

	ListInValue := make([][]interface{}, 0, 0)
	for _, v := range q.ListIn {
		ListInValue = append(ListInValue, []interface{}{v.ProductType, v.Brand})
	}
	db = db.Where("(product_type, product_model) IN ?", ListInValue)
	return db
}
