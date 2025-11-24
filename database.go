package ql

import (
	"fmt"
	"regexp"
	"strings"
)

func NewDatabase[T any](data []T) *Database[T] {
	db := &Database[T]{
		data:     data,
		strs:     make([]string, 0),
		strMap:   make(map[string]int),
		tables:   make([]*Table[T], 0),
		tableMap: make(map[string]*Table[T]),
	}
	table := NewTable[T]("", db, nil)
	for i := range data {
		table.AddRecord(NewRecord[T](table, i, nil))
	}
	db.tableMap[""] = table
	db.tables = append(db.tables, table)
	return db
}

type Database[T any] struct {
	data     []T
	strs     []string
	strMap   map[string]int
	tables   []*Table[T]
	tableMap map[string]*Table[T]
}

func (db *Database[T]) GetString(i int) string {
	return db.strs[i]
}

func (db *Database[T]) StoreString(s string) int {
	if ret, ok := db.strMap[s]; ok {
		return ret
	} else {
		ret = len(db.strs)
		db.strs = append(db.strs, s)
		db.strMap[s] = ret
		return ret
	}
}

func (db *Database[T]) GetTable(name string) *Table[T] {
	return db.tableMap[name]
}

func (db *Database[T]) AddTable(baseTableName, tableName string, fields []string, fn func(t T) []string) error {
	baseTable, ok := db.tableMap[baseTableName]
	if !ok {
		return fmt.Errorf("table %s not found", baseTableName)
	}
	table := NewTable[T](tableName, db, fields)
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
	ret := regexp.MustCompile(`^from (.+) (?:where (.+))? select (.+)$`).FindStringSubmatch(q)
	if len(ret) != 4 {
		return nil, fmt.Errorf("invalid query statement: %s", q)
	}
	from := ret[1]
	where := ret[2]
	select_ := ret[3]
	variables := make(map[string]*Table[T])
	for _, d := range strings.Split(from, ",") {
		s := strings.Fields(d)
		if len(s) != 2 {
			return nil, fmt.Errorf("invalid query statement: %s", q)
		}
		table := db.GetTable(s[0])
		if table == nil {
			return nil, fmt.Errorf("table %s not found", s[0])
		}
		variables[s[1]] = table
	}
	records, err := eval[T](variables, where, select_)
	if err != nil {
		return nil, fmt.Errorf("fail to eval: %w", err)
	}
	result := make([]T, 0)
	for _, r := range records {
		result = append(result, db.data[r.this])
	}
	return result, nil
}
