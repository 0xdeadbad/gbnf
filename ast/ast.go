package ast

type AST struct {
	Root []Node
}

type Node interface {
	String() string
}
