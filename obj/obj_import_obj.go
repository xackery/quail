package obj

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

func importObjFile(obj *ObjData, objPath string) error {
	var lastMaterial *common.Material

	objCache := &objCache{
		vertexLookup: make(map[string]int),
	}

	lineNumber := 0
	ro, err := os.Open(objPath)
	if err != nil {
		return err
	}
	defer ro.Close()

	scanner := bufio.NewScanner(ro)
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.HasPrefix(line, "usemtl ") {
			line = strings.TrimPrefix(line, "usemtl ")
			lastMaterial = materialByName(line, obj)
			if lastMaterial == nil {
				lastMaterial = &common.Material{
					Name: line,
				}
				obj.Materials = append(obj.Materials, lastMaterial)
				fmt.Printf("warning: obj line %d refers to material %s, which isn't declared\n", lineNumber, line)
			}
			continue
		}

		if strings.HasPrefix(line, "f ") {
			line = strings.TrimPrefix(line, "f ")
			records := strings.Split(line, " ")
			if len(records) != 3 {
				return fmt.Errorf("line %d has %d records, expected 3", lineNumber, len(records))
			}
			faces := []int{}
			for i, record := range records {
				entries := strings.Split(record, "/")
				if len(entries) != 3 {
					return fmt.Errorf("line %d has %d entries, expected 3", lineNumber, len(entries))
				}
				for j, entry := range entries {
					val, err := strconv.Atoi(entry)
					if err != nil {
						return fmt.Errorf("line %d parse face index %d %d %s: %w", lineNumber, i, j, entry, err)
					}
					faces = append(faces, val)
				}
			}
			index1, err := face(faces[0], faces[1], faces[2], objCache, obj)
			if err != nil {
				return fmt.Errorf("face 1: line %d: %w", lineNumber, err)
			}
			index2, err := face(faces[3], faces[4], faces[5], objCache, obj)
			if err != nil {
				return fmt.Errorf("face 2: line %d: %w", lineNumber, err)
			}
			index3, err := face(faces[6], faces[7], faces[8], objCache, obj)
			if err != nil {
				return fmt.Errorf("face 3: line %d: %w", lineNumber, err)
			}
			obj.Triangles = append(obj.Triangles, &common.Triangle{
				Index:        math32.Vector3{X: float32(index1), Y: float32(index2), Z: float32(index3)},
				MaterialName: lastMaterial.Name,
				Flag:         lastMaterial.Flag,
			})
		}
		if strings.HasPrefix(line, "v ") {
			position := math32.Vector3{}
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
			objCache.vertices = append(objCache.vertices, position)
		}
		if strings.HasPrefix(line, "vn ") {
			normal := math32.Vector3{}
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
			objCache.normals = append(objCache.normals, normal)
		}
		if strings.HasPrefix(line, "vt ") {
			uv := math32.Vector2{}
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
			objCache.uvs = append(objCache.uvs, uv)
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read obj %s: %w", objPath, err)
	}

	//for i := 0; i < len(positions); i++ {
	//obj.Vertices = append(obj.Vertices, )
	//err = e.AddVertex(positions[i], rotations[i], uvs[i])

	//}
	return nil
}

func face(v, t, n int, objCache *objCache, obj *ObjData) (int, error) {

	index, ok := objCache.vertexLookup[fmt.Sprintf("%d/%d/%d", v, t, n)]
	if ok {
		return index, nil
	}

	v -= 1
	t -= 1
	n -= 1

	if len(objCache.vertices) <= v {
		return 0, fmt.Errorf("wanted vertex %d, but objCache is %d", v, len(objCache.vertices))
	}
	vert := objCache.vertices[v]

	if len(objCache.normals) <= n {
		return 0, fmt.Errorf("wanted normal %d, but objCache is %d", n, len(objCache.normals))
	}
	norm := objCache.normals[n]

	uv := math32.Vector2{}
	if t != 0 && len(objCache.uvs) >= t {
		uv = objCache.uvs[t]
	}

	obj.Vertices = append(obj.Vertices, &common.Vertex{
		Position: math32.Vector3{X: vert.X, Y: vert.Y, Z: vert.Z},
		Normal:   math32.Vector3{X: norm.X, Y: norm.Y, Z: norm.Z},
		Uv:       math32.Vector2{X: uv.X, Y: uv.Y},
	})

	index = len(obj.Vertices) - 1
	objCache.vertexLookup[fmt.Sprintf("%d/%d/%d", v, t, n)] = index
	return index, nil
}
