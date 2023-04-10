package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// DmRGBTrack, Referenced by an ObjectInstance fragment.
type DmRGBTrack struct {
	VertexColor *DmRGBTrackDef
	Reference   uint32
	hashIndex   uint32
}

func LoadDmRGBTrack(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &DmRGBTrack{}
	err := parseDmRGBTrack(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse DmRGBTrack: %w", err)
	}
	return v, nil
}

func parseDmRGBTrack(r io.ReadSeeker, v *DmRGBTrack) error {
	if v == nil {
		return fmt.Errorf("DmRGBTrack is nil")
	}
	err := binary.Read(r, binary.LittleEndian, &v.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	return nil
}

func (v *DmRGBTrack) FragmentType() string {
	return "DmRGBTrack"
}

func (e *DmRGBTrack) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
