package wad

type nameBuilder struct {
	names   map[int32]string
	nameBuf []byte
}

func NewNameBuilder() *nameBuilder {
	return &nameBuilder{
		names:   make(map[int32]string),
		nameBuf: []byte{},
	}
}

// Build builds a name builder
func (n *nameBuilder) Build(nameData string) {
	chunk := []rune{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			n.names[int32(lastOffset)] = string(chunk)
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}
}

// Set sets the names within a buffer
func (n *nameBuilder) Set(newNames map[int32]string) {
	if newNames == nil {
		n.names = make(map[int32]string)
		return
	}
	for k, v := range newNames {
		n.names[k] = v
	}
}

// Name returns the Name of an id
func (n *nameBuilder) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if n.names == nil {
		return "!UNK"
	}
	name := n.names[id]

	return name
}

// Add is used when building a world file, appending new names
func (n *nameBuilder) Add(name string) int32 {
	if n.names == nil {
		n.names = make(map[int32]string)
	}
	if name == "" {
		return 0
	}

	if id := n.Index(name); id != -1 {
		return id
	}
	n.names[int32(len(n.nameBuf))] = name
	n.nameBuf = append(n.nameBuf, []byte(name)...)
	n.nameBuf = append(n.nameBuf, 0)
	return int32(len(n.nameBuf) - len(name) - 1)
}

// Index returns the index of a name, or -1 if not found
func (n *nameBuilder) Index(name string) int32 {
	if n.names == nil {
		return -1
	}
	for k, v := range n.names {
		if v == name {
			return k
		}
	}
	return -1
}

// Data returns the name buffer
func (n *nameBuilder) Data() []byte {
	return n.nameBuf
}

// Clear clears the name buffer
func (n *nameBuilder) Clear() {
	n.names = make(map[int32]string)
	n.nameBuf = []byte{}
}
