package data

import (
	"context"
	json1 "encoding/json"

	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (s *tracing) Create(ctx context.Context, createRequest CreateRequest) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Create", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		createRequestByte, _ := json1.Marshal(createRequest)
		createRequestJson := string(createRequestByte)

		span.LogKV("createRequest", createRequestJson, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Create(ctx, createRequest)
}
func (s *tracing) CreateToPmServiceTopic(ctx context.Context, uuid string) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CreateToPmServiceTopic", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		span.LogKV("uuid", uuid)
		span.Finish()
	}()
	s.next.CreateToPmServiceTopic(ctx, uuid)
}
func (s *tracing) Delete(ctx context.Context, uuid string) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Delete", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		span.LogKV("uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Delete(ctx, uuid)
}
func (s *tracing) GetByUUID(ctx context.Context, uuid string) (res ListResponse, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GetByUUID", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		span.LogKV("uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.GetByUUID(ctx, uuid)
}
func (s *tracing) Impi(ctx context.Context, uuid string, action string) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Impi", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		span.LogKV("uuid", uuid, "action", action, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Impi(ctx, uuid, action)
}
func (s *tracing) Include(ctx context.Context, body []byte) (res []byte, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Include", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		bodyByte, _ := json1.Marshal(body)
		bodyJson := string(bodyByte)

		span.LogKV("body", bodyJson, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Include(ctx, body)
}
func (s *tracing) List(ctx context.Context, page int, pageSize int, order string, computerRoomUuid string, keyword string, brandUuid string, brand string) (list []ListResponse, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "List", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		span.LogKV("page", page, "pageSize", pageSize, "order", order, "computerRoomUuid", computerRoomUuid, "keyword", keyword, "brandUuid", brandUuid, "brand", brand, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.List(ctx, page, pageSize, order, computerRoomUuid, keyword, brandUuid, brand)
}
func (s *tracing) ProjectService(ctx context.Context, ips []string, namespace string, name string) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ProjectService", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		ipsByte, _ := json1.Marshal(ips)
		ipsJson := string(ipsByte)

		span.LogKV("ips", ipsJson, "namespace", namespace, "name", name, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ProjectService(ctx, ips, namespace, name)
}
func (s *tracing) ProjectServiceAvailable(ctx context.Context, page int, pageSize int, ip string) (list []ListResponse, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ProjectServiceAvailable", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		span.LogKV("page", page, "pageSize", pageSize, "ip", ip, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.ProjectServiceAvailable(ctx, page, pageSize, ip)
}
func (s *tracing) PutInStorage(ctx context.Context, uuid string, body PutInStorageRequestBody) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PutInStorage", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		bodyByte, _ := json1.Marshal(body)
		bodyJson := string(bodyByte)

		span.LogKV("uuid", uuid, "body", bodyJson, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.PutInStorage(ctx, uuid, body)
}
func (s *tracing) PutInUse(ctx context.Context, uuid string, physicalMachine CreateRequest) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PutInUse", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		physicalMachineByte, _ := json1.Marshal(physicalMachine)
		physicalMachineJson := string(physicalMachineByte)

		span.LogKV("uuid", uuid, "physicalMachine", physicalMachineJson, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.PutInUse(ctx, uuid, physicalMachine)
}
func (s *tracing) RemoveAllAssociations(ctx context.Context, uuid string) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RemoveAllAssociations", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		span.LogKV("uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.RemoveAllAssociations(ctx, uuid)
}
func (s *tracing) RemoveProjectService(ctx context.Context, uuid string) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RemoveProjectService", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		span.LogKV("uuid", uuid, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.RemoveProjectService(ctx, uuid)
}
func (s *tracing) Update(ctx context.Context, uuid string, physicalMachine CreateRequest) (res bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Update", opentracing.Tag{Key: string(ext.Component), Value: "github.com/fitan/genx/test/data"})
	defer func() {
		physicalMachineByte, _ := json1.Marshal(physicalMachine)
		physicalMachineJson := string(physicalMachineByte)

		span.LogKV("uuid", uuid, "physicalMachine", physicalMachineJson, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return s.next.Update(ctx, uuid, physicalMachine)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{next: next, tracer: otTracer}
	}
}
