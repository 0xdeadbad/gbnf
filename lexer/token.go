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

func (t *Token) TokenType() TokenType {
	return t.Type
}

func (t *Token) TokenLexeme() string {
	return t.Lexeme
}

type TokenType uint

const (
	TerminalSymbol TokenType = iota
	NonTerminalSymbol
	ProdRule
	Or
	And
	String
	Action
	ParenLeft
	ParenRight
	BracketLeft
	BracketRight
	ActionArg
	Assign
	Sequence
	EndMark
	EndOfRule
	Not
)

func (t TokenType) String() string {
	switch t {
	case TerminalSymbol:
		return "TerminalSymbol"
	case NonTerminalSymbol:
		return "NonTerminalSymbol"
	case ProdRule:
		return "ProdRule"
	case Or:
		return "Or"
	case And:
		return "And"
	case String:
		return "String"
	case Action:
		return "Action"
	case ParenLeft:
		return "ParenLeft"
	case ParenRight:
		return "ParenRight"
	case BracketLeft:
		return "BracketLeft"
	case BracketRight:
		return "BracketRight"
	case ActionArg:
		return "ActionArg"
	case Assign:
		return "Assign"
	case Sequence:
		return "Sequence"
	case EndMark:
		return "EndMark"
	case EndOfRule:
		return "EndOfRule"
	case Not:
		return "Not"
	default:
		return "Unknown"
	}
}
