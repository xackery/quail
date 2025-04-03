package raw

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type DatWtr struct {
	MetaFileName      string
	Watersheets       []*DatWtrSheet
	WatersheetEntries []*DatWtrSheetEntry
}

type DatWtrSheet struct {
	Tag              string
	MinX             float64
	MaxX             float64
	MinY             float64
	MaxY             float64
	ZHeight          float64
	FresnelBias      float64
	FresnelPower     float64
	ReflectionAmount float64
	UVScale          float64
	ReflectionColor  [4]float64
	WaterColor1      [4]float64
	WaterColor2      [4]float64
	NormalMap        string
	EnvironmentMap   string
}

type DatWtrSheetEntry struct {
	Index            int
	FresnelBias      float64
	FresnelPower     float64
	ReflectionAmount float64
	UVScale          float64
	ReflectionColor  [4]float64
	WaterColor1      [4]float64
	WaterColor2      [4]float64
	NormalMap        string
	EnvironmentMap   string
}

func (e *DatWtr) Identity() string {
	return "dat"
}

// Decode reads a v4 zone dat file
// https://github.com/EQEmu/zone-utilities/blob/master/src/common/eqg_v4_loader.cpp#L115
func (e *DatWtr) Read(r io.ReadSeeker) error {
	var err error

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		records := strings.Fields(scanner.Text())
		if len(records) < 1 {
			continue
		}

		if records[0] == "*WATERSHEET" {
			if len(records) < 2 {
				return fmt.Errorf("parse %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet := &DatWtrSheet{
				Tag: records[1],
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("minx: expected 1 value, got %d", len(records)-1)
			}
			sheet.MinX, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}

			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("maxx: expected 1 value, got %d", len(records)-1)
			}
			sheet.MaxX, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}

			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("miny: expected 1 value, got %d", len(records)-1)
			}

			sheet.MinY, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}

			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("maxy: expected 1 value, got %d", len(records)-1)
			}

			sheet.MaxY, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("zheight: expected 1 value, got %d", len(records)-1)
			}
			sheet.ZHeight, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("fresnelbias: expected 1 value, got %d", len(records)-1)
			}

			sheet.FresnelBias, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("fresnelpower: expected 1 value, got %d", len(records)-1)
			}

			sheet.FresnelPower, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("reflectionamount: expected 1 value, got %d", len(records)-1)
			}

			sheet.ReflectionAmount, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("uvscale: expected 1 value, got %d", len(records)-1)
			}

			sheet.UVScale, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			sheet.UVScale, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}

			records = nextLineFields(scanner)
			if len(records) < 5 {
				return fmt.Errorf("parse %s %s: expected 4 values, got %d", records[0], records[1], len(records)-1)
			}
			for i := 0; i < 4; i++ {
				sheet.ReflectionColor[i], err = strconv.ParseFloat(records[i+1], 32)
				if err != nil {
					return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
				}
			}
			records = nextLineFields(scanner)
			if len(records) < 5 {
				return fmt.Errorf("parse %s %s: expected 4 values, got %d", records[0], records[1], len(records)-1)
			}

			for i := 0; i < 4; i++ {
				sheet.WaterColor1[i], err = strconv.ParseFloat(records[i+1], 32)
				if err != nil {
					return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
				}
			}

			records = nextLineFields(scanner)
			if len(records) < 5 {
				return fmt.Errorf("parse %s %s: expected 4 values, got %d", records[0], records[1], len(records)-1)
			}

			for i := 0; i < 4; i++ {
				sheet.WaterColor2[i], err = strconv.ParseFloat(records[i+1], 32)
				if err != nil {
					return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
				}
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("parse %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}

			sheet.NormalMap = records[1]

			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("parse %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet.EnvironmentMap = records[1]
			records = nextLineFields(scanner)
			if len(records) < 1 {
				return fmt.Errorf("parse %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			if !strings.HasPrefix("*END_SHEET", records[0]) {
				return fmt.Errorf("parse %s %s: expected *END_SHEET, got %s", records[0], records[1], records[0])
			}
			e.Watersheets = append(e.Watersheets, sheet)
			continue
		}

		if records[0] == "*WATERSHEETDATA" {
			sheet := &DatWtrSheetEntry{}

			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("parse %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet.Index, err = strconv.Atoi(records[1])
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}

			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("fresnelbias %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet.FresnelBias, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}

			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("fresnelpower %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet.FresnelPower, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("reflectionamount %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet.ReflectionAmount, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("uvscale %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet.UVScale, err = strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
			}
			records = nextLineFields(scanner)
			if len(records) < 5 {
				return fmt.Errorf("reflectionColor: expected 4 values, got %d", len(records)-1)
			}
			for i := 0; i < 4; i++ {
				sheet.ReflectionColor[i], err = strconv.ParseFloat(records[i+1], 32)
				if err != nil {
					return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
				}
			}
			records = nextLineFields(scanner)
			if len(records) < 5 {
				return fmt.Errorf("waterColor1: expected 4 values, got %d", len(records)-1)
			}
			for i := 0; i < 4; i++ {
				sheet.WaterColor1[i], err = strconv.ParseFloat(records[i+1], 32)
				if err != nil {
					return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
				}
			}
			records = nextLineFields(scanner)
			if len(records) < 5 {
				return fmt.Errorf("waterColor2: expected 4 values, got %d", len(records)-1)
			}
			for i := 0; i < 4; i++ {
				sheet.WaterColor2[i], err = strconv.ParseFloat(records[i+1], 32)
				if err != nil {
					return fmt.Errorf("parse %s %s: %w", records[0], records[1], err)
				}
			}
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("normalmap %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet.NormalMap = records[1]
			records = nextLineFields(scanner)
			if len(records) < 2 {
				return fmt.Errorf("environmentmap %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			sheet.EnvironmentMap = records[1]
			records = nextLineFields(scanner)
			if len(records) < 1 {
				return fmt.Errorf("parse %s %s: expected 1 value, got %d", records[0], records[1], len(records)-1)
			}
			if records[0] != "*ENDWATERSHEETDATA" {
				return fmt.Errorf("expected *END_WATERSHEETDATA, got %s", records[0])
			}

			e.WatersheetEntries = append(e.WatersheetEntries, sheet)
			continue
		}
	}
	return nil
}

// SetName sets the name of the file
func (e *DatWtr) SetFileName(name string) {
	e.MetaFileName = name
}

func (e *DatWtr) FileName() string {
	return e.MetaFileName
}

func nextLineFields(scanner *bufio.Scanner) []string {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil
		}
		return nil
	}
	return strings.Fields(scanner.Text())
}
