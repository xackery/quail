package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragTrack is a bone in a skeleton. It is Track in libeq, Mob Skeleton Piece Track Reference in openzone, TRACKINSTANCE in wld, TrackDefFragment in lantern
type WldFragTrack struct {
	NameRef int32  `yaml:"name_ref"`
	Track   int32  `yaml:"track_ref"`
	Flags   uint32 `yaml:"flags"`
	Sleep   uint32 `yaml:"sleep"` // if 0x01 is set, this is the number of milliseconds to sleep before starting the animation
}

func (e *WldFragTrack) FragCode() int {
	return FragCodeTrack
}

func (e *WldFragTrack) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.Track)
	enc.Uint32(e.Flags)
	if e.Flags&0x01 == 0x01 {
		enc.Uint32(e.Sleep)
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragTrack) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Track = dec.Int32()
	e.Flags = dec.Uint32()
	if e.Flags&0x01 == 0x01 {
		e.Sleep = dec.Uint32()
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
