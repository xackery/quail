package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragSphereListDef is SphereListDef in libeq, empty in openzone, SPHERELISTDEFINITION in wld
type WldFragSphereListDef struct {
	NameRef     int32         `yaml:"name_ref"`
	Flags       uint32        `yaml:"flags"`
	SphereCount uint32        `yaml:"sphere_count"`
	Radius      float32       `yaml:"radius"`
	Scale       float32       `yaml:"scale"`
	Spheres     []model.Quad4 `yaml:"spheres"`
}

func (e *WldFragSphereListDef) FragCode() int {
	return FragCodeSphereListDef
}

func (e *WldFragSphereListDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SphereCount)
	enc.Float32(e.Radius)
	enc.Float32(e.Scale)
	for _, sphere := range e.Spheres {
		enc.Float32(sphere.X)
		enc.Float32(sphere.Y)
		enc.Float32(sphere.Z)
		enc.Float32(sphere.W)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSphereListDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.SphereCount = dec.Uint32()
	e.Radius = dec.Float32()
	e.Scale = dec.Float32()
	for i := uint32(0); i < e.SphereCount; i++ {
		var sphere model.Quad4
		sphere.X = dec.Float32()
		sphere.Y = dec.Float32()
		sphere.Z = dec.Float32()
		sphere.W = dec.Float32()
		e.Spheres = append(e.Spheres, sphere)
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
