package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

func (ani *Ani) Write(w io.Writer) error {
	ani.NameClear()
	for _, bone := range ani.Bones {
		ani.NameAdd(bone.Name)
	}
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.StringFixed("EQGA", 4)
	enc.Uint32(ani.Version)
	nameData := ani.NameData()
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
		enc.Int32(ani.NameIndex(bone.Name))
		for _, frame := range bone.Frames {
			enc.Uint32(frame.Milliseconds)
			enc.Float32(frame.Translation.X)
			enc.Float32(frame.Translation.Y)
			enc.Float32(frame.Translation.Z)
			enc.Float32(frame.Rotation.X)
			enc.Float32(frame.Rotation.Y)
			enc.Float32(frame.Rotation.Z)
			enc.Float32(frame.Rotation.W)
			enc.Float32(frame.Scale.X)
			enc.Float32(frame.Scale.Y)
			enc.Float32(frame.Scale.Z)
		}
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
