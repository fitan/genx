package kithttp

import (
	"github.com/fitan/genx/common"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const (
	KitHttp          = "@kit-http"
	KitService       = "@kit-http-service"
	kitParam         = "@kit-http-param"
	KitHttpRequest   = "@kit-http-request"
	KitHttpResponse  = "@kit-http-response"
	kitEndpointParam = "@kit-endpoint-param"
)

type KitCommentConf struct {
	Url       string
	UrlMethod string

	HttpRequestName string
	HttpRequestBody bool

	//HttpResponseName string

	HttpParams      map[string]HttpParam
	KitServiceParam KitServiceParam
}

type HttpParam struct {
	MethodParamName string
	SourceType      string
	Validate        string
	Annotation      string
}

type KitServiceParam struct {
	EndpointName string
	DecodeName   string
	EncodeName   string
}

func KitComment(doc common.Doc) (kitConf KitCommentConf, err error) {
	has := doc.ByFuncNameAndArgs(KitHttp, &kitConf.Url, &kitConf.UrlMethod)
	if has {
		if kitConf.Url == "" || kitConf.UrlMethod == "" {
			err = errors.New("must format: @kit-http url method")
			return
		}
	}

	doc.ByFuncNameAndArgs(KitService, &kitConf.KitServiceParam.EndpointName, &kitConf.KitServiceParam.DecodeName, &kitConf.KitServiceParam.EncodeName)

	var isRequestBody string

	doc.ByFuncNameAndArgs(KitHttpRequest, &kitConf.HttpRequestName, &isRequestBody)
	kitConf.HttpRequestBody = lo.Ternary(isRequestBody != "" && isRequestBody != "false", true, false)

	return
}

//func (m *KitCommentConf) ParamKitHttpResponse(s []string) (err error) {
//	if len(s) < 3 {
//		err = errors.New("must format: @kit-http-response responseName")
//		return
//	}
//	m.HttpResponseName = s[2]
//	return
//}

type Kit struct {
	//Comment KitConf
	Conf KitCommentConf
}

type InterfaceMethodParam struct {
	ParamName string
	ParamType string
}

func NewKit(doc common.Doc) (res Kit, err error) {
	commentConf, err := KitComment(doc)
	if err != nil {
		return
	}

	res.Conf = commentConf

	return

}
