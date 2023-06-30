package data

import (
	// @extra "gitlab.creditease.corp/paas/paas-assets/src/dnsmiddleware"
	// @extra "gitlab.creditease.corp/paas/paas-assets/src/encode"
	"context"
)

type Middleware func(Service) Service

// @KitHttp
// @log
// @trace
type Service interface {
	// List 获取所有物理机
	// @kitUrl("/device/physicalmachine", "GET")
	// @kitRequest("ListRequest", "true")
	List(ctx context.Context, page, pageSize int, order string, computerRoomUuid, keyword, brandUuid, brand string) (list []ListResponse, total int64, err error)
	// GetByUUID 根据uuid查询物理机
	// @kit-http /device/physicalmachine/{uuid} GET
	// @kit-http-request GetByUUIDRequest
	GetByUUID(ctx context.Context, uuid string) (res ListResponse, err error)
	// Create 创建物理机
	// @kit-http /device/physicalmachine POST
	// @kit-http-request CreateRequest true
	Create(ctx context.Context, createRequest CreateRequest) (res bool, err error)
	// Update 更新物理机
	// @kit-http /device/physicalmachine/{uuid} PUT
	// @kit-http-request UpdateRequest
	Update(ctx context.Context, uuid string, physicalMachine CreateRequest) (res bool, err error)
	// Delete 根据uuid删除物理机
	// @kit-http /device/physicalmachine/{uuid} DELETE
	// @kit-http-request DeleteRequest
	Delete(ctx context.Context, uuid string) (res bool, err error)
	// 解除所有关联关系
	// @kit-http /device/physicalmachine/{uuid}/unlink DELETE
	// @kit-http-request DeleteRequest
	RemoveAllAssociations(ctx context.Context, uuid string) (res bool, err error)
	// 入库
	// @kit-http /device/physicalmachine/{uuid}/save PUT
	// @kit-http-request PutInStorageRequest
	PutInStorage(ctx context.Context, uuid string, body PutInStorageRequestBody) (res bool, err error)
	// 上架
	// @kit-http /device/physicalmachine/{uuid}/up PUT
	// @kit-http-request UpdateRequest
	PutInUse(ctx context.Context, uuid string, physicalMachine CreateRequest) (res bool, err error)

	// 批量录入物理机
	// @kit-http /device/physicalmachine/include POST
	// @kit-http-request IncludeRequest
	// @kit-http-service "" "" encode.CsvResponse
	Include(ctx context.Context, body []byte) (res []byte, err error)
	// 解除项目服务
	// @kit-http /device/physicalmachine/{uuid}/projectservice DELETE
	// @kit-http-request RemoveProjectServiceRequest
	RemoveProjectService(ctx context.Context, uuid string) (res bool, err error)

	// 给物理机分配项目服务
	// @kit-http /device/physicalmachine/projectservice PUT
	// @kit-http-request ProjectServiceRequest
	ProjectService(ctx context.Context, ips []string, namespace, name string) (res bool, err error)
	// 可以分配项目服务的物理机
	// @kit-http /device/physicalmachine/projectservice/available GET
	// @kit-http-request ProjectServiceAvailableRequest
	ProjectServiceAvailable(ctx context.Context, page, pageSize int, ip string) (list []ListResponse, total int64, err error)
	// impi 操作
	// @kit-http /device/physicalmachine/{uuid}/impi/{action} PUT
	// @kit-http-request ImpiRequest
	Impi(ctx context.Context, uuid, action string) (res bool, err error)

	// 把消息放到topic中
	CreateToPmServiceTopic(ctx context.Context, uuid string)
}
