package ql

func NewTable[T any](name string, db *Database[T], fields []string) *Table[T] {
	fieldMap := make(map[string]int, len(fields))
	for i, field := range fields {
		fieldMap[field] = i
	}
	return &Table[T]{
		name:      name,
		db:        db,
		fields:    fields,
		fieldMap:  fieldMap,
		records:   make([]*Record[T], 0),
		recordMap: make(map[int]*Record[T]),
	}
}

type Table[T any] struct {
	name      string
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
