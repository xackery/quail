package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// 0x10 skeletonTrackDef
type skeletonTrackDef struct {
	nameRef            int32
	flags              uint32
	animCount          uint32
	collisionVolumeRef uint32
	centerOffset       geo.Vector3
	radius             float32
	anims              []anim
	meshCount          uint32
	meshRefs           []int32
	skinToAnimRefs     []int32
}

type anim struct {
	NameRef         int32
	Flags           uint32
	TrackRef        int32
	MeshOrSpriteRef int32
	SubAnimCount    uint32
	SubAnims        []uint32
}

func (e *WLD) skeletonTrackDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &skeletonTrackDef{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.nameRef = dec.Int32()
	def.flags = dec.Uint32()
	def.animCount = dec.Uint32()
	def.collisionVolumeRef = dec.Uint32()
	if def.flags&0x1 != 0 {
		def.centerOffset.X = dec.Float32()
		def.centerOffset.Y = dec.Float32()
		def.centerOffset.Z = dec.Float32()
	}

	if def.flags&0x2 != 0 {
		def.radius = dec.Float32()
	}

	// TODO: figure out why anim lenght is off alignment
	// decode animCount
	/*if def.flags&0x4 != 0 {
		def.anims = make([]anim, def.animCount)
		for i := 0; i < int(def.animCount); i++ {
			var anim anim
			anim.NameRef = dec.Int32()
			anim.Flags = dec.Uint32()
			anim.TrackRef = dec.Int32()
			anim.MeshOrSpriteRef = dec.Int32()
			anim.SubAnimCount = dec.Uint32()
			anim.SubAnims = make([]uint32, anim.SubAnimCount)
			for j := 0; j < int(anim.SubAnimCount); j++ {
				anim.SubAnims[j] = dec.Uint32()
			}
			def.anims[i] = anim
		}
	}

	if def.flags&0x200 != 0 {
		def.meshCount = dec.Uint32()
		for i := 0; i < int(def.meshCount); i++ {
			def.meshRefs = append(def.meshRefs, dec.Int32())
			def.skinToAnimRefs = append(def.skinToAnimRefs, dec.Int32())
		}
	}
	*/
	if dec.Error() != nil {
		return fmt.Errorf("skeletonTrackDefRead: %w", dec.Error())
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *skeletonTrackDef) build(e *WLD) error {
	return nil
}
