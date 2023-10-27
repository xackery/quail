package common

// Zone is a zone
type Zone struct {
	Header  *Header  `yaml:"header,omitempty"`
	Models  []string `yaml:"models,omitempty"`
	Objects []Object `yaml:"objects,omitempty"`
	Regions []Region `yaml:"regions,omitempty"`
	Lights  []Light  `yaml:"lights,omitempty"`
	Lits    []*RGBA  `yaml:"lits,omitempty"`
	V4Info  V4Info   `yaml:"v4info,omitempty"`
	names   map[int32]string
	nameBuf []byte
}

func NewZone(name string) *Zone {
	return &Zone{
		Header: &Header{
			Name: name,
		},
	}
}

type V4Info struct {
	MinLng               int
	MinLat               int
	MaxLng               int
	MaxLat               int
	MinExtents           Vector3
	MaxExtents           Vector3
	UnitsPerVert         float32
	QuadsPerTile         int
	CoverMapInputSize    int
	LayeringMapInputSize int
}

// Object is an object
type Object struct {
	Name      string
	ModelName string
	Position  Vector3
	Rotation  Vector3
	Scale     float32
}

// Region is a region
type Region struct {
	Name    string
	Center  Vector3
	Unknown Vector3
	Extent  Vector3
}

// Light is a light
type Light struct {
	Name     string
	Position Vector3
	Color    Vector3
	Radius   float32
}

// SetNames sets the names within a buffer
func (e *Zone) SetNames(names map[int32]string) {
	if e.names == nil {
		e.names = make(map[int32]string)
	}
	for k, v := range names {
		e.names[k] = v
	}
}

// Name returns the name of an id
func (e *Zone) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if e.names == nil {
		return "!UNK"
	}
	//fmt.Println("name: [", e.names[id], "]")
	return e.names[id]
}

// NameAdd is used when building a world file, appending new names
func (e *Zone) NameAdd(name string) int32 {
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

// NameIndex returns the index of a name, or -1 if not found
func (e *Zone) NameIndex(name string) int32 {
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

// NameData dumps the name cache
func (e *Zone) NameData() []byte {
	//os.WriteFile("dst.txt", []byte(fmt.Sprintf("%+v", e.names)), 0644)
	return e.nameBuf
}

// NameClear purges names and namebuf, called when encode starts
func (e *Zone) NameClear() {
	e.names = make(map[int32]string)
	e.nameBuf = nil
}
