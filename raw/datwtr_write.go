package raw

import (
	"bufio"
	"fmt"
	"io"
)

func (e *DatWtr) Write(w io.Writer) error {
	var err error
	writer := bufio.NewWriter(w)

	if len(e.Watersheets) == 0 && len(e.WatersheetEntries) == 0 {
		writer.WriteString("*BEGIN_WATERSHEETS\r\n")
		writer.WriteString("*END_WATERSHEETS\r\n")
	}

	// Write watersheets section
	if len(e.Watersheets) > 0 {
		writer.WriteString("*BEGIN_WATERSHEETS\r\n")
		for _, sheet := range e.Watersheets {
			writer.WriteString(fmt.Sprintf("\t*WATERSHEET %s\r\n", sheet.Tag))
			writer.WriteString(fmt.Sprintf("\t\t*MINX\t%0.6f\r\n", sheet.MinX))
			writer.WriteString(fmt.Sprintf("\t\t*MAXX\t%0.6f\r\n", sheet.MaxX))
			writer.WriteString(fmt.Sprintf("\t\t*MINY\t%0.6f\r\n", sheet.MinY))
			writer.WriteString(fmt.Sprintf("\t\t*MAXY\t%0.6f\r\n", sheet.MaxY))
			writer.WriteString(fmt.Sprintf("\t\t*ZHEIGHT\t%0.6f\r\n", sheet.ZHeight))
			writer.WriteString(fmt.Sprintf("\t\t*FRESNELBIAS\t%0.3f\r\n", sheet.FresnelBias))
			writer.WriteString(fmt.Sprintf("\t\t*FRESNELPOWER\t%0.3f\r\n", sheet.FresnelPower))
			writer.WriteString(fmt.Sprintf("\t\t*REFLECTIONAMOUNT\t%0.3f\r\n", sheet.ReflectionAmount))
			writer.WriteString(fmt.Sprintf("\t\t*UVSCALE\t%0.3f\r\n", sheet.UVScale))
			writer.WriteString(fmt.Sprintf("\t\t*REFLECTIONCOLOR\t%0.3f\t%0.3f\t%0.3f\t%0.3f\r\n",
				sheet.ReflectionColor[0], sheet.ReflectionColor[1], sheet.ReflectionColor[2], sheet.ReflectionColor[3]))
			writer.WriteString(fmt.Sprintf("\t\t*WATERCOLOR1\t%0.3f\t%0.3f\t%0.3f\t%0.3f\r\n",
				sheet.WaterColor1[0], sheet.WaterColor1[1], sheet.WaterColor1[2], sheet.WaterColor1[3]))
			writer.WriteString(fmt.Sprintf("\t\t*WATERCOLOR2\t%0.3f\t%0.3f\t%0.3f\t%0.3f\r\n",
				sheet.WaterColor2[0], sheet.WaterColor2[1], sheet.WaterColor2[2], sheet.WaterColor2[3]))
			writer.WriteString(fmt.Sprintf("\t\t*NORMALMAP\t%s\r\n", sheet.NormalMap))
			writer.WriteString(fmt.Sprintf("\t\t*ENVIRONMENTMAP\t%s\r\n", sheet.EnvironmentMap))
			writer.WriteString(fmt.Sprintf("\t*END_SHEET\r\n"))
		}

		writer.WriteString("*END_WATERSHEETS\r\n")
	}

	// Write watersheetdata section
	if len(e.WatersheetEntries) > 0 {
		writer.WriteString("*BEGIN_WATERSHEETDATA\r\n")
		for _, entry := range e.WatersheetEntries {
			writer.WriteString(fmt.Sprintf("\t*WATERSHEETDATA\r\n"))
			writer.WriteString(fmt.Sprintf("\t\t*INDEX\t%d\r\n", entry.Index))
			writer.WriteString(fmt.Sprintf("\t\t*FRESNELBIAS\t%0.3f\r\n", entry.FresnelBias))
			writer.WriteString(fmt.Sprintf("\t\t*FRESNELPOWER\t%0.3f\r\n", entry.FresnelPower))
			writer.WriteString(fmt.Sprintf("\t\t*REFLECTIONAMOUNT\t%0.3f\r\n", entry.ReflectionAmount))
			writer.WriteString(fmt.Sprintf("\t\t*UVSCALE\t%0.3f\r\n", entry.UVScale))
			writer.WriteString(fmt.Sprintf("\t\t*REFLECTIONCOLOR\t%0.3f\t%0.3f\t%0.3f\t%0.3f\r\n",
				entry.ReflectionColor[0], entry.ReflectionColor[1], entry.ReflectionColor[2], entry.ReflectionColor[3]))
			writer.WriteString(fmt.Sprintf("\t\t*WATERCOLOR1\t%0.3f\t%0.3f\t%0.3f\t%0.3f\r\n",
				entry.WaterColor1[0], entry.WaterColor1[1], entry.WaterColor1[2], entry.WaterColor1[3]))
			writer.WriteString(fmt.Sprintf("\t\t*WATERCOLOR2\t%0.3f\t%0.3f\t%0.3f\t%0.3f\r\n",
				entry.WaterColor2[0], entry.WaterColor2[1], entry.WaterColor2[2], entry.WaterColor2[3]))
			writer.WriteString(fmt.Sprintf("\t\t*NORMALMAP\t%s\r\n", entry.NormalMap))
			writer.WriteString(fmt.Sprintf("\t\t*ENVIRONMENTMAP\t%s\r\n", entry.EnvironmentMap))
			writer.WriteString(fmt.Sprintf("\t*ENDWATERSHEETDATA\r\n"))
		}

		writer.WriteString("*END_WATERSHEETDATA\r\n")
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("flush writer: %w", err)
	}

	return nil
}
