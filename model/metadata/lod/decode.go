package lod

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func (e *LOD) Decode(r io.ReadSeeker) error {

	// read all of r as a string
	buf := make([]byte, 1024)
	var str string
	for {
		n, err := r.Read(buf)
		if err != nil {
			break
		}
		str += string(buf[:n])
	}

	// split by newline
	lines := strings.Split(str, "\n")
	if strings.TrimSpace(lines[0]) != "EQLOD" {
		return fmt.Errorf("header does not match EQLOD, got %s", lines[0])
	}

	e.lods = make([]*LODEntry, 0)
	// parse each line
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// split into commas
		records := strings.Split(line, ",")
		if len(records) < 3 {
			return fmt.Errorf("line %d expected 3 entries, got %d", i, len(records))
		}
		le := &LODEntry{
			Category:   records[0],
			ObjectName: records[1],
		}
		val, err := strconv.Atoi(records[2])
		if err != nil {
			return fmt.Errorf("line %d lod %d is not a number", i, val)
		}
		le.Distance = val
		e.lods = append(e.lods, le)
	}

	return nil
}
