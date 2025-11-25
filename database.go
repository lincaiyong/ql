package ql

import (
	"fmt"
	"regexp"
	"strings"
)

func NewDatabase[T any](entities []*T) *Database[T] {
	db := &Database[T]{
		entities: entities,
		strs:     make([]string, 0),
		strMap:   make(map[string]int),
		tables:   make([]*Table[T], 0),
		tableMap: make(map[string]*Table[T]),
	}
	table := NewTable[T]("Entity", db, nil, nil)
	for i := range entities {
		table.AddRecord(NewRecord[T](table, i, nil))
	}
	db.tableMap["Entity"] = table
	db.tables = append(db.tables, table)
	return db
}

type Database[T any] struct {
	entities []*T
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

func (db *Database[T]) GetBaseTable() *Table[T] {
	return db.tableMap["Entity"]
}

func (db *Database[T]) GetTable(name string) *Table[T] {
	return db.tableMap[name]
}

func (db *Database[T]) AddTable(baseTableName, tableName string, fields []string, fn func(t *T) []string) (*Table[T], error) {
	baseTable, ok := db.tableMap[baseTableName]
	if !ok {
		return nil, fmt.Errorf("table %s not found", baseTableName)
	}
	table := NewTable[T](tableName, db, fields, baseTable.getters)
	db.tableMap[tableName] = table
	db.tables = append(db.tables, table)
	for _, record := range baseTable.records {
		if ret := fn(db.entities[record.id]); ret != nil {
			r := NewRecord[T](table, record.id, ret)
			table.AddRecord(r)
		}
	}
	return table, nil
}

func (db *Database[T]) Query(q string) ([]*T, error) {
	ret := regexp.MustCompile(`^select (.+?)(?:where (.+))?$`).FindStringSubmatch(q)
	if len(ret) != 3 {
		return nil, fmt.Errorf("invalid query statement: %s", q)
	}
	select_ := ret[1]
	where := ret[2]
	s := strings.Fields(select_)
	if len(s) != 2 {
		return nil, fmt.Errorf("invalid query statement: %s", q)
	}
	table := db.GetTable(s[0])
	if table == nil {
		return nil, fmt.Errorf("table %s not found", s[0])
	}
	varName := s[1]
	records, err := eval[T](table, varName, where)
	if err != nil {
		return nil, fmt.Errorf("fail to eval: %w", err)
	}
	result := make([]*T, 0)
	for _, r := range records {
		result = append(result, db.entities[r.id])
	}
	return result, nil
}
