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
	NumVertices    uint32                  // size1 in libeq
	NumFaces       uint32                  // size2 in libeq
	BoundingRadius float32                 // params1 in libeq
	ScaleFactor    float32                 // params2 in libeq
	Vertices       [][3]float32            // entries1 in libeq
	Faces          []WldFragPolyhedronFace // entries2 in libeq
}

type WldFragPolyhedronFace struct {
	NumVertices uint32
	Vertices    []uint32
}

func (e *WldFragPolyhedronDef) FragCode() int {
	return FragCodePolyhedronDef
}

func (e *WldFragPolyhedronDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.NumVertices)
	enc.Uint32(e.NumFaces)
	enc.Float32(e.BoundingRadius)
	enc.Float32(e.ScaleFactor)
	for _, entry := range e.Vertices {
		enc.Float32(entry[0])
		enc.Float32(entry[1])
		enc.Float32(entry[2])
	}
	for _, entry := range e.Faces {
		enc.Uint32(entry.NumVertices)
		for _, unk2 := range entry.Vertices {
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
	e.NumVertices = dec.Uint32()
	e.NumFaces = dec.Uint32()
	e.BoundingRadius = dec.Float32()
	e.ScaleFactor = dec.Float32()
	for i := uint32(0); i < e.NumVertices; i++ {
		e.Vertices = append(e.Vertices, [3]float32{dec.Float32(), dec.Float32(), dec.Float32()})
	}
	for i := uint32(0); i < e.NumFaces; i++ {
		entry := WldFragPolyhedronFace{}
		entry.NumVertices = dec.Uint32()
		for j := uint32(0); j < entry.NumVertices; j++ {
			entry.Vertices = append(entry.Vertices, dec.Uint32())
		}
		e.Faces = append(e.Faces, entry)
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
