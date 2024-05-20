package raw

import (
	"encoding/base64"
	"io"
)

// Bmp takes a raw BMP type and converts it to png
type Bmp struct {
	MetaFileName string `yaml:"file_name"`
	Data         string `yaml:"data"`
}

// Identity returns the type of the struct
func (bmp *Bmp) Identity() string {
	return "bmp"
}

func (bmp *Bmp) Read(r io.ReadSeeker) error {
	//var err error
	//dec := encdec.NewDecoder(r, binary.LittleEndian)

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	bmp.Data = base64.StdEncoding.EncodeToString(data)

	/*var img image.Image
	header := dec.StringFixed(4)
	r.Seek(0, io.SeekStart)
	if header[0:3] == "DDS" {
		img, err = dds.Decode(r)
		if err != nil {
			return fmt.Errorf("dds decode: %w", err)
		}
	} else {
		img, err = gbmp.Decode(r)
		if err != nil {
			return fmt.Errorf("bmp decode: %w", err)
		}
	}

	buf := bytes.NewBuffer(nil)
	benc := base64.NewEncoder(base64.StdEncoding, buf)
	err = png.Encode(benc, img)
	if err != nil {
		return fmt.Errorf("png encode: %w", err)
	}

	bmp.Data = buf.String()
	*/

	return nil
}

// SetFileName sets the name of the file
func (bmp *Bmp) SetFileName(name string) {
	bmp.MetaFileName = name
}

// FileName returns the name of the file
func (bmp *Bmp) FileName() string {
	return bmp.MetaFileName
}
