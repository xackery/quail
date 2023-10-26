package ani

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/tag"
)

// Encode writes a ani file
func Encode(anim *common.Animation, version uint32, w io.Writer) error {
	var err error
	modelNames := []string{}

	if len(anim.Bones) > 0 {
		modelNames = append(modelNames, anim.Header.Name)
	}

	names, nameData, err := anim.NameBuild(modelNames)
	if err != nil {
		return fmt.Errorf("nameBuild: %w", err)
	}

	boneData, err := anim.BoneBuild(version, true, names)
	if err != nil {
		return fmt.Errorf("boneBuild: %w", err)
	}

	tag.New()
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGA")
	enc.Uint32(version)
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(len(anim.Bones)))
	if version > 1 {
		if anim.IsStrict {
			enc.Uint32(1)
		} else {
			enc.Uint32(0)
		}
	}
	enc.Bytes(nameData)
	enc.Bytes(boneData)

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	//log.Debugf("%s encoded %d bones, bone 0 had %d frames", anim.Header.Name, len(anim.Bones), anim.Bones[0].FrameCount)
	return nil
}
