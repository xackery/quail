package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/g3n/engine/math32"
)

// LightSourceInstance information
type LightSourceInstance struct {
	HashIndex uint32
	Flags     uint32
	Reference uint32
	Position  math32.Vector3
	Radius    float32
}

func loadLightSourceInstance(r io.ReadSeeker) (Fragment, error) {
	l := &LightSourceInstance{}
	err := parseLightSourceInstance(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse LightSourceInstance: %w", err)
	}
	return l, nil
}

func parseLightSourceInstance(r io.ReadSeeker, l *LightSourceInstance) error {
	if l == nil {
		return fmt.Errorf("lightsourceInstance is nil")
	}
	err := binary.Read(r, binary.LittleEndian, &l)
	if err != nil {
		return fmt.Errorf("read light source instance: %w", err)
	}
	return nil
}

func (l *LightSourceInstance) FragmentType() string {
	return "Light Source Instance"
}
