package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// 0x12 trackDef
type trackDef struct {
	nameRef            int32
	flags              uint32
	skeletonTransforms []boneTransform
	//data2              []uint8 // potentailly can be ignored
}

type boneTransform struct {
	translation geo.Vector3
	rotation    geo.Quad4
	scale       float32
	matrix      [16]float32
}

func (e *WLD) trackDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &trackDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	skeletonCount := dec.Uint32()
	for i := 0; i < int(skeletonCount); i++ {
		rotDenom := dec.Int16()
		rotX := dec.Int16()
		rotY := dec.Int16()
		rotZ := dec.Int16()
		shiftX := dec.Int16()
		shiftY := dec.Int16()
		shiftZ := dec.Int16()
		shiftDenom := dec.Int16()
		ft := boneTransform{}
		if shiftDenom != 0 {
			ft.scale = float32(shiftDenom) / 256
			ft.translation.X = float32(shiftX) / 256
			ft.translation.Y = float32(shiftY) / 256
			ft.translation.Z = float32(shiftZ) / 256
		}
		ft.rotation.X = float32(rotX)
		ft.rotation.Y = float32(rotY)
		ft.rotation.Z = float32(rotZ)
		ft.rotation.W = float32(rotDenom)
		ft.rotation = geo.Normalize(ft.rotation)
		def.skeletonTransforms = append(def.skeletonTransforms, ft)
	}

	if dec.Error() != nil {
		return fmt.Errorf("trackDefRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *trackDef) build(e *WLD) error {
	return nil
}

func (e *WLD) trackDefWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
