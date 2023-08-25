package common

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"golang.org/x/tools/go/packages"
)

var (
	//enumLexer = lexer.MustSimple([]lexer.SimpleRule{
	//	{"whitespace", `\s+`},
	//	{"Punct", `[,)(]`},
	//	{"FuncName", `^@[a-zA-Z][a-zA-Z_\d]*`},
	//	//{"FuncTag", `^@[a-zA-Z][a-zA-Z_\d]*\(`},
	//	{"String", `"(\\.|[^"])*"|'(\\.|[^'])*'`},
	//	{"Ident", `[^ \f\n\r\t\v,)]+`},
	//})
	enumLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"whitespace", `\s+`},
		{`String`, `"(?:\\.|[^"])*"|'(?:\\.|[^'])*'`},
		{"Punct", `[)(,]`},
		{"Name", `^@[a-zA-Z][a-zA-Z_\d]*`},
		{"Comment", `^[^@].+`},
	})
	parser = participle.MustBuild[Doc](
		participle.Lexer(enumLexer),
	)
)

type Doc struct {
	Funcs []*Func `@@*`
}

type Func struct {
	Others *string `@Comment | @Name (?! "(")`
	Func   *F      `| @@`
}

type F struct {
	FuncName string   `@Name`
	Args     []string `( "(" ( @String ("," @String)* )? ")"  )`
}

func (f *F) UpFunkName() string {
	return strings.ToUpper(f.FuncName)
}

func (d *Doc) ByFuncName(name string) *F {
	for _, v := range d.Funcs {
		if v.Func != nil && strings.EqualFold(v.Func.FuncName, name) {
			return v.Func
		}
	}
	return nil
}

func (d *Doc) ByFuncNameAndArgs(name string, args ...*string) bool {
	f := d.ByFuncName(name)
	if f == nil {
		return false
	}
	record := make([]string, len(f.Args), len(f.Args))
	for i, arg := range f.Args {
		record[i] = arg
	}
	for i, _ := range args {
		args[i] = &record[i]
	}

	return true
}

func ParseDoc(s string) (*Doc, error) {
	//slog.Info("parse doc: ", slog.String("doc", s))
	return parser.ParseString("", s)
}

// 获取struct字段里的注释
func GetCommentByTokenPos(pkg *packages.Package, pos token.Pos) *ast.CommentGroup {
	fieldFileName := pkg.Fset.Position(pos).Filename
	fieldLine := pkg.Fset.Position(pos).Line
	var fieldComment *ast.CommentGroup
	for _, syntax := range pkg.Syntax {
		fileName := pkg.Fset.Position(syntax.Pos()).Filename
		if fieldFileName == fileName {
			for _, c := range syntax.Comments {
				if pkg.Fset.Position(c.End()).Line+1 == fieldLine {
					fieldComment = c
				}
			}
			break
		}
	}
	return fieldComment
}
