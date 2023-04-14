package pts

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model/geo"
)

func (e *PTS) Decode(r io.ReadSeeker) error {
	var err error
	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'E', 'Q', 'P', 'T'} {
		return fmt.Errorf("header does not match EQPT, got %s", header)
	}

	particleCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &particleCount)
	if err != nil {
		return fmt.Errorf("read particleCount: %w", err)
	}
	dump.Hex(particleCount, "particleCount=%d", particleCount)

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("read header version: %w", err)
	}
	dump.Hex(version, "version=%d", version)
	if version != 1 {
		return fmt.Errorf("version is %d, wanted 1", version)
	}

	for i := 0; i < int(particleCount); i++ {
		entry := geo.NewParticlePoint()

		entry.Name, err = helper.ReadFixedString(r, 64)
		if err != nil {
			return fmt.Errorf("%d read name: %w", i, err)
		}
		dump.Hex(entry.Name, "%dname=%s", i, entry.Name)

		entry.Bone, err = helper.ReadFixedString(r, 64)
		if err != nil {
			return fmt.Errorf("%d read bone: %w", i, err)
		}
		dump.Hex(entry.Bone, "%dbone=%s", i, entry.Bone)

		err = binary.Read(r, binary.LittleEndian, entry.Translation)
		if err != nil {
			return fmt.Errorf("%d read translation: %w", i, err)
		}
		dump.Hex(entry.Translation, "%dtranslation=%0.0f", i, entry.Translation)

		err = binary.Read(r, binary.LittleEndian, entry.Rotation)
		if err != nil {
			return fmt.Errorf("%d read rotation: %w", i, err)
		}
		dump.Hex(entry.Rotation, "%drotation=%0.0f", i, entry.Rotation)

		err = binary.Read(r, binary.LittleEndian, entry.Scale)
		if err != nil {
			return fmt.Errorf("%d read scale: %w", i, err)
		}
		dump.Hex(entry.Scale, "%dscale=%0.0f", i, entry.Scale)

		e.particleManager.PointAdd(entry)
	}
	return nil
}
