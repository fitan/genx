package kithttp

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"strings"
	"unicode"

	"github.com/fitan/genx/common"
	"github.com/fitan/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

const (
	RequestParamTagName string = "param"
	QueryTag            string = "query"
	HeaderTag           string = "header"
	PathTag             string = "path"
	BodyTag             string = "body"
	CtxTag              string = "ctx"

	DocKitHttpParamMark string = "@kit-http-param"
)

type KitRequest struct {
	pkg *packages.Package

	ServiceName   string
	RequestTypeOf *types.Struct
	RequestName   string
	RequestIsBody bool
	RequestIsNil  bool

	NamedMap map[string]string
	Query    map[string]RequestParam
	Path     map[string]RequestParam
	Body     map[string]RequestParam
	Header   map[string]RequestParam
	Ctx      map[string]RequestParam
	Empty    map[string]RequestParam
}

func NewKitRequest(pkg *packages.Package, serviceName, requestName string, requestIsBody bool) *KitRequest {
	return &KitRequest{
		pkg:           pkg,
		ServiceName:   serviceName,
		RequestName:   requestName,
		RequestIsBody: requestIsBody,
		NamedMap:      make(map[string]string),
		Query:         make(map[string]RequestParam),
		Path:          make(map[string]RequestParam),
		Body:          make(map[string]RequestParam),
		Header:        make(map[string]RequestParam),
		Ctx:           make(map[string]RequestParam),
		Empty:         make(map[string]RequestParam),
	}
}

type RequestParam struct {
	ParamDoc *ast.CommentGroup

	ParamPath string

	FieldName string

	ParamName string
	// path, query, header, body, empty
	ParamSource string
	// basic, map, slice,ptr
	ParamType string
	// time [time.Time]
	ParamTypeName string
	// int,string,bool,float
	BasicType string

	HasPtr bool
}

func (r RequestParam) Annotations() string {
	if r.ParamDoc == nil {
		return ""
	}
	var annotations string
	format := &AstDocFormat{r.ParamDoc}
	format.MarkValuesMapping(r.FieldName, &annotations)
	if annotations != "" {
		return annotations
	}
	return fmt.Sprintf(`"%s"`, strings.TrimPrefix(r.ParamDoc.List[0].Text, "// "))
}

func (r RequestParam) ToVal() jen.Code {
	return jen.Var().Id("_" + r.ParamName).Id(r.ParamTypeName)
}

