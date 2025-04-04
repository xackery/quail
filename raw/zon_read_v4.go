package raw

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Decode reads a v4 ZON file
// https://github.com/EQEmu/zone-utilities/blob/master/src/common/eqg_v4_loader.cpp#L736
func (zon *Zon) ReadV4(r io.ReadSeeker) error {
	var err error
	scanner := bufio.NewScanner(r)
	lineNumber := 1
	for scanner.Scan() {

		line := scanner.Text()
		if strings.HasPrefix(line, "*NAME") {
			zon.V4Info.Name = strings.TrimPrefix(line, "*NAME ")
			continue
		}
		if strings.HasPrefix(line, "*MINLNG") {
			vals := strings.Split(line, " ")
			if len(vals) < 4 {
				return fmt.Errorf("line %d: MINLNG: not enough values", lineNumber)
			}

			zon.V4Info.MinLng, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: MINLNG: %w", lineNumber, err)
			}
			zon.V4Info.MaxLng, err = strconv.Atoi(vals[3])
			if err != nil {
				return fmt.Errorf("line %d: MAXLNG: %w", lineNumber, err)
			}
			continue
		}

		if strings.HasPrefix(line, "*MINLAT") {
			vals := strings.Split(line, " ")
			if len(vals) < 4 {
				return fmt.Errorf("line %d: MINLAT: not enough values", lineNumber)
			}

			zon.V4Info.MinLat, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: MINLAT: %w", lineNumber, err)
			}
			zon.V4Info.MaxLat, err = strconv.Atoi(vals[3])
			if err != nil {
				return fmt.Errorf("line %d: MAXLAT: %w", lineNumber, err)
			}
			continue
		}

		if strings.HasPrefix(line, "*MIN_EXTENTS") {
			val := strings.TrimPrefix(line, "*MIN_EXTENTS ")
			vals := strings.Split(val, " ")
			fval := float64(0)
			fval, err = strconv.ParseFloat(vals[0], 32)
			if err != nil {
				return fmt.Errorf("line %d: MIN_EXTENTS X: %w", lineNumber, err)
			}
			zon.V4Info.MinExtents[0] = float32(fval)
			fval, err = strconv.ParseFloat(vals[1], 32)
			if err != nil {
				return fmt.Errorf("line %d: MIN_EXTENTS Y: %w", lineNumber, err)
			}
			zon.V4Info.MinExtents[1] = float32(fval)
			fval, err = strconv.ParseFloat(vals[2], 32)
			if err != nil {
				return fmt.Errorf("line %d: MIN_EXTENTS Z: %w", lineNumber, err)
			}
			zon.V4Info.MinExtents[2] = float32(fval)
			continue
		}

		if strings.HasPrefix(line, "*MAX_EXTENTS") {
			val := strings.TrimPrefix(line, "*MAX_EXTENTS ")
			vals := strings.Split(val, " ")
			fval := float64(0)
			fval, err = strconv.ParseFloat(vals[0], 32)
			if err != nil {
				return fmt.Errorf("line %d: MAX_EXTENTS X: %w", lineNumber, err)
			}
			zon.V4Info.MaxExtents[0] = float32(fval)
			fval, err = strconv.ParseFloat(vals[1], 32)
			if err != nil {
				return fmt.Errorf("line %d: MAX_EXTENTS Y: %w", lineNumber, err)
			}
			zon.V4Info.MaxExtents[1] = float32(fval)
			fval, err = strconv.ParseFloat(vals[2], 32)
			if err != nil {
				return fmt.Errorf("line %d: MAX_EXTENTS Z: %w", lineNumber, err)
			}
			zon.V4Info.MaxExtents[2] = float32(fval)
			continue
		}

		if strings.HasPrefix(line, "*UNITSPERVERT") {
			vals := strings.Split(line, " ")
			fval := float64(0)
			fval, err = strconv.ParseFloat(vals[1], 32)
			if err != nil {
				return fmt.Errorf("line %d: UNITSPERVERT: %w", lineNumber, err)
			}
			zon.V4Info.UnitsPerVert = float32(fval)
			continue
		}

		if strings.HasPrefix(line, "*QUADSPERTILE") {
			vals := strings.Split(line, " ")
			zon.V4Info.QuadsPerTile, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: QUADSPERTILE: %w", lineNumber, err)
			}
			continue
		}

		if strings.HasPrefix(line, "*COVERMAPINPUTSIZE") {
			vals := strings.Split(line, " ")
			zon.V4Info.CoverMapInputSize, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: COVERMAPINPUTSIZE: %w", lineNumber, err)
			}
			continue
		}

		if strings.HasPrefix(line, "*LAYERINGMAPINPUTSIZE") {
			vals := strings.Split(line, " ")
			zon.V4Info.LayeringMapInputSize, err = strconv.Atoi(vals[1])
			if err != nil {
				return fmt.Errorf("line %d: LAYERINGMAPINPUTSIZE: %w", lineNumber, err)
			}
			continue
		}

		lineNumber++
	}

	zon.Version = 4

	return nil
}
