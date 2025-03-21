package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

func (ani *Ani) Write(w io.Writer) error {
	if ani.name == nil {
		ani.name = &eqgName{}
	}

	ani.name.clear()
	for _, bone := range ani.Bones {
		ani.name.add(bone.Name)
	}
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.StringFixed("EQGA", 4)
	enc.Uint32(ani.Version)
	nameData := ani.name.data()
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(len(ani.Bones)))

	if ani.Version > 1 {
		if ani.IsStrict {
			enc.Uint32(1)
		} else {
			enc.Uint32(0)
		}
	}
	enc.Bytes(nameData)

	for _, bone := range ani.Bones {
		enc.Uint32(uint32(len(bone.Frames)))
		enc.Int32(ani.name.offsetByName(bone.Name))
		for _, frame := range bone.Frames {
			enc.Uint32(frame.Milliseconds)
			enc.Float32(frame.Translation[0])
			enc.Float32(frame.Translation[1])
			enc.Float32(frame.Translation[2])
			enc.Float32(frame.Rotation[0])
			enc.Float32(frame.Rotation[1])
			enc.Float32(frame.Rotation[2])
			enc.Float32(frame.Rotation[3])
			enc.Float32(frame.Scale[0])
			enc.Float32(frame.Scale[1])
			enc.Float32(frame.Scale[2])
		}
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
