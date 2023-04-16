package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs/archive"
)

// BitmapInfo information
type BitmapInfo struct {
	name         string
	textureCount uint32
	textureName  string
}

// LoadBitmapInfo loads a Bitmap Information
func LoadBitmapInfo(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &BitmapInfo{}
	err := parseBitmapInfoImage(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse BitmapInfo: %w", err)
	}
	return e, nil
}

func parseBitmapInfoImage(r io.ReadSeeker, e *BitmapInfo) error {
	if e == nil {
		return fmt.Errorf("BitmapInfo is nil")
	}
	var err error

	e.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.textureCount)
	if err != nil {
		return fmt.Errorf("read textureCount: %w", err)
	}

	for i := 0; i < int(e.textureCount); i++ {
		e.textureName, err = nameFromHashIndex(r)
		if err != nil {
			return fmt.Errorf("nameFromHashIndex: %w", err)
		}
		log.Debugf("textureName: %s", e.textureName)
	}

	return nil
}

func (e *BitmapInfo) FragmentType() string {
	return "BitmapInfo"
}

func (e *BitmapInfo) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
