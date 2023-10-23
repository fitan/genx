package common

import (
	"go/ast"
	"go/types"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
	"strings"
)

type InterfaceSerialize struct {
	pkg *packages.Package
}

func NewInterfaceSerialize(pkg *packages.Package) *InterfaceSerialize {
	return &InterfaceSerialize{pkg: pkg}
}

func (i *InterfaceSerialize) Parse(it *types.Interface) (InterfaceMetaDate, error) {
	impl := InterfaceMetaDate{
		Methods: make([]InterfaceMethod, 0),
	}

	for n := 0; n < it.NumMethods(); n++ {
		var returnsError bool
		var acceptsContext bool
		var params MethodParamSlice
		var results MethodParamSlice
		var comment *ast.CommentGroup
		m := it.Method(n)
		comment = GetCommentByTokenPos(i.pkg, m.Pos())
		doc, err := ParseDoc(comment.Text())
		if err != nil {
			slog.Error("ParseDoc", err, slog.String("comment", comment.Text()))
			return impl, err
		}
		methodName := m.Name()
		signature := m.Type().(*types.Signature)
		ps := signature.Params()
		for pi := 0; pi < ps.Len(); pi++ {
			mParam := MethodParam{}
			if pi == 0 {
				if ps.At(pi).Type().String() == "context.Context" {
					acceptsContext = true
				}
			}
			p := ps.At(pi)
			t := p.Type()
			pName := p.Name()

			mParam.Name = pName
			mParam.Type = t
			mParam.XType = TypeOf(t)
			mParam.ID = mParam.XType.TypeAsJenComparePkgNameString(i.pkg)
			params = append(params, mParam)
		}

		rs := m.Type().(*types.Signature).Results()
		for ri := 0; ri < rs.Len(); ri++ {
			mParam := MethodParam{}
			if ri == rs.Len()-1 {
				if rs.At(ri).Type().String() == "error" {
					returnsError = true
				}
			}

			r := rs.At(ri)
			rName := r.Name()
			mParam.Name = rName

			t := r.Type()
			mParam.Type = t
			mParam.XType = TypeOf(t)
			mParam.ID = mParam.XType.TypeAsJenComparePkgNameString(i.pkg)

			results = append(results, mParam)
		}

		implMethod := InterfaceMethod{
			Name:           methodName,
			Doc:            doc,
			Params:         params,
			Results:        results,
			ReturnsError:   returnsError,
			AcceptsContext: acceptsContext,
			Variadic:       signature.Variadic(),
		}
		impl.Methods = append(impl.Methods, implMethod)
	}

	return impl, nil
}

type InterfaceMetaDate struct {
	Imports []*ast.ImportSpec
	Doc     *Doc
	Methods []InterfaceMethod
}

type InterfaceMethod struct {
	Name           string
	Doc            Doc
	Params         MethodParamSlice
	Results        MethodParamSlice
	ReturnsError   bool
	AcceptsContext bool
	Variadic       bool
}

func (m InterfaceMethod) HasParams() bool {
	return len(m.Params) > 0
}

// HasResults returns true if method has results
func (m InterfaceMethod) HasResults() bool {
	return len(m.Results) > 0
}

func (m InterfaceMethod) ResultsExcludeErr() MethodParamSlice {
	tmp := make(MethodParamSlice, 0, 0)
	for _, p := range m.Results {
		if p.ID == "error" {
			continue
		}

		tmp = append(tmp, p)
	}
	return tmp
}

func (m InterfaceMethod) ParamsExcludeCtx() MethodParamSlice {
	tmp := make(MethodParamSlice, 0, 0)
	for _, p := range m.Params {
		if p.Type.String() == "context.Context" {
			continue
		}

		tmp = append(tmp, p)
	}
	return tmp
}

func (m InterfaceMethod) SwagRespObjData() string {
	if len(m.ResultsExcludeErr()) == 0 {
		return "data=string"
	}

	if len(m.ResultsExcludeErr()) == 1 {
		return "data=" + m.ResultsExcludeErr()[0].ID
	}

	var s []string
	for _, v := range m.ResultsExcludeErr() {
		s = append(s, "data."+v.Name+"="+v.ID)
	}
	return strings.Join(s, ",")
}

func (m InterfaceMethod) ResultsMapExcludeErr() string {
	ss := []string{}
	for _, r := range m.ResultsExcludeErr() {
		ss = append(ss, `"`+r.Name+`": `+r.Name)
	}
	return "map[string]interface{}{\n" + strings.Join(ss, ",\n ") + "}"
}

type MethodParamSlice []MethodParam

type MethodParam struct {
	Comment []*ast.Comment
	Name    string
	Type    types.Type
	ID      string
	XType   *Type
}

func (m MethodParam) Basic() bool {
	return m.XType.Basic
}
