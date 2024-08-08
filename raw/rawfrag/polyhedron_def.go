package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragPolyhedronDef is PolyhedronDef in libeq, Polygon animation in openzone, POLYHEDRONDEFINITION in wld, Fragment17 in lantern
type WldFragPolyhedronDef struct {
	NameRef        int32
	Flags          uint32
	BoundingRadius float32      // params1 in libeq
	ScaleFactor    float32      // params2 in libeq
	Vertices       [][3]float32 // entries1 in libeq
	Faces          [][]uint32   // entries2 in libeq
}

func (e *WldFragPolyhedronDef) FragCode() int {
	return FragCodePolyhedronDef
}

func (e *WldFragPolyhedronDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Vertices)))
	enc.Uint32(uint32(len(e.Faces)))
	enc.Float32(e.BoundingRadius)
	enc.Float32(e.ScaleFactor)
	for _, entry := range e.Vertices {
		enc.Float32(entry[0])
		enc.Float32(entry[1])
		enc.Float32(entry[2])
	}
	for _, faceEntries := range e.Faces {
		enc.Uint32(uint32(len(faceEntries)))
		for _, unk2 := range faceEntries {
			enc.Uint32(unk2)
		}
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPolyhedronDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	vertexCount := dec.Uint32()
	faceCount := dec.Uint32()
	e.BoundingRadius = dec.Float32()
	e.ScaleFactor = dec.Float32()
	e.Vertices = make([][3]float32, vertexCount)
	for i := range e.Vertices {
		e.Vertices[i] = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	}
	e.Faces = make([][]uint32, faceCount)
	for i := range e.Faces {
		entryCount := dec.Uint32()
		e.Faces[i] = make([]uint32, entryCount)
		for j := range e.Faces[i] {
			e.Faces[i][j] = dec.Uint32()
		}
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
