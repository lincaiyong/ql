package main

import (
	"github.com/lincaiyong/log"
	"github.com/lincaiyong/ql"
	"strconv"
)

func main() {
	var numbers []string
	for i := 0; i < 50; i++ {
		numbers = append(numbers, strconv.Itoa(i))
	}
	db := ql.NewDatabase[string](numbers)
	err := db.AddTable("", "OddNumber", nil, func(t string) []string {
		if i, _ := strconv.Atoi(t); i%2 == 1 {
			return []string{}
		}
		return nil
	})
	if err != nil {
		log.ErrorLog("fail to add table: %v", err)
	}
	err = db.AddTable("", "EvenNumber", []string{"string"}, func(t string) []string {
		if i, _ := strconv.Atoi(t); i%2 == 0 {
			return []string{t + "_string"}
		}
		return nil
	})
	err = db.AddTable("", "DividableBy4", nil, func(t string) []string {
		if i, _ := strconv.Atoi(t); i%2 == 0 {
			return []string{}
		}
		return nil
	})
	ret, err := db.Query(`from EvenNumber n where n > 40 select n`)
	if err != nil {
		log.ErrorLog("fail to query: %v", err)
		return
	}
	for _, n := range ret {
		log.InfoLog("%s", n)
	}
	log.InfoLog("done")
}
