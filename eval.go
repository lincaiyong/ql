package ql

import (
	"fmt"
	"github.com/lincaiyong/log"
	"github.com/lincaiyong/ql/parser"
	"strconv"
	"strings"
)

func eval[T any](table *Table[T], varName, where string) (result []*Record[T], err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	if where == "" {
		return table.Records(), nil
	}
	node, err := parse(where)
	if err != nil {
		return nil, err
	}
	v := Evaluator[T]{
		varName: varName,
		table:   table,
	}
	result = v.EvalSet(node, table.Records())
	return result, nil
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

type Evaluator[T any] struct {
	table   *Table[T]
	varName string
}

func (v *Evaluator[T]) EvalSet(node *parser.Node, all []*Record[T]) []*Record[T] {
	switch node.Type() {
	case parser.NodeTypeIdent:
		name := node.Ident()
		if name == v.varName {
			log.FatalLog("invalid identifier %s", name)
			return nil
		}
		return all
	case parser.NodeTypeUnary:
		s := v.EvalSet(node.UnaryTarget(), all)
		m := make(map[int]struct{})
		result := make([]*Record[T], 0, len(all)-len(s))
		for _, record := range s {
			if _, ok := m[record.id]; !ok {
				result = append(result, record)
			}
		}
		return result
	case parser.NodeTypeBinary:
		if node.Op() == "and" {
			lhs := v.EvalSet(node.BinaryLhs(), all)
			result := v.EvalSet(node.BinaryRhs(), lhs)
			return result
		} else if node.Op() == "or" {
			lhs := v.EvalSet(node.BinaryLhs(), all)
			rhs := v.EvalSet(node.BinaryRhs(), all)
			m := make(map[int]struct{})
			result := make([]*Record[T], 0, len(lhs)+len(rhs))
			for _, record := range lhs {
				m[record.id] = struct{}{}
				result = append(result, record)
			}
			for _, record := range rhs {
				if _, ok := m[record.id]; !ok {
					result = append(result, record)
				}
			}
			return result
		} else if node.Op() == ">" || node.Op() == "<" || node.Op() == ">=" || node.Op() == "<=" || node.Op() == "==" || node.Op() == "!=" {
			lhs := v.EvalValue(node.BinaryLhs())
			rhs := v.EvalValue(node.BinaryRhs())
			if lhs == nil {
				log.FatalLog("invalid lhs")
				return nil
			}
			if rhs == nil {
				log.FatalLog("invalid rhs")
				return nil
			}
			result := make([]*Record[T], 0, len(all))
			for _, record := range all {
				lhsValue := lhs(record.Entity())
				rhsValue := rhs(record.Entity())
				if v.compare(node.Op(), lhsValue, rhsValue) {
					result = append(result, record)
				}
			}
			return result
		} else {
			log.FatalLog("invalid operator %s", node.Op())
			return nil
		}
	default:
		log.FatalLog("invalid node type %s", node.Type())
		return nil
	}
}

func (v *Evaluator[T]) EvalValue(node *parser.Node) func(*T) *Value {
	if node.Type() == parser.NodeTypeSelector {
		if node.SelectorTarget().Type() == parser.NodeTypeIdent {
			n := node.SelectorTarget().Ident()
			if n != v.varName {
				log.FatalLog("invalid identifier %s", n)
				return nil
			}
			getter := v.table.Getter(node.SelectorKey())
			return getter
		} else {
			log.FatalLog("invalid selector target %s", node.SelectorTarget().Type())
			return nil
		}
	} else if node.Type() == parser.NodeTypeString {
		return func(entity *T) *Value {
			s := strings.Trim(node.String(), "'")
			s = strings.ReplaceAll(s, "\\'", "'")
			return NewStringValue(s)
		}
	} else if node.Type() == parser.NodeTypeNumber {
		return func(entity *T) *Value {
			i, _ := strconv.Atoi(node.String())
			return NewIntValue(i)
		}
	}
	log.FatalLog("invalid node type %s", node.Type())
	return nil
}

func (v *Evaluator[T]) compare(op string, lhs, rhs *Value) bool {
	if lhs.type_ != rhs.type_ {
		log.FatalLog("invalid lhs, rhs %s %s %s", lhs.type_, rhs.type_, op)
		return false
	}
	if lhs.type_ == ValueTypeInt {
		switch op {
		case ">":
			return lhs.IntValue() > rhs.IntValue()
		case "<":
			return lhs.IntValue() < rhs.IntValue()
		case ">=":
			return lhs.IntValue() >= rhs.IntValue()
		case "<=":
			return lhs.IntValue() <= rhs.IntValue()
		case "==":
			return lhs.IntValue() == rhs.IntValue()
		default:
			log.FatalLog("invalid op %s", op)
			return false
		}
	}
	if lhs.type_ == ValueTypeString {
		if op == "==" {
			return lhs.StringValue() == rhs.StringValue()
		} else if op == "!=" {
			return lhs.StringValue() != rhs.StringValue()
		}
	}
	log.FatalLog("invalid op %s", op)
	return false
}
