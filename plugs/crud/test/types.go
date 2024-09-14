package test

import "github.com/fitan/genx/plugs/gormq/test/model"

type CreateRequest struct {
	Name string
}

type CreateListRequest struct {
	Body []*model.PhysicalMachine `json:"body" param:"body,body"`
}

type DeleteRequest struct {
	Id int `json:"id" param:"path,id"`
}

type GetRequest struct {
	Id int `json:"id" param:"path,id"`
}

type ListRequest struct {
	// @gq-op =
	Name string
}

type SaveRequest struct {
	Name string
}
