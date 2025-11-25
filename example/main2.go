package main

import (
	"github.com/lincaiyong/log"
	"github.com/lincaiyong/ql"
	"github.com/lincaiyong/ql/parser"
)

func main() {
	var entities []*parser.Node
	expr := "a.op == '+' and a.lhs.code == 'x'"
	tokens, err := parser.Tokenize(expr)
	if err != nil {
		log.ErrorLog("%v", err)
		return
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		log.ErrorLog("%v", err)
		return
	}
	node.Visit(func(node *parser.Node) {
		entities = append(entities, node)
	})
	type Entity = parser.Node
	db := ql.NewDatabase[Entity](entities)
	tbl := db.GetTable("")
	tbl.Define("op", func(e *Entity) *ql.Value {
		return ql.NewStringValue(e.Op())
	})
	_, err = db.AddTable("", "BinaryExpr", nil, func(t *Entity) []string {
		if t.Type() == parser.NodeTypeBinary {
			return []string{}
		}
		return nil
	})
	if err != nil {
		log.ErrorLog("fail to add table: %v", err)
		return
	}
	_, err = db.AddTable("", "UnaryExpr", nil, func(t *Entity) []string {
		if t.Type() == parser.NodeTypeUnary {
			return []string{}
		}
		return nil
	})
	if err != nil {
		log.ErrorLog("fail to add table: %v", err)
		return
	}
	ret, err := db.Query(`select BinaryExpr n`)
	if err != nil {
		log.ErrorLog("fail to query: %v", err)
		return
	}
	for _, n := range ret {
		log.InfoLog("%s", n.Op())
	}
	log.InfoLog("done")
}
