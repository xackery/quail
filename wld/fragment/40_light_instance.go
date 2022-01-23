package fragment

import (
	"bytes"
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
	e := &LightInstance{}
	err := parseLightInstance(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse kight instance: %w", err)
	}
	return e, nil
}

func parseLightInstance(r io.ReadSeeker, e *LightInstance) error {
	if e == nil {
		return fmt.Errorf("light instance is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &e.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.Position.X)
	if err != nil {
		return fmt.Errorf("read x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.Position.Y)
	if err != nil {
		return fmt.Errorf("read y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.Position.Z)
	if err != nil {
		return fmt.Errorf("read z: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.Radius)
	if err != nil {
		return fmt.Errorf("read radius: %w", err)
	}

	return nil
}

func (e *LightInstance) FragmentType() string {
	return "Light Instance"
}

func (e *LightInstance) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