func (k *KitRequest) ParamPath(paramName string) (res string) {
	defer func() {
		if res != "" {
			res = "." + res
		}
	}()

	if strings.ToUpper(paramName) == strings.ToUpper(k.RequestName) {
		return ""
	}
	param, ok := k.Query[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Path[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Header[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Body[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Ctx[paramName]
	if ok {
		return param.ParamPath
	}

	param, ok = k.Empty[paramName]
	if ok {
		return param.ParamPath
	}

	panic("param not found: " + paramName)
}

func (k *KitRequest) SetParam(param RequestParam) {
	switch param.ParamSource {
	case QueryTag:
		k.Query[param.ParamName] = param
	case PathTag:
		k.Path[param.ParamName] = param
	case HeaderTag:
		k.Header[param.ParamName] = param
	case BodyTag:
		k.Body[param.ParamName] = param
	case CtxTag:
		k.Ctx[param.ParamName] = param
	case "":
		k.Empty[param.ParamName] = param

	default:
		panic("param source error: " + param.ParamSource + "," + param.ParamName)
	}
}

func (k *KitRequest) ParseParamTag(fieldName, tag string) (paramSource string, paramName string) {

	split := strings.Split(tag, ",")
	if len(split) == 1 {
		return split[0], downFirst(fieldName)
	}

	if len(split) == 2 {
		return split[0], split[1]
	}

	return "", ""

}

func (k *KitRequest) Statement() *jen.Statement {
	listCode := make([]jen.Code, 0, 0)
	// req := Request{}
	listCode = append(listCode, jen.Id("req").Op(":=").Id(k.RequestName).Block())
	listCode = append(listCode, k.DefineVal()...)
	listCode = append(listCode, k.BindPathParam()...)
	listCode = append(listCode, k.BindQueryParam()...)
	listCode = append(listCode, k.BindHeaderParam()...)
	listCode = append(listCode, k.BindBodyParam()...)
	listCode = append(listCode, k.BindCtxParam()...)
	listCode = append(listCode, k.BindRequest()...)
	listCode = append(listCode, k.Validate()...)
	listCode = append(listCode, jen.Return(jen.Id("req"), jen.Id("err")))
	var LineListCode []jen.Code
	for _, v := range listCode {
		LineListCode = append(LineListCode, jen.Line(), v)
	}

	fn := jen.Func().Id("decode"+upFirst(k.ServiceName)+"Request").Params(
		jen.Id("ctx").Id("context.Context"),
		jen.Id("r").Op("*").Qual("net/http", "Request"),
	).Call(
		jen.Id("res").Interface(),
		jen.Id("err").Id("error"),
	).Block(
		LineListCode...,
	)

	return fn
}

func (k *KitRequest) DecodeRequest() (s string) {
	return k.Statement().GoString()
}

func (k *KitRequest) DefineVal() []jen.Code {
	listCode := make([]jen.Code, 0, 0)
	for _, v := range k.Query {
		listCode = append(listCode, v.ToVal())
	}
	for _, v := range k.Path {
		listCode = append(listCode, v.ToVal())
	}
	for _, v := range k.Header {
		listCode = append(listCode, v.ToVal())
	}

	for _, v := range k.Ctx {
		listCode = append(listCode, v.ToVal())
	}
	return listCode
}

func (k *KitRequest) Validate() []jen.Code {
	list := make([]jen.Code, 0, 0)
	list = append(
		list,
		jen.List(jen.Id("validReq"), jen.Id("err")).Op(":=").Qual("github.com/asaskevich/govalidator", "ValidateStruct").Call(jen.Id("req")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Err().Op("=").Qual("github.com/pkg/errors", "Wrap").Call(jen.Id("err"), jen.Lit("govalidator.ValidateStruct")),
			jen.Return(),
		),
		jen.If(jen.Id("!validReq")).Block(
			jen.Err().Op("=").Qual("github.com/pkg/errors", "Wrap").Call(jen.Id("err"), jen.Lit("valid false")),
			jen.Return(),
		),
	)
	return list
}

func (k *KitRequest) BindBodyParam() []jen.Code {
	listCode := make([]jen.Code, 0, 0)
	if k.RequestIsNil {
		return listCode
	}
	returnCode := jen.If(jen.Err().Op("!=").Nil()).Block(
		jen.Err().Op("=").Qual("github.com/pkg/errors", "Wrap").Call(jen.Id("err"), jen.Lit("json.Decode")),
		jen.Return(),
	)
	if k.RequestIsBody {
		// err = json.NewDecoder(r.Body).Decode(&req)
		decode := jen.Id("err").Op("=").Qual("encoding/json", "NewDecoder").Call(jen.Id("r.Body")).Dot("Decode").Parens(jen.Id("&req"))
		listCode = append(listCode, decode, returnCode)

		return listCode
	}

	if len(k.Body) == 0 {
		return listCode
	}

	if len(k.Body) != 1 {
		panic("body param count error " + fmt.Sprint(len(k.Body)))
	}

	for _, v := range k.Body {
		decode := jen.Id("err").Op("=").Qual("encoding/json", "NewDecoder").Call(jen.Id("r.Body")).Dot("Decode").Parens(jen.Id("&req." + v.ParamPath))
		listCode = append(listCode, decode, returnCode)
	}

	return listCode
}

func (k *KitRequest) BindHeaderParam() []jen.Code {
	list := make([]jen.Code, 0, 0)

	for _, v := range k.Header {
		//r.Header.Get("project")
		varBind := jen.Id("r.Header.Get").Call(jen.Lit(v.ParamName))
		if v.BasicType != "string" {
			// cast.ToInt(vars["id"])
			varBind = jen.Id("cast").Dot("To" + upFirst(v.BasicType) + "E").Call(varBind)
			// id, err := cast.ToIntE(vars["id"])
			varBind = jen.List(jen.Id("_"+v.ParamName), jen.Err()).Op("=").Add(varBind)
			// if err != nil {
			// 	return err
			// }
			returnCode := jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(),
			)
			list = append(list, varBind, returnCode)
			continue
		}
		// id = vars["id"]
		val := jen.Id("_" + v.ParamName).Op("=").Add(varBind)
		list = append(list, val)
	}

	return list

}

func (k *KitRequest) BindQueryParam() []jen.Code {
	list := make([]jen.Code, 0, 0)

	if len(k.Query) == 0 {
		return list
	}

	for _, v := range k.Query {
		//r.URL.Query().Get("project")
		varBind := jen.Id("r.URL.Query().Get").Call(jen.Lit(v.ParamName))
		if !(v.ParamType == "basic" && v.BasicType == "string") {
			castCode, err := CastMap(v.ParamName, v.ParamTypeName, varBind)
			if err != nil {
				panic(err)
			}
			list = append(list, castCode...)
			continue
		}
		// id = vars["id"]
		val := jen.Id("_" + v.ParamName).Op("=").Add(varBind)
		list = append(list, val)
	}

	return list
}

func (k *KitRequest) BindPathParam() []jen.Code {
	list := make([]jen.Code, 0, 0)

	if len(k.Path) == 0 {
		return list
	}

	// vars := mux.Vars(r)
	vars := jen.Id("vars").Op(":=").Qual("github.com/gorilla/mux", "Vars").Call(jen.Id("r"))
	list = append(list, vars)
	for _, v := range k.Path {
		// vars["id"]
		varBind := jen.Id("vars").Index(jen.Lit(v.ParamName))
		if v.BasicType != "string" {
			// cast.ToInt(vars["id"])
			varBind = jen.Id("cast").Dot("To" + upFirst(v.BasicType) + "E").Call(varBind)
			// id, err := cast.ToIntE(vars["id"])
			varBind = jen.List(jen.Id("_"+v.ParamName), jen.Err()).Op("=").Add(varBind)
			// if err != nil {
			// 	return err
			// }
			returnCode := jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(),
			)
			list = append(list, varBind, returnCode)
			continue
		}
		// id = vars["id"]
		val := jen.Id("_" + v.ParamName).Op("=").Add(varBind)
		list = append(list, val)
	}

	return list
}

func (k *KitRequest) BindCtxParam() []jen.Code {
	list := make([]jen.Code, 0, 0)
	for _, v := range k.Ctx {
		var ctxKey string

		if v.ParamDoc == nil {
			panic("ctx param doc is nil")
		}
		for _, d := range v.ParamDoc.List {
			fields := strings.Fields(d.Text)
			if fields[1] == DocKitHttpParamMark {
				if len(fields) < 3 {
					panic("ctx param doc error: " + d.Text)
				}

				if fields[2] == "ctx" {
					ctxKey = fields[3]
				}

			}
		}
		if ctxKey == "" {
			panic("not find ctx param doc error: " + v.ParamDoc.Text())
		}
		ctxVal := jen.Var().Id(v.ParamName + "OK").Bool()
		varBind := jen.List(jen.Id("_"+v.ParamName), jen.Id(v.ParamName+"OK")).Op("=").Id("ctx.Value").Call(jen.Id(ctxKey)).Assert(jen.Id(v.ParamTypeName))
		ifBind := jen.If(jen.Id(v.ParamName+"OK")).Op("==").False().Block(
			jen.Err().Op("=").Id("errors.New").Call(jen.Lit("ctx param "+v.ParamName+" is not found")),
			jen.Return(),
		)
		list = append(list, ctxVal, varBind, ifBind)
	}
	return list
}

func (k *KitRequest) BindRequest() []jen.Code {
	list := make([]jen.Code, 0, 0)
	for _, v := range k.Path {
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Id("_" + v.ParamName)
		list = append(list, reqBindVal)
	}

	for _, v := range k.Query {
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Id("_" + v.ParamName)
		list = append(list, reqBindVal)
	}

	for _, v := range k.Header {
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Id("_" + v.ParamName)
		list = append(list, reqBindVal)
	}

	for _, v := range k.Ctx {
		reqBindVal := jen.Id("req").Dot(v.ParamPath).Op("=").Id("_" + v.ParamName)
		list = append(list, reqBindVal)
	}
	return list
}

func (k *KitRequest) ParseRequest() {
	var hasFindRequest bool
	for _, s := range k.pkg.Syntax {
		ast.Inspect(s, func(node ast.Node) bool {
			switch nodeT := node.(type) {
			case *ast.GenDecl:
				for _, spec := range nodeT.Specs {
					if specT, ok := spec.(*ast.TypeSpec); ok {
						if specT.Name.Name == k.RequestName {
							hasFindRequest = true
							doc := nodeT.Doc
							k.Doc(doc)
							t := k.pkg.TypesInfo.TypeOf(specT.Type).(*types.Struct)
							k.RequestTypeOf = t
							k.RequestType([]string{}, k.RequestName, t, "", doc)
							k.CheckRequestIsNil()

							return false
						}
					}
				}

			}
			return true
		})
	}
	if !hasFindRequest {
		panic("not find request" + k.RequestName)
	}
}

func (k *KitRequest) Doc(doc *ast.CommentGroup) {
}

func (k *KitRequest) ParseFieldComment(pos token.Pos) (s *ast.CommentGroup) {
	fieldFileName := k.pkg.Fset.Position(pos).Filename
	fieldLine := k.pkg.Fset.Position(pos).Line
	var fieldComment *ast.CommentGroup
	for _, syntax := range k.pkg.Syntax {
		fileName := k.pkg.Fset.Position(syntax.Pos()).Filename
		if fieldFileName == fileName {
			for _, c := range syntax.Comments {
				if k.pkg.Fset.Position(c.End()).Line+1 == fieldLine {
					fieldComment = c
				}
			}
			break
		}
	}
	return fieldComment

	//if fieldComment == nil {
	//	return ""
	//}
	//
	//for _, c := range fieldComment {
	//	commentField := strings.Fields(c.Text)
	//	if len(commentField) < 3 {
	//		panic("comment error: " + c.Text)
	//	}
	//	fmt.Println("commentField", commentField)
	//	if commentField[0] == "@kit-request" && commentField[1] == "ctx" {
	//		return commentField[2]
	//	}
	//}
	//return ""
}

func (k *KitRequest) CheckRequestIsNil() {
	if k.RequestIsBody {
		if k.RequestTypeOf.NumFields() == 0 {
			k.RequestIsNil = true
		}
	}
}

func (k *KitRequest) RequestType(prefix []string, requestName string, requestType types.Type, requestParamTagTypeTag string, doc *ast.CommentGroup) {
	rawParamType := common.TypeOf(requestType).TypeAsJenComparePkgNameString(k.pkg)
	paramSource, paramName := k.ParseParamTag(requestName, requestParamTagTypeTag)

	switch rt := requestType.(type) {
	case *types.Named:
		k.NamedMap[paramName] = rawParamType
		k.RequestType(prefix, requestName, rt.Underlying(), requestParamTagTypeTag, doc)
	case *types.Struct:
		var paramTypeName string
		var ok bool
		if paramTypeName, ok = k.NamedMap[paramName]; !ok {
			paramTypeName = common.TypeOf(requestType).TypeAsJenComparePkgNameString(k.pkg)
		}

		k.SetParam(RequestParam{
			FieldName:     requestName,
			ParamDoc:      doc,
			ParamPath:     strings.Join(prefix, "."),
			ParamName:     paramName,
			ParamSource:   paramSource,
			ParamType:     "struct",
			ParamTypeName: paramTypeName,
			BasicType:     paramTypeName,
			HasPtr:        false,
		})
		for i := 0; i < rt.NumFields(); i++ {
			field := rt.Field(i)
			fieldName := field.Name()
			fieldType := field.Type()
			tag, _ := reflect.StructTag(rt.Tag(i)).Lookup(RequestParamTagName)
			k.RequestType(append(prefix, fieldName), fieldName, fieldType, tag, k.ParseFieldComment(field.Pos()))
		}
	case *types.Pointer:
		k.RequestType(prefix, requestName, rt.Elem().Underlying(), requestParamTagTypeTag, doc)
	case *types.Slice:
		var paramTypeName string
		var ok bool
		if paramTypeName, ok = k.NamedMap[paramName]; !ok {
			paramTypeName = common.TypeOf(requestType).TypeAsJenComparePkgNameString(k.pkg)
		}
		k.SetParam(RequestParam{
			FieldName:     requestName,
			ParamDoc:      doc,
			ParamPath:     strings.Join(prefix, "."),
			ParamName:     paramName,
			ParamSource:   paramSource,
			ParamType:     "slice",
			ParamTypeName: paramTypeName,
			BasicType:     rt.Elem().Underlying().String(),
			HasPtr:        false,
		})
	case *types.Map:
		//var paramTypeName string
		//var ok bool
		//if paramTypeName, ok = k.NamedMap[paramName]; !ok {
		//	paramTypeName = rt.Elem().Underlying().String()
		//}
		//k.SetParam(RequestParam{
		//	FieldName:     requestName,
		//	ParamDoc:      doc,
		//	ParamPath:     strings.Join(prefix, "."),
		//	ParamName:     paramName,
		//	ParamSource:   paramSource,
		//	ParamType:     "map",
		//	ParamTypeName: paramTypeName,
		//	ParamTypeName:  rawParamType,
		//	BasicType:     rt.Elem().Underlying().String(),
		//	HasPtr:        false,
		//})
	case *types.Basic:
		var paramTypeName string
		var ok bool
		if paramTypeName, ok = k.NamedMap[paramName]; !ok {
			paramTypeName = rt.String()
		}
		k.SetParam(RequestParam{
			FieldName:     requestName,
			ParamDoc:      doc,
			ParamPath:     strings.Join(prefix, "."),
			ParamName:     paramName,
			ParamSource:   paramSource,
			ParamType:     "basic",
			ParamTypeName: paramTypeName,
			BasicType:     rt.Name(),
			HasPtr:        false,
		})
	default:
		return
	}

	return
}

func downFirst(s string) string {
	for _, v := range s {
		return string(unicode.ToLower(v)) + s[len(string(v)):]
	}
	return ""
}

func upFirst(s string) string {
	for _, v := range s {
		return string(unicode.ToUpper(v)) + s[len(string(v)):]
	}
	return ""
}

func CastMap(paramName, paramTypeName string, code jen.Code) (res []jen.Code, err error) {
	if paramTypeName == "[]string" {
		res = append(res, jen.Id("_"+paramName).Op("=").Qual("strings", "Split").Call(code, jen.Lit(",")))
		return
	}
	var m = map[string]string{
		"int":     "ToIntE",
		"int8":    "ToInt8E",
		"int16":   "ToInt16E",
		"int32":   "ToInt32E",
		"int64":   "ToInt64E",
		"uint":    "ToUintE",
		"uint8":   "ToUint8E",
		"uint16":  "ToUint16E",
		"uint32":  "ToUint32E",
		"uint64":  "ToUint64E",
		"float32": "ToFloat32E",
		"float64": "ToFloat64E",
		"string":  "ToStringE",
		"bool":    "ToBoolE",

		"[]int":  "ToIntSliceE",
		"[]bool": "ToBoolSliceE",

		"map.int":   "ToStringMapIntE",
		"map.int64": "ToStringMapInt64E",
		"map.bool":  "ToStringMapBoolE",

		"time.Time":     "ToTimeE",
		"time.Duration": "ToDurationE",
	}
	var ok bool
	fnStr, ok := m[paramTypeName]
	if !ok {
		err = fmt.Errorf("CastMap not found %s", paramTypeName)
		return
	}

	paramStr := paramName + "Str"
	varParamStr := jen.Id(paramStr).Op(":=").Add(code)
	paramStrCode := jen.Id(paramStr)
	if strings.HasPrefix(paramTypeName, "[]") {
		paramStrCode = jen.Qual("strings", "Split").Call(paramStrCode, jen.Lit(","))
	}
	ifParamStr := jen.If(jen.Id(paramStr).Op("!=").Lit("")).Block(
		jen.List(jen.Id("_"+paramName), jen.Err()).Op("=").Qual("github.com/spf13/cast", fnStr).Call(paramStrCode),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(),
		),
	)

	res = append(res, varParamStr, ifParamStr)

	return
}
