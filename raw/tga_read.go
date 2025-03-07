package raw

import (
	"encoding/base64"
	"io"
)

// Tga takes a raw TGA type and converts it to an image.Image friendly format
type Tga struct {
	MetaFileName string
	Data         string
}

// Identity returns the type of the struct
func (tga *Tga) Identity() string {
	return "tga"
}

func (tga *Tga) Read(r io.ReadSeeker) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	tga.Data = base64.StdEncoding.EncodeToString(data)
	/* img, err := gdds.Decode(r)
	if err != nil {
		return fmt.Errorf("tga decode: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	benc := base64.NewEncoder(base64.StdEncoding, buf)
	err = png.Encode(benc, img)
	if err != nil {
		return fmt.Errorf("png encode: %w", err)
	}

	tga.Data = buf.String() */
	return nil
}

// SetFileName sets the name of the file
func (tga *Tga) SetFileName(name string) {
	tga.MetaFileName = name
}

// FileName returns the name of the file
func (tga *Tga) FileName() string {
	return tga.MetaFileName
}
