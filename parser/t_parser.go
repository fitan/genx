// Code generated from TParser.g4 by ANTLR 4.13.0. DO NOT EDIT.

package parser // TParser

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

import "strings"

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type TParser struct {
	*antlr.BaseParser
}

var TParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func tparserParserInit() {
	staticData := &TParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "'='", "','",
	}
	staticData.SymbolicNames = []string{
		"", "ATID", "FieldFuncName", "LPAREN", "NEWLINE", "WS", "INSET", "EQ",
		"Comma", "PARENWS", "ID", "String", "RPAREN", "S", "CLOSE", "OLDFUNCCLOSE",
		"OLDFUNCWS", "FIELD",
	}
	staticData.RuleNames = []string{
		"doc", "line", "func", "argument",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 17, 78, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 1, 0, 5,
		0, 10, 8, 0, 10, 0, 12, 0, 13, 9, 0, 1, 1, 1, 1, 1, 1, 3, 1, 18, 8, 1,
		1, 1, 1, 1, 3, 1, 22, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 3, 2, 32, 8, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 40, 8, 2,
		5, 2, 42, 8, 2, 10, 2, 12, 2, 45, 9, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 61, 8, 2, 10,
		2, 12, 2, 64, 9, 2, 3, 2, 66, 8, 2, 1, 2, 1, 2, 1, 2, 3, 2, 71, 8, 2, 1,
		3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 0, 0, 4, 0, 2, 4, 6, 0, 0, 85, 0, 11,
		1, 0, 0, 0, 2, 21, 1, 0, 0, 0, 4, 70, 1, 0, 0, 0, 6, 72, 1, 0, 0, 0, 8,
		10, 3, 2, 1, 0, 9, 8, 1, 0, 0, 0, 10, 13, 1, 0, 0, 0, 11, 9, 1, 0, 0, 0,
		11, 12, 1, 0, 0, 0, 12, 1, 1, 0, 0, 0, 13, 11, 1, 0, 0, 0, 14, 22, 3, 4,
		2, 0, 15, 17, 5, 6, 0, 0, 16, 18, 5, 13, 0, 0, 17, 16, 1, 0, 0, 0, 17,
		18, 1, 0, 0, 0, 18, 19, 1, 0, 0, 0, 19, 22, 5, 14, 0, 0, 20, 22, 5, 4,
		0, 0, 21, 14, 1, 0, 0, 0, 21, 15, 1, 0, 0, 0, 21, 20, 1, 0, 0, 0, 22, 3,
		1, 0, 0, 0, 23, 24, 5, 1, 0, 0, 24, 25, 6, 2, -1, 0, 25, 31, 5, 3, 0, 0,
		26, 27, 5, 11, 0, 0, 27, 32, 6, 2, -1, 0, 28, 29, 3, 6, 3, 0, 29, 30, 6,
		2, -1, 0, 30, 32, 1, 0, 0, 0, 31, 26, 1, 0, 0, 0, 31, 28, 1, 0, 0, 0, 32,
		43, 1, 0, 0, 0, 33, 39, 5, 8, 0, 0, 34, 35, 5, 11, 0, 0, 35, 40, 6, 2,
		-1, 0, 36, 37, 3, 6, 3, 0, 37, 38, 6, 2, -1, 0, 38, 40, 1, 0, 0, 0, 39,
		34, 1, 0, 0, 0, 39, 36, 1, 0, 0, 0, 40, 42, 1, 0, 0, 0, 41, 33, 1, 0, 0,
		0, 42, 45, 1, 0, 0, 0, 43, 41, 1, 0, 0, 0, 43, 44, 1, 0, 0, 0, 44, 46,
		1, 0, 0, 0, 45, 43, 1, 0, 0, 0, 46, 47, 5, 12, 0, 0, 47, 48, 1, 0, 0, 0,
		48, 71, 5, 4, 0, 0, 49, 50, 5, 1, 0, 0, 50, 51, 6, 2, -1, 0, 51, 52, 5,
		3, 0, 0, 52, 53, 5, 12, 0, 0, 53, 71, 5, 4, 0, 0, 54, 55, 5, 2, 0, 0, 55,
		65, 6, 2, -1, 0, 56, 57, 5, 17, 0, 0, 57, 62, 6, 2, -1, 0, 58, 59, 5, 17,
		0, 0, 59, 61, 6, 2, -1, 0, 60, 58, 1, 0, 0, 0, 61, 64, 1, 0, 0, 0, 62,
		60, 1, 0, 0, 0, 62, 63, 1, 0, 0, 0, 63, 66, 1, 0, 0, 0, 64, 62, 1, 0, 0,
		0, 65, 56, 1, 0, 0, 0, 65, 66, 1, 0, 0, 0, 66, 67, 1, 0, 0, 0, 67, 71,
		5, 15, 0, 0, 68, 69, 5, 1, 0, 0, 69, 71, 6, 2, -1, 0, 70, 23, 1, 0, 0,
		0, 70, 49, 1, 0, 0, 0, 70, 54, 1, 0, 0, 0, 70, 68, 1, 0, 0, 0, 71, 5, 1,
		0, 0, 0, 72, 73, 5, 10, 0, 0, 73, 74, 5, 7, 0, 0, 74, 75, 5, 11, 0, 0,
		75, 76, 6, 3, -1, 0, 76, 7, 1, 0, 0, 0, 9, 11, 17, 21, 31, 39, 43, 62,
		65, 70,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// TParserInit initializes any static state used to implement TParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewTParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func TParserInit() {
	staticData := &TParserParserStaticData
	staticData.once.Do(tparserParserInit)
}

