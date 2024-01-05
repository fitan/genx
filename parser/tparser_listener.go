// Code generated from TParser.g4 by ANTLR 4.13.0. DO NOT EDIT.

package parser // TParser

import "github.com/antlr4-go/antlr/v4"

// TParserListener is a complete listener for a parse tree produced by TParser.
type TParserListener interface {
	antlr.ParseTreeListener

	// EnterDoc is called when entering the doc production.
	EnterDoc(c *DocContext)

	// EnterLine is called when entering the line production.
	EnterLine(c *LineContext)

	// EnterFunc is called when entering the func production.
	EnterFunc(c *FuncContext)

	// EnterArgument is called when entering the argument production.
	EnterArgument(c *ArgumentContext)

	// ExitDoc is called when exiting the doc production.
	ExitDoc(c *DocContext)

	// ExitLine is called when exiting the line production.
	ExitLine(c *LineContext)

	// ExitFunc is called when exiting the func production.
	ExitFunc(c *FuncContext)

	// ExitArgument is called when exiting the argument production.
	ExitArgument(c *ArgumentContext)
}
