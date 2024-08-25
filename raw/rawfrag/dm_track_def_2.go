package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDmTrackDef2 is DmTrackDef2 in libeq, Mesh Animated Vertices in openzone, DMTRACKDEF in wld, MeshAnimatedVertices in lantern
type WldFragDmTrackDef2 struct {
	NameRef int32
	Flags   uint32
	Param1  uint16
	Param2  uint16
	Scale   uint16
	Frames  [][][3]int16
	Size6   uint16
}

func (e *WldFragDmTrackDef2) FragCode() int {
	return FragCodeDmTrackDef2
}

func (e *WldFragDmTrackDef2) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	if len(e.Frames) < 1 {
		return fmt.Errorf("no frames found")
	}
	enc.Uint16(uint16(len(e.Frames[0])))
	enc.Uint16(uint16(len(e.Frames)))
	enc.Uint16(e.Param1)
	enc.Uint16(e.Param2)
	enc.Uint16(e.Scale)
	for _, frame := range e.Frames {
		for _, vert := range frame {
			enc.Int16(vert[0])
			enc.Int16(vert[1])
			enc.Int16(vert[2])
		}
	}
	enc.Uint16(e.Size6)
	err := enc.Error()
	if err != nil {
		return err
	}
	return nil
}

func (e *WldFragDmTrackDef2) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	vertexCount := dec.Uint16()
	frameCount := dec.Uint16()
	e.Param1 = dec.Uint16()
	e.Param2 = dec.Uint16()
	e.Scale = dec.Uint16()
	e.Frames = make([][][3]int16, frameCount)
	for i := range e.Frames {
		e.Frames[i] = make([][3]int16, vertexCount)
		for j := range e.Frames[i] {
			e.Frames[i][j][0] = dec.Int16()
			e.Frames[i][j][1] = dec.Int16()
			e.Frames[i][j][2] = dec.Int16()
		}
	}
	e.Size6 = dec.Uint16()
	err := dec.Error()
	if err != nil {
		return err
	}
	return nil
}
