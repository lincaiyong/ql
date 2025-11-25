package main

import (
	"github.com/lincaiyong/log"
	"github.com/lincaiyong/ql"
)

func main() {
	type Entity struct {
		num int
	}
	var entities []*Entity
	for i := 0; i < 100; i++ {
		entities = append(entities, &Entity{num: i})
	}
	db := ql.NewDatabase[Entity](entities)
	tbl := db.GetTable("")
	tbl.Define("num", func(e *Entity) *ql.Value {
		return ql.NewIntValue(e.num)
	})
	_, err := db.AddTable("", "OddNumber", nil, func(t *Entity) []string {
		if t.num%2 == 1 {
			return []string{}
		}
		return nil
	})
	if err != nil {
		log.ErrorLog("fail to add table: %v", err)
	}
	_, err = db.AddTable("", "EvenNumber", []string{"string"}, func(t *Entity) []string {
		if t.num%2 == 0 {
			return []string{}
		}
		return nil
	})
	_, err = db.AddTable("", "DividableBy4", nil, func(t *Entity) []string {
		if t.num%2 == 0 {
			return []string{}
		}
		return nil
	})
	ret, err := db.Query(`select EvenNumber n where n.num > 40 and n.num <= 50`)
	if err != nil {
		log.ErrorLog("fail to query: %v", err)
		return
	}
	for _, n := range ret {
		log.InfoLog("%d", n.num)
	}
	log.InfoLog("done")
}
