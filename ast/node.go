package ast

import "gbnf/lexer"

type Assignment struct {
	Left  *lexer.Token
	Right []*lexer.Token
}
