package raw

import "fmt"

// eqgName stores and manages names structures for eqg related files
type eqgName struct {
	lastIndex int
	names     []*nameEntry
	nameBuf   []byte
}

type nameEntry struct {
	name   string
	offset int
}

func (e *nameEntry) String() string {
	return fmt.Sprintf("%s(%d)", e.name, e.offset)
}

func (e *eqgName) parse(nameData []byte) error {
	names := make(map[int32]string)
	chunk := []byte{}
	lastOffset := 0
	//lastElement := ""
	for i, b := range nameData {
		if b == 0 {
			names[int32(lastOffset)] = string(chunk)
			//	lastElement = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	e.set(names)
	return nil
}

func (e *eqgName) byOffset(id int32) string {
	if id < 0 {
		id = -id
	}
	if e.names == nil {
		return fmt.Sprintf("!UNK(%d)", id)
	}
	//fmt.Println("name: [", names[id], "]")

	for _, v := range e.names {
		if int32(v.offset) == id {
			return v.name
		}
	}
	return fmt.Sprintf("!UNK(%d)", id)
}

func (e *eqgName) offsetByName(name string) int32 {
	if e.names == nil {
		return -1
	}
	for _, v := range e.names {
		if v.name == name {
			return int32(v.offset)
		}
	}
	return -1
}

// indexByName is used when reading, returns the index of a name, or -1 if not found
func (e *eqgName) indexByName(name string) int32 {
	if e.names == nil {
		return -1
	}
	for k, v := range e.names {
		if v.name != name {
			continue
		}
		return int32(k)
	}
	return -1
}

// nameSet is used during reading, sets the names within a buffer
func (e *eqgName) set(newNames map[int32]string) {
	if newNames == nil {
		e.names = []*nameEntry{}
		return
	}
	for k, v := range newNames {
		e.names = append(e.names, &nameEntry{offset: int(k), name: v})
	}
	e.nameBuf = []byte{0x00}

	for _, v := range e.names {
		e.nameBuf = append(e.nameBuf, []byte(v.name)...)
		e.nameBuf = append(e.nameBuf, 0)
	}
}

// NameAdd is used when writing, appending new names
func (e *eqgName) add(name string) int32 {

	if e.names == nil {
		e.names = []*nameEntry{}
		//			{offset: 0, name: ""},
		//		}
		//		e.nameBuf = []byte{0x00}
	}
	if name == "" {
		return 0
	}

	/* if name[len(e.name)-1:] != "\x00" {
		name += "\x00"
	}
	*/
	// if id := e.NameOffset(name); id != -1 {
	// 	return -id
	// }
	e.names = append(e.names, &nameEntry{offset: len(e.nameBuf), name: name})
	lastRef := int32(len(e.nameBuf))
	e.nameBuf = append(e.nameBuf, []byte(name)...)
	e.nameBuf = append(e.nameBuf, 0)
	return int32(-lastRef)
}

// NameData is used during writing, dumps the name cache
func (e *eqgName) data() []byte {
	if len(e.nameBuf) == 0 {
		return nil
	}
	return e.nameBuf
}

// NameClear purges names and namebuf, called when encode starts
func (e *eqgName) clear() {
	e.names = nil
	e.nameBuf = nil
}

// func (e *eqgName) slice() []string {
// 	if e.names == nil {
// 		return nil
// 	}
// 	names := []string{}
// 	for _, v := range e.names {
// 		names = append(names, v.name)
// 	}
// 	return names
// }

func (e *eqgName) len() int {
	if e.names == nil {
		return 0
	}
	return len(e.names)
}
