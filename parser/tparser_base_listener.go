// Code generated from TParser.g4 by ANTLR 4.13.0. DO NOT EDIT.

package parser // TParser

import "github.com/antlr4-go/antlr/v4"

// BaseTParserListener is a complete listener for a parse tree produced by TParser.
type BaseTParserListener struct{}

var _ TParserListener = &BaseTParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseTParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseTParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseTParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseTParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterDoc is called when production doc is entered.
func (s *BaseTParserListener) EnterDoc(ctx *DocContext) {}

// ExitDoc is called when production doc is exited.
func (s *BaseTParserListener) ExitDoc(ctx *DocContext) {}

// EnterLine is called when production line is entered.
func (s *BaseTParserListener) EnterLine(ctx *LineContext) {}

// ExitLine is called when production line is exited.
func (s *BaseTParserListener) ExitLine(ctx *LineContext) {}

// EnterFunc is called when production func is entered.
func (s *BaseTParserListener) EnterFunc(ctx *FuncContext) {}

// ExitFunc is called when production func is exited.
func (s *BaseTParserListener) ExitFunc(ctx *FuncContext) {}

// EnterArgument is called when production argument is entered.
func (s *BaseTParserListener) EnterArgument(ctx *ArgumentContext) {}

// ExitArgument is called when production argument is exited.
func (s *BaseTParserListener) ExitArgument(ctx *ArgumentContext) {}
