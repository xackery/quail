package prt

import (
	"encoding/binary"
	"fmt"
	"io"

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
		entry := &particleEntry{}

		err = binary.Read(r, binary.LittleEndian, &entry.id)
		if err != nil {
			return fmt.Errorf("%d read id: %w", i, err)
		}
		dump.Hex(entry.id, "%did=%d", i, entry.id)

		if version == 5 {
			err = binary.Read(r, binary.LittleEndian, &entry.id2)
			if err != nil {
				return fmt.Errorf("%d read id2: %w", i, err)
			}
			dump.Hex(entry.id2, "%did2=%d", i, entry.id2)
		}

		entry.name, err = helper.ReadFixedString(r, 64)
		if err != nil {
			return fmt.Errorf("%d read name: %w", i, err)
		}
		dump.Hex(entry.name, "%dname=%s", i, entry.name)

		err = binary.Read(r, binary.LittleEndian, &entry.unknownA)
		if err != nil {
			return fmt.Errorf("%d read unknownA: %w", i, err)
		}
		dump.Hex(entry.unknownA, "%dunknownA=%d", i, entry.unknownA)

		err = binary.Read(r, binary.LittleEndian, &entry.duration)
		if err != nil {
			return fmt.Errorf("%d read duration: %w", i, err)
		}
		dump.Hex(entry.duration, "%dduration=%d", i, entry.duration)

		err = binary.Read(r, binary.LittleEndian, &entry.unknownB)
		if err != nil {
			return fmt.Errorf("%d read unknownB: %w", i, err)
		}
		dump.Hex(entry.unknownB, "%dunknownB=%d", i, entry.unknownB)

		err = binary.Read(r, binary.LittleEndian, &entry.unknownFFFFFFFF)
		if err != nil {
			return fmt.Errorf("%d read unknownFFFFFFFF: %w", i, err)
		}
		dump.Hex(entry.unknownFFFFFFFF, "%dunknownFFFFFFFF=%d", i, entry.unknownFFFFFFFF)

		err = binary.Read(r, binary.LittleEndian, &entry.unknownC)
		if err != nil {
			return fmt.Errorf("%d read unknownC: %w", i, err)
		}
		dump.Hex(entry.unknownC, "%dunknownC=%d", i, entry.unknownC)
		e.particles = append(e.particles, entry)
	}
	return nil
}
