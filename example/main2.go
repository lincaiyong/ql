package main

import (
	"github.com/lincaiyong/log"
	"github.com/lincaiyong/ql"
	"github.com/lincaiyong/ql/parser"
)

func main() {
	var entities []*parser.Node
	expr := "a.op == '+' and a.lhs.code != 'x'"
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
	tbl := db.GetBaseTable()
	tbl.Define("op", func(e *Entity) *ql.Value {
		return ql.NewStringValue(e.Op())
	})
	tbl.Define("type", func(e *Entity) *ql.Value {
		return ql.NewStringValue(e.Type())
	})
	_, err = db.AddTable("Entity", "BinaryExpr", nil, func(t *Entity) []string {
		if t.Type() == parser.NodeTypeBinary {
			return []string{}
		}
		return nil
	})
	if err != nil {
		log.ErrorLog("fail to add table: %v", err)
		return
	}
	_, err = db.AddTable("Entity", "UnaryExpr", nil, func(t *Entity) []string {
		if t.Type() == parser.NodeTypeUnary {
			return []string{}
		}
		return nil
	})
	if err != nil {
		log.ErrorLog("fail to add table: %v", err)
		return
	}
	ret, err := db.Query(`select Entity n where n.type == 'binary' and n.op == '!='`)
	if err != nil {
		log.ErrorLog("fail to query: %v", err)
		return
	}
	for _, n := range ret {
		log.InfoLog("%s", n.Dump())
	}
	log.InfoLog("done")
}
