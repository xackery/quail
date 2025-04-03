package raw

import (
	"fmt"
	"io"
)

func (e *Dat) Write(w io.Writer) error {
	switch e.DatType {
	case DatTypeZon:
		return e.DatZon.Write(w)
	case DatTypeInvisibleWall:
		return e.DatIw.Write(w)
	case DatTypeWater:
		return e.DatWtr.Write(w)
	default:
		return fmt.Errorf("unknown dat type")
	}
}
