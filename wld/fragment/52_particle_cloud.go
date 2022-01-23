package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// ParticleCloud information
type ParticleCloud struct {
	hashIndex uint32
}

func LoadParticleCloud(r io.ReadSeeker) (common.WldFragmenter, error) {
	v := &ParticleCloud{}
	err := parseParticleCloud(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse particle cloud: %w", err)
	}
	return v, nil
}

func parseParticleCloud(r io.ReadSeeker, v *ParticleCloud) error {
	if v == nil {
		return fmt.Errorf("particle cloud is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &v.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}
	if value != 4 {
		return fmt.Errorf("flags wanted 4, got %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value8: %w", err)
	}
	if value != 3 {
		return fmt.Errorf("value8 wanted 3, got %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value12: %w", err)
	}
	if value != 1 && value != 3 && value != 4 {
		return fmt.Errorf("value12 wanted 1 3 or 4, got %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value16: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value20: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value24: %w", err)
	}
	if value != 0 {
		return fmt.Errorf("value24 wanted 0, got %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value28: %w", err)
	}
	if value != 0 {
		return fmt.Errorf("value28 wanted 0, got %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value32: %w", err)
	}
	if value != 0 {
		return fmt.Errorf("value32 wanted 0, got %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value36: %w", err)
	}
	if value != 0 {
		return fmt.Errorf("value36 wanted 0, got %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value40: %w", err)
	}
	if value != 0 {
		return fmt.Errorf("value40 wanted 0, got %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value44: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value48: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value52: %w", err)
	}
	//s/b 100's, e.g. 1000, 100, 750, 500, 1600, 2500.

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value56: %w", err)
	}
	//s/b low numbers, e.g. 4, 5, 8, 10, 0

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value60: %w", err)
	}
	if value != 0 && value != 1 {
		return fmt.Errorf("value60 want 0 or 1, got %d", value)
	}

	var intvalue int32
	err = binary.Read(r, binary.LittleEndian, &intvalue)
	if err != nil {
		return fmt.Errorf("read value64: %w", err)
	}
	if intvalue != 0 && intvalue != -1 {
		return fmt.Errorf("value64 want 0 or -1, got %d", intvalue)
	}

	err = binary.Read(r, binary.LittleEndian, &intvalue)
	if err != nil {
		return fmt.Errorf("read value68: %w", err)
	}
	if intvalue != 0 && intvalue != -1 {
		return fmt.Errorf("value68 want 0 or -1, got %d", intvalue)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value72: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value76: %w", err)
	}
	//float?

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read value80: %w", err)
	}
	//float?
	return nil
}

func (v *ParticleCloud) FragmentType() string {
	return "Particle Cloud"
}
