package data

import (
	"github.com/fitan/genx/test/data/types"
	mygorm "github.com/fitan/mykit/mygorm"
	"gorm.io/gorm"
)

func (g PMListReq) GormQScopes(res []func(*gorm.DB) *gorm.DB, err error) {
	req := make([]mygorm.GenxScopesReq, 0)

	req = append(req, mygorm.GenxScopesReq{
		Field: "UUID",
		Op:    "\"=\"",
		Value: g.UUID,
	})
	req = append(req, mygorm.GenxScopesReq{
		Field: "Brand.UUID",
		Op:    "\"?=\"",
		Value: g.BrandUUID,
	})
	req = append(req, mygorm.GenxScopesReq{
		Field: "Brand.Users.Name",
		Op:    "=",
		Value: g.BrandUsersName,
	})
	req = append(req, mygorm.GenxScopesReq{
		Field: "Brand.ID",
		Op:    "\">\"",
		Value: g.BrandID,
	})
	req = append(req, mygorm.GenxScopesReq{
		Field: "Brand.CreatedAt",
		Op:    "\"><\"",
		Value: g.BrandCreatedAt,
	})

	return mygorm.GenxScopes(types.PhysicalMachine, req)
}
