package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw/rawfrag"
)

type Wld struct {
	MetaFileName string
	Version      uint32
	IsNewWorld   bool
	IsZone       bool
	Fragments    []model.FragmentReadWriter
	Unk2         uint32
	Unk3         uint32
	names        []*nameEntry
	nameBuf      []byte
}

func (wld *Wld) Identity() string {
	return "wld"
}

// Read reads a wld file that was prepped by Load
func (wld *Wld) Read(r io.ReadSeeker) error {
	if wld.Fragments == nil {
		wld.Fragments = []model.FragmentReadWriter{&rawfrag.WldFragDefault{}}
	}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		return fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	wld.Version = dec.Uint32()

	wld.IsNewWorld = false
	switch wld.Version {
	case 0x00015500:
		wld.IsNewWorld = false
	case 0x1000C800:
		wld.IsNewWorld = true
	default:
		return fmt.Errorf("unknown wld version %d", wld.Version)
	}

	fragmentCount := dec.Uint32()

	bspRegionCount := dec.Uint32() //bspRegionCount
	maxFragSize := dec.Uint32()    // max fragment size
	hashSize := dec.Uint32()
	stringCount := dec.Uint32() // string count
	hashRaw := dec.Bytes(int(hashSize))
	nameData := helper.ReadStringHash(hashRaw)

	wld.names = []*nameEntry{}
	chunk := []rune{}
	lastOffset := 0
	//nameBuf = []byte{}
	for i, b := range nameData {
		if b == 0 {
			wld.names = append(wld.names, &nameEntry{name: string(chunk), offset: lastOffset})
			//nameBuf = append(wld.nameBuf, []byte(string(chunk))...)
			//nameBuf = append(wld.nameBuf, 0)
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		if i == len(nameData)-1 {
			break // some times there's garbage at the end
		}
		chunk = append(chunk, b)
	}

	if len(wld.names) != int(stringCount)+1 {
		fmt.Printf("name count mismatch, wanted %d, got %d (ignoring, openzone?)\n", stringCount, len(wld.names))
		//return fmt.Errorf("name count mismatch, wanted %d, got %d", stringCount, len(wld.names))
	}

	wld.nameBuf = hashRaw

	fragments, err := readFragments(fragmentCount, r)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	totalRegions := 0
	for i := uint32(0); i < fragmentCount; i++ {
		data := fragments[i]
		if len(data) > int(maxFragSize+4) {
			return fmt.Errorf("fragment %d (size: %d) exceeds max size %d", i, len(data), maxFragSize)
		}
		r := bytes.NewReader(data)

		reader := NewFrag(r)
		if reader == nil {
			return fmt.Errorf("unknown fragment at offset %d", i)
		}

		err = reader.Read(r, wld.IsNewWorld)
		if err != nil {
			return fmt.Errorf("frag %d 0x%x (%s) read: %w", i, reader.FragCode(), FragName(int(reader.FragCode())), err)
		}
		wld.Fragments = append(wld.Fragments, reader)

		pos, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("fragment %d (size: %d) seek: %w", i, len(data), err)
		}
		if pos != int64(len(data)) {
			isNonZero := false
			for i, bdata := range data {
				if int64(i) > pos {
					continue
				}
				if bdata != 0 {
					isNonZero = true
					break
				}
			}
			if !isNonZero {
				fmt.Printf("fragment %d seek mismatch (%d/%d) (%T)\n", i, pos, len(data), reader)

				fmt.Printf("fragment %d data: %x\n", i, data[pos:])
			}
		}
		_, ok := reader.(*rawfrag.WldFragRegion)
		if ok {
			totalRegions++
		}

	}

	if totalRegions != int(bspRegionCount) {
		return fmt.Errorf("region count mismatch, wanted %d, got %d", bspRegionCount, totalRegions)
	}
	return nil
}

// rawFrags is user by tests to compare for writer
func (wld *Wld) rawFrags(r io.ReadSeeker) ([][]byte, error) {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		return nil, fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	wld.Version = dec.Uint32()

	wld.IsNewWorld = true
	switch wld.Version {
	case 0x00015500:
		wld.IsNewWorld = false
	case 0x1000C800:
		wld.IsNewWorld = true
	default:
		return nil, fmt.Errorf("unknown wld version %d", wld.Version)
	}

	fragmentCount := dec.Uint32()

	_ = dec.Uint32() //bspRegionCount
	_ = dec.Uint32() // max_fragment_size
	hashSize := dec.Uint32()
	_ = dec.Uint32() // string count
	hashRaw := dec.Bytes(int(hashSize))
	nameData := helper.ReadStringHash(hashRaw)

	wld.names = []*nameEntry{}
	chunk := []rune{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			wld.names = append(wld.names, &nameEntry{name: string(chunk), offset: lastOffset})
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		if i == len(nameData)-1 {
			break // some times there's garbage at the end
		}
		chunk = append(chunk, b)
	}

	wld.nameBuf = hashRaw

	fragments, err := readFragments(fragmentCount, r)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	return fragments, nil
}

