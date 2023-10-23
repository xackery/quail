package ani

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

// Decode decodes an ANI file
func Decode(animation *common.Animation, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGA" {
		return fmt.Errorf("invalid header %s, wanted EQGA", header)
	}
	tag.New()

	version := dec.Uint32()
	animation.Header.Version = int(version)
	nameLength := int(dec.Uint32())
	boneCount := dec.Uint32()

	animation.IsStrict = false
	if version > 1 {
		isStrict := dec.Uint32()
		if isStrict > 0 {
			animation.IsStrict = true
		}
	}
	tag.Add(0, dec.Pos(), "red", "header")

	nameData := dec.Bytes(int(nameLength))
	tag.Add(tag.LastPos(), dec.Pos(), "green", "names")

	names := make(map[uint32]string)
	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	for i := 0; i < int(boneCount); i++ {
		bone := &common.BoneAnimation{}
		bone.FrameCount = dec.Uint32()

		nameOffset := dec.Uint32()
		bone.Name = names[nameOffset]

		for j := 0; j < int(bone.FrameCount); j++ {
			frame := &common.BoneAnimationFrame{}

			frame.Milliseconds = dec.Uint32()
			frame.Translation = common.Vector3{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
			}

			frame.Rotation = common.Quad4{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
				W: dec.Float32(),
			}

			frame.Scale = common.Vector3{
				X: dec.Float32(),
				Y: dec.Float32(),
				Z: dec.Float32(),
			}
			bone.Frames = append(bone.Frames, frame)
		}

		dump.HexRange([]byte{0x01, 0x02}, int(bone.FrameCount*44), "%dbone frames(%d bytes)", i, bone.FrameCount*44)
		animation.Bones = append(animation.Bones, bone)
	}

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	log.Debugf("%s (ani) decoded %d bones, bone 0 had %d frames", animation.Header.Name, len(animation.Bones), animation.Bones[0].FrameCount)
	return nil
}
