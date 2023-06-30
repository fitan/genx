package kithttp

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

func Gen(pkg *packages.Package, methods []common.InterfaceMethod) {

}

func genEndpointConst(methodNameList []string) jen.Code {
	j := jen.Null()

	for _, methodName := range methodNameList {
		j.Const().Id(methodName + "MethodName").Op("=").Lit(methodName).Line()
	}

	j.Var().Id("MethodNameList").Op("=").Index().String().ValuesFunc(func(g *jen.Group) {
		for _, methodName := range methodNameList {
			g.Id(methodName + "MethodName")
		}
	})

	return j
}

func genEndpoints(methodNameList []string) jen.Code {
	listCode := make([]jen.Code, 0, len(methodNameList))
	for _, methodName := range methodNameList {
		listCode = append(listCode, jen.Id(methodName+"Endpoint").Qual("github.com/go-kit/kit/endpoint", "Endpoint"))
	}
	return jen.Null().Type().Id("Endpoints").Struct(
		listCode...,
	)
}

func genNewEndpoint(methodNameList []string) jen.Code {
	endpointVarList := make([]jen.Code, 0, len(methodNameList))
	endpointForList := make([]jen.Code, 0, len(methodNameList))

	for _, methodName := range methodNameList {
		endpointVarList = append(endpointVarList, jen.Id(methodName+"Endpoint").Op(":").Id("make"+methodName+"Endpoint").Call(jen.Id("s")))

		endpointForList = append(endpointForList, jen.For(jen.List(jen.Id("_"), jen.Id("m")).Op(":=").Range().Id("dmw").Index(jen.Id(methodName+"MethodName"))).Block(jen.Id("eps").Dot(methodName+"Endpoint").Op("=").Id("m").Call(jen.Id("eps").Dot(methodName+"Endpoint"))).Line())
	}

	endpointForListStatement := jen.Statement(endpointForList)

	return jen.Func().Id("NewEndpoint").Params(jen.Id("s").Id("Service"), jen.Id("dmw").Map(jen.Id("string")).Index().Qual("github.com/go-kit/kit/endpoint", "Middleware")).Params(jen.Id("Endpoints")).Block(
		jen.Id("eps").Op(":=").Id("Endpoints").Values(
			endpointVarList...,
		),
		&endpointForListStatement,
		jen.Return().Id("eps"),
	)
}
