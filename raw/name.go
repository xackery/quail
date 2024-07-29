package raw

import (
	"github.com/xackery/quail/helper"
)

var (
	names   = []*nameEntry{}
	nameBuf = []byte{}
)

type nameEntry struct {
	name   string
	offset int
}

// NameSet is used during reading, sets the names within a buffer
func NameSet(newNames map[int32]string) {
	if newNames == nil {
		names = []*nameEntry{}
		return
	}
	for k, v := range newNames {
		names = append(names, &nameEntry{offset: int(k), name: v})
	}
	nameBuf = []byte{0x00}

	for _, v := range names {
		nameBuf = append(nameBuf, []byte(v.name)...)
		nameBuf = append(nameBuf, 0)
	}
}

// Name is used during reading, returns the Name of an id
func Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if names == nil {
		return "!UNK"
	}
	//fmt.Println("name: [", names[id], "]")

	for _, v := range names {
		if int32(v.offset) == id {
			return v.name
		}
	}

	return "!UNK"
}

// NameAdd is used when writing, appending new names
func NameAdd(name string) int32 {
	if names == nil {
		names = []*nameEntry{
			{offset: 0, name: ""},
		}
		nameBuf = []byte{0x00}
	}
	/* if !strings.HasSuffix(name, "\000") {
		name += "\000"
	} */
	if name == "" {
		return 0
	}

	if id := NameIndex(name); id != -1 {
		return id
	}
	names = append(names, &nameEntry{offset: len(nameBuf), name: name})
	lastRef := int32(len(nameBuf))
	nameBuf = append(nameBuf, []byte(name)...)
	nameBuf = append(nameBuf, 0)
	return int32(-lastRef)
}

// NameIndex is used when reading, returns the index of a name, or -1 if not found
func NameIndex(name string) int32 {
	if names == nil {
		return -1
	}
	for k, v := range names {
		if v.name == name {
			return int32(k)
		}
	}
	return -1
}

// NameData is used during writing, dumps the name cache
func NameData() []byte {
	return helper.WriteStringHash(string(nameBuf))
}

// NameClear purges names and namebuf, called when encode starts
func NameClear() {
	names = nil
	nameBuf = nil
}
