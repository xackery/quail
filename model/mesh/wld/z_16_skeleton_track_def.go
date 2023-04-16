package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

type skeletonTrackDef struct {
	NameRef            int32
	Flags              uint32
	AnimCount          uint32
	CollisionVolumeRef uint32
	centerOffset       geo.Vector3 `bin:"-"`
	radius             float32     `bin:"-"`
	anims              []anim      `bin:"-"`
	meshCount          uint32      `bin:"-"`
	meshRefs           []int32     `bin:"-"`
	skinToAnimRefs     []int32     `bin:"-"`
}

type anim struct {
	NameRef         int32
	Flags           uint32
	TrackRef        int32
	MeshOrSpriteRef int32
	SubAnimCount    uint32
	SubAnims        []uint32 `bin:"-"`
}

func (e *WLD) skeletonTrackDefRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &skeletonTrackDef{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	if def.Flags&0x1 != 0 {
		err = binary.Read(r, binary.LittleEndian, &def.centerOffset)
		if err != nil {
			return fmt.Errorf("read centerOffset: %w", err)
		}
	}

	if def.Flags&0x2 != 0 {
		err = binary.Read(r, binary.LittleEndian, &def.radius)
		if err != nil {
			return fmt.Errorf("read radius: %w", err)
		}
	}

	// decode animCount
	if def.Flags&0x4 != 0 {
		def.anims = make([]anim, def.AnimCount)
		for i := 0; i < int(def.AnimCount); i++ {
			var anim anim
			err = dec.Decode(&anim)
			if err != nil {
				return fmt.Errorf("decode anim: %w", err)
			}
		}
	}

	if def.Flags&0x200 != 0 {
		err = binary.Read(r, binary.LittleEndian, &def.meshCount)
		if err != nil {
			return fmt.Errorf("read meshCount: %w", err)
		}

		for i := 0; i < int(def.meshCount); i++ {
			var meshRef int32
			err = binary.Read(r, binary.LittleEndian, &meshRef)
			if err != nil {
				return fmt.Errorf("read meshCount: %w", err)
			}
			def.meshRefs = append(def.meshRefs, meshRef)

			var skinToAnimRef int32
			err = binary.Read(r, binary.LittleEndian, &skinToAnimRef)
			if err != nil {
				return fmt.Errorf("read skinToAnimRef: %w", err)
			}
			def.skinToAnimRefs = append(def.skinToAnimRefs, skinToAnimRef)
		}
	}

	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *skeletonTrackDef) build(e *WLD) error {
	return nil
}