// NewTParser produces a new parser instance for the optional input antlr.TokenStream.
func NewTParser(input antlr.TokenStream) *TParser {
	TParserInit()
	this := new(TParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &TParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "TParser.g4"

	return this
}

// Note that '@members' cannot be changed now, but this should have been 'globals'
// If you are looking to have variables for each instance, use '@structmembers'

func trimQuotation(s string) string {
	if strings.HasPrefix(s, "\"") {
		return strings.Trim(s, "\"")
	}
	if strings.HasPrefix(s, "'") {
		return strings.Trim(s, "'")
	}
	return s
}

func GenFuncArg(name, value string) (res FuncArg) {
	res.Name = name
	res.Value = value
	return
}

type FuncArg struct {
	Name  string
	Value string
}

// TParser tokens.
const (
	TParserEOF           = antlr.TokenEOF
	TParserATID          = 1
	TParserFieldFuncName = 2
	TParserLPAREN        = 3
	TParserNEWLINE       = 4
	TParserWS            = 5
	TParserINSET         = 6
	TParserEQ            = 7
	TParserComma         = 8
	TParserPARENWS       = 9
	TParserID            = 10
	TParserString_       = 11
	TParserRPAREN        = 12
	TParserS             = 13
	TParserCLOSE         = 14
	TParserOLDFUNCCLOSE  = 15
	TParserOLDFUNCWS     = 16
	TParserFIELD         = 17
)

// TParser rules.
const (
	TParserRULE_doc      = 0
	TParserRULE_line     = 1
	TParserRULE_func     = 2
	TParserRULE_argument = 3
)

// IDocContext is an interface to support dynamic dispatch.
type IDocContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLine() []ILineContext
	Line(i int) ILineContext

	// IsDocContext differentiates from other interfaces.
	IsDocContext()
}

type DocContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDocContext() *DocContext {
	var p = new(DocContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_doc
	return p
}

func InitEmptyDocContext(p *DocContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_doc
}

func (*DocContext) IsDocContext() {}

func NewDocContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DocContext {
	var p = new(DocContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_doc

	return p
}

func (s *DocContext) GetParser() antlr.Parser { return s.parser }

func (s *DocContext) AllLine() []ILineContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILineContext); ok {
			len++
		}
	}

	tst := make([]ILineContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILineContext); ok {
			tst[i] = t.(ILineContext)
			i++
		}
	}

	return tst
}

func (s *DocContext) Line(i int) ILineContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILineContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILineContext)
}

func (s *DocContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DocContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DocContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterDoc(s)
	}
}

func (s *DocContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitDoc(s)
	}
}

func (p *TParser) Doc() (localctx IDocContext) {
	localctx = NewDocContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, TParserRULE_doc)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(11)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&86) != 0 {
		{
			p.SetState(8)
			p.Line()
		}

		p.SetState(13)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILineContext is an interface to support dynamic dispatch.
type ILineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Func_() IFuncContext
	INSET() antlr.TerminalNode
	CLOSE() antlr.TerminalNode
	S() antlr.TerminalNode
	NEWLINE() antlr.TerminalNode

	// IsLineContext differentiates from other interfaces.
	IsLineContext()
}

type LineContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLineContext() *LineContext {
	var p = new(LineContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_line
	return p
}

func InitEmptyLineContext(p *LineContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_line
}

func (*LineContext) IsLineContext() {}

func NewLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineContext {
	var p = new(LineContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_line

	return p
}

func (s *LineContext) GetParser() antlr.Parser { return s.parser }

func (s *LineContext) Func_() IFuncContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFuncContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFuncContext)
}

