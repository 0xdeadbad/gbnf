package lexer

import "fmt"

type Token struct {
	Lexeme string
	Type   TokenType
	Line   uint
	Column uint
}

func NewToken(lexeme string, line, column uint, typ TokenType) *Token {
	return &Token{
		Lexeme: lexeme,
		Type:   typ,
		Line:   line,
		Column: column,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("[%s:%s:%d:%d]", t.Lexeme, t.Type, t.Line, t.Column)
}

type TokenType uint

const (
	TerminalSymbol TokenType = iota
	NonTerminalSymbol
	Assignment
	Or
	String
	Action
	ParenLeft
	ParenRight
	ActionArg
)

func (t TokenType) String() string {
	switch t {
	case TerminalSymbol:
		return "TerminalSymbol"
	case NonTerminalSymbol:
		return "NonTerminalSymbol"
	case Assignment:
		return "Assignment"
	case Or:
		return "Or"
	case String:
		return "String"
	case Action:
		return "Action"
	case ParenLeft:
		return "ParenLeft"
	case ParenRight:
		return "ParenRight"
	case ActionArg:
		return "ActionArg"
	default:
		return "Unknown"
	}
}
