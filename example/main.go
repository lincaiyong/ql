package main

import (
	"fmt"
	"github.com/lincaiyong/log"
	"regexp"
	"strconv"
	"strings"
)

func NewDatabase[T any](data []T) *Database[T] {
	db := &Database[T]{
		data:     data,
		tables:   make([]*Table[T], 0),
		tableMap: make(map[string]*Table[T]),
	}
	table := NewTable[T](db, nil)
	for i := range data {
		table.AddRecord(NewRecord[T](table, i, nil))
	}
	db.tableMap[""] = table
	db.tables = append(db.tables, table)
	return db
}

type Database[T any] struct {
	data     []T
	tables   []*Table[T]
	tableMap map[string]*Table[T]
}

func (db *Database[T]) GetTable(name string) *Table[T] {
	return db.tableMap[name]
}

func (db *Database[T]) AddTable(baseTableName, tableName string, fields []string, fn func(t T) []string) error {
	baseTable, ok := db.tableMap[baseTableName]
	if !ok {
		return fmt.Errorf("table %s not found", baseTableName)
	}
	table := NewTable[T](db, fields)
	db.tableMap[tableName] = table
	db.tables = append(db.tables, table)
	for _, record := range baseTable.records {
		if ret := fn(db.data[record.this]); ret != nil {
			r := NewRecord[T](table, record.this, ret)
			table.AddRecord(r)
		}
	}
	return nil
}

func (db *Database[T]) Query(q string) ([]T, error) {
	ret := regexp.MustCompile(`^from (.+) select (.+)$`).FindStringSubmatch(q)
	if len(ret) != 3 {
		return nil, fmt.Errorf("invalid query statement: %s", q)
	}
	from := ret[1]
	select_ := ret[2]
	fromTypeNames := make([][2]string, 0)
	for _, d := range strings.Split(from, ",") {
		s := strings.Fields(d)
		if len(s) != 2 {
			return nil, fmt.Errorf("invalid query statement: %s", q)
		}
		fromTypeNames = append(fromTypeNames, [2]string{s[0], s[1]})
	}
	log.InfoLog("%v %s", fromTypeNames, select_)
	tableName := fromTypeNames[0][0]
	table := db.GetTable(tableName)
	if table == nil {
		return nil, fmt.Errorf("table %s not found", select_)
	}
	result := make([]T, 0)
	for _, r := range table.records {
		result = append(result, db.data[r.this])
	}
	return result, nil
}

func NewTable[T any](db *Database[T], fields []string) *Table[T] {
	fieldMap := make(map[string]int, len(fields))
	for i, field := range fields {
		fieldMap[field] = i
	}
	return &Table[T]{
		db:        db,
		fields:    fields,
		fieldMap:  fieldMap,
		records:   make([]*Record[T], 0),
		recordMap: make(map[int]*Record[T]),
	}
}

type Table[T any] struct {
	db        *Database[T]
	fields    []string
	fieldMap  map[string]int
	records   []*Record[T]
	recordMap map[int]*Record[T]
}

func (t *Table[T]) AddRecord(r *Record[T]) {
	t.recordMap[r.this] = r
	t.records = append(t.records, r)
}

func NewRecord[T any](table *Table[T], this int, data []string) *Record[T] {
	return &Record[T]{
		table: table,
		this:  this,
		data:  data,
	}
}

type Record[T any] struct {
	table *Table[T]
	this  int
	data  []string
}

func main() {
	// https://codeql.github.com/docs/codeql-language-guides/analyzing-data-flow-in-go/
	//	q := `from Function osOpen, CallExpr call
	//where
	//  osOpen.hasQualifiedName("os", "Open") and
	//  call.getTarget() = osOpen
	//select call.getArgument(0)
	//`
	//	ret := ql.Query(q)
	//	fmt.Println(ret)
	var numbers []string
	for i := 0; i < 50; i++ {
		numbers = append(numbers, strconv.Itoa(i))
	}
	db := NewDatabase[string](numbers)
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
	ret, err := db.Query(`from EvenNumber n select n`)
	if err != nil {
		log.ErrorLog("fail to query: %v", err)
		return
	}
	for _, n := range ret {
		log.InfoLog("%s", n)
	}
	log.InfoLog("done")
}
