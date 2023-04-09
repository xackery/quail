package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// LightInstance information
type LightInstance struct {
	name      string
	Reference uint32
	Position  [3]float32
	Radius    float32
}

func LoadLightInstance(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &LightInstance{}
	err := parseLightInstance(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse kight instance: %w", err)
	}
	return e, nil
}

func parseLightInstance(r io.ReadSeeker, v *LightInstance) error {
	if v == nil {
		return fmt.Errorf("light instance is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Reference)
	if err != nil {
		return fmt.Errorf("read reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Position[0])
	if err != nil {
		return fmt.Errorf("read x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Position[1])
	if err != nil {
		return fmt.Errorf("read y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Position[2])
	if err != nil {
		return fmt.Errorf("read z: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Radius)
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