func (s *LineContext) INSET() antlr.TerminalNode {
	return s.GetToken(TParserINSET, 0)
}

func (s *LineContext) CLOSE() antlr.TerminalNode {
	return s.GetToken(TParserCLOSE, 0)
}

func (s *LineContext) S() antlr.TerminalNode {
	return s.GetToken(TParserS, 0)
}

func (s *LineContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(TParserNEWLINE, 0)
}

func (s *LineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterLine(s)
	}
}

func (s *LineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitLine(s)
	}
}

func (p *TParser) Line() (localctx ILineContext) {
	localctx = NewLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, TParserRULE_line)
	var _la int

	p.SetState(21)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TParserATID, TParserFieldFuncName:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(14)
			p.Func_()
		}

	case TParserINSET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(15)
			p.Match(TParserINSET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(17)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == TParserS {
			{
				p.SetState(16)
				p.Match(TParserS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(19)
			p.Match(TParserCLOSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case TParserNEWLINE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(20)
			p.Match(TParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFuncContext is an interface to support dynamic dispatch.
type IFuncContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_ATID returns the _ATID token.
	Get_ATID() antlr.Token

	// Get_String_ returns the _String_ token.
	Get_String_() antlr.Token

	// Get_FieldFuncName returns the _FieldFuncName token.
	Get_FieldFuncName() antlr.Token

	// Get_FIELD returns the _FIELD token.
	Get_FIELD() antlr.Token

	// Set_ATID sets the _ATID token.
	Set_ATID(antlr.Token)

	// Set_String_ sets the _String_ token.
	Set_String_(antlr.Token)

	// Set_FieldFuncName sets the _FieldFuncName token.
	Set_FieldFuncName(antlr.Token)

	// Set_FIELD sets the _FIELD token.
	Set_FIELD(antlr.Token)

	// Get_argument returns the _argument rule contexts.
	Get_argument() IArgumentContext

	// Set_argument sets the _argument rule contexts.
	Set_argument(IArgumentContext)

	// GetFuncArgs returns the FuncArgs attribute.
	GetFuncArgs() []FuncArg

	// GetFuncName returns the FuncName attribute.
	GetFuncName() string

	// SetFuncArgs sets the FuncArgs attribute.
	SetFuncArgs([]FuncArg)

	// SetFuncName sets the FuncName attribute.
	SetFuncName(string)

	// Getter signatures
	ATID() antlr.TerminalNode
	NEWLINE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllString_() []antlr.TerminalNode
	String_(i int) antlr.TerminalNode
	AllArgument() []IArgumentContext
	Argument(i int) IArgumentContext
	AllComma() []antlr.TerminalNode
	Comma(i int) antlr.TerminalNode
	FieldFuncName() antlr.TerminalNode
	OLDFUNCCLOSE() antlr.TerminalNode
	AllFIELD() []antlr.TerminalNode
	FIELD(i int) antlr.TerminalNode

	// IsFuncContext differentiates from other interfaces.
	IsFuncContext()
}

type FuncContext struct {
	antlr.BaseParserRuleContext
	parser         antlr.Parser
	FuncArgs       []FuncArg
	FuncName       string
	_ATID          antlr.Token
	_String_       antlr.Token
	_argument      IArgumentContext
	_FieldFuncName antlr.Token
	_FIELD         antlr.Token
}

func NewEmptyFuncContext() *FuncContext {
	var p = new(FuncContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_func
	return p
}

func InitEmptyFuncContext(p *FuncContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_func
}

func (*FuncContext) IsFuncContext() {}

func NewFuncContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncContext {
	var p = new(FuncContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_func

	return p
}

func (s *FuncContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncContext) Get_ATID() antlr.Token { return s._ATID }

func (s *FuncContext) Get_String_() antlr.Token { return s._String_ }

func (s *FuncContext) Get_FieldFuncName() antlr.Token { return s._FieldFuncName }

func (s *FuncContext) Get_FIELD() antlr.Token { return s._FIELD }

func (s *FuncContext) Set_ATID(v antlr.Token) { s._ATID = v }

func (s *FuncContext) Set_String_(v antlr.Token) { s._String_ = v }

func (s *FuncContext) Set_FieldFuncName(v antlr.Token) { s._FieldFuncName = v }

func (s *FuncContext) Set_FIELD(v antlr.Token) { s._FIELD = v }

func (s *FuncContext) Get_argument() IArgumentContext { return s._argument }

func (s *FuncContext) Set_argument(v IArgumentContext) { s._argument = v }

func (s *FuncContext) GetFuncArgs() []FuncArg { return s.FuncArgs }

func (s *FuncContext) GetFuncName() string { return s.FuncName }

func (s *FuncContext) SetFuncArgs(v []FuncArg) { s.FuncArgs = v }

func (s *FuncContext) SetFuncName(v string) { s.FuncName = v }

func (s *FuncContext) ATID() antlr.TerminalNode {
	return s.GetToken(TParserATID, 0)
}

func (s *FuncContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(TParserNEWLINE, 0)
}

func (s *FuncContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(TParserLPAREN, 0)
}

func (s *FuncContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(TParserRPAREN, 0)
}

func (s *FuncContext) AllString_() []antlr.TerminalNode {
	return s.GetTokens(TParserString_)
}

func (s *FuncContext) String_(i int) antlr.TerminalNode {
	return s.GetToken(TParserString_, i)
}

func (s *FuncContext) AllArgument() []IArgumentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IArgumentContext); ok {
			len++
		}
	}

	tst := make([]IArgumentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IArgumentContext); ok {
			tst[i] = t.(IArgumentContext)
			i++
		}
	}

	return tst
}

func (s *FuncContext) Argument(i int) IArgumentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgumentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgumentContext)
}

func (s *FuncContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(TParserComma)
}

func (s *FuncContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(TParserComma, i)
}

func (s *FuncContext) FieldFuncName() antlr.TerminalNode {
	return s.GetToken(TParserFieldFuncName, 0)
}

func (s *FuncContext) OLDFUNCCLOSE() antlr.TerminalNode {
	return s.GetToken(TParserOLDFUNCCLOSE, 0)
}

func (s *FuncContext) AllFIELD() []antlr.TerminalNode {
	return s.GetTokens(TParserFIELD)
}

func (s *FuncContext) FIELD(i int) antlr.TerminalNode {
	return s.GetToken(TParserFIELD, i)
}

func (s *FuncContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FuncContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterFunc(s)
	}
}

