package raw

import (
	"bytes"
	"fmt"
	"io"
)

type DatType int

const (
	DatTypeUnknown DatType = iota
	DatTypeZon
	DatTypeInvisibleWall
	DatTypeWater
	DatTypeFloraExclude
)

// Dat is a generic type that turns into a derived dat type
type Dat struct {
	DatType DatType
	DatIw   *DatIw
	DatZon  *DatZon
	DatWtr  *DatWtr
	DatFe   *DatFe
}

func (dat *Dat) Identity() string {
	switch dat.DatType {
	case DatTypeZon:
		return "datzon"
	case DatTypeInvisibleWall:
		return "datiw"
	case DatTypeWater:
		return "datwtr"
	default:
		return "dat"
	}
}

func (e *Dat) Read(r io.ReadSeeker) error {

	header := make([]byte, 4)
	_, err := r.Read(header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	r.Seek(0, io.SeekStart)

	switch string(header) {
	case "*BEG":
		e.DatType = DatTypeFloraExclude
		e.DatFe = &DatFe{}
		return e.DatFe.Read(r)

	//case DatTypeUnknown:
	//	return fmt.Errorf("unknown dat type")
	//case DatTypeZon:
	//	return e.DatZon.Read(r)
	//case DatTypeInvisibleWall:
	//	return e.DatIw.Read(r)
	//case DatTypeWater:
	//	return e.DatWtr.Read(r)
	default:
		if bytes.Equal(header[1:], []byte{0x0, 0x0, 0x0}) {
			e.DatType = DatTypeInvisibleWall
			e.DatIw = &DatIw{}
			return e.DatIw.Read(r)
		}
		fmt.Printf("unknown dat type: type: %x (%s)\n", header, header)
		e.DatType = DatTypeZon
		e.DatZon = &DatZon{}
		return e.DatZon.Read(r)
		//		return fmt.Errorf("unknown dat type: %x (%s)", header, header)
	}
}

func (e *Dat) SetType(datType DatType) {
	e.DatType = datType
}

// SetName sets the name of the file
func (e *Dat) SetFileName(name string) {
	switch e.DatType {
	case DatTypeZon:
		e.DatZon.SetFileName(name)
	case DatTypeInvisibleWall:
		e.DatIw.SetFileName(name)
	case DatTypeWater:
		e.DatWtr.SetFileName(name)
	}
}

func (e *Dat) FileName() string {
	switch e.DatType {
	case DatTypeZon:
		return e.DatZon.FileName()
	case DatTypeInvisibleWall:
		return e.DatIw.FileName()
	case DatTypeWater:
		return e.DatWtr.FileName()
	default:
		return ""
	}
}
