package main

import (
	"github.com/lincaiyong/log"
	"github.com/lincaiyong/ql"
	"strconv"
)

type Entity struct {
	num string
}

func main() {
	var entities []*Entity
	for i := 0; i < 50; i++ {
		entities = append(entities, &Entity{num: strconv.Itoa(i)})
	}
	db := ql.NewDatabase[*Entity](entities)
	err := db.AddTable("", "OddNumber", nil, func(t *Entity) []string {
		if i, _ := strconv.Atoi(t.num); i%2 == 1 {
			return []string{}
		}
		return nil
	})
	if err != nil {
		log.ErrorLog("fail to add table: %v", err)
	}
	err = db.AddTable("", "EvenNumber", []string{"string"}, func(t *Entity) []string {
		if i, _ := strconv.Atoi(t.num); i%2 == 0 {
			return []string{t.num + "_string"}
		}
		return nil
	})
	err = db.AddTable("", "DividableBy4", nil, func(t *Entity) []string {
		if i, _ := strconv.Atoi(t.num); i%2 == 0 {
			return []string{}
		}
		return nil
	})
	ret, err := db.Query(`from EvenNumber n, DividableBy4 d where n.num > 40 and n in d select n`)
	if err != nil {
		log.ErrorLog("fail to query: %v", err)
		return
	}
	for _, n := range ret {
		log.InfoLog("%s", n)
	}
	log.InfoLog("done")
}
