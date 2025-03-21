package raw

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// LOD is level of detail information
// Typical usaeg is like so:
/*
EQLOD
LOD,OBJ_FIREPIT_STMFT,150
LOD,OBJ_FIREPIT_STMFT_LOD1,250
LOD,OBJ_FIREPIT_STMFT_LOD2,400
LOD,OBJ_FIREPIT_STMFT_LOD3,1000
*/
type Lod struct {
	MetaFileName string
	Entries      []*LodEntry
}

type LodEntry struct {
	Category   string
	ObjectName string
	Distance   float32
}

// Identity returns the type of the struct
func (lod *Lod) Identity() string {
	return "lod"
}

func (lod *Lod) Read(r io.ReadSeeker) error {
	var err error
	line := ""

	lineNumber := 0

	lod.Entries = []*LodEntry{}

	hasHeader := false
	isNewline := false
	for {
		if isNewline {
			line = ""
			isNewline = false
		}
		chunk := make([]byte, 1)
		_, err = r.Read(chunk)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			break
		}
		line += string(chunk)
		if chunk[0] != '\n' {
			continue
		}
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")

		lineNumber++
		isNewline = true
		if !hasHeader {
			if line != "EQLOD" {
				return fmt.Errorf("invalid header %s, wanted EQLOD", line)
			}
			hasHeader = true
			continue
		}

		records := strings.Split(line, ",")
		if len(records) == 1 {
			continue
		}
		entry := &LodEntry{
			Distance: 0,
		}
		if len(records) > 2 {
			val, err := strconv.ParseFloat(records[2], 32)
			if err != nil {
				return fmt.Errorf("line %d lod %0.3f (%s) is not a number", lineNumber, val, records[2])
			}
			entry.Distance = float32(val)
		}
		if len(records) < 2 {
			return fmt.Errorf("line %d expected at least 2 entries, got %d", lineNumber, len(records))
		}

		entry.Category = records[0]
		entry.ObjectName = records[1]

		lod.Entries = append(lod.Entries, entry)
	}

	return nil
}

// SetFileName sets the name of the file
func (lod *Lod) SetFileName(name string) {
	lod.MetaFileName = name
}

// FileName returns the name of the file
func (lod *Lod) FileName() string {
	return lod.MetaFileName
}
