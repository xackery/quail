package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// RegionType information
type RegionType struct {
	name    string
	indices []int32
}

func LoadRegionType(r io.ReadSeeker) (common.WldFragmenter, error) {
	v := &RegionType{}
	err := parseRegionType(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse RegionType: %w", err)
	}
	return v, nil
}

// based on https://github.com/danwilkins/LanternExtractor/blob/development/0.2.0/LanternExtractor/EQ/Wld/Fragments/BspRegionType.cs
func parseRegionType(r io.ReadSeeker, v *RegionType) error {
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

	regionTypeString, err := helper.ReadString(r)
	if err != nil {
		return fmt.Errorf("read regionTypeString: %w", err)
	}

	regionTypes := []int{} //0: water, 1: lava, 2: zone line
	if strings.HasPrefix(regionTypeString, "wtn_") ||
		strings.HasPrefix(regionTypeString, "wt_") {
		regionTypes = append(regionTypes, 0)
	}

	if strings.HasPrefix(regionTypeString, "wtntp") {
		regionTypes = append(regionTypes, 0)
		regionTypes = append(regionTypes, 1)
		//TODO: decode zoneline
	}
	if len(regionTypes) > 0 {

	}
	// TODO: tons of more variants

	return nil
}

func (v *RegionType) FragmentType() string {
	return "RegionType"
}

func (e *RegionType) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
