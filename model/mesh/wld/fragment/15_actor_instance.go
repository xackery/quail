package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// ActorInstance information
type ActorInstance struct {
	name     string
	Position [3]float32
	Rotation [3]float32
	Scale    [3]float32
}

func LoadActorInstance(r io.ReadSeeker) (archive.WldFragmenter, error) {
	v := &ActorInstance{}
	err := parseActorInstance(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse Actor instance: %w", err)
	}
	return v, nil
}

func parseActorInstance(r io.ReadSeeker, v *ActorInstance) error {
	if v == nil {
		return fmt.Errorf("Actor instance is nil")
	}
	var value uint32
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	var flags uint32
	err = binary.Read(r, binary.LittleEndian, &flags)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}
	// Main zone: 0x2E, Actors: 0x32E
	//TODO if flags != 0x2E && flags != 0x32E {
	//	return fmt.Errorf("unknown flags want 0x2E or 0x32E, got 0x%x", flags)
	//}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown2: %w", err)
	}

	if flags == 0x2E && value != 0x16 {
		return fmt.Errorf("expected unknown2 to be 0x16, got 0x%x", value)
	}

	if flags == 0x32E && value != 0 {
		return fmt.Errorf("expected unknown2 to be 0, got 0x%x", value)
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

	var rotX, rotY, rotZ float32
	err = binary.Read(r, binary.LittleEndian, &rotX)
	if err != nil {
		return fmt.Errorf("read rotX: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &rotY)
	if err != nil {
		return fmt.Errorf("read rotY: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &rotZ)
	if err != nil {
		return fmt.Errorf("read rotZ: %w", err)
	}

	modifier := float32(float32(1) / float32(512) * 360)
	v.Rotation[0] = 0
	v.Rotation[1] = rotY * modifier
	v.Rotation[2] = -(rotX * modifier)

	err = binary.Read(r, binary.LittleEndian, &rotX)
	if err != nil {
		return fmt.Errorf("read scaleX: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &rotY)
	if err != nil {
		return fmt.Errorf("read scaleY: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &rotZ)
	if err != nil {
		return fmt.Errorf("read scaleZ: %w", err)
	}

	v.Scale[0], v.Scale[1], v.Scale[2] = rotY, rotY, rotY

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read colorFragment: %w", err)
	}
	//if value != 0 {
	//TODO: look up vertexcolorreference
	//}

	return nil
}

func (v *ActorInstance) FragmentType() string {
	return "Actor Instance"
}

func (e *ActorInstance) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
