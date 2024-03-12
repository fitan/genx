package gormq

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/jennifer/jen"
	"github.com/samber/lo"
	"gorm.io/gorm/schema"
)

func Gen(j *jen.File, option gen.Option, implGoTypeMetes []gen.StructGoTypeMeta) {
	j.ImportAlias("gorm.io/gorm", "gorm")
	for _, v := range implGoTypeMetes {

		p := NewParser(option, v)
		p.Parse([]string{}, true)

		j.Add(p.GormQuery.GenScope())
	}
}

func NewParser(option gen.Option, s gen.StructGoTypeMeta) *Parser {
	p := &Parser{
		Option: option,
		Meta:   s,
	}
	return p
}

type Parser struct {
	Option    gen.Option
	GormQuery *GormQuery
	Meta      gen.StructGoTypeMeta
}

func (p *Parser) Parse(prePath []string, parseModel bool) {
	var modelName string
	if parseModel {
		p.Meta.Doc.ByFuncNameAndArgs(FuncName, &modelName)
		if !p.Meta.Doc.ByFuncNameAndArgs(FuncName, &modelName) {
			panic("not found @gq model")
		}
	}

	q := &GormQuery{
		StructName:           p.Meta.Name,
		ModelName:            modelName,
		FieldQueryList:       []FieldQuery{},
		StructQueryList:      []StructQuery{},
		SubQueryList:         []SubQuery{},
		GroupQueryList:       []GroupQuery{},
		PointStructQueryList: []PointStructQuery{},
	}

	p.GormQuery = q

	for i := 0; i < p.Meta.Obj.NumFields(); i++ {
		f := p.Meta.Obj.Field(i)

		if !f.IsField() {
			continue
		}

		tag := p.Meta.Obj.Tag(i)

		p.parse(prePath, f, tag)
	}

}

func (p *Parser) embed(pre []string, xt *common.Type, parentPoint bool) bool {
	if xt.Pointer {
		return p.embed(pre, common.TypeOf(xt.PointerInner.T), true)
	}

	if xt.Struct && !parentPoint {
		for i := 0; i < xt.StructType.NumFields(); i++ {
			f := xt.StructType.Field(i)
			if !f.IsField() {
				continue
			}

			tag := xt.StructType.Tag(i)
			p.parse(pre, f, tag)
		}
		return true
	} else if xt.Struct && parentPoint {
		meta := gen.StructGoTypeMeta{
			Name:   p.Meta.Name,
			Doc:    p.Meta.Doc,
			Params: p.Meta.Params,
			Obj:    xt.StructType,
		}
		np := NewParser(p.Option, meta)
		np.Parse(pre, false)

		p.GormQuery.PointStructQueryList = append(p.GormQuery.PointStructQueryList, PointStructQuery{
			Path:   pre,
			Parser: np,
		})
		return true
	}

	if xt.Named {
		return p.embed(pre, common.TypeOf(xt.NamedType.Underlying()), parentPoint)
	}

	return false
}

func (p *Parser) parse(pre []string, f *types.Var, tag string) {
	xt := common.TypeOf(f.Type())

	var clause, op, foreignKey, references string
	var dbNameList []string

	parseDoc, err := common.ParseDoc(common.GetCommentByTokenPos(p.Option.Pkg, f.Pos()).Text())
	if err != nil {
		panic(err)
	}

	parseDoc.ListByFuncName("@gq-column", &dbNameList)
	parseDoc.ByFuncNameAndArgs("@gq-clause", &clause)
	parseDoc.ByFuncNameAndArgs("@gq-op", &op)

	if len(dbNameList) == 0 {
		gormTag, ok := reflect.StructTag(tag).Lookup("gorm")
		if ok {
			for _, v := range strings.Split(gormTag, ";") {
				vv := strings.Split(v, ":")
				if len(vv) == 2 {
					if vv[0] == "column" {
						dbNameList = append(dbNameList, vv[1])
					}
				}
			}
		}
		name := schema.NamingStrategy{}
		dbNameList = append(dbNameList, name.ColumnName("", f.Name()))
	}

	if clause == "" {
		clause = "where"
	}
	clause = common.UpFirst(clause)

	if op == "" {
		op = "="
	}

	if parseDoc.ByFuncNameAndArgs("@gq-sub", &foreignKey, &references) {
		sq := SubQuery{
			Name:       f.Name(),
			XType:      xt,
			Path:       append(pre, f.Name()),
			Clause:     clause,
			ForeignKey: foreignKey,
			References: references,
		}
		p.GormQuery.SubQueryList = append(p.GormQuery.SubQueryList, sq)
		return
	}

	var gqStructValue []string
	if parseDoc.ListByFuncName("@gq-struct", &gqStructValue) {

		sq := StructQuery{
			Name:   f.Name(),
			XType:  xt,
			Clause: clause,
			Path:   append(pre, f.Name()),
			Values: gqStructValue,
		}

		p.GormQuery.StructQueryList = append(p.GormQuery.StructQueryList, sq)
		return
	}

	if parseDoc.ByFuncNameAndArgs("@gq-group") {
		gt := common.TypeOf(f.Type())
		if !gt.Struct {
			panic(fmt.Sprintf("not support type %s", f.String()))
		}

		gtDoc, _ := common.ParseDoc(common.GetCommentByTokenPos(p.Option.Pkg, f.Pos()).Text())
		gtParser := NewParser(p.Option, gen.StructGoTypeMeta{
			Name:   f.Name(),
			Doc:    gtDoc,
			Params: []string{},
			Obj:    gt.StructType,
		})

		gtParser.Parse(append(pre, f.Name()), false)

		p.GormQuery.GroupQueryList = append(p.GormQuery.GroupQueryList, GroupQuery{
			XType:     xt,
			Path:      append(pre, f.Name()),
			Clause:    clause,
			GormQuery: gtParser.GormQuery,
		})
		return
	}

	if p.embed(append(pre, f.Name()), xt, false) {
		return
	}

	qf := FieldQuery{
		Name:       f.Name(),
		XType:      xt,
		Op:         Op(op),
		DbNameList: dbNameList,
		Clause:     clause,
		Path:       append(pre, f.Name()),
	}

	qf.Parse(p.Option)
	p.GormQuery.FieldQueryList = append(p.GormQuery.FieldQueryList, qf)
}

