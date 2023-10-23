package common

type Wld struct {
	Version    int
	IsOldWorld bool
	Fragments  map[int]FragmentReader
	names      map[int32]string
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
