package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDmRGBTrackDef is a list of colors, one per vertex, for baked lighting. It is DmRGBTrackDef in libeq, Vertex Color in openzone, empty in wld, VertexColors in lantern
type WldFragDmRGBTrackDef struct {
	NameRef int32
	Data1   uint32 // usually contains 1
	Data2   uint32 // usually contains 1
	Sleep   uint32 // usually contains 200
	Data4   uint32 // usually contains 0
	RGBAs   [][4]uint8
}

func (e *WldFragDmRGBTrackDef) FragCode() int {
	return FragCodeDmRGBTrackDef
}

func (e *WldFragDmRGBTrackDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.Int32(e.NameRef)
	if e.Data1 == 0 {
		e.Data1 = 1
	}
	enc.Uint32(e.Data1)
	enc.Uint32(uint32(len(e.RGBAs)))
	if e.Data2 == 0 {
		e.Data2 = 1
	}
	enc.Uint32(e.Data2)
	if e.Sleep == 0 {
		e.Sleep = 200
	}
	enc.Uint32(e.Sleep)
	enc.Uint32(e.Data4)

	for _, rgba := range e.RGBAs {
		enc.Uint8(rgba[0])
		enc.Uint8(rgba[1])
		enc.Uint8(rgba[2])
		enc.Uint8(rgba[3])
	}

	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func (e *WldFragDmRGBTrackDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	e.NameRef = dec.Int32()
	e.Data1 = dec.Uint32()
	if e.Data1 != 1 {
		fmt.Printf("Data1 on rgbtrack is not 1 (%d), tell xack you found this!\n", e.Data1)
	}
	numRGBA := dec.Uint32()
	e.Data2 = dec.Uint32()
	if e.Data2 != 1 {
		fmt.Printf("Data2 on rgbtrack is not 1 (%d), tell xack you found this!\n", e.Data2)
	}
	e.Sleep = dec.Uint32()
	e.Data4 = dec.Uint32()
	if e.Data4 != 0 {
		fmt.Printf("Data4 (NumVertices) on rgbtrack is not 0 (%d), tell xack you found this!\n", e.Data4)
	}

	e.RGBAs = make([][4]uint8, numRGBA)
	for i := range e.RGBAs {
		e.RGBAs[i][0] = dec.Uint8()
		e.RGBAs[i][1] = dec.Uint8()
		e.RGBAs[i][2] = dec.Uint8()
		e.RGBAs[i][3] = dec.Uint8()
	}
	return nil
}
