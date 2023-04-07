package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/pfs/archive"
)

// SkeletonHierarchy information
type SkeletonHierarchy struct {
	name string
}

func LoadSkeletonHierarchy(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &SkeletonHierarchy{}
	err := parseSkeletonHierarchy(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse SkeletonHiearchy: %w", err)
	}
	return v, nil
}

// based on https://github.com/danwilkins/LanternExtractor/blob/development/0.2.0/LanternExtractor/EQ/Wld/Fragments/SkeletonHierarchy.cs
func parseSkeletonHierarchy(r io.ReadSeeker, v *SkeletonHierarchy) error {
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
	}

	// Always 2 when used in main zone, and object files.
	// This means, it has a bounding radius
	// Some differences in character + model archives
	// Confirmed

	flags := int32(0)
	err = binary.Read(r, binary.LittleEndian, &flags)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}
	dump.Hex(flags, "flags=%d", flags)

	boneCount := int32(0)
	err = binary.Read(r, binary.LittleEndian, &boneCount)
	if err != nil {
		return fmt.Errorf("read boneCount: %w", err)
	}
	dump.Hex(boneCount, "boneCount=%d", boneCount)

	polygonAnimationReference := int32(0)
	err = binary.Read(r, binary.LittleEndian, &polygonAnimationReference)
	if err != nil {
		return fmt.Errorf("read polygonAnimationReference: %w", err)
	}
	dump.Hex(polygonAnimationReference, "polygonAnimationReference=%d(0x%x)", polygonAnimationReference, polygonAnimationReference)

	if flags == 0 {
		unk1 := int32(0)
		err = binary.Read(r, binary.LittleEndian, &unk1)
		if err != nil {
			return fmt.Errorf("read unk1: %w", err)
		}
		dump.Hex(unk1, "unk1=%d", unk1)

		unk2 := int32(0)
		err = binary.Read(r, binary.LittleEndian, &unk2)
		if err != nil {
			return fmt.Errorf("read unk2: %w", err)
		}
		dump.Hex(unk2, "unk2=%d", unk2)

		unk3 := int32(0)
		err = binary.Read(r, binary.LittleEndian, &unk3)
		if err != nil {
			return fmt.Errorf("read unk3: %w", err)
		}
		dump.Hex(unk3, "unk3=%d", unk3)
	}
	if flags&1 > 0 {
		boundingRadius := float32(0)
		err = binary.Read(r, binary.LittleEndian, &boundingRadius)
		if err != nil {
			return fmt.Errorf("read boundingRadius: %w", err)
		}
		dump.Hex(boundingRadius, "boundingRadius=%0.2f", boundingRadius)
	}

	for i := 0; i < int(boneCount); i++ {
		_, err = nameFromHashIndex(r)
		if err != nil {
			return fmt.Errorf("%d boneNameFromHasIndex: %w", i, err)
		}

		boneFlags := int32(0)
		err = binary.Read(r, binary.LittleEndian, &boneFlags)
		if err != nil {
			return fmt.Errorf("read %d boneFlags: %w", i, err)
		}
		dump.Hex(boneFlags, "%dboneFlags=%d", i, boneFlags)

		trackReference := int32(0)
		err = binary.Read(r, binary.LittleEndian, &trackReference)
		if err != nil {
			return fmt.Errorf("read %d trackReference: %w", i, err)
		}
		dump.Hex(trackReference, "%dtrackReference=%d", i, trackReference)

		meshReference := int32(0)
		err = binary.Read(r, binary.LittleEndian, &meshReference)
		if err != nil {
			return fmt.Errorf("read %d meshReference: %w", i, err)
		}
		dump.Hex(meshReference, "%dmeshReference=%d", i, meshReference)

		childCount := int32(0)
		err = binary.Read(r, binary.LittleEndian, &childCount)
		if err != nil {
			return fmt.Errorf("read %d childCount: %w", i, err)
		}
		dump.Hex(childCount, "%dchildCount=%d", i, childCount)
		for j := 0; j < int(childCount); j++ {
			childReference := int32(0)
			err = binary.Read(r, binary.LittleEndian, &childReference)
			if err != nil {
				return fmt.Errorf("read %d %d childReference: %w", i, j, err)
			}
		}

		// All meshes will have vertex bone assignments
		if flags&9 > 0 { //hasMeshReferences
			meshCount := int32(0)
			err = binary.Read(r, binary.LittleEndian, &meshCount)
			if err != nil {
				return fmt.Errorf("read %d meshCount: %w", i, err)
			}
			dump.Hex(meshCount, "%dmeshCount=%d", i, meshCount)
			for i := 0; i < int(meshCount); i++ {
				meshReference := int32(0)
				err = binary.Read(r, binary.LittleEndian, &meshReference)
				if err != nil {
					return fmt.Errorf("read %d meshReference: %w", i, err)
				}
				dump.Hex(meshReference, "%dmeshReference=%d", i, meshReference)
			}

			for i := 0; i < int(meshCount); i++ {
				unkownReference := int32(0)
				err = binary.Read(r, binary.LittleEndian, &unkownReference)
				if err != nil {
					return fmt.Errorf("read %d unkownReference: %w", i, err)
				}
				dump.Hex(unkownReference, "%dunkownReference=%d", i, unkownReference)
			}

		}
	}

	return nil
}

func (v *SkeletonHierarchy) FragmentType() string {
	return "SkeletonHierarchy"
}

func (e *SkeletonHierarchy) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
