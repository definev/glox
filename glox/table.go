package glox

type Table struct {
	Count    int
	capacity int
	Entries  []*Entry
}

type Entry struct {
	Key   *ObjString
	Value Value
}

func NewTable() *Table {
	return &Table{
		Count:    0,
		capacity: 0,
		Entries:  make([]*Entry, 0),
	}
}

func (t *Table) Init() {
	t.Count = 0
	t.capacity = 0
	t.Entries = nil
}

func (t *Table) Free() {
	t.Init()
}

func (t *Table) Set(key *ObjString, value Value) {}
