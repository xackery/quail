package raw

import (
	"fmt"
	"io"
)

type DatType int

const (
	DatTypeUnknown DatType = iota
	DatTypeZon
	DatTypeInvisibleWall
	DatTypeWater
)

// Dat is a generic type that turns into a derived dat type
type Dat struct {
	DatType DatType
	DatIw   *DatIw
	DatZon  *DatZon
	DatWtr  *DatWtr
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
	switch e.DatType {
	case DatTypeUnknown:
		return fmt.Errorf("unknown dat type")
	case DatTypeZon:
		return e.DatZon.Read(r)
	case DatTypeInvisibleWall:
		return e.DatIw.Read(r)
	case DatTypeWater:
		return e.DatWtr.Read(r)
	default:
		return fmt.Errorf("unknown dat type")
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