func (s *FuncContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitFunc(s)
	}
}

func (p *TParser) Func_() (localctx IFuncContext) {
	localctx = NewFuncContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, TParserRULE_func)
	var _la int

	p.SetState(70)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(23)

			var _m = p.Match(TParserATID)

			localctx.(*FuncContext)._ATID = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*FuncContext).FuncName = (func() string {
			if localctx.(*FuncContext).Get_ATID() == nil {
				return ""
			} else {
				return localctx.(*FuncContext).Get_ATID().GetText()
			}
		}())

		{
			p.SetState(25)
			p.Match(TParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(31)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case TParserString_:
			{
				p.SetState(26)

				var _m = p.Match(TParserString_)

				localctx.(*FuncContext)._String_ = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, GenFuncArg("", trimQuotation((func() string {
				if localctx.(*FuncContext).Get_String_() == nil {
					return ""
				} else {
					return localctx.(*FuncContext).Get_String_().GetText()
				}
			}()))))

		case TParserID:
			{
				p.SetState(28)

				var _x = p.Argument()

				localctx.(*FuncContext)._argument = _x
			}
			localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, localctx.(*FuncContext).Get_argument().GetRes())

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		p.SetState(43)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == TParserComma {
			{
				p.SetState(33)
				p.Match(TParserComma)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(39)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case TParserString_:
				{
					p.SetState(34)

					var _m = p.Match(TParserString_)

					localctx.(*FuncContext)._String_ = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, GenFuncArg("", trimQuotation((func() string {
					if localctx.(*FuncContext).Get_String_() == nil {
						return ""
					} else {
						return localctx.(*FuncContext).Get_String_().GetText()
					}
				}()))))

			case TParserID:
				{
					p.SetState(36)

					var _x = p.Argument()

					localctx.(*FuncContext)._argument = _x
				}
				localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, localctx.(*FuncContext).Get_argument().GetRes())

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(45)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(46)
			p.Match(TParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		{
			p.SetState(48)
			p.Match(TParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(49)

			var _m = p.Match(TParserATID)

			localctx.(*FuncContext)._ATID = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*FuncContext).FuncName = (func() string {
			if localctx.(*FuncContext).Get_ATID() == nil {
				return ""
			} else {
				return localctx.(*FuncContext).Get_ATID().GetText()
			}
		}())
		{
			p.SetState(51)
			p.Match(TParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(52)
			p.Match(TParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(53)
			p.Match(TParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(54)

			var _m = p.Match(TParserFieldFuncName)

			localctx.(*FuncContext)._FieldFuncName = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*FuncContext).FuncName = strings.TrimSpace((func() string {
			if localctx.(*FuncContext).Get_FieldFuncName() == nil {
				return ""
			} else {
				return localctx.(*FuncContext).Get_FieldFuncName().GetText()
			}
		}()))
		p.SetState(65)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == TParserFIELD {
			{
				p.SetState(56)

				var _m = p.Match(TParserFIELD)

				localctx.(*FuncContext)._FIELD = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, GenFuncArg("", (func() string {
				if localctx.(*FuncContext).Get_FIELD() == nil {
					return ""
				} else {
					return localctx.(*FuncContext).Get_FIELD().GetText()
				}
			}())))
			p.SetState(62)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == TParserFIELD {
				{
					p.SetState(58)

					var _m = p.Match(TParserFIELD)

					localctx.(*FuncContext)._FIELD = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, GenFuncArg("", (func() string {
					if localctx.(*FuncContext).Get_FIELD() == nil {
						return ""
					} else {
						return localctx.(*FuncContext).Get_FIELD().GetText()
					}
				}())))

				p.SetState(64)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(67)
			p.Match(TParserOLDFUNCCLOSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(68)

			var _m = p.Match(TParserATID)

			localctx.(*FuncContext)._ATID = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*FuncContext).FuncName = (func() string {
			if localctx.(*FuncContext).Get_ATID() == nil {
				return ""
			} else {
				return localctx.(*FuncContext).Get_ATID().GetText()
			}
		}())

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArgumentContext is an interface to support dynamic dispatch.
type IArgumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_ID returns the _ID token.
	Get_ID() antlr.Token

	// Get_String_ returns the _String_ token.
	Get_String_() antlr.Token

	// Set_ID sets the _ID token.
	Set_ID(antlr.Token)

	// Set_String_ sets the _String_ token.
	Set_String_(antlr.Token)

	// GetRes returns the res attribute.
	GetRes() FuncArg

	// SetRes sets the res attribute.
	SetRes(FuncArg)

	// Getter signatures
	ID() antlr.TerminalNode
	EQ() antlr.TerminalNode
	String_() antlr.TerminalNode

	// IsArgumentContext differentiates from other interfaces.
	IsArgumentContext()
}

type ArgumentContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	res      FuncArg
	_ID      antlr.Token
	_String_ antlr.Token
}

func NewEmptyArgumentContext() *ArgumentContext {
	var p = new(ArgumentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_argument
	return p
}

func InitEmptyArgumentContext(p *ArgumentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_argument
}

func (*ArgumentContext) IsArgumentContext() {}

func NewArgumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentContext {
	var p = new(ArgumentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_argument

	return p
}

func (s *ArgumentContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentContext) Get_ID() antlr.Token { return s._ID }

func (s *ArgumentContext) Get_String_() antlr.Token { return s._String_ }

func (s *ArgumentContext) Set_ID(v antlr.Token) { s._ID = v }

func (s *ArgumentContext) Set_String_(v antlr.Token) { s._String_ = v }

func (s *ArgumentContext) GetRes() FuncArg { return s.res }

func (s *ArgumentContext) SetRes(v FuncArg) { s.res = v }

func (s *ArgumentContext) ID() antlr.TerminalNode {
	return s.GetToken(TParserID, 0)
}

func (s *ArgumentContext) EQ() antlr.TerminalNode {
	return s.GetToken(TParserEQ, 0)
}

func (s *ArgumentContext) String_() antlr.TerminalNode {
	return s.GetToken(TParserString_, 0)
}

func (s *ArgumentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterArgument(s)
	}
}

func (s *ArgumentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitArgument(s)
	}
}

func (p *TParser) Argument() (localctx IArgumentContext) {
	localctx = NewArgumentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, TParserRULE_argument)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(72)

		var _m = p.Match(TParserID)

		localctx.(*ArgumentContext)._ID = _m
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(73)
		p.Match(TParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(74)

		var _m = p.Match(TParserString_)

		localctx.(*ArgumentContext)._String_ = _m
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

	localctx.(*ArgumentContext).SetRes(GenFuncArg((func() string {
		if localctx.(*ArgumentContext).Get_ID() == nil {
			return ""
		} else {
			return localctx.(*ArgumentContext).Get_ID().GetText()
		}
	}()), trimQuotation((func() string {
		if localctx.(*ArgumentContext).Get_String_() == nil {
			return ""
		} else {
			return localctx.(*ArgumentContext).Get_String_().GetText()
		}
	}()))))

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