type GormQuery struct {
	StructName string
	// @gq model.User
	ModelName string
	// @gq-op
	FieldQueryList []FieldQuery

	StructQueryList []StructQuery

	SubQueryList []SubQuery

	GroupQueryList []GroupQuery

	PointStructQueryList []PointStructQuery
}

func (g *GormQuery) GenScope() jen.Code {
	return jen.Func().Params(jen.Id("q").Op("*").Id(g.StructName)).Id("Scope").Params(jen.Id("db").Id("*gorm.DB")).Id("*gorm.DB").BlockFunc(
		func(group *jen.Group) {
			group.Id("db").Op("=").Id("db.Model").Call(jen.Id("&" + g.ModelName + "{}"))

			g.GenWhere(group)

			//for _, v := range g.FieldQueryList {
			//	group.Id("db").Op("=").Id("db.").Add(v.GenClause())
			//}
			//
			//for _, v := range g.StructQueryList {
			//	group.Id("db").Op("=").Id("db.").Add(v.GenClause())
			//}
			//
			//for _, v := range g.SubQueryList {
			//	group.Id("db").Op("=").Id("db.").Add(v.GenClause())
			//}
			//
			//for _, v := range g.GroupQueryList {
			//	group.Id("db").Op("=").Id("db." + v.Clause).Call(v.GormQuery.GenWhereScope())
			//}

			group.Return(jen.Id("db"))
		},
	)
}

func (g *GormQuery) GenWhere(code *jen.Group) {
	for _, v := range g.FieldQueryList {
		before, ql := v.GenClause()
		genCode := jen.Line().Add(before).Line()
		genCode.Id("db").Op("=").Id("db.").Add(ql)
		if v.XType.Pointer {
			code.If(jen.Id("q." + strings.Join(v.Path, ".")).Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
				group.Add(genCode)
			})
		} else {
			code.Add(genCode)
		}
	}

	for _, v := range g.StructQueryList {
		genCode := jen.Id("db").Op("=").Id("db.").Add(v.GenClause())
		if v.XType.Pointer {
			code.If(jen.Id("q." + strings.Join(v.Path, ".")).Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
				group.Add(genCode)
			})
		} else {
			code.Add(genCode)
		}
	}

	for _, v := range g.SubQueryList {
		genCode := jen.Id("db").Op("=").Id("db.").Add(v.GenClause())
		if v.XType.Pointer {
			code.If(jen.Id("q." + strings.Join(v.Path, ".")).Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
				group.Add(genCode)
			})
		} else {
			code.Add(genCode)
		}
	}

	for _, v := range g.GroupQueryList {
		genCode := jen.Id("db").Op("=").Id("db." + v.Clause).Call(v.GormQuery.GenWhereScope().Call(jen.Id("db.Session(&gorm.Session{NewDB: true})")))
		if v.XType.Pointer {
			code.If(jen.Id("q." + strings.Join(v.Path, ".")).Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
				group.Add(genCode)
			})
		} else {
			code.Add(genCode)
		}
	}

	for _, v := range g.PointStructQueryList {
		code.If(jen.Id("q." + strings.Join(v.Path, ".")).Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
			v.Parser.GormQuery.GenWhere(group)
		})
	}
}

func (g *GormQuery) GenWhereScope() *jen.Statement {
	return jen.Func().Call(jen.Id("db *gorm.DB")).Params(jen.Id("*gorm.DB")).BlockFunc(
		func(group *jen.Group) {
			g.GenWhere(group)

			group.Return(jen.Id("db"))
		})
}

