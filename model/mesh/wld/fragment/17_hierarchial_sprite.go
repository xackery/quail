package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// skeletonTrack information
type skeletonTrack struct {
	name      string
	Reference uint32
	FrameMs   uint32
}

func LoadskeletonTrack(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &skeletonTrack{}
	err := parseskeletonTrack(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse skeletonTrack: %w", err)
	}
	return v, nil
}

func parseskeletonTrack(r io.ReadSeeker, v *skeletonTrack) error {
	if v == nil {
		return fmt.Errorf("skeletonTrack is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flag: %w", err)
	}

	if value&1 == 1 {
		err = binary.Read(r, binary.LittleEndian, &v.FrameMs)
		if err != nil {
			return fmt.Errorf("read frame ms: %w", err)
		}
	}

	return nil
}

func (v *skeletonTrack) FragmentType() string {
	return "skeletonTrack"
}

func (e *skeletonTrack) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
