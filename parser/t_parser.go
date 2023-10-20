// Code generated from .\TParser.g4 by ANTLR 4.13.0. DO NOT EDIT.

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
		"", "", "", "", "','", "'('", "')'",
	}
	staticData.SymbolicNames = []string{
		"", "ID", "String", "FieldFuncName", "Comma", "LPAREN", "RPAREN", "IGNORE_NEWLINE",
		"NEWLINE", "WS", "INSET", "S", "CLOSE", "OLDFUNCCLOSE", "FIELD", "OLDFUNCWS",
	}
	staticData.RuleNames = []string{
		"doc", "line", "func",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 15, 60, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 1, 0, 4, 0, 8, 8, 0,
		11, 0, 12, 0, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 17, 8, 1, 1, 2, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 27, 8, 2, 10, 2, 12, 2, 30,
		9, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 47, 8, 2, 10, 2, 12, 2, 50, 9, 2, 3, 2, 52,
		8, 2, 1, 2, 3, 2, 55, 8, 2, 1, 2, 3, 2, 58, 8, 2, 1, 2, 0, 0, 3, 0, 2,
		4, 0, 0, 65, 0, 7, 1, 0, 0, 0, 2, 16, 1, 0, 0, 0, 4, 57, 1, 0, 0, 0, 6,
		8, 3, 2, 1, 0, 7, 6, 1, 0, 0, 0, 8, 9, 1, 0, 0, 0, 9, 7, 1, 0, 0, 0, 9,
		10, 1, 0, 0, 0, 10, 1, 1, 0, 0, 0, 11, 17, 3, 4, 2, 0, 12, 13, 5, 10, 0,
		0, 13, 14, 5, 11, 0, 0, 14, 17, 5, 12, 0, 0, 15, 17, 5, 8, 0, 0, 16, 11,
		1, 0, 0, 0, 16, 12, 1, 0, 0, 0, 16, 15, 1, 0, 0, 0, 17, 3, 1, 0, 0, 0,
		18, 19, 5, 1, 0, 0, 19, 20, 6, 2, -1, 0, 20, 21, 5, 5, 0, 0, 21, 22, 5,
		2, 0, 0, 22, 28, 6, 2, -1, 0, 23, 24, 5, 4, 0, 0, 24, 25, 5, 2, 0, 0, 25,
		27, 6, 2, -1, 0, 26, 23, 1, 0, 0, 0, 27, 30, 1, 0, 0, 0, 28, 26, 1, 0,
		0, 0, 28, 29, 1, 0, 0, 0, 29, 31, 1, 0, 0, 0, 30, 28, 1, 0, 0, 0, 31, 32,
		5, 6, 0, 0, 32, 33, 1, 0, 0, 0, 33, 58, 5, 8, 0, 0, 34, 35, 5, 1, 0, 0,
		35, 36, 6, 2, -1, 0, 36, 37, 5, 5, 0, 0, 37, 38, 5, 6, 0, 0, 38, 58, 5,
		8, 0, 0, 39, 40, 5, 3, 0, 0, 40, 51, 6, 2, -1, 0, 41, 42, 5, 14, 0, 0,
		42, 48, 6, 2, -1, 0, 43, 44, 5, 15, 0, 0, 44, 45, 5, 14, 0, 0, 45, 47,
		6, 2, -1, 0, 46, 43, 1, 0, 0, 0, 47, 50, 1, 0, 0, 0, 48, 46, 1, 0, 0, 0,
		48, 49, 1, 0, 0, 0, 49, 52, 1, 0, 0, 0, 50, 48, 1, 0, 0, 0, 51, 41, 1,
		0, 0, 0, 51, 52, 1, 0, 0, 0, 52, 54, 1, 0, 0, 0, 53, 55, 5, 15, 0, 0, 54,
		53, 1, 0, 0, 0, 54, 55, 1, 0, 0, 0, 55, 56, 1, 0, 0, 0, 56, 58, 5, 13,
		0, 0, 57, 18, 1, 0, 0, 0, 57, 34, 1, 0, 0, 0, 57, 39, 1, 0, 0, 0, 58, 5,
		1, 0, 0, 0, 7, 9, 16, 28, 48, 51, 54, 57,
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

// TParser tokens.
const (
	TParserEOF            = antlr.TokenEOF
	TParserID             = 1
	TParserString_        = 2
	TParserFieldFuncName  = 3
	TParserComma          = 4
	TParserLPAREN         = 5
	TParserRPAREN         = 6
	TParserIGNORE_NEWLINE = 7
	TParserNEWLINE        = 8
	TParserWS             = 9
	TParserINSET          = 10
	TParserS              = 11
	TParserCLOSE          = 12
	TParserOLDFUNCCLOSE   = 13
	TParserFIELD          = 14
	TParserOLDFUNCWS      = 15
)

// TParser rules.
const (
	TParserRULE_doc  = 0
	TParserRULE_line = 1
	TParserRULE_func = 2
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
	p.SetState(7)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1290) != 0) {
		{
			p.SetState(6)
			p.Line()
		}

		p.SetState(9)
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
	S() antlr.TerminalNode
	CLOSE() antlr.TerminalNode
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

func (s *LineContext) S() antlr.TerminalNode {
	return s.GetToken(TParserS, 0)
}

func (s *LineContext) CLOSE() antlr.TerminalNode {
	return s.GetToken(TParserCLOSE, 0)
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
	p.SetState(16)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TParserID, TParserFieldFuncName:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(11)
			p.Func_()
		}

	case TParserINSET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(12)
			p.Match(TParserINSET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(13)
			p.Match(TParserS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(14)
			p.Match(TParserCLOSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case TParserNEWLINE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(15)
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

	// Get_ID returns the _ID token.
	Get_ID() antlr.Token

	// Get_String_ returns the _String_ token.
	Get_String_() antlr.Token

	// Get_FieldFuncName returns the _FieldFuncName token.
	Get_FieldFuncName() antlr.Token

	// Get_FIELD returns the _FIELD token.
	Get_FIELD() antlr.Token

	// Set_ID sets the _ID token.
	Set_ID(antlr.Token)

	// Set_String_ sets the _String_ token.
	Set_String_(antlr.Token)

	// Set_FieldFuncName sets the _FieldFuncName token.
	Set_FieldFuncName(antlr.Token)

	// Set_FIELD sets the _FIELD token.
	Set_FIELD(antlr.Token)

	// GetFuncArgs returns the FuncArgs attribute.
	GetFuncArgs() []string

	// GetFuncName returns the FuncName attribute.
	GetFuncName() string

	// SetFuncArgs sets the FuncArgs attribute.
	SetFuncArgs([]string)

	// SetFuncName sets the FuncName attribute.
	SetFuncName(string)

	// Getter signatures
	ID() antlr.TerminalNode
	NEWLINE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllString_() []antlr.TerminalNode
	String_(i int) antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllComma() []antlr.TerminalNode
	Comma(i int) antlr.TerminalNode
	FieldFuncName() antlr.TerminalNode
	OLDFUNCCLOSE() antlr.TerminalNode
	AllFIELD() []antlr.TerminalNode
	FIELD(i int) antlr.TerminalNode
	AllOLDFUNCWS() []antlr.TerminalNode
	OLDFUNCWS(i int) antlr.TerminalNode

	// IsFuncContext differentiates from other interfaces.
	IsFuncContext()
}

type FuncContext struct {
	antlr.BaseParserRuleContext
	parser         antlr.Parser
	FuncArgs       []string
	FuncName       string
	_ID            antlr.Token
	_String_       antlr.Token
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

func (s *FuncContext) Get_ID() antlr.Token { return s._ID }

func (s *FuncContext) Get_String_() antlr.Token { return s._String_ }

func (s *FuncContext) Get_FieldFuncName() antlr.Token { return s._FieldFuncName }

func (s *FuncContext) Get_FIELD() antlr.Token { return s._FIELD }

func (s *FuncContext) Set_ID(v antlr.Token) { s._ID = v }

func (s *FuncContext) Set_String_(v antlr.Token) { s._String_ = v }

func (s *FuncContext) Set_FieldFuncName(v antlr.Token) { s._FieldFuncName = v }

func (s *FuncContext) Set_FIELD(v antlr.Token) { s._FIELD = v }

func (s *FuncContext) GetFuncArgs() []string { return s.FuncArgs }

func (s *FuncContext) GetFuncName() string { return s.FuncName }

func (s *FuncContext) SetFuncArgs(v []string) { s.FuncArgs = v }

func (s *FuncContext) SetFuncName(v string) { s.FuncName = v }

func (s *FuncContext) ID() antlr.TerminalNode {
	return s.GetToken(TParserID, 0)
}

func (s *FuncContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(TParserNEWLINE, 0)
}

func (s *FuncContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(TParserLPAREN, 0)
}

func (s *FuncContext) AllString_() []antlr.TerminalNode {
	return s.GetTokens(TParserString_)
}

func (s *FuncContext) String_(i int) antlr.TerminalNode {
	return s.GetToken(TParserString_, i)
}

func (s *FuncContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(TParserRPAREN, 0)
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

func (s *FuncContext) AllOLDFUNCWS() []antlr.TerminalNode {
	return s.GetTokens(TParserOLDFUNCWS)
}

func (s *FuncContext) OLDFUNCWS(i int) antlr.TerminalNode {
	return s.GetToken(TParserOLDFUNCWS, i)
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

	var _alt int

	p.SetState(57)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(18)

			var _m = p.Match(TParserID)

			localctx.(*FuncContext)._ID = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*FuncContext).FuncName = (func() string {
			if localctx.(*FuncContext).Get_ID() == nil {
				return ""
			} else {
				return localctx.(*FuncContext).Get_ID().GetText()
			}
		}())

		{
			p.SetState(20)
			p.Match(TParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(21)

			var _m = p.Match(TParserString_)

			localctx.(*FuncContext)._String_ = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, trimQuotation((func() string {
			if localctx.(*FuncContext).Get_String_() == nil {
				return ""
			} else {
				return localctx.(*FuncContext).Get_String_().GetText()
			}
		}())))
		p.SetState(28)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == TParserComma {
			{
				p.SetState(23)
				p.Match(TParserComma)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(24)

				var _m = p.Match(TParserString_)

				localctx.(*FuncContext)._String_ = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, trimQuotation((func() string {
				if localctx.(*FuncContext).Get_String_() == nil {
					return ""
				} else {
					return localctx.(*FuncContext).Get_String_().GetText()
				}
			}())))

			p.SetState(30)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(31)
			p.Match(TParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		{
			p.SetState(33)
			p.Match(TParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(34)

			var _m = p.Match(TParserID)

			localctx.(*FuncContext)._ID = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*FuncContext).FuncName = (func() string {
			if localctx.(*FuncContext).Get_ID() == nil {
				return ""
			} else {
				return localctx.(*FuncContext).Get_ID().GetText()
			}
		}())
		{
			p.SetState(36)
			p.Match(TParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(37)
			p.Match(TParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(38)
			p.Match(TParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(39)

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
		p.SetState(51)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == TParserFIELD {
			{
				p.SetState(41)

				var _m = p.Match(TParserFIELD)

				localctx.(*FuncContext)._FIELD = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, (func() string {
				if localctx.(*FuncContext).Get_FIELD() == nil {
					return ""
				} else {
					return localctx.(*FuncContext).Get_FIELD().GetText()
				}
			}()))
			p.SetState(48)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
			for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
				if _alt == 1 {
					{
						p.SetState(43)
						p.Match(TParserOLDFUNCWS)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}
					{
						p.SetState(44)

						var _m = p.Match(TParserFIELD)

						localctx.(*FuncContext)._FIELD = _m
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}
					localctx.(*FuncContext).FuncArgs = append(localctx.(*FuncContext).FuncArgs, (func() string {
						if localctx.(*FuncContext).Get_FIELD() == nil {
							return ""
						} else {
							return localctx.(*FuncContext).Get_FIELD().GetText()
						}
					}()))

				}
				p.SetState(50)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
				if p.HasError() {
					goto errorExit
				}
			}

		}
		p.SetState(54)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == TParserOLDFUNCWS {
			{
				p.SetState(53)
				p.Match(TParserOLDFUNCWS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(56)
			p.Match(TParserOLDFUNCCLOSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

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
