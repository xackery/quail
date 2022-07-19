package prt

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/helper"
)

func (e *PRT) Load(r io.ReadSeeker) error {
	var err error
	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'P', 'T', 'C', 'L'} {
		return fmt.Errorf("header does not match PTCL")
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
	if version < 4 {
		return fmt.Errorf("version is %d, wanted > 3", version)
	}

	for i := 0; i < int(particleCount); i++ {
		entry := &common.ParticleEntry{}

		err = binary.Read(r, binary.LittleEndian, &entry.Id)
		if err != nil {
			return fmt.Errorf("%d read id: %w", i, err)
		}
		dump.Hex(entry.Id, "%did=%d", i, entry.Id)

		if version == 5 {
			err = binary.Read(r, binary.LittleEndian, &entry.Id2)
			if err != nil {
				return fmt.Errorf("%d read id2: %w", i, err)
			}
			dump.Hex(entry.Id2, "%did2=%d", i, entry.Id2)
		}

		entry.Bone, err = helper.ReadFixedString(r, 64)
		if err != nil {
			return fmt.Errorf("%d read name: %w", i, err)
		}
		dump.Hex(entry.Bone, "%dbone=%s", i, entry.Bone)

		err = binary.Read(r, binary.LittleEndian, &entry.UnknownA)
		if err != nil {
			return fmt.Errorf("%d read unknownA: %w", i, err)
		}
		dump.Hex(entry.UnknownA, "%dunknownA=%d", i, entry.UnknownA)

		err = binary.Read(r, binary.LittleEndian, &entry.Duration)
		if err != nil {
			return fmt.Errorf("%d read duration: %w", i, err)
		}
		dump.Hex(entry.Duration, "%dduration=%d", i, entry.Duration)

		err = binary.Read(r, binary.LittleEndian, &entry.UnknownB)
		if err != nil {
			return fmt.Errorf("%d read unknownB: %w", i, err)
		}
		dump.Hex(entry.UnknownB, "%dunknownB=%d", i, entry.UnknownB)

		err = binary.Read(r, binary.LittleEndian, &entry.UnknownFFFFFFFF)
		if err != nil {
			return fmt.Errorf("%d read unknownFFFFFFFF: %w", i, err)
		}
		dump.Hex(entry.UnknownFFFFFFFF, "%dunknownFFFFFFFF=%d", i, entry.UnknownFFFFFFFF)

		err = binary.Read(r, binary.LittleEndian, &entry.UnknownC)
		if err != nil {
			return fmt.Errorf("%d read unknownC: %w", i, err)
		}
		dump.Hex(entry.UnknownC, "%dunknownC=%d", i, entry.UnknownC)
		e.particles = append(e.particles, entry)
	}
	return nil
}
