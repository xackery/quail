package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// GlobalAmbientLight information
type GlobalAmbientLight struct {
	Unk1 uint32
}

func LoadGlobalAmbientLight(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &GlobalAmbientLight{}
	err := parseGlobalAmbientLight(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse GlobalAmbientLight: %w", err)
	}
	return e, nil
}

func parseGlobalAmbientLight(r io.ReadSeeker, e *GlobalAmbientLight) error {
	if e == nil {
		return fmt.Errorf("globalAmbientLight is nil")
	}

	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("read light source : %w", err)
	}
	return nil
}

func (e *GlobalAmbientLight) FragmentType() string {
	return "Global Ambient Light"
}

func (e *GlobalAmbientLight) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
