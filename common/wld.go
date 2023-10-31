package common

type Wld struct {
	Header     *Header                `yaml:"header"`
	IsOldWorld bool                   `yaml:"is_old_world"`
	Fragments  map[int]FragmentReader `yaml:"fragments,omitempty"`
	names      map[int32]string       `yaml:"-"`
	nameBuf    []byte                 `yaml:"-"`
	Models     []*Model               `yaml:"models,omitempty"`
}

func NewWld(name string) *Wld {
	e := &Wld{
		Header: &Header{
			Name: name,
		},
		Fragments: make(map[int]FragmentReader),
		names:     make(map[int32]string),
	}
	return e
}

// SetNames sets the names within a buffer
func (e *Wld) SetNames(names map[int32]string) {
	if e.names == nil {
		e.names = make(map[int32]string)
	}
	for k, v := range names {
		e.names[k] = v
	}
}

// Name returns the name of an id
func (e *Wld) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if e.names == nil {
		return "!UNK"
	}
	//fmt.Println("name: [", e.names[id], "]")
	return e.names[id]
}

// NameIndex returns the index of a name, or -1 if not found
func (e *Wld) NameIndex(name string) int32 {
	if e.names == nil {
		return -1
	}
	for k, v := range e.names {
		if v == name {
			return k
		}
	}
	return -1
}

// NameAdd is used when building a world file, appending new names
func (e *Wld) NameAdd(name string) int32 {
	if e.names == nil {
		e.names = make(map[int32]string)
	}

	if id := e.NameIndex(name); id != -1 {
		return id
	}
	e.names[int32(len(e.nameBuf))] = name
	e.nameBuf = append(e.nameBuf, []byte(name)...)
	e.nameBuf = append(e.nameBuf, 0)
	return int32(len(e.nameBuf) - len(name) - 1)
}