type FieldQuery struct {
	Name  string
	XType *common.Type
	Op    Op
	// @gq-column name name1
	DbNameList []string
	// @gq-clause where or not
	Clause string

	Path []string

	InStruct InStruct
}

type InStruct struct {
	Is     bool
	Parser *Parser
}

func (q *FieldQuery) Parse(op gen.Option) {
	if q.XType.List {
		if q.XType.ListInner.Struct {
			q.InStruct.Is = true
			q.InStruct.Parser = NewParser(op, gen.StructGoTypeMeta{
				Name:   q.Name,
				Params: []string{},
				Obj:    q.XType.ListInner.StructType,
			})
			q.InStruct.Parser.Parse([]string{}, false)
		} else {
			q.Op = "in"
		}
	}
}

func (q FieldQuery) GenQuery() string {
	//if strings.Contains(q.Op, "better") || strings.Contains(q.Op, "><") {
	//	return q.DbName + " BETWEEN ? AND ?"
	//}
	if q.InStruct.Is {
		structDbNameList := []string{}
		for _, v := range q.InStruct.Parser.GormQuery.FieldQueryList {
			dbName, _ := lo.Last(v.DbNameList)
			structDbNameList = append(structDbNameList, dbName)
		}
		return "(" + strings.Join(structDbNameList, ", ") + ")" + " IN ?"
	}

	if len(q.DbNameList) == 1 {
		return q.DbNameList[0] + " " + q.Op.Convert() + " ?"
	} else {
		return strings.Join(lo.Map(q.DbNameList, func(item string, index int) string {
			return item + " " + q.Op.Convert() + " @key"
		}), " OR ")
	}
}

// Where("...", v)
func (q FieldQuery) GenClause() (before, ql jen.Code) {
	if q.InStruct.Is {
		id := q.Name + "Value"

		b := jen.Id(id).Op(":=").Id("make([][]interface{},0,0)").Line()
		b.Id("for _, v := range").Id("q." + strings.Join(q.Path, ".")).BlockFunc(func(group *jen.Group) {
			var valuesCode []jen.Code
			for _, v := range q.InStruct.Parser.GormQuery.FieldQueryList {
				valuesCode = append(valuesCode, jen.Id("v."+strings.Join(v.Path, ".")))
			}
			group.Id(id).Op("=").Append(jen.Id(id), jen.Index().Interface().Values(valuesCode...))
		})
		//for _, v := range q.InStruct.Parser.GormQuery.FieldQueryList {
		//	b.Add(jen.Id(q.Name).Op("=").Append(jen.Id(q.Name), jen.Id("q."+strings.Join(v.Path, "."))).Line())
		//}

		code := []jen.Code{jen.Lit(q.GenQuery())}
		code = append(code, jen.Id(id))
		return b, jen.Id(q.Clause).Call(code...)

	}
	header := "q."
	if q.XType.Pointer {
		header = "*" + header
	}

	if len(q.DbNameList) == 1 {
		code := []jen.Code{jen.Lit(q.GenQuery())}
		code = append(code, q.Op.ConvertValue(header+strings.Join(q.Path, "."))...)
		return nil, jen.Id(q.Clause).Call(code...)
	} else {

		first := q.DbNameList[0] + " " + q.Op.Convert() + " ?"
		value := q.Op.ConvertValue(header + strings.Join(q.Path, "."))
		code := []jen.Code{jen.Lit(first)}
		code = append(code, value...)
		statement := jen.Id(q.Clause).Call(code...)
		for _, v := range q.DbNameList[1:] {
			c := []jen.Code{jen.Lit(v + " " + q.Op.Convert() + " ?")}
			c = append(c, value...)
			statement.Dot("Or").Call(c...)
		}
		return nil, statement
	}
}

type SubQuery struct {
	// @gq-clause where or not
	Name       string
	XType      *common.Type
	Path       []string
	Clause     string
	ForeignKey string
	References string
}

// Where("id in (?)", scope(db.Select(id)))
func (s SubQuery) GenClause() jen.Code {
	return jen.Id(s.Clause).Call(jen.Lit(s.ForeignKey+" in (?)"), jen.Id("q."+strings.Join(s.Path, ".")+".Scope").Call(jen.Id(`db.Session(&gorm.Session{NewDB: true}).Select("`+s.References+`")`)))
}

type StructQuery struct {
	Name   string
	XType  *common.Type
	Clause string
	Path   []string
	Values []string
}

func (s StructQuery) GenClause() jen.Code {
	code := []jen.Code{jen.Id("&q." + strings.Join(s.Path, "."))}
	lo.ForEach(s.Values, func(item string, index int) {
		code = append(code, jen.Lit(item))
	})
	return jen.Id(s.Clause).Call(code...)
}

type EmbedQuery struct {
	SubQuery
}

type GroupQuery struct {
	XType     *common.Type
	Path      []string
	Clause    string
	GormQuery *GormQuery
}

type PointStructQuery struct {
	Path   []string
	Parser *Parser
}
