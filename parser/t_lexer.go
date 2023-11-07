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
		"DEFAULT_MODE", "PAREN", "INSIDE", "OLDFUNC",
	}
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "'='", "','",
	}
	staticData.SymbolicNames = []string{
		"", "ATID", "FieldFuncName", "LPAREN", "NEWLINE", "WS", "INSET", "EQ",
		"Comma", "PARENWS", "ID", "String", "RPAREN", "S", "CLOSE", "OLDFUNCCLOSE",
		"OLDFUNCWS", "FIELD",
	}
	staticData.RuleNames = []string{
		"ATID", "FieldFuncName", "LPAREN", "NEWLINE", "WS", "INSET", "EQ", "Comma",
		"PARENWS", "ID", "String", "RPAREN", "S", "CLOSE", "OLDFUNCCLOSE", "OLDFUNCWS",
		"FIELD",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 17, 179, 6, -1, 6, -1, 6, -1, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2,
		7, 2, 2, 3, 7, 3, 2, 4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8,
		7, 8, 2, 9, 7, 9, 2, 10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13,
		2, 14, 7, 14, 2, 15, 7, 15, 2, 16, 7, 16, 1, 0, 1, 0, 1, 0, 5, 0, 42, 8,
		0, 10, 0, 12, 0, 45, 9, 0, 1, 1, 1, 1, 4, 1, 49, 8, 1, 11, 1, 12, 1, 50,
		1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 3, 3, 61, 8, 3, 1, 3, 1,
		3, 3, 3, 65, 8, 3, 1, 4, 4, 4, 68, 8, 4, 11, 4, 12, 4, 69, 1, 4, 1, 4,
		1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 4, 8, 83, 8, 8, 11,
		8, 12, 8, 84, 1, 8, 1, 8, 1, 9, 1, 9, 5, 9, 91, 8, 9, 10, 9, 12, 9, 94,
		9, 9, 1, 10, 1, 10, 1, 10, 1, 10, 5, 10, 100, 8, 10, 10, 10, 12, 10, 103,
		9, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 5, 10, 110, 8, 10, 10, 10, 12,
		10, 113, 9, 10, 1, 10, 3, 10, 116, 8, 10, 1, 11, 1, 11, 1, 11, 1, 11, 1,
		11, 1, 12, 4, 12, 124, 8, 12, 11, 12, 12, 12, 125, 1, 13, 3, 13, 129, 8,
		13, 1, 13, 1, 13, 3, 13, 133, 8, 13, 1, 13, 1, 13, 1, 14, 3, 14, 138, 8,
		14, 1, 14, 1, 14, 3, 14, 142, 8, 14, 1, 14, 1, 14, 1, 15, 4, 15, 147, 8,
		15, 11, 15, 12, 15, 148, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 16, 5, 16,
		157, 8, 16, 10, 16, 12, 16, 160, 9, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1,
		16, 5, 16, 167, 8, 16, 10, 16, 12, 16, 170, 9, 16, 1, 16, 1, 16, 4, 16,
		174, 8, 16, 11, 16, 12, 16, 175, 3, 16, 178, 8, 16, 0, 0, 17, 4, 1, 6,
		2, 8, 3, 10, 4, 12, 5, 14, 6, 16, 7, 18, 8, 20, 9, 22, 10, 24, 11, 26,
		12, 28, 13, 30, 14, 32, 15, 34, 16, 36, 17, 4, 0, 1, 2, 3, 8, 3, 0, 65,
		90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 2, 0, 9, 9,
		32, 32, 1, 0, 64, 64, 1, 0, 39, 39, 1, 0, 34, 34, 2, 0, 10, 10, 13, 13,
		3, 0, 9, 10, 13, 13, 32, 32, 200, 0, 4, 1, 0, 0, 0, 0, 6, 1, 0, 0, 0, 0,
		8, 1, 0, 0, 0, 0, 10, 1, 0, 0, 0, 0, 12, 1, 0, 0, 0, 0, 14, 1, 0, 0, 0,
		1, 16, 1, 0, 0, 0, 1, 18, 1, 0, 0, 0, 1, 20, 1, 0, 0, 0, 1, 22, 1, 0, 0,
		0, 1, 24, 1, 0, 0, 0, 1, 26, 1, 0, 0, 0, 2, 28, 1, 0, 0, 0, 2, 30, 1, 0,
		0, 0, 3, 32, 1, 0, 0, 0, 3, 34, 1, 0, 0, 0, 3, 36, 1, 0, 0, 0, 4, 38, 1,
		0, 0, 0, 6, 46, 1, 0, 0, 0, 8, 54, 1, 0, 0, 0, 10, 64, 1, 0, 0, 0, 12,
		67, 1, 0, 0, 0, 14, 73, 1, 0, 0, 0, 16, 77, 1, 0, 0, 0, 18, 79, 1, 0, 0,
		0, 20, 82, 1, 0, 0, 0, 22, 88, 1, 0, 0, 0, 24, 115, 1, 0, 0, 0, 26, 117,
		1, 0, 0, 0, 28, 123, 1, 0, 0, 0, 30, 132, 1, 0, 0, 0, 32, 141, 1, 0, 0,
		0, 34, 146, 1, 0, 0, 0, 36, 177, 1, 0, 0, 0, 38, 39, 5, 64, 0, 0, 39, 43,
		7, 0, 0, 0, 40, 42, 7, 1, 0, 0, 41, 40, 1, 0, 0, 0, 42, 45, 1, 0, 0, 0,
		43, 41, 1, 0, 0, 0, 43, 44, 1, 0, 0, 0, 44, 5, 1, 0, 0, 0, 45, 43, 1, 0,
		0, 0, 46, 48, 3, 4, 0, 0, 47, 49, 5, 32, 0, 0, 48, 47, 1, 0, 0, 0, 49,
		50, 1, 0, 0, 0, 50, 48, 1, 0, 0, 0, 50, 51, 1, 0, 0, 0, 51, 52, 1, 0, 0,
		0, 52, 53, 6, 1, 0, 0, 53, 7, 1, 0, 0, 0, 54, 55, 5, 40, 0, 0, 55, 56,
		6, 2, 1, 0, 56, 57, 1, 0, 0, 0, 57, 58, 6, 2, 2, 0, 58, 9, 1, 0, 0, 0,
		59, 61, 5, 13, 0, 0, 60, 59, 1, 0, 0, 0, 60, 61, 1, 0, 0, 0, 61, 62, 1,
		0, 0, 0, 62, 65, 5, 10, 0, 0, 63, 65, 5, 0, 0, 1, 64, 60, 1, 0, 0, 0, 64,
		63, 1, 0, 0, 0, 65, 11, 1, 0, 0, 0, 66, 68, 7, 2, 0, 0, 67, 66, 1, 0, 0,
		0, 68, 69, 1, 0, 0, 0, 69, 67, 1, 0, 0, 0, 69, 70, 1, 0, 0, 0, 70, 71,
		1, 0, 0, 0, 71, 72, 6, 4, 3, 0, 72, 13, 1, 0, 0, 0, 73, 74, 8, 3, 0, 0,
		74, 75, 1, 0, 0, 0, 75, 76, 6, 5, 4, 0, 76, 15, 1, 0, 0, 0, 77, 78, 5,
		61, 0, 0, 78, 17, 1, 0, 0, 0, 79, 80, 5, 44, 0, 0, 80, 19, 1, 0, 0, 0,
		81, 83, 7, 2, 0, 0, 82, 81, 1, 0, 0, 0, 83, 84, 1, 0, 0, 0, 84, 82, 1,
		0, 0, 0, 84, 85, 1, 0, 0, 0, 85, 86, 1, 0, 0, 0, 86, 87, 6, 8, 3, 0, 87,
		21, 1, 0, 0, 0, 88, 92, 7, 0, 0, 0, 89, 91, 7, 1, 0, 0, 90, 89, 1, 0, 0,
		0, 91, 94, 1, 0, 0, 0, 92, 90, 1, 0, 0, 0, 92, 93, 1, 0, 0, 0, 93, 23,
		1, 0, 0, 0, 94, 92, 1, 0, 0, 0, 95, 101, 5, 39, 0, 0, 96, 97, 5, 39, 0,
		0, 97, 100, 5, 39, 0, 0, 98, 100, 8, 4, 0, 0, 99, 96, 1, 0, 0, 0, 99, 98,
		1, 0, 0, 0, 100, 103, 1, 0, 0, 0, 101, 99, 1, 0, 0, 0, 101, 102, 1, 0,
		0, 0, 102, 104, 1, 0, 0, 0, 103, 101, 1, 0, 0, 0, 104, 116, 5, 39, 0, 0,
		105, 111, 5, 34, 0, 0, 106, 107, 5, 34, 0, 0, 107, 110, 5, 34, 0, 0, 108,
		110, 8, 5, 0, 0, 109, 106, 1, 0, 0, 0, 109, 108, 1, 0, 0, 0, 110, 113,
		1, 0, 0, 0, 111, 109, 1, 0, 0, 0, 111, 112, 1, 0, 0, 0, 112, 114, 1, 0,
		0, 0, 113, 111, 1, 0, 0, 0, 114, 116, 5, 34, 0, 0, 115, 95, 1, 0, 0, 0,
		115, 105, 1, 0, 0, 0, 116, 25, 1, 0, 0, 0, 117, 118, 5, 41, 0, 0, 118,
		119, 6, 11, 5, 0, 119, 120, 1, 0, 0, 0, 120, 121, 6, 11, 6, 0, 121, 27,
		1, 0, 0, 0, 122, 124, 8, 6, 0, 0, 123, 122, 1, 0, 0, 0, 124, 125, 1, 0,
		0, 0, 125, 123, 1, 0, 0, 0, 125, 126, 1, 0, 0, 0, 126, 29, 1, 0, 0, 0,
		127, 129, 5, 13, 0, 0, 128, 127, 1, 0, 0, 0, 128, 129, 1, 0, 0, 0, 129,
		130, 1, 0, 0, 0, 130, 133, 5, 10, 0, 0, 131, 133, 5, 0, 0, 1, 132, 128,
		1, 0, 0, 0, 132, 131, 1, 0, 0, 0, 133, 134, 1, 0, 0, 0, 134, 135, 6, 13,
		6, 0, 135, 31, 1, 0, 0, 0, 136, 138, 5, 13, 0, 0, 137, 136, 1, 0, 0, 0,
		137, 138, 1, 0, 0, 0, 138, 139, 1, 0, 0, 0, 139, 142, 5, 10, 0, 0, 140,
		142, 5, 0, 0, 1, 141, 137, 1, 0, 0, 0, 141, 140, 1, 0, 0, 0, 142, 143,
		1, 0, 0, 0, 143, 144, 6, 14, 6, 0, 144, 33, 1, 0, 0, 0, 145, 147, 7, 2,
		0, 0, 146, 145, 1, 0, 0, 0, 147, 148, 1, 0, 0, 0, 148, 146, 1, 0, 0, 0,
		148, 149, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150, 151, 6, 15, 3, 0, 151,
		35, 1, 0, 0, 0, 152, 158, 5, 39, 0, 0, 153, 154, 5, 39, 0, 0, 154, 157,
		5, 39, 0, 0, 155, 157, 8, 4, 0, 0, 156, 153, 1, 0, 0, 0, 156, 155, 1, 0,
		0, 0, 157, 160, 1, 0, 0, 0, 158, 156, 1, 0, 0, 0, 158, 159, 1, 0, 0, 0,
		159, 161, 1, 0, 0, 0, 160, 158, 1, 0, 0, 0, 161, 178, 5, 39, 0, 0, 162,
		168, 5, 34, 0, 0, 163, 164, 5, 34, 0, 0, 164, 167, 5, 34, 0, 0, 165, 167,
		8, 5, 0, 0, 166, 163, 1, 0, 0, 0, 166, 165, 1, 0, 0, 0, 167, 170, 1, 0,
		0, 0, 168, 166, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169, 171, 1, 0, 0, 0,
		170, 168, 1, 0, 0, 0, 171, 178, 5, 34, 0, 0, 172, 174, 8, 7, 0, 0, 173,
		172, 1, 0, 0, 0, 174, 175, 1, 0, 0, 0, 175, 173, 1, 0, 0, 0, 175, 176,
		1, 0, 0, 0, 176, 178, 1, 0, 0, 0, 177, 152, 1, 0, 0, 0, 177, 162, 1, 0,
		0, 0, 177, 173, 1, 0, 0, 0, 178, 37, 1, 0, 0, 0, 28, 0, 1, 2, 3, 43, 50,
		60, 64, 69, 84, 92, 99, 101, 109, 111, 115, 125, 128, 132, 137, 141, 148,
		156, 158, 166, 168, 175, 177, 7, 5, 3, 0, 1, 2, 0, 5, 1, 0, 6, 0, 0, 5,
		2, 0, 1, 11, 1, 4, 0, 0,
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
	TLexerATID          = 1
	TLexerFieldFuncName = 2
	TLexerLPAREN        = 3
	TLexerNEWLINE       = 4
	TLexerWS            = 5
	TLexerINSET         = 6
	TLexerEQ            = 7
	TLexerComma         = 8
	TLexerPARENWS       = 9
	TLexerID            = 10
	TLexerString_       = 11
	TLexerRPAREN        = 12
	TLexerS             = 13
	TLexerCLOSE         = 14
	TLexerOLDFUNCCLOSE  = 15
	TLexerOLDFUNCWS     = 16
	TLexerFIELD         = 17
)

// TLexer modes.
const (
	TLexerPAREN = iota + 1
	TLexerINSIDE
	TLexerOLDFUNC
)

var nesting int

func (l *TLexer) Action(localctx antlr.RuleContext, ruleIndex, actionIndex int) {
	switch ruleIndex {
	case 2:
		l.LPAREN_Action(localctx, actionIndex)

	case 11:
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