// readFragments convert frag data to structs
func readFragments(fragmentCount uint32, r io.ReadSeeker) (fragments [][]byte, err error) {

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(fragmentCount); fragOffset++ {
		fragSize := dec.Uint32()
		totalFragSize += fragSize

		fragCode := dec.Bytes(4)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("frag position seek %d/%d: %w", fragOffset, fragmentCount, err)
		}

		data := make([]byte, fragSize)
		_, err = r.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read frag %d/%d: %w", fragOffset, fragmentCount, err)
		}

		data = append(fragCode, data...)

		fragments = append(fragments, data)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek end of frag %d/%d: %w", fragOffset, fragmentCount, err)
		}
	}

	if dec.Error() != nil {
		return nil, fmt.Errorf("read: %w", dec.Error())
	}
	return fragments, nil
}

// SetFileName sets the name of the file
func (wld *Wld) SetFileName(name string) {
	wld.MetaFileName = name
}

// FileName returns the name of the file
func (wld *Wld) FileName() string {
	return wld.MetaFileName
}

// Name is used during reading, returns the Name of an id
func (wld *Wld) Name(id int32) string {
	if id < 0 {
		id = -id
	}
	if wld.names == nil {
		return fmt.Sprintf("!UNK(%d)", id)
	}
	//fmt.Println("name: [", names[id], "]")

	for _, v := range wld.names {
		if int32(v.offset) == id {
			return v.name
		}
	}
	return fmt.Sprintf("!UNK(%d)", id)
}

// NameSet is used during reading, sets the names within a buffer
func (wld *Wld) NameSet(newNames map[int32]string) {
	if newNames == nil {
		wld.names = []*nameEntry{}
		return
	}
	for k, v := range newNames {
		wld.names = append(wld.names, &nameEntry{offset: int(k), name: v})
	}
	wld.nameBuf = []byte{0x00}

	for _, v := range wld.names {
		wld.nameBuf = append(wld.nameBuf, []byte(v.name)...)
		wld.nameBuf = append(wld.nameBuf, 0)
	}
}

// NameAdd is used when writing, appending new names
func (wld *Wld) NameAdd(name string) int32 {

	if wld.names == nil {
		wld.names = []*nameEntry{
			{offset: 0, name: ""},
		}
		wld.nameBuf = []byte{0x00}
	}
	if name == "" {
		return 0
	}

	/* if name[len(wld.name)-1:] != "\x00" {
		name += "\x00"
	}
	*/
	if id := wld.NameOffset(name); id != -1 {
		return -id
	}
	wld.names = append(wld.names, &nameEntry{offset: len(wld.nameBuf), name: name})
	lastRef := int32(len(wld.nameBuf))
	wld.nameBuf = append(wld.nameBuf, []byte(name)...)
	wld.nameBuf = append(wld.nameBuf, 0)
	return int32(-lastRef)
}

func (wld *Wld) NameOffset(name string) int32 {
	if wld.names == nil {
		return -1
	}
	for _, v := range wld.names {
		if v.name == name {
			return int32(v.offset)
		}
	}
	return -1
}

// NameIndex is used when reading, returns the index of a name, or -1 if not found
func (wld *Wld) NameIndex(name string) int32 {
	if wld.names == nil {
		return -1
	}
	for k, v := range wld.names {
		if v.name == name {
			return int32(k)
		}
	}
	return -1
}

// NameData is used during writing, dumps the name cache
func (wld *Wld) NameData() []byte {

	return helper.WriteStringHash(string(wld.nameBuf))
}

// NameClear purges names and namebuf, called when encode starts
func (wld *Wld) NameClear() {
	wld.names = nil
	wld.nameBuf = nil
}

func (wld *Wld) TagByFrag(srcFrag interface{}) string {
	switch frag := srcFrag.(type) {
	case *rawfrag.WldFragActorDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragActor:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragBlitSpriteDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragBlitSprite:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragBMInfo:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragCompositeSpriteDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragCompositeSprite:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragDmRGBTrack:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragDmRGBTrackDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragDmSpriteDef2:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragDMSpriteDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragDmTrackDef2:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragLight:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragHierarchicalSpriteDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragLightDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragMaterialDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragMaterialPalette:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragParticleCloudDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragParticleSpriteDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragParticleSprite:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragPointLight:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragPolyhedron:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragPolyhedronDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragRegion:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSimpleSpriteDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSimpleSprite:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSoundDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSound:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSphereListDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSphereList:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSphere:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSprite2D:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSprite2DDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSprite3D:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSprite3DDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSprite4DDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragSprite4D:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragTrack:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragTrackDef:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragWorldTree:
		return wld.Name(frag.NameRef)
	case *rawfrag.WldFragZone:
		return wld.Name(frag.NameRef)
	}

	return ""
}
