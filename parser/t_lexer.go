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
		"", "ID", "String", "FieldFuncName", "Comma", "LPAREN", "RPAREN", "NEWLINE",
		"WS", "INSET", "S", "CLOSE", "OLDFUNCCLOSE", "FIELD", "OLDFUNCWS",
	}
	staticData.RuleNames = []string{
		"ID", "String", "FieldFuncName", "Comma", "LPAREN", "RPAREN", "NEWLINE",
		"WS", "INSET", "S", "CLOSE", "OLDFUNCCLOSE", "FIELD", "OLDFUNCWS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 14, 122, 6, -1, 6, -1, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2,
		2, 3, 7, 3, 2, 4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8,
		2, 9, 7, 9, 2, 10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 1,
		0, 1, 0, 1, 0, 5, 0, 35, 8, 0, 10, 0, 12, 0, 38, 9, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 5, 1, 44, 8, 1, 10, 1, 12, 1, 47, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 5, 1, 54, 8, 1, 10, 1, 12, 1, 57, 9, 1, 1, 1, 3, 1, 60, 8, 1, 1, 2,
		1, 2, 4, 2, 64, 8, 2, 11, 2, 12, 2, 65, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1,
		4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 3, 6, 79, 8, 6, 1, 6, 1, 6, 1, 7, 4, 7,
		84, 8, 7, 11, 7, 12, 7, 85, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 4,
		9, 95, 8, 9, 11, 9, 12, 9, 96, 1, 10, 3, 10, 100, 8, 10, 1, 10, 1, 10,
		1, 10, 1, 10, 1, 11, 3, 11, 107, 8, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1,
		12, 4, 12, 114, 8, 12, 11, 12, 12, 12, 115, 1, 13, 4, 13, 119, 8, 13, 11,
		13, 12, 13, 120, 0, 0, 14, 3, 1, 5, 2, 7, 3, 9, 4, 11, 5, 13, 6, 15, 7,
		17, 8, 19, 9, 21, 10, 23, 11, 25, 12, 27, 13, 29, 14, 3, 0, 1, 2, 8, 3,
		0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 1, 0,
		39, 39, 1, 0, 34, 34, 2, 0, 9, 9, 32, 32, 1, 0, 64, 64, 2, 0, 10, 10, 13,
		13, 3, 0, 9, 10, 13, 13, 32, 32, 133, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0,
		0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0,
		0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 1, 21, 1,
		0, 0, 0, 1, 23, 1, 0, 0, 0, 2, 25, 1, 0, 0, 0, 2, 27, 1, 0, 0, 0, 2, 29,
		1, 0, 0, 0, 3, 31, 1, 0, 0, 0, 5, 59, 1, 0, 0, 0, 7, 61, 1, 0, 0, 0, 9,
		69, 1, 0, 0, 0, 11, 71, 1, 0, 0, 0, 13, 74, 1, 0, 0, 0, 15, 78, 1, 0, 0,
		0, 17, 83, 1, 0, 0, 0, 19, 89, 1, 0, 0, 0, 21, 94, 1, 0, 0, 0, 23, 99,
		1, 0, 0, 0, 25, 106, 1, 0, 0, 0, 27, 113, 1, 0, 0, 0, 29, 118, 1, 0, 0,
		0, 31, 32, 5, 64, 0, 0, 32, 36, 7, 0, 0, 0, 33, 35, 7, 1, 0, 0, 34, 33,
		1, 0, 0, 0, 35, 38, 1, 0, 0, 0, 36, 34, 1, 0, 0, 0, 36, 37, 1, 0, 0, 0,
		37, 4, 1, 0, 0, 0, 38, 36, 1, 0, 0, 0, 39, 45, 5, 39, 0, 0, 40, 41, 5,
		39, 0, 0, 41, 44, 5, 39, 0, 0, 42, 44, 8, 2, 0, 0, 43, 40, 1, 0, 0, 0,
		43, 42, 1, 0, 0, 0, 44, 47, 1, 0, 0, 0, 45, 43, 1, 0, 0, 0, 45, 46, 1,
		0, 0, 0, 46, 48, 1, 0, 0, 0, 47, 45, 1, 0, 0, 0, 48, 60, 5, 39, 0, 0, 49,
		55, 5, 34, 0, 0, 50, 51, 5, 34, 0, 0, 51, 54, 5, 34, 0, 0, 52, 54, 8, 3,
		0, 0, 53, 50, 1, 0, 0, 0, 53, 52, 1, 0, 0, 0, 54, 57, 1, 0, 0, 0, 55, 53,
		1, 0, 0, 0, 55, 56, 1, 0, 0, 0, 56, 58, 1, 0, 0, 0, 57, 55, 1, 0, 0, 0,
		58, 60, 5, 34, 0, 0, 59, 39, 1, 0, 0, 0, 59, 49, 1, 0, 0, 0, 60, 6, 1,
		0, 0, 0, 61, 63, 3, 3, 0, 0, 62, 64, 5, 32, 0, 0, 63, 62, 1, 0, 0, 0, 64,
		65, 1, 0, 0, 0, 65, 63, 1, 0, 0, 0, 65, 66, 1, 0, 0, 0, 66, 67, 1, 0, 0,
		0, 67, 68, 6, 2, 0, 0, 68, 8, 1, 0, 0, 0, 69, 70, 5, 44, 0, 0, 70, 10,
		1, 0, 0, 0, 71, 72, 5, 40, 0, 0, 72, 73, 6, 4, 1, 0, 73, 12, 1, 0, 0, 0,
		74, 75, 5, 41, 0, 0, 75, 76, 6, 5, 2, 0, 76, 14, 1, 0, 0, 0, 77, 79, 5,
		13, 0, 0, 78, 77, 1, 0, 0, 0, 78, 79, 1, 0, 0, 0, 79, 80, 1, 0, 0, 0, 80,
		81, 5, 10, 0, 0, 81, 16, 1, 0, 0, 0, 82, 84, 7, 4, 0, 0, 83, 82, 1, 0,
		0, 0, 84, 85, 1, 0, 0, 0, 85, 83, 1, 0, 0, 0, 85, 86, 1, 0, 0, 0, 86, 87,
		1, 0, 0, 0, 87, 88, 6, 7, 3, 0, 88, 18, 1, 0, 0, 0, 89, 90, 8, 5, 0, 0,
		90, 91, 1, 0, 0, 0, 91, 92, 6, 8, 4, 0, 92, 20, 1, 0, 0, 0, 93, 95, 8,
		6, 0, 0, 94, 93, 1, 0, 0, 0, 95, 96, 1, 0, 0, 0, 96, 94, 1, 0, 0, 0, 96,
		97, 1, 0, 0, 0, 97, 22, 1, 0, 0, 0, 98, 100, 5, 13, 0, 0, 99, 98, 1, 0,
		0, 0, 99, 100, 1, 0, 0, 0, 100, 101, 1, 0, 0, 0, 101, 102, 5, 10, 0, 0,
		102, 103, 1, 0, 0, 0, 103, 104, 6, 10, 5, 0, 104, 24, 1, 0, 0, 0, 105,
		107, 5, 13, 0, 0, 106, 105, 1, 0, 0, 0, 106, 107, 1, 0, 0, 0, 107, 108,
		1, 0, 0, 0, 108, 109, 5, 10, 0, 0, 109, 110, 1, 0, 0, 0, 110, 111, 6, 11,
		5, 0, 111, 26, 1, 0, 0, 0, 112, 114, 8, 7, 0, 0, 113, 112, 1, 0, 0, 0,
		114, 115, 1, 0, 0, 0, 115, 113, 1, 0, 0, 0, 115, 116, 1, 0, 0, 0, 116,
		28, 1, 0, 0, 0, 117, 119, 7, 4, 0, 0, 118, 117, 1, 0, 0, 0, 119, 120, 1,
		0, 0, 0, 120, 118, 1, 0, 0, 0, 120, 121, 1, 0, 0, 0, 121, 30, 1, 0, 0,
		0, 17, 0, 1, 2, 36, 43, 45, 53, 55, 59, 65, 78, 85, 96, 99, 106, 115, 120,
		6, 5, 2, 0, 1, 4, 0, 1, 5, 1, 6, 0, 0, 5, 1, 0, 4, 0, 0,
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
	TLexerID            = 1
	TLexerString_       = 2
	TLexerFieldFuncName = 3
	TLexerComma         = 4
	TLexerLPAREN        = 5
	TLexerRPAREN        = 6
	TLexerNEWLINE       = 7
	TLexerWS            = 8
	TLexerINSET         = 9
	TLexerS             = 10
	TLexerCLOSE         = 11
	TLexerOLDFUNCCLOSE  = 12
	TLexerFIELD         = 13
	TLexerOLDFUNCWS     = 14
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

	default:
		panic("No registered action for: " + fmt.Sprint(ruleIndex))
	}
}

func (l *TLexer) LPAREN_Action(localctx antlr.RuleContext, actionIndex int) {
	switch actionIndex {
	case 0:
		nesting++

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
func (l *TLexer) RPAREN_Action(localctx antlr.RuleContext, actionIndex int) {
	switch actionIndex {
	case 1:
		nesting--

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
