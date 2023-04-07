package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// TrackReference information
type TrackReference struct {
	name      string
	Reference uint32
	FrameMs   uint32
}

func LoadTrackReference(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &TrackReference{}
	err := parseTrackReference(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse track reference: %w", err)
	}
	return v, nil
}

func parseTrackReference(r io.ReadSeeker, v *TrackReference) error {
	if v == nil {
		return fmt.Errorf("track reference is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flag: %w", err)
	}

	//TODO: review
	// Either 4 or 5 - maybe something to look into
	// Bits are set 0, or 2. 0 has the extra field for delay.
	// 2 doesn't have any additional fields.
	if value&1 == 1 {
		err = binary.Read(r, binary.LittleEndian, &v.FrameMs)
		if err != nil {
			return fmt.Errorf("read frame ms: %w", err)
		}
	}

	return nil
}

func (v *TrackReference) FragmentType() string {
	return "Track Reference"
}

func (e *TrackReference) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
