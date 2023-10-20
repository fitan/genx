// Code generated from .\TLexer.g4 by ANTLR 4.13.0. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type TLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var TLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func tlexerLexerInit() {
	staticData := &TLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE", "INSIDE", "OLDFUNC",
	}
	staticData.LiteralNames = []string{
		"", "", "", "", "','", "'('", "')'",
	}
	staticData.SymbolicNames = []string{
		"", "ID", "String", "FieldFuncName", "Comma", "LPAREN", "RPAREN", "IGNORE_NEWLINE",
		"NEWLINE", "WS", "INSET", "S", "CLOSE", "OLDFUNCCLOSE", "FIELD", "OLDFUNCWS",
	}
	staticData.RuleNames = []string{
		"ID", "String", "FieldFuncName", "Comma", "LPAREN", "RPAREN", "IGNORE_NEWLINE",
		"NEWLINE", "WS", "INSET", "S", "CLOSE", "OLDFUNCCLOSE", "FIELD", "OLDFUNCWS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 15, 135, 6, -1, 6, -1, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2,
		2, 3, 7, 3, 2, 4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8,
		2, 9, 7, 9, 2, 10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2,
		14, 7, 14, 1, 0, 1, 0, 1, 0, 5, 0, 37, 8, 0, 10, 0, 12, 0, 40, 9, 0, 1,
		1, 1, 1, 1, 1, 1, 1, 5, 1, 46, 8, 1, 10, 1, 12, 1, 49, 9, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 5, 1, 56, 8, 1, 10, 1, 12, 1, 59, 9, 1, 1, 1, 3, 1, 62,
		8, 1, 1, 2, 1, 2, 4, 2, 66, 8, 2, 11, 2, 12, 2, 67, 1, 2, 1, 2, 1, 3, 1,
		3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 3, 6, 81, 8, 6, 1, 6, 1, 6,
		1, 6, 1, 6, 1, 6, 1, 7, 3, 7, 89, 8, 7, 1, 7, 1, 7, 1, 8, 4, 8, 94, 8,
		8, 11, 8, 12, 8, 95, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 4,
		10, 106, 8, 10, 11, 10, 12, 10, 107, 1, 11, 3, 11, 111, 8, 11, 1, 11, 1,
		11, 1, 11, 1, 11, 1, 11, 1, 12, 3, 12, 119, 8, 12, 1, 12, 1, 12, 1, 12,
		1, 12, 1, 12, 1, 13, 4, 13, 127, 8, 13, 11, 13, 12, 13, 128, 1, 14, 4,
		14, 132, 8, 14, 11, 14, 12, 14, 133, 0, 0, 15, 3, 1, 5, 2, 7, 3, 9, 4,
		11, 5, 13, 6, 15, 7, 17, 8, 19, 9, 21, 10, 23, 11, 25, 12, 27, 13, 29,
		14, 31, 15, 3, 0, 1, 2, 8, 3, 0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57,
		65, 90, 95, 95, 97, 122, 1, 0, 39, 39, 1, 0, 34, 34, 2, 0, 9, 9, 32, 32,
		1, 0, 64, 64, 2, 0, 10, 10, 13, 13, 3, 0, 9, 10, 13, 13, 32, 32, 147, 0,
		3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0,
		11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0,
		0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 1, 23, 1, 0, 0, 0, 1, 25, 1, 0, 0,
		0, 2, 27, 1, 0, 0, 0, 2, 29, 1, 0, 0, 0, 2, 31, 1, 0, 0, 0, 3, 33, 1, 0,
		0, 0, 5, 61, 1, 0, 0, 0, 7, 63, 1, 0, 0, 0, 9, 71, 1, 0, 0, 0, 11, 73,
		1, 0, 0, 0, 13, 76, 1, 0, 0, 0, 15, 80, 1, 0, 0, 0, 17, 88, 1, 0, 0, 0,
		19, 93, 1, 0, 0, 0, 21, 99, 1, 0, 0, 0, 23, 105, 1, 0, 0, 0, 25, 110, 1,
		0, 0, 0, 27, 118, 1, 0, 0, 0, 29, 126, 1, 0, 0, 0, 31, 131, 1, 0, 0, 0,
		33, 34, 5, 64, 0, 0, 34, 38, 7, 0, 0, 0, 35, 37, 7, 1, 0, 0, 36, 35, 1,
		0, 0, 0, 37, 40, 1, 0, 0, 0, 38, 36, 1, 0, 0, 0, 38, 39, 1, 0, 0, 0, 39,
		4, 1, 0, 0, 0, 40, 38, 1, 0, 0, 0, 41, 47, 5, 39, 0, 0, 42, 43, 5, 39,
		0, 0, 43, 46, 5, 39, 0, 0, 44, 46, 8, 2, 0, 0, 45, 42, 1, 0, 0, 0, 45,
		44, 1, 0, 0, 0, 46, 49, 1, 0, 0, 0, 47, 45, 1, 0, 0, 0, 47, 48, 1, 0, 0,
		0, 48, 50, 1, 0, 0, 0, 49, 47, 1, 0, 0, 0, 50, 62, 5, 39, 0, 0, 51, 57,
		5, 34, 0, 0, 52, 53, 5, 34, 0, 0, 53, 56, 5, 34, 0, 0, 54, 56, 8, 3, 0,
		0, 55, 52, 1, 0, 0, 0, 55, 54, 1, 0, 0, 0, 56, 59, 1, 0, 0, 0, 57, 55,
		1, 0, 0, 0, 57, 58, 1, 0, 0, 0, 58, 60, 1, 0, 0, 0, 59, 57, 1, 0, 0, 0,
		60, 62, 5, 34, 0, 0, 61, 41, 1, 0, 0, 0, 61, 51, 1, 0, 0, 0, 62, 6, 1,
		0, 0, 0, 63, 65, 3, 3, 0, 0, 64, 66, 5, 32, 0, 0, 65, 64, 1, 0, 0, 0, 66,
		67, 1, 0, 0, 0, 67, 65, 1, 0, 0, 0, 67, 68, 1, 0, 0, 0, 68, 69, 1, 0, 0,
		0, 69, 70, 6, 2, 0, 0, 70, 8, 1, 0, 0, 0, 71, 72, 5, 44, 0, 0, 72, 10,
		1, 0, 0, 0, 73, 74, 5, 40, 0, 0, 74, 75, 6, 4, 1, 0, 75, 12, 1, 0, 0, 0,
		76, 77, 5, 41, 0, 0, 77, 78, 6, 5, 2, 0, 78, 14, 1, 0, 0, 0, 79, 81, 5,
		13, 0, 0, 80, 79, 1, 0, 0, 0, 80, 81, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82,
		83, 5, 10, 0, 0, 83, 84, 4, 6, 0, 0, 84, 85, 1, 0, 0, 0, 85, 86, 6, 6,
		3, 0, 86, 16, 1, 0, 0, 0, 87, 89, 5, 13, 0, 0, 88, 87, 1, 0, 0, 0, 88,
		89, 1, 0, 0, 0, 89, 90, 1, 0, 0, 0, 90, 91, 5, 10, 0, 0, 91, 18, 1, 0,
		0, 0, 92, 94, 7, 4, 0, 0, 93, 92, 1, 0, 0, 0, 94, 95, 1, 0, 0, 0, 95, 93,
		1, 0, 0, 0, 95, 96, 1, 0, 0, 0, 96, 97, 1, 0, 0, 0, 97, 98, 6, 8, 3, 0,
		98, 20, 1, 0, 0, 0, 99, 100, 8, 5, 0, 0, 100, 101, 6, 9, 4, 0, 101, 102,
		1, 0, 0, 0, 102, 103, 6, 9, 5, 0, 103, 22, 1, 0, 0, 0, 104, 106, 8, 6,
		0, 0, 105, 104, 1, 0, 0, 0, 106, 107, 1, 0, 0, 0, 107, 105, 1, 0, 0, 0,
		107, 108, 1, 0, 0, 0, 108, 24, 1, 0, 0, 0, 109, 111, 5, 13, 0, 0, 110,
		109, 1, 0, 0, 0, 110, 111, 1, 0, 0, 0, 111, 112, 1, 0, 0, 0, 112, 113,
		5, 10, 0, 0, 113, 114, 6, 11, 6, 0, 114, 115, 1, 0, 0, 0, 115, 116, 6,
		11, 7, 0, 116, 26, 1, 0, 0, 0, 117, 119, 5, 13, 0, 0, 118, 117, 1, 0, 0,
		0, 118, 119, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120, 121, 5, 10, 0, 0, 121,
		122, 6, 12, 8, 0, 122, 123, 1, 0, 0, 0, 123, 124, 6, 12, 7, 0, 124, 28,
		1, 0, 0, 0, 125, 127, 8, 7, 0, 0, 126, 125, 1, 0, 0, 0, 127, 128, 1, 0,
		0, 0, 128, 126, 1, 0, 0, 0, 128, 129, 1, 0, 0, 0, 129, 30, 1, 0, 0, 0,
		130, 132, 7, 4, 0, 0, 131, 130, 1, 0, 0, 0, 132, 133, 1, 0, 0, 0, 133,
		131, 1, 0, 0, 0, 133, 134, 1, 0, 0, 0, 134, 32, 1, 0, 0, 0, 18, 0, 1, 2,
		38, 45, 47, 55, 57, 61, 67, 80, 88, 95, 107, 110, 118, 128, 133, 9, 5,
		2, 0, 1, 4, 0, 1, 5, 1, 6, 0, 0, 1, 9, 2, 5, 1, 0, 1, 11, 3, 4, 0, 0, 1,
		12, 4,
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

// TLexerInit initializes any static state used to implement TLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewTLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func TLexerInit() {
	staticData := &TLexerLexerStaticData
	staticData.once.Do(tlexerLexerInit)
}

// NewTLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewTLexer(input antlr.CharStream) *TLexer {
	TLexerInit()
	l := new(TLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &TLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "TLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// TLexer tokens.
const (
	TLexerID             = 1
	TLexerString_        = 2
	TLexerFieldFuncName  = 3
	TLexerComma          = 4
	TLexerLPAREN         = 5
	TLexerRPAREN         = 6
	TLexerIGNORE_NEWLINE = 7
	TLexerNEWLINE        = 8
	TLexerWS             = 9
	TLexerINSET          = 10
	TLexerS              = 11
	TLexerCLOSE          = 12
	TLexerOLDFUNCCLOSE   = 13
	TLexerFIELD          = 14
	TLexerOLDFUNCWS      = 15
)

// TLexer modes.
const (
	TLexerINSIDE = iota + 1
	TLexerOLDFUNC
)

var nesting int

func (l *TLexer) Action(localctx antlr.RuleContext, ruleIndex, actionIndex int) {
	switch ruleIndex {
	case 4:
		l.LPAREN_Action(localctx, actionIndex)

	case 5:
		l.RPAREN_Action(localctx, actionIndex)

	case 9:
		l.INSET_Action(localctx, actionIndex)

	case 11:
		l.CLOSE_Action(localctx, actionIndex)

	case 12:
		l.OLDFUNCCLOSE_Action(localctx, actionIndex)

	default:
		panic("No registered action for: " + fmt.Sprint(ruleIndex))
	}
}

func (l *TLexer) LPAREN_Action(localctx antlr.RuleContext, actionIndex int) {
	switch actionIndex {
	case 0:
		fmt.Println(nesting)
		nesting++

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *TLexer) RPAREN_Action(localctx antlr.RuleContext, actionIndex int) {
	switch actionIndex {
	case 1:
		fmt.Println(nesting)
		nesting--

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *TLexer) INSET_Action(localctx antlr.RuleContext, actionIndex int) {
	switch actionIndex {
	case 2:
		fmt.Println("ent INSIDE")

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *TLexer) CLOSE_Action(localctx antlr.RuleContext, actionIndex int) {
	switch actionIndex {
	case 3:
		fmt.Println("out INSIDE")

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *TLexer) OLDFUNCCLOSE_Action(localctx antlr.RuleContext, actionIndex int) {
	switch actionIndex {
	case 4:
		fmt.Println("out Fied")

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}

func (l *TLexer) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 6:
		return l.IGNORE_NEWLINE_Sempred(localctx, predIndex)

	default:
		panic("No registered predicate for: " + fmt.Sprint(ruleIndex))
	}
}

func (p *TLexer) IGNORE_NEWLINE_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return nesting > 0

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
