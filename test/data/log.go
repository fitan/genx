package data

import (
	"context"
	json1 "encoding/json"
	"time"

	log "github.com/go-kit/kit/log"
	level "github.com/go-kit/kit/log/level"
)

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Create(ctx context.Context, createRequest CreateRequest) (res bool, err error) {
	createRequestByte, _ := json1.Marshal(createRequest)
	createRequestJson := string(createRequestByte)

	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "Create", "createRequest", createRequestJson, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Create(ctx, createRequest)
}
func (s *logging) CreateToPmServiceTopic(ctx context.Context, uuid string) {
	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "CreateToPmServiceTopic", "uuid", uuid, "took", time.Since(begin))
	}(time.Now())
	s.next.CreateToPmServiceTopic(ctx, uuid)
}
func (s *logging) Delete(ctx context.Context, uuid string) (res bool, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "Delete", "uuid", uuid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Delete(ctx, uuid)
}
func (s *logging) GetByUUID(ctx context.Context, uuid string) (res ListResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "GetByUUID", "uuid", uuid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.GetByUUID(ctx, uuid)
}
func (s *logging) Impi(ctx context.Context, uuid string, action string) (res bool, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "Impi", "uuid", uuid, "action", action, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Impi(ctx, uuid, action)
}
func (s *logging) Include(ctx context.Context, body []byte) (res []byte, err error) {
	bodyByte, _ := json1.Marshal(body)
	bodyJson := string(bodyByte)

	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "Include", "body", bodyJson, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Include(ctx, body)
}
func (s *logging) List(ctx context.Context, page int, pageSize int, order string, computerRoomUuid string, keyword string, brandUuid string, brand string) (list []ListResponse, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "List", "page", page, "pageSize", pageSize, "order", order, "computerRoomUuid", computerRoomUuid, "keyword", keyword, "brandUuid", brandUuid, "brand", brand, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.List(ctx, page, pageSize, order, computerRoomUuid, keyword, brandUuid, brand)
}
func (s *logging) ProjectService(ctx context.Context, ips []string, namespace string, name string) (res bool, err error) {
	ipsByte, _ := json1.Marshal(ips)
	ipsJson := string(ipsByte)

	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "ProjectService", "ips", ipsJson, "namespace", namespace, "name", name, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.ProjectService(ctx, ips, namespace, name)
}
func (s *logging) ProjectServiceAvailable(ctx context.Context, page int, pageSize int, ip string) (list []ListResponse, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "ProjectServiceAvailable", "page", page, "pageSize", pageSize, "ip", ip, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.ProjectServiceAvailable(ctx, page, pageSize, ip)
}
func (s *logging) PutInStorage(ctx context.Context, uuid string, body PutInStorageRequestBody) (res bool, err error) {
	bodyByte, _ := json1.Marshal(body)
	bodyJson := string(bodyByte)

	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "PutInStorage", "uuid", uuid, "body", bodyJson, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.PutInStorage(ctx, uuid, body)
}
func (s *logging) PutInUse(ctx context.Context, uuid string, physicalMachine CreateRequest) (res bool, err error) {
	physicalMachineByte, _ := json1.Marshal(physicalMachine)
	physicalMachineJson := string(physicalMachineByte)

	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "PutInUse", "uuid", uuid, "physicalMachine", physicalMachineJson, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.PutInUse(ctx, uuid, physicalMachine)
}
func (s *logging) RemoveAllAssociations(ctx context.Context, uuid string) (res bool, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "RemoveAllAssociations", "uuid", uuid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.RemoveAllAssociations(ctx, uuid)
}
func (s *logging) RemoveProjectService(ctx context.Context, uuid string) (res bool, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "RemoveProjectService", "uuid", uuid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.RemoveProjectService(ctx, uuid)
}
func (s *logging) Update(ctx context.Context, uuid string, physicalMachine CreateRequest) (res bool, err error) {
	physicalMachineByte, _ := json1.Marshal(physicalMachine)
	physicalMachineJson := string(physicalMachineByte)

	defer func(begin time.Time) {
		_ = s.logger.Log(s.traceId, ctx.Value(s.traceId), "method", "Update", "uuid", uuid, "physicalMachine", physicalMachineJson, "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Update(ctx, uuid, physicalMachine)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "github.com/fitan/genx/test/data", "logging")
	return func(next Service) Service {
		return &logging{logger: level.Info(logger), next: next, traceId: traceId}
	}
}
