package raw

var (
	names   = make(map[int32]string)
	nameBuf = []byte{}
)

// NamesSet sets the names within a buffer
func NamesSet(newNames map[int32]string) {
	if newNames == nil {
		names = make(map[int32]string)
		return
	}
	for k, v := range newNames {
		names[k] = v
	}
}

// Name returns the Name of an id
func Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if names == nil {
		return "!UNK"
	}
	//fmt.Println("name: [", names[id], "]")
	name := names[id]

	return name
}

// NameAdd is used when building a world file, appending new names
func NameAdd(name string) int32 {
	if names == nil {
		names = make(map[int32]string)
	}

	if id := NameIndex(name); id != -1 {
		return id
	}
	names[int32(len(nameBuf))] = name
	nameBuf = append(nameBuf, []byte(name)...)
	nameBuf = append(nameBuf, 0)
	return int32(len(nameBuf) - len(name) - 1)
}

// NameIndex returns the index of a name, or -1 if not found
func NameIndex(name string) int32 {
	if names == nil {
		return -1
	}
	for k, v := range names {
		if v == name {
			return k
		}
	}
	return -1
}

// NameData dumps the name cache
func NameData() []byte {
	//os.WriteFile("dst.txt", []byte(fmt.Sprintf("%+v", names)), 0644)
	return nameBuf
}

// NameClear purges names and namebuf, called when encode starts
func NameClear() {
	names = make(map[int32]string)
	nameBuf = nil
}
