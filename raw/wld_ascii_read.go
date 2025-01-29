package raw

import (
	"fmt"
	"io"
)

type WldAscii struct {
	MetaFileName string
	Data         string
}

func (wld *WldAscii) Identity() string {
	return "wld.ascii"
}

// Read reads a wld file that was prepped by Load
func (wld *WldAscii) Read(r io.ReadSeeker) error {

	data := make([]byte, 0)
	for {
		b := make([]byte, 1)
		_, err := r.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read: %w", err)
		}
		data = append(data, b...)
	}
	wld.Data = string(data)

	return nil
}

// SetFileName sets the name of the file
func (wld *WldAscii) SetFileName(name string) {
	wld.MetaFileName = name
}

// FileName returns the name of the file
func (wld *WldAscii) FileName() string {
	return wld.MetaFileName
}
