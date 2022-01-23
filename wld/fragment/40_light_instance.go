package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

// LightInstance information
type LightInstance struct {
	hashIndex uint32
	Reference uint32
	Position  math32.Vector3
	Radius    float32
}

func LoadLightInstance(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &LightInstance{}
	err := parseLightInstance(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse kight instance: %w", err)
	}
	return l, nil
}

func parseLightInstance(r io.ReadSeeker, l *LightInstance) error {
	if l == nil {
		return fmt.Errorf("light instance is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &l.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &l.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &l.Position.X)
	if err != nil {
		return fmt.Errorf("read x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &l.Position.Y)
	if err != nil {
		return fmt.Errorf("read y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &l.Position.Z)
	if err != nil {
		return fmt.Errorf("read z: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &l.Radius)
	if err != nil {
		return fmt.Errorf("read radius: %w", err)
	}

	return nil
}

func (l *LightInstance) FragmentType() string {
	return "Light Instance"
}
