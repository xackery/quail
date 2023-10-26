package zon

import (
	"bufio"
	"io"
	"strconv"

	"github.com/xackery/quail/common"
)

func EncodeV4(zone *common.Zone, w io.Writer) error {
	var err error
	writer := bufio.NewWriter(w)
	writer.WriteString("*NAME ")
	writer.WriteString(zone.Header.Name)
	writer.WriteString("\n")
	writer.WriteString("*MINLNG ")
	writer.WriteString(strconv.Itoa(zone.V4Info.MinLng))
	writer.WriteString(" MAXLNG ")
	writer.WriteString(strconv.Itoa(zone.V4Info.MaxLng))
	writer.WriteString("\n")
	writer.WriteString("*MINLAT ")
	writer.WriteString(strconv.Itoa(zone.V4Info.MinLat))
	writer.WriteString(" MAXLAT ")
	writer.WriteString(strconv.Itoa(zone.V4Info.MaxLat))
	writer.WriteString("\n")
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
