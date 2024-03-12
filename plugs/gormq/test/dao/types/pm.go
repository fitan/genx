package types

//go:generate gowrap gen -g -p ./

// @gq model.PhysicalMachine
type PmList struct {
	// @gq-column uuid
	UUID string `json:"uuid" param:"query,uuid"`
	// @gq-column sn
	Sn string `json:"sn" param:"query,sn"`
	// @gq-sub brand_uuid uuid
	BrandSub PmListBrand
	// @gq-sub os_brand_uuid uuid
	OsBrandSub PmListOsBrand
	// @gq-column maintain_status
	MaintainStatus int `json:"maintainStatus" param:"query,maintainStatus"`
	// @gq-column sn uuid
	//Keyword string `json:"keyword" param:"query,keyword"

	Name []string `json:"name"`

	// @gq-column sn
	Sn1 *string `json:"sn1" param:"query,sn"`

	Keyword PmListKeyword

	ListIn []PmListBrand

	// @gq-group
	// @gq-clause or
	PmListKeyword

	PmListNest

	PointPmListNest *PmListNest

	// @gq-struct MaintainStatus
	Match PmListNest
}

type PmListNest struct {
	// @gq-sub os_brand_uuid uuid
	OsBrandSub PmListOsBrand
	// @gq-column maintain_status
	MaintainStatus *int `json:"maintainStatus" param:"query,maintainStatus"`
}

type PmListKeyword struct {
	// @gq-column sn uuid
	Keyword string `json:"keyword" param:"query,keyword"`
}

// @gq model.Brand
type PmListBrand struct {
	// @gq-column product_type
	// @gq-op like
	ProductType string
	// @gq-op like
	// @gq-column product_model
	Brand string
}

// @gq model.Brand
type PmListOsBrand struct {
	// @gq-column product_type
	// @gq-op like
	ProductType *string
	// @gq-op like
	// @gq-column product_model
	Brand string

	ListIn []PmListBrand
}
