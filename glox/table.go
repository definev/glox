package glox

const TABLE_MAX_LOAD = 0.75

type Table struct {
	Count    int
	capacity int
	Entries  []Entry
}

type Entry struct {
	Key   *ObjString
	Value Value
}

func NewTable() *Table {
	return &Table{
		Count:    0,
		capacity: 0,
		Entries:  make([]Entry, 0),
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

func (t *Table) Set(key *ObjString, value Value) bool {
	if float64(t.Count+1) > float64(t.capacity)*TABLE_MAX_LOAD {
		capacity := GROW_CAPACITY(t.capacity)
		t.adjustEntries(capacity)
	}

	entry := t.findEntry(t.Entries, key, t.capacity)

	isNewEntry := entry.Key == nil
	if isNewEntry && entry.Value.IsNil() {
		t.Count += 1
	}

	entry.Key = key
	entry.Value = value

	return isNewEntry
}

func (t *Table) SetAll(other *Table) {
	for index := 0; index < other.capacity; index++ {
		entry := other.Entries[index]

		if entry.Key != nil {
			t.Set(entry.Key, entry.Value)
		}
	}
}

func (t *Table) Get(key *ObjString) Value {
	if t.capacity == 0 {
		return NewNilVal()
	}

	entry := t.findEntry(t.Entries, key, t.capacity)
	if entry.Key == nil {
		return NewNilVal()
	}

	return entry.Value
}

func (t *Table) Delete(key *ObjString) bool {
	if t.capacity == 0 {
		return false
	}

	entry := t.findEntry(t.Entries, key, t.capacity)
	if entry.Key == nil {
		return false
	}

	entry.Key = nil
	entry.Value = NewBoolVal(true)

	return true
}

func (t *Table) adjustEntries(capacity int) {
	entries := make([]Entry, capacity)
	for index := 0; index < capacity; index++ {
		entries[index].Key = nil
		entries[index].Value = NewNilVal()
	}

	t.Count = 0

	for index := 0; index < t.capacity; index++ {
		entry := t.Entries[index]
		if entry.Key == nil {
			continue
		}

		dest := t.findEntry(entries, entry.Key, capacity)
		dest.Key = entry.Key
		dest.Value = entry.Value
		t.Count += 1
	}

	t.Entries = entries
	t.capacity = capacity
}

func (t *Table) findEntry(entries []Entry, key *ObjString, capacity int) *Entry {
	index := key.Hash % uint32(capacity)

	var tombstone *Entry = nil

	for {
		entry := entries[index]
		if entry.Key == nil {
			if entry.Value.IsNil() {
				// Treat tombstone as empty slot
				if tombstone != nil {
					return tombstone
				}
				return &entry
			} else {
				if tombstone == nil {
					tombstone = &entry
				}
			}
		}
		if entry.Key.IsEqual(*key) {
			return &entry
		}

		index = (index + 1) % uint32(capacity)
	}
}
