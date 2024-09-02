package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

func Gen(pkg *packages.Package, doc common.Doc, methods []common.InterfaceMethod) {
	slog.Info("gen cache", slog.Any("pkg", pkg), slog.Any("methods", methods))

	var keyPrefix string
	_ = doc.ByFuncNameAndArgs("@cache", &keyPrefix)

	j := jen.NewFile(pkg.Name)
	common.JenAddImports(pkg, j)
	j.AddImport("context", "")
	j.AddImport("encoding/json", "")
	j.AddImport("log/slog", "")
	j.AddImport("time", "")
	j.AddImport("github.com/samber/do/v2", "")
	j.AddImport("github.com/redis/go-redis/v9", "")
	j.AddImport("github.com/pkg/errors", "")
	j.Add(genCacheStruct())
	j.Add()

	funcList := make([]jen.Code, 0)

	for _, method := range methods {
		funcList = append(funcList, genCacheFunc(pkg.Name, keyPrefix, method))
	}

	j.Add(funcList...)
	j.Add(genNewCache(pkg.Name))

	common.WriteGO("cache.go", j.GoString())
}

func genCacheFunc(serviceName string, keyPrefix string, method common.InterfaceMethod) jen.Code {
	methodParamCode := make([]jen.Code, 0)
	methodParamCode = append(methodParamCode, jen.Id("ctx").Qual("context", "Context"))
	methodResultCode := make([]jen.Code, 0)

	nextMethodParamCode := make([]jen.Code, 0)
	for _, param := range method.Params {
		nextMethodParamCode = append(nextMethodParamCode, jen.Id(param.Name))
	}

	for _, param := range method.ParamsExcludeCtx() {
		methodParamCode = append(methodParamCode, jen.Id(param.Name).Id(param.ID))
	}

	for _, param := range method.Results {
		methodResultCode = append(methodResultCode, jen.Id(param.Name).Id(param.ID))
	}

	return jen.Func().Params(jen.Id("s").Op("*").Id("cache")).Id(method.Name).Params(
		methodParamCode...,
	).Params(
		methodResultCode...,
	).BlockFunc(func(g *jen.Group) {
		var key string
		var exp string
		has := method.Doc.ByFuncNameAndArgs("@cache-get", &key, &exp)
		if has {
			g.Add(cacheGet(method, serviceName, keyPrefix, key, exp))
		}

		has = method.Doc.ByFuncNameAndArgs("@cache-del", &key)
		if has {
			g.Add(cacheDel(method, serviceName, keyPrefix, key))
		}

	}).Line()
}

