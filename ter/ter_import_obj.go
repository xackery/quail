package ter

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/g3n/engine/math32"
)

func (e *TER) ImportObj(objPath string, mtlPath string) error {
	var err error
	rm, err := os.Open(mtlPath)
	if err != nil {
		return err
	}
	defer rm.Close()

	scanner := bufio.NewScanner(rm)
	lineNumber := 0
	lastMaterial := ""
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.HasPrefix(line, "newmtl") {
			if len(line) < 8 {
				return fmt.Errorf("%s line %d: newmtl is too short", mtlPath, lineNumber)
			}
			lastMaterial = line[7:]
			err = e.AddMaterial(lastMaterial, "")
			if err != nil {
				return fmt.Errorf("addMaterial line %d: %w", lineNumber, err)
			}
			continue
		}
		if strings.HasPrefix(line, "map_Kd") {
			if lastMaterial == "" {
				return fmt.Errorf("map_kd line %d found before material definition", lineNumber)
			}

			err = e.AddMaterialProperty(lastMaterial, "e_TextureDiffuse0", 2, 0, 0)
			if err != nil {
				return fmt.Errorf("addMaterialProperty map_Kd line %d: %w", lineNumber, err)
			}
		}
		if strings.HasPrefix(line, "map_Bump") {
			if lastMaterial == "" {
				return fmt.Errorf("map_Bump line %d found before material definition", lineNumber)
			}

			err = e.AddMaterialProperty(lastMaterial, "e_TextureNormal0", 2, 0, 0)
			if err != nil {
				return fmt.Errorf("addMaterialProperty map_Bump line %d: %w", lineNumber, err)
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read mtl %s: %w", objPath, err)
	}

	positions := []*math32.Vector3{}
	normals := []*math32.Vector3{}
	uvs := []*math32.Vector2{}

	ro, err := os.Open(objPath)
	if err != nil {
		return err
	}
	defer rm.Close()

	fReg := regexp.MustCompile(`f ([0-9]+)\/([0-9]+)\/([0-9]+) ([0-9]+)\/([0-9]+)\/([0-9]+) ([0-9]+)\/([0-9]+)\/([0-9]+)`)

	scanner = bufio.NewScanner(ro)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "usemtl ") {
			line = strings.TrimPrefix(line, "usemtl ")
			isKnownMaterial := false
			for _, m := range e.materials {
				if m.Name != line {
					continue
				}
				isKnownMaterial = true
				lastMaterial = m.Name
				break
			}
			if !isKnownMaterial {
				return fmt.Errorf("obj line %d refers to material %s, which isn't in mtl", lineNumber, line)
			}
			continue
		}
		matches := fReg.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 {
			if len(matches[0]) < 9 {
				return fmt.Errorf("line %d pattern face should be 9, got %d", len(matches[0]))
			}
			faces := []float32{}
			for i := 0; i < 9; i++ {
				val, err := strconv.ParseFloat(matches[0][i], 32)
				if err != nil {
					return fmt.Errorf("line %d parse face index %d %s: %w", lineNumber, i, matches[0][i], err)
				}
				faces = append(faces, float32(val))
			}
			e.AddTriangle(index math32.Vector3, materialName string, flag uint32)
		}
		if strings.HasPrefix(line, "v ") {
			position := &math32.Vector3{}
			records := strings.Split(line, " ")
			if len(records) != 4 {
				return fmt.Errorf("line %d split v expected 4, got %d", lineNumber, len(records))
			}
			val, err := strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("line %d parse x %s failed: %w", lineNumber, records[1], err)
			}
			position.X = float32(val)
			val, err = strconv.ParseFloat(records[2], 32)
			if err != nil {
				return fmt.Errorf("line %d parse y %s failed: %w", lineNumber, records[2], err)
			}
			position.Y = float32(val)
			val, err = strconv.ParseFloat(records[3], 32)
			if err != nil {
				return fmt.Errorf("line %d parse z %s failed: %w", lineNumber, records[3], err)
			}
			position.Z = float32(val)
			positions = append(positions, position)
		}
		if strings.HasPrefix(line, "vn ") {
			normal := &math32.Vector3{}
			records := strings.Split(line, " ")
			if len(records) != 4 {
				return fmt.Errorf("line %d split v expected 4, got %d", lineNumber, len(records))
			}
			val, err := strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("line %d parse normal x %s failed: %w", lineNumber, records[1], err)
			}
			normal.X = float32(val)
			val, err = strconv.ParseFloat(records[2], 32)
			if err != nil {
				return fmt.Errorf("line %d parse normal y %s failed: %w", lineNumber, records[2], err)
			}
			normal.Y = float32(val)
			val, err = strconv.ParseFloat(records[3], 32)
			if err != nil {
				return fmt.Errorf("line %d parse normal z %s failed: %w", lineNumber, records[3], err)
			}
			normal.Z = float32(val)
			normals = append(normals, normal)
		}
		if strings.HasPrefix(line, "vt ") {
			uv := &math32.Vector2{}
			records := strings.Split(line, " ")
			if len(records) != 3 {
				return fmt.Errorf("line %d split vt expected 3, got %d", lineNumber, len(records))
			}
			val, err := strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("line %d parse x %s failed: %w", lineNumber, records[1], err)
			}
			uv.X = float32(val)
			val, err = strconv.ParseFloat(records[2], 32)
			if err != nil {
				return fmt.Errorf("line %d parse y %s failed: %w", lineNumber, records[2], err)
			}
			uv.Y = float32(val)
			uvs = append(uvs, uv)
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read obj %s: %w", objPath, err)
	}

	for i := 0; i < len(positions); i++ {

		err = e.AddVertex(positions[i], rotations[i], uvs[i])
		if err != nil {
			return fmt.Errorf("addVertex %d: %w", i, err)
		}

	}

	return nil
}
