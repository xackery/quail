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
	skeletonCount      uint32
	skeletonTransforms []boneTransform
	//data2              []uint8 // potentailly can be ignored
}

type boneTransform struct {
	translation geo.Vector3
	rotation    geo.Quad4
	scale       float32
}

func (e *WLD) trackDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &trackDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.skeletonCount = dec.Uint32()
	for i := 0; i < int(def.skeletonCount); i++ {
		var bone boneTransform
		bone.translation.X = dec.Float32()
		bone.translation.Y = dec.Float32()
		bone.translation.Z = dec.Float32()
		bone.rotation.X = dec.Float32()
		bone.rotation.Y = dec.Float32()
		bone.rotation.Z = dec.Float32()
		bone.rotation.W = dec.Float32()
		bone.scale = dec.Float32()
		def.skeletonTransforms = append(def.skeletonTransforms, bone)
	}

	if dec.Error() != nil {
		return fmt.Errorf("trackDefRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *trackDef) build(e *WLD) error {
	return nil
}