func cacheGet(method common.InterfaceMethod, serviceName, keyPrefix, key, exp string) jen.Code {
	td, err := time.ParseDuration(exp)
	if err != nil {
		slog.Error("parse duration error", "method", method.Name, "exp", exp, "error", err)
		panic(err)
	}

	j := jen.Null()
	j.Var().Id("_key").Op("=").Id(`fmt.Sprintf`).Call(jen.Lit("%s.%d"), jen.Lit(keyPrefix), jen.Id(key)).Line()

	j.Id("log").Op(":=").Id("s.logger.With").Call(
		jen.Lit("middleware"), jen.Lit("cache"),
		jen.Lit("serviceName"), jen.Lit(serviceName),
		jen.Lit("method"), jen.Lit(method.Name),
		jen.Lit("key"), jen.Id("_key"),
	).Line()

	j.Id("_val, err").Op(":=").Id("s.redis.Get").Call(jen.Id("ctx"), jen.Lit(key)).Dot("Result").Call().Line()

	j.Id("if err != nil").BlockFunc(func(g *jen.Group) {
		g.Line()
		g.Id("if err == redis.Nil").Block(
			jen.Id(`log.WarnContext(ctx, "redis get nil", "error", err)`)).Else().Block(
			jen.Id(`log.ErrorContext(ctx, "redis get error", "error", err)`),
		)

		results := lo.Map(method.Results, func(item common.MethodParam, index int) string {
			return item.Name
		})

		params := lo.Map(method.Params, func(item common.MethodParam, index int) string {
			return item.Name
		})
		g.Id(strings.Join(results, ",")).Op("=").Id("s.next." + method.Name).Call(jen.Id(strings.Join(params, ",")))

		if method.ReturnsError {
			g.Id("if err != nil").Block(
				jen.Id("return"),
			)
		}

		g.Var().Id("_b").Id("[]byte")

		if len(method.ResultsExcludeErr()) > 1 {
			g.Id("_res").Op(":=").Id(method.ResultsMapValPointExcludeErr())
			g.Id("_b, err").Op("=").Id("json.Marshal").Call(jen.Id("_res"))
		} else {
			g.Id("_b, err").Op("=").Id("json.Marshal").Call(jen.Id(results[0]))
		}

		g.Id(`if err != nil {
			log.ErrorContext(ctx, "set redis before marshal error", "error", err)
			err = nil
			return
		}`)

		g.Id("err").Op("=").Id("s.redis.Set").Call(jen.Id("ctx"), jen.Lit(key), jen.Id("string(_b)"), jen.Id(fmt.Sprintf("%v", td.Seconds())+"*time.Second")).Dot("Err").Call()

		g.Id(`if err != nil {
			log.ErrorContext(ctx, "set redis error", "error", err)
			err = nil
			return
		}

		return`)
	}).Line()

	if len(method.ResultsExcludeErr()) > 1 {
		j.Id("_res").Op(":=").Id(method.ResultsMapValPointExcludeErr()).Line()
		j.Id("err").Op("=").Id("json.Unmarshal").Call(jen.Id("[]byte(_val)"), jen.Id("&_res")).Line()
	} else {
		j.Id("err").Op("=").Id("json.Unmarshal").Call(jen.Id("[]byte(_val)"), jen.Id("&"+method.Results[0].Name)).Line()
	}
	j.Id(`if err != nil {
			log.ErrorContext(ctx, "set redis error", "error", err)
			err = nil
			return
		}`).Line()
	j.Id("return")

	return j
}

func cacheDel(method common.InterfaceMethod, serviceName, keyPrefix, key string) jen.Code {
	t := `
	defer func() {
		if err == nil {
			delErr := s.redis.Del(ctx, %s).Err()
			if delErr != nil {
				log.ErrorContext(ctx, "del redis error", "error", delErr)
			}
		}
	}()
	`

	params := lo.Map(method.Params, func(item common.MethodParam, index int) string {
		return item.Name
	})
	j := jen.Null()
	j.Var().Id("_key").Op("=").Id(`fmt.Sprintf`).Call(jen.Lit("%s.%d"), jen.Lit(keyPrefix), jen.Id(key)).Line()
	j.Id("log").Op(":=").Id("s.logger.With").Call(
		jen.Lit("middleware"), jen.Lit("cache"),
		jen.Lit("serviceName"), jen.Lit(serviceName),
		jen.Lit("method"), jen.Lit(method.Name),
		jen.Lit("key"), jen.Id("_key"),
	).Line()

	j.Id(fmt.Sprintf(t, "_key")).Line()

	return j.Return().Id("s.next." + method.Name).Call(jen.Id(strings.Join(params, ","))).Line()

}

func genCacheStruct() jen.Code {
	return jen.Null().Type().Id("cache").Struct(
		jen.Id("pkgName").String(),
		jen.Id("next").Id("Service"),
		jen.Id("logger").Id("*slog.Logger"),
		jen.Id("redis").Id("*redis.Client"),
	)
}

func genNewCache(pkgName string) jen.Code {
	return jen.Func().Id("NewCache").Params(jen.Id("i").Id("do.Injector")).Params(jen.Id("Middleware")).BlockFunc(
		func(g *jen.Group) {
			g.Return().Func().Params(jen.Id("next").Id("Service")).Params(jen.Id("Service")).Block(
				jen.Return().Op("&").Id("cache").Values(
					jen.Id("next").Op(":").Id("next"),
					jen.Id("logger").Op(":").Id(`do.MustInvoke[*slog.Logger](i)`),
					jen.Id("redis").Op(":").Id(`do.MustInvoke[*redis.Client](i)`),
				),
			)
		},
	)
}
