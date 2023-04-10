package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// DmRGBTrackDef information
type DmRGBTrackDef struct {
	// Colors of the vertex, if applicable
	Colors []color.RGBA
	name   string
}

func LoadDmRGBTrackDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &DmRGBTrackDef{}
	err := parseDmRGBTrackDef(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse DmRGBTrackDef: %w", err)
	}
	return v, nil
}

func parseDmRGBTrackDef(r io.ReadSeeker, v *DmRGBTrackDef) error {
	if v == nil {
		return fmt.Errorf("DmRGBTrackDef is nil")
	}
	var value uint32
	var err error

	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	//unknown
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown: %w", err)
	}

	var DmRGBTrackDefCount int32
	err = binary.Read(r, binary.LittleEndian, &DmRGBTrackDefCount)
	if err != nil {
		return fmt.Errorf("read vertex color count: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown 2: %w", err)
	}
	//if value != 1 {
	//	return fmt.Errorf("unknown 2 expected %d, got %d", 1, value)
	//}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown 3: %w", err)
	}
	//also got 70? with nexus
	//if value != 200 {
	//	return fmt.Errorf("unknown 3 expected %d, got %d", 200, value)
	//}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown 4: %w", err)
	}
	//also got 74? with nexus
	//if value != 0 {
	//	return fmt.Errorf("unknown 4 expected %d, got %d", 0, value)
	//}
	for i := 0; i < int(DmRGBTrackDefCount)/4; i++ {
		rgba := color.RGBA{}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read r: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read g: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read b: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read a: %w", err)
		}
		v.Colors = append(v.Colors, rgba)
	}

	return nil
}

func (v *DmRGBTrackDef) FragmentType() string {
	return "DmRGBTrackDef"
}

func (e *DmRGBTrackDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
