package raw

import (
	"bufio"
	"fmt"
	"io"
)

func (zon *Zon) WriteV4(srcW io.Writer) error {
	var err error
	w := bufio.NewWriter(srcW)
	w.WriteString("EQTZP\r\n")

	w.WriteString(fmt.Sprintf("*NAME %s\r\n", zon.V4Info.Name))
	w.WriteString(fmt.Sprintf("*MINLNG %d *MAXLNG %d\r\n", zon.V4Info.MinLng, zon.V4Info.MaxLng))
	w.WriteString(fmt.Sprintf("*MINLAT %d *MAXLAT %d\r\n", zon.V4Info.MinLat, zon.V4Info.MaxLat))

	w.WriteString(fmt.Sprintf("*MIN_EXTENTS %0.3f %0.3f %0.3f\r\n",
		zon.V4Info.MinExtents[0],
		zon.V4Info.MinExtents[1],
		zon.V4Info.MinExtents[2]))

	w.WriteString(fmt.Sprintf("*MAX_EXTENTS %0.3f %0.3f %0.3f\r\n",
		zon.V4Info.MaxExtents[0],
		zon.V4Info.MaxExtents[1],
		zon.V4Info.MaxExtents[2]))

	w.WriteString(fmt.Sprintf("*UNITSPERVERT %0.2f\r\n", zon.V4Info.UnitsPerVert))
	w.WriteString(fmt.Sprintf("*QUADSPERTILE %d\r\n", zon.V4Info.QuadsPerTile))
	w.WriteString(fmt.Sprintf("*COVERMAPINPUTSIZE %d\r\n", zon.V4Info.CoverMapInputSize))
	w.WriteString(fmt.Sprintf("*LAYERINGMAPINPUTSIZE %d\r\n", zon.V4Info.LayeringMapInputSize))
	w.WriteString(fmt.Sprintf("*VERSION %d\r\n", zon.Version))
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}
