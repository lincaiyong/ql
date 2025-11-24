package ql

func NewTable[T any](name string, db *Database[T], fields []string, getters map[string]func(*T) *Value) *Table[T] {
	fieldMap := make(map[string]int, len(fields))
	for i, field := range fields {
		fieldMap[field] = i
	}
	if getters == nil {
		getters = make(map[string]func(*T) *Value)
	}
	return &Table[T]{
		name:      name,
		db:        db,
		fields:    fields,
		fieldMap:  fieldMap,
		records:   make([]*Record[T], 0),
		recordMap: make(map[int]*Record[T]),
		getters:   getters,
	}
}

type Table[T any] struct {
	name      string
	db        *Database[T]
	fields    []string
	fieldMap  map[string]int
	records   []*Record[T]
	recordMap map[int]*Record[T]
	getters   map[string]func(*T) *Value
}

func (t *Table[T]) AddRecord(r *Record[T]) {
	t.recordMap[r.id] = r
	t.records = append(t.records, r)
}

func (t *Table[T]) Records() []*Record[T] {
	return t.records
}

func (t *Table[T]) Getter(n string) func(*T) *Value {
	return t.getters[n]
}

func (t *Table[T]) Define(n string, getter func(*T) *Value) {
	t.getters[n] = getter
}
