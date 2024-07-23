package wld

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (wld *Wld) ReadAscii(path string) error {
	wld.mu.Lock()
	defer wld.mu.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for lineNumber := 0; lineNumber < len(lines); lineNumber++ {
		line := lines[lineNumber]
		line = strings.TrimSpace(line)
		lineUpper := strings.ToUpper(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(lineUpper, "//") {
			continue
		}

		//if strings.HasPrefix(lineUpper, "INCLUDE") {
		//return fmt.Errorf("include not yet supported")
		//}

		if strings.HasPrefix(lineUpper, "DMSPRITEDEF2") {

			dmSprite := &DMSpriteDef2{}
			lineNumber, err = dmSprite.ReadAscii(lineNumber, lines)
			if err != nil {
				return err
			}
			wld.DMSpriteDef2s = append(wld.DMSpriteDef2s, dmSprite)
			continue
		}
	}

	return nil
}

func (d *DMSpriteDef2) ReadAscii(lineNumber int, lines []string) (int, error) {
	var err error

	scope := "DMSPRITEDEF2"
	for i := lineNumber; i < len(lines); i++ {
		line := lines[i]
		line = strings.TrimSpace(line)
		lineUpper := strings.ToUpper(line)
		if line == "" {
			continue
		}

		token := "ENDDMSPRITEDEF2"
		if strings.HasPrefix(lineUpper, token) {
			return i, nil
		}

		token = "TAG" // TAG "tag"
		if strings.HasPrefix(lineUpper, token) {
			lineIndex := strings.Index(line, token+" \"")
			if lineIndex == -1 {
				return i, fmt.Errorf("line %d %s: missing space quote after %s", i+1, scope, token)
			}
			line = line[lineIndex+5:]
			lineIndex = strings.Index(line, "\"")
			if lineIndex == -1 {
				return i, fmt.Errorf("line %d %s: missing quote after %s", i+1, scope, token)
			}
			d.Tag = line[:lineIndex]
			continue
		}

		token = "CENTEROFFSET" // CENTEROFFSET %0.2f %0.2f %0.2f
		if strings.HasPrefix(lineUpper, token) {
			lineIndex := strings.Index(line, token+" ")
			if lineIndex == -1 {
				return i, fmt.Errorf("line %d %s: missing space after %s", i+1, scope, token)
			}
			line = line[lineIndex+13:]
			splices := strings.Split(line, " ")
			if len(splices) < 3 {
				return i, fmt.Errorf("line %d %s %s: expected 3 slices", i+1, scope, token)
			}
			for j := 0; j < 3; j++ {
				val := float64(0)
				val, err = strconv.ParseFloat(splices[j], 32)
				if err != nil {
					return i, fmt.Errorf("line %d %s %s element %d: %w", i+1, scope, token, j, err)
				}
				d.CenterOffset[j] = float32(val)
			}
			continue
		}

		token = "NUMVERTICES" // NUMVERTICES %d
		if strings.HasPrefix(lineUpper, token) {
			lineIndex := strings.Index(line, token+" ")
			if lineIndex == -1 {
				return i, fmt.Errorf("line %d %s: missing space after %s", i+1, scope, token)
			}
			line = line[lineIndex+12:]
			numVertices := 0
			_, err := fmt.Sscanf(line, "%d", &numVertices)
			if err != nil {
				return i, fmt.Errorf("line %d %s %s: %w", i+1, scope, token, err)
			}

			for j := 0; j < numVertices; j++ {
				i++
				if i >= len(lines) {
					return i, fmt.Errorf("line %d %s %s %d: unexpected end of file", i+1, scope, token, j)
				}
				line = lines[i]
				line = strings.TrimSpace(line)
				lineUpper = strings.ToUpper(line)
				if line == "" {
					return i, fmt.Errorf("line %d %s %s %d: unexpected empty line", i+1, scope, token, j)
				}
				if !strings.HasPrefix(lineUpper, "XYZ") {
					return i, fmt.Errorf("line %d %s %s %d: expected XYZ", i+1, scope, token, j)
				}
				// XYZ %0.2f %0.2f %0.2f
				lineIndex := strings.Index(line, "XYZ ")
				if lineIndex == -1 {
					return i, fmt.Errorf("line %d %s %s %d: missing XYZ", i+1, scope, token, j)
				}
				line = line[lineIndex+5:]
				d.Vertices = append(d.Vertices, [3]float32{})
				splices := strings.Split(line, " ")
				if len(splices) < 3 {
					return i, fmt.Errorf("line %d %s %s %d: expected 3 slices", i+1, scope, token, j)
				}
				for k := 0; k < 3; k++ {
					val := float64(0)
					val, err = strconv.ParseFloat(splices[k], 32)
					if err != nil {
						return i, fmt.Errorf("line %d %s %s %d element %d: %w", i+1, scope, token, j, k, err)
					}
					d.Vertices[j][k] = float32(val)
				}
			}

			continue
		}

		token = "NUMUVS" // NUMUVS %d
		if strings.HasPrefix(lineUpper, token) {
			lineIndex := strings.Index(line, token+" ")
			if lineIndex == -1 {
				return i, fmt.Errorf("line %d %s: missing space after %s", i+1, scope, token)
			}
			line = line[lineIndex+7:]
			numUVs := 0
			_, err := fmt.Sscanf(line, "%d", &numUVs)
			if err != nil {
				return i, fmt.Errorf("line %d %s %s: %w", i+1, scope, token, err)
			}

			for j := 0; j < numUVs; j++ {
				i++
				if i >= len(lines) {
					return i, fmt.Errorf("line %d %s %s %d: unexpected end of file", i+1, scope, token, j)
				}
				line = lines[i]
				line = strings.TrimSpace(line)
				lineUpper = strings.ToUpper(line)
				if line == "" {
					return i, fmt.Errorf("line %d %s %s %d: unexpected empty line", i+1, scope, token, j)
				}
				if !strings.HasPrefix(lineUpper, "UV") {
					return i, fmt.Errorf("line %d %s %s %d: expected UV", i+1, scope, token, j)
				}
				// UV %0.2f %0.2f
				lineIndex := strings.Index(line, "UV ")
				if lineIndex == -1 {
					return i, fmt.Errorf("line %d %s %s %d: missing UV", i+1, scope, token, j)
				}
				line = line[lineIndex+4:]
				if line[0] == ' ' {
					line = line[1:]
				}
				d.UVs = append(d.UVs, [2]float32{})
				splices := strings.Split(line, " ")
				if len(splices) < 2 {
					return i, fmt.Errorf("line %d %s %s %d: expected 2 slices", i+1, scope, token, j)
				}
				for k := 0; k < 2; k++ {
					val := float64(0)
					splice := splices[k]
					splice = strings.ReplaceAll(splice, ",", "")
					val, err = strconv.ParseFloat(splice, 32)
					if err != nil {
						return i, fmt.Errorf("line %d %s %s %d element %d: %w", i+1, scope, token, j, k, err)
					}
					d.UVs[j][k] = float32(val)
				}
			}
			continue
		}

	}

	return len(lines), fmt.Errorf("unexpected end of file, expecetd ENDDMSPRITEDEF2")
}
