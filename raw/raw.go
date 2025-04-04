package raw

import (
	"fmt"
	"io"
	"strings"
)

type Reader interface {
	Identity() string
	Read(r io.ReadSeeker) error
	SetFileName(name string)
}

type Writer interface {
	Identity() string
	FileName() string
	Write(w io.Writer) error
}

type ReadWriter interface {
	Reader
	Writer
}

// New takes an extension and returns a ReadWriter that can parse it
func New(ext string) ReadWriter {
	switch ext {
	case ".ani":
		return &Ani{}
	case ".bmp":
		return &Bmp{}
	case ".dat":
		return &Dat{}
	case ".dds":
		return &Dds{}
	case ".edd":
		return &Edd{}
	case ".lay":
		return &Lay{}
	case ".lit":
		return &Lit{}
	case ".lod":
		return &Lod{}
	case ".mds":
		return &Mds{}
	case ".mod":
		return &Mod{}
	case ".png":
		return &Png{}
	case ".prt":
		return &Prt{}
	case ".pts":
		return &Pts{}
	case ".ter":
		return &Ter{}
	case ".tog":
		return &Tog{}
	case ".wld":
		return &Wld{}
	case ".zon":
		return &Zon{}
	case ".env":
		return &Unk{}
	case ".txt":
		return &Txt{}
	case ".sps": // map file, safely ignored
		return &Txt{}
	case ".sph": // map file, safely ignored
		return &Txt{}
	case ".fx": // raw fxo files
		return &Txt{}
	case ".spk": // map file, safely ignored
		return &Txt{}
	case ".spm": // map file, safely ignored
		return &Txt{}
	case ".ms": // map file, safely ignored
		return &Txt{}
	case ".mdf": // model definition file, safely ignored
		return &Txt{}
	case ".tga":
		return &Tga{}
	case ".def":
		return &Def{}
	case ".jpg":
		return &Jpg{}
	case ".bat":
		return &Txt{}
	case ".db":
		return &Txt{}
	case ".max":
		return &Txt{}
	case ".uvw":
		return &Txt{}
	case ".eco":
		return &Eco{}
	case ".rfd":
		return &Rfd{}
	case ".prj":
		return &Txt{}
	case ".obg":
		return &Obg{}
	default:
		return nil
	}
}

// Read takes an extension and a reader and returns a ReadWriter that can parse it
func Read(ext string, r io.ReadSeeker) (ReadWriter, error) {
	reader := New(ext)
	if reader == nil {
		return nil, fmt.Errorf("unknown extension %s", ext)
	}
	err := reader.Read(r)
	if err != nil {
		if ext == ".wld" && strings.Contains(err.Error(), "header wanted 0x023d5054") {
			r.Seek(0, io.SeekStart)
			reader = &WldAscii{}
			err = reader.Read(r)
			if err != nil {
				return nil, err
			}
			return reader, nil
		}
		return nil, err
	}
	return reader, nil
}

// Write takes an extension and a writer and returns a ReadWriter that can parse it
func Write(ext string, w io.Writer) (ReadWriter, error) {
	writer := New(ext)
	if writer == nil {
		return nil, fmt.Errorf("unknown extension %s", ext)
	}
	err := writer.Write(w)
	if err != nil {
		return nil, err
	}
	return writer, nil
}
