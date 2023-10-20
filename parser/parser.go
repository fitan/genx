package parser

import "github.com/antlr4-go/antlr/v4"

type DocLine struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

func ParserStr(s string) []DocLine {
	input := antlr.NewInputStream(s)
	lexer := NewTLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewTParser(stream)
	//p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.Doc()
	doc := NewTreeShapeListener()
	antlr.ParseTreeWalkerDefault.Walk(doc, tree)
	return doc.Docs
}

type TreeShapeListener struct {
	*BaseTParserListener
	Docs []DocLine `json:"docs"`
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (t *TreeShapeListener) EnterFunc(ctx FuncContext) {
	t.Docs = append(t.Docs, DocLine{
		Name: ctx.FuncName,
		Args: ctx.FuncArgs,
	})
}
