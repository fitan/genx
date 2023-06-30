package common

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"go/ast"
	"go/token"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
	"strings"
)

var (
	enumLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"whitespace", `\s+`},
		{"Punct", `[,)(]`},
		{"FuncName", `^@[a-zA-Z][a-zA-Z_\d]*`},
		//{"FuncTag", `^@[a-zA-Z][a-zA-Z_\d]*\(`},
		{"String", `"(\\.|[^"])*"|'(\\.|[^'])*'`},
		{"Ident", `[^ \f\n\r\t\v,)]+`},
	})
	parser = participle.MustBuild[Doc](
		participle.Lexer(enumLexer),
	)
)

type Doc struct {
	Funcs []*Func `@@*`
}

type Func struct {
	Others *string `@~FuncName`
	Func   *F      `| @@`
}

type F struct {
	FuncName string   `@FuncName`
	Args     []string `( "(" (@String | @Ident) ("," (@String | @Ident))* ")" )?`
}

func (d *Doc) ByFuncName(f F) *Func {
	for _, v := range d.Funcs {
		if v.Func != nil && strings.EqualFold(v.Func.FuncName, f.FuncName) {
			return v
		}
	}
	return nil
}

func ParseDoc(s string) (*Doc, error) {
	slog.Info("parse doc: ", slog.String("doc", s))
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
