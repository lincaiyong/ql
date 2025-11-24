package ql

func NewRecord[T any](table *Table[T], id int, data []string) *Record[T] {
	si := make([]int, len(data))
	for i, d := range data {
		si[i] = table.db.StoreString(d)
	}
	return &Record[T]{
		table: table,
		id:    id,
		data:  si,
	}
}

type Record[T any] struct {
	table *Table[T]
	id    int
	data  []int
}

func (r *Record[T]) Entity() *T {
	return r.table.db.entities[r.id]
}
