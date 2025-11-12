package ql

func NewRecord[T any](table *Table[T], this int, data []string) *Record[T] {
	si := make([]int, len(data))
	for i, d := range data {
		si[i] = table.db.StoreString(d)
	}
	return &Record[T]{
		table: table,
		this:  this,
		data:  si,
	}
}

type Record[T any] struct {
	table *Table[T]
	this  int
	data  []int
}
