package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// TextureImage information
type TextureImage struct {
	name         string
	textureCount uint32
	textureName  string
}

// LoadTextureImage loads a TextureImage
func LoadTextureImage(r io.ReadSeeker) (archive.WldFragmenter, error) {
	e := &TextureImage{}
	err := parseTextureImage(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse TextureImage: %w", err)
	}
	return e, nil
}

func parseTextureImage(r io.ReadSeeker, e *TextureImage) error {
	if e == nil {
		return fmt.Errorf("TextureImage is nil")
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
		fmt.Println("textureName", e.textureName)
	}

	return nil
}

func (e *TextureImage) FragmentType() string {
	return "TextureImage"
}

func (e *TextureImage) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
