package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs/archive"
)

// Zone information
type Zone struct {
	name    string
	indices []int32
}

func LoadZone(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &Zone{}
	err := parseZone(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse Zone: %w", err)
	}
	return v, nil
}

// based on https://github.com/danwilkins/LanternExtractor/blob/development/0.2.0/LanternExtractor/EQ/Wld/Fragments/BspZone.cs
func parseZone(r io.ReadSeeker, v *Zone) error {
	var err error

	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	flags := int32(0)
	err = binary.Read(r, binary.LittleEndian, &flags)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}
	//dump.Hex(flags, "flags=%d", flags)

	regionCount := int32(0)
	err = binary.Read(r, binary.LittleEndian, &regionCount)
	if err != nil {
		return fmt.Errorf("read regionCount: %w", err)
	}
	//dump.Hex(regionCount, "regionCount=%d", regionCount)

	for i := 0; i < int(regionCount); i++ {
		index := int32(0)
		err = binary.Read(r, binary.LittleEndian, &index)
		if err != nil {
			return fmt.Errorf("read index: %w", err)
		}
		//dump.Hex(index, "index=%d", index)
		v.indices = append(v.indices, index)
	}

	ZoneString, err := helper.ReadString(r)
	if err != nil {
		return fmt.Errorf("read ZoneString: %w", err)
	}

	Zones := []int{} //0: water, 1: lava, 2: zone line
	if strings.HasPrefix(ZoneString, "wtn_") ||
		strings.HasPrefix(ZoneString, "wt_") {
		Zones = append(Zones, 0)
	}

	if strings.HasPrefix(ZoneString, "wtntp") {
		Zones = append(Zones, 0)
		Zones = append(Zones, 1)
	}

	return nil
}

func (v *Zone) FragmentType() string {
	return "Zone"
}

func (e *Zone) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
