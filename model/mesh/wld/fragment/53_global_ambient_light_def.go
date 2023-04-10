package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// GlobalAmbientLightDef information
type GlobalAmbientLightDef struct {
	Unk1 uint32
}

func LoadGlobalAmbientLightDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &GlobalAmbientLightDef{}
	err := parseGlobalAmbientLightDef(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse GlobalAmbientLightDef: %w", err)
	}
	return e, nil
}

func parseGlobalAmbientLightDef(r io.ReadSeeker, e *GlobalAmbientLightDef) error {
	if e == nil {
		return fmt.Errorf("GlobalAmbientLightDef is nil")
	}

	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("read lightDef source : %w", err)
	}
	return nil
}

func (e *GlobalAmbientLightDef) FragmentType() string {
	return "GlobalAmbientLightDef"
}

func (e *GlobalAmbientLightDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
