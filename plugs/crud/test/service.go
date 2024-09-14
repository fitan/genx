package test

import (
	"context"

	"github.com/fitan/genx/plugs/gormq/test/model"
	"gorm.io/gorm"
)

// @crud model.PhysicalMachine
// @crud-list-op-query("")
// @crud-list-op-preload("")
// @crud-create-op-column(select="", omit="")
// @crud-update-method(name="update", select="UUID" ,desc="更新物理机信息")
// @crud-create-method(name="createUser", omit="UUID", desc="创建用户")
type service struct {
	HttpCrud *httpCrud
	DbCrud   *dbCrud
}

type DdCrud interface {
	Create(ctx context.Context, req *model.PhysicalMachine) (err error)
	CreateList(ctx context.Context, req []*model.PhysicalMachine) (err error)
	Delete(ctx context.Context, id []int) (err error)
	Get(ctx context.Context, id int64) (res model.PhysicalMachine, err error)
	List(ctx context.Context, page int, pageSize int, orded string, query query) (list []model.PhysicalMachine, total int64, err error)
	Save(ctx context.Context, id int, req model.PhysicalMachine) (err error)
}

type dbCrud struct {
	db *gorm.DB
}

func (d *dbCrud) Create(ctx context.Context, req *model.PhysicalMachine) (err error) {
	return d.db.WithContext(ctx).Create(req).Error
}

func (d *dbCrud) CreateList(ctx context.Context, req []*model.PhysicalMachine) (err error) {
	return d.db.WithContext(ctx).Create(req).Error
}

func (d *dbCrud) Get(ctx context.Context, id int64) (res model.PhysicalMachine, err error) {
	return res, d.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
}

type query struct {
}

func (d *dbCrud) List(ctx context.Context, page, pageSize int, orded string, query query) (list []model.PhysicalMachine, total int64, err error) {
	return list, total, d.db.WithContext(ctx).Model(&model.PhysicalMachine{}).Count(&total).Error
}

func (d *dbCrud) Save(ctx context.Context, id int, req model.PhysicalMachine) (err error) {
	return d.db.WithContext(ctx).Where("id = ?", id).Save(&req).Error
}

func (d *dbCrud) Delete(ctx context.Context, id []int) (err error) {
	return d.db.WithContext(ctx).Where("id in ?", id).Delete(&model.PhysicalMachine{}).Error
}

type HttpCrud interface {
	// @kit-http / POST
	// @kit-http-request CreateRequest
	Create(ctx context.Context, req CreateRequest) (err error)
	CreateList(ctx context.Context, body []*model.PhysicalMachine) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Get(ctx context.Context, id int64) (res model.PhysicalMachine, err error)
	List(ctx context.Context, page int, pageSize int, orded string, query query) (list []model.PhysicalMachine, total int64, err error)
	Save(ctx context.Context, id int, req model.PhysicalMachine) (err error)
}

type httpCrud struct {
	dbCrud DdCrud
}

func (s *httpCrud) Create(ctx context.Context, req model.PhysicalMachine) (err error) {
	return s.dbCrud.Create(ctx, &req)
}

func (s *httpCrud) CreateList(ctx context.Context, req []*model.PhysicalMachine) (err error) {
	return s.dbCrud.CreateList(ctx, req)
}

func (s *httpCrud) Get(ctx context.Context, id int64) (res model.PhysicalMachine, err error) {
	return s.dbCrud.Get(ctx, id)
}

func (s *httpCrud) List(ctx context.Context, page, pageSize int, orded string, query query) (list []model.PhysicalMachine, total int64, err error) {
	return s.dbCrud.List(ctx, page, pageSize, order, query)
}

func (s *httpCrud) Save(ctx context.Context, id int, req model.PhysicalMachine) (err error) {
	return s.dbCrud.Save(ctx, id, req)
}

func (s *httpCrud) Delete(ctx context.Context, id []int) (err error) {
	return s.dbCrud.Delete(ctx, id)
}
