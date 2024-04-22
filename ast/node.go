package ast

import (
	"fmt"
	"gbnf/lexer"
)

type ProdRule struct {
	Left  *lexer.Token
	Right []*lexer.Token
}

func (p *ProdRule) String() string {
	return fmt.Sprintf("%s ::= %s", p.Left, p.Right)
}

type Action struct {
	Action *lexer.Token
	Args   []*lexer.Token
}

func (a *Action) String() string {
	return fmt.Sprintf("%s(%s)", a.Action, a.Args)
}
