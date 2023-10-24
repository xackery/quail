package common

type Wld struct {
	Header     *Header                `yaml:"header"`
	IsOldWorld bool                   `yaml:"is_old_world"`
	Fragments  map[int]FragmentReader `yaml:"-"`
	names      map[int32]string       `yaml:"-"`
	Models     []*Model               `yaml:"models"`
}

func NewWld(name string) *Wld {
	e := &Wld{
		Header: &Header{
			Name: name,
		},
	}
	return e
}

func (e *Wld) SetNames(names map[int32]string) {
	if e.names == nil {
		e.names = make(map[int32]string)
	}
	for k, v := range names {
		e.names[k] = v
	}
}

func (e *Wld) Name(id int32) string {
	if e.names == nil {
		return "!UNK"
	}
	return e.names[id]
}
