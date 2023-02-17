package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
)

// AmbientLight information
type AmbientLight struct {
	Unk1 uint32
}

func LoadAmbientLight(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &AmbientLight{}
	err := parseAmbientLight(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse AmbientLight: %w", err)
	}
	return e, nil
}

func parseAmbientLight(r io.ReadSeeker, e *AmbientLight) error {
	if e == nil {
		return fmt.Errorf("AmbientLight is nil")
	}

	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("read light source : %w", err)
	}
	return nil
}

func (e *AmbientLight) FragmentType() string {
	return "Ambient Light"
}

func (e *AmbientLight) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
