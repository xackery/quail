package raw

import (
	"bufio"
	"io"
	"strconv"
)

func (zon *Zon) WriteV4(w io.Writer) error {
	var err error
	writer := bufio.NewWriter(w)
	writer.WriteString("*NAME ")
	writer.WriteString(zon.Name)
	writer.WriteString("\n")
	writer.WriteString("*MINLNG ")
	writer.WriteString(strconv.Itoa(zon.V4Info.MinLng))
	writer.WriteString(" MAXLNG ")
	writer.WriteString(strconv.Itoa(zon.V4Info.MaxLng))
	writer.WriteString("\n")
	writer.WriteString("*MINLAT ")
	writer.WriteString(strconv.Itoa(zon.V4Info.MinLat))
	writer.WriteString(" MAXLAT ")
	writer.WriteString(strconv.Itoa(zon.V4Info.MaxLat))
	writer.WriteString("\n")
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
