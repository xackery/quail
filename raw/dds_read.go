package raw

import (
	"encoding/base64"
	"io"
)

// Dds takes a raw DDS type and converts it to an image.Image friendly format
type Dds struct {
	MetaFileName string `yaml:"file_name"`
	Data         string `yaml:"data"`
}

func (dds *Dds) Read(r io.ReadSeeker) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	dds.Data = base64.StdEncoding.EncodeToString(data)
	/* img, err := gdds.Decode(r)
	if err != nil {
		return fmt.Errorf("dds decode: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	benc := base64.NewEncoder(base64.StdEncoding, buf)
	err = png.Encode(benc, img)
	if err != nil {
		return fmt.Errorf("png encode: %w", err)
	}

	dds.Data = buf.String() */
	return nil
}

// SetFileName sets the name of the file
func (dds *Dds) SetFileName(name string) {
	dds.MetaFileName = name
}

// FileName returns the name of the file
func (dds *Dds) FileName() string {
	return dds.MetaFileName
}
