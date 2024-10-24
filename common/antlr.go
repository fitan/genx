package common

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/cockroachdb/errors"
	"github.com/fitan/genx/parser"
	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

type Doc []DocLine

func (d *Doc) ByFuncName(name string) *DocLine {
	for _, doc := range *d {
		if doc.Name == name {
			return &doc
		}
	}
	return nil
}

func (d *Doc) ByFuncNameAndArgName(funcName, argName string) (string, bool) {
	if d == nil {
		return "", false
	}
	f := d.ByFuncName(funcName)
	if f == nil {
		return "", false
	}
	for _, arg := range f.Args {
		if arg.Name == argName {
			return arg.Value, true
		}
	}
	return "", false
}

func (d *Doc) ListByFuncName(name string, list *[]string) bool {
	if d == nil {
		return false
	}
	f := d.ByFuncName(name)
	if f == nil {
		return false
	}
	for _, arg := range f.Args {
		*list = append(*list, arg.Value)
	}
	return true
}

func (d *Doc) ByFuncNameAndArgs(name string, args ...*string) bool {
	if d == nil {
		return false
	}
	f := d.ByFuncName(name)
	if f == nil {
		return false
	}
	max := lo.Max([]int{len(f.Args), len(args)})
	record := make([]string, max, max)
	for i, arg := range f.Args {
		record[i] = arg.Value
	}
	for i, _ := range args {
		*args[i] = record[i]
	}

	return true
}

type DocLine struct {
	Name string           `json:"name"`
	Args []parser.FuncArg `json:"args"`
}

func (d *DocLine) UpFuncName() string {
	return strings.ToUpper(d.Name)
}

func ParseDoc(s string) (Doc, error) {
	input := antlr.NewInputStream(s)
	lexer := parser.NewTLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewTParser(stream)
	customErrorListener := NewCustomErrorListener()
	p.AddErrorListener(customErrorListener)
	tree := p.Doc()
	doc := NewTreeShapeListener()
	antlr.ParseTreeWalkerDefault.Walk(doc, tree)
	return doc.Doc, errors.Join(customErrorListener.Errors...)
}

type TreeShapeListener struct {
	*parser.BaseTParserListener
	Doc Doc `json:"docs"`
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (t *TreeShapeListener) EnterFunc(ctx *parser.FuncContext) {
	t.Doc = append(t.Doc, DocLine{
		Name: ctx.FuncName,
		Args: ctx.FuncArgs,
	})
}

type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []error
}

func NewCustomErrorListener() *CustomErrorListener {
	return &CustomErrorListener{
		Errors: make([]error, 0, 0),
	}
}
func (c *CustomErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	c.Errors = append(c.Errors, fmt.Errorf("语法错误：行 %d，列 %d，消息：%s\n", line, column, msg))
}

// 获取struct字段里的注释
func GetCommentByTokenPos(pkg *packages.Package, pos token.Pos) *ast.CommentGroup {
	fieldFileName := pkg.Fset.Position(pos).Filename
	fieldLine := pkg.Fset.Position(pos).Line
	fieldComment := ast.CommentGroup{}
	for _, syntax := range pkg.Syntax {
		fileName := pkg.Fset.Position(syntax.Pos()).Filename
		if fieldFileName == fileName {
			for _, c := range syntax.Comments {
				if pkg.Fset.Position(c.End()).Line+1 == fieldLine {
					fieldComment = *c
				}
			}
			break
		}
	}
	return &fieldComment
}
