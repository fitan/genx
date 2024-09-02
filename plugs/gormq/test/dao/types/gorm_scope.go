package types

import (
	"github.com/fitan/genx/plugs/gormq/test/model"
	"gorm.io/gorm"
)

func (q *PmList) Scope(db *gorm.DB) *gorm.DB {
	db = db.Model(&model.PhysicalMachine{})
	BrandSubDB := q.BrandSub.Scope(db.Session(&gorm.Session{NewDB: true}))
	OsBrandSubDB := q.OsBrandSub.Scope(db.Session(&gorm.Session{NewDB: true}))
	PmListNest_OsBrandSubDB := q.PmListNest.OsBrandSub.Scope(db.Session(&gorm.Session{NewDB: true}))

	if q.UUID == "" && q.Sn == "" && q.MaintainStatus == 0 && len(q.Name) == 0 && q.Sn1 == nil && q.Keyword.Keyword == "" && len(q.ListIn) == 0 && q.PmListNest.MaintainStatus == nil && BrandSubDB == nil && OsBrandSubDB == nil && PmListNest_OsBrandSubDB == nil && q.PointPmListNest == nil {
		return nil
	}

	if q.UUID != "" {

		db = db.Where("uuid = ?", q.UUID)
	}
	if q.Sn != "" {

		db = db.Where("sn = ?", q.Sn)
	}
	if q.MaintainStatus != 0 {

		db = db.Where("maintain_status = ?", q.MaintainStatus)
	}
	if len(q.Name) != 0 {

		db = db.Where("uuid in ?", q.Name)
	}
	if q.Sn1 != nil {

		db = db.Where("sn = ?", *q.Sn1)
	}
	if q.Keyword.Keyword != "" {

		db = db.Where("sn = ?", q.Keyword.Keyword).Or("uuid = ?", q.Keyword.Keyword)
	}
	if len(q.ListIn) != 0 {

		ListInValue := make([][]interface{}, 0, 0)
		for _, v := range q.ListIn {
			ListInValue = append(ListInValue, []interface{}{v.ProductType, v.Brand})
		}
		db = db.Where("(product_type, product_model) IN ?", ListInValue)
	}
	if q.PmListNest.MaintainStatus != nil {

		db = db.Where("maintain_status = ?", *q.PmListNest.MaintainStatus)
	}
	db = db.Where(&q.Match, "maintain_status")

	if BrandSubDB != nil {
		db = db.Where("brand_uuid in (?)", BrandSubDB.Select("uuid"))
	}

	if OsBrandSubDB != nil {
		db = db.Where("os_brand_uuid in (?)", OsBrandSubDB.Select("uuid"))
	}

	if PmListNest_OsBrandSubDB != nil {
		db = db.Where("os_brand_uuid in (?)", PmListNest_OsBrandSubDB.Select("uuid"))
	}
	db = db.Or(func(db *gorm.DB) *gorm.DB {

		if q.PmListKeyword.Keyword == "" {
			return nil
		}

		if q.PmListKeyword.Keyword != "" {

			db = db.Where("sn = ?", q.PmListKeyword.Keyword).Or("uuid = ?", q.PmListKeyword.Keyword)
		}
		return db
	}(db.Session(&gorm.Session{NewDB: true})))
	if q.PointPmListNest != nil {
		PointPmListNest_OsBrandSubDB := q.PointPmListNest.OsBrandSub.Scope(db.Session(&gorm.Session{NewDB: true}))

		if q.PointPmListNest.MaintainStatus == nil && PointPmListNest_OsBrandSubDB == nil {
			return nil
		}

		if q.PointPmListNest.MaintainStatus != nil {

			db = db.Where("maintain_status = ?", *q.PointPmListNest.MaintainStatus)
		}

		if PointPmListNest_OsBrandSubDB != nil {
			db = db.Where("os_brand_uuid in (?)", PointPmListNest_OsBrandSubDB.Select("uuid"))
		}
	}

	return db
}

func (q *PmListBrand) Scope(db *gorm.DB) *gorm.DB {
	db = db.Model(&model.Brand{})

	if q.ProductType == "" && q.Brand == "" {
		return nil
	}

	if q.ProductType != "" {

		db = db.Where("product_type like ?", "%"+q.ProductType+"%")
	}
	if q.Brand != "" {

		db = db.Where("product_model like ?", "%"+q.Brand+"%")
	}

	return db
}

func (q *PmListOsBrand) Scope(db *gorm.DB) *gorm.DB {
	db = db.Model(&model.Brand{})

	if q.ProductType == nil && q.Brand == "" && len(q.ListIn) == 0 {
		return nil
	}

	if q.ProductType != nil {

		db = db.Where("product_type like ?", "%"+*q.ProductType+"%")
	}
	if q.Brand != "" {

		db = db.Where("product_model like ?", "%"+q.Brand+"%")
	}
	if len(q.ListIn) != 0 {

		ListInValue := make([][]interface{}, 0, 0)
		for _, v := range q.ListIn {
			ListInValue = append(ListInValue, []interface{}{v.ProductType, v.Brand})
		}
		db = db.Where("(product_type, product_model) IN ?", ListInValue)
	}

	return db
}
