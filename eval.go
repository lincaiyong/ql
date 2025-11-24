package ql

import (
	"fmt"
	"github.com/lincaiyong/ql/parser"
)

func eval[T any](variables map[string]*Table[T], where string, select_ string) ([]*Record[T], error) {
	node, err := parse(where)
	if err != nil {
		return nil, err
	}
	v := Visitor[T]{
		input: variables,
	}
	v.visit(node)
	return nil, nil
}

func parse(expr string) (*parser.Node, error) {
	tokens, err := parser.Tokenize(expr)
	if err != nil {
		return nil, fmt.Errorf("fail to tokenize: %w", err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("fail to parse: %w", err)
	}
	return node, nil
}

type Visitor[T any] struct {
	input       map[string]*Table[T]
	constraints map[string]*Table[T]
}

func (v *Visitor[T]) visit(node *parser.Node) {

}
