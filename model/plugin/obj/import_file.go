package obj

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xackery/quail/common"
)

func importFile(req *ObjRequest) error {
	var lastMaterial *common.Material

	objCache := &objCache{
		vertexLookup: make(map[string]int),
	}

	lineNumber := 0
	ro, err := os.Open(req.ObjPath)
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
			lastMaterial = materialByName(line, req.Data)
			if lastMaterial == nil {
				lastMaterial = &common.Material{
					Name: line,
				}
				req.Data.Materials = append(req.Data.Materials, lastMaterial)
				fmt.Printf("warning: obj line %d refers to material '%s', which is not declared\n", lineNumber, line)
			}
			continue
		}

		if strings.HasPrefix(line, "f ") {
			line = strings.TrimPrefix(line, "f ")
			records := strings.Split(line, " ")
			if len(records) != 3 {
				return fmt.Errorf("line %d has %d records, expected 3", lineNumber, len(records))
			}
			triangles := []int{}
			for i, record := range records {
				entries := strings.Split(record, "/")
				if len(entries) != 3 {
					return fmt.Errorf("line %d has %d entries, expected 3", lineNumber, len(entries))
				}
				for j, entry := range entries {
					val, err := strconv.Atoi(entry)
					if err != nil {
						return fmt.Errorf("line %d parse triangle index %d %d %s: %w", lineNumber, i, j, entry, err)
					}
					triangles = append(triangles, val)
				}
			}
			index1, err := triangle(triangles[0], triangles[1], triangles[2], objCache, req.Data)
			if err != nil {
				return fmt.Errorf("triangle 1: line %d: %w", lineNumber, err)
			}
			index2, err := triangle(triangles[3], triangles[4], triangles[5], objCache, req.Data)
			if err != nil {
				return fmt.Errorf("triangle 2: line %d: %w", lineNumber, err)
			}
			index3, err := triangle(triangles[6], triangles[7], triangles[8], objCache, req.Data)
			if err != nil {
				return fmt.Errorf("triangle 3: line %d: %w", lineNumber, err)
			}
			req.Data.Triangles = append(req.Data.Triangles, &common.Triangle{
				Index:        [3]uint32{uint32(index1), uint32(index2), uint32(index3)},
				MaterialName: lastMaterial.Name,
				Flag:         lastMaterial.Flag,
			})
		}

		if strings.HasPrefix(line, "v ") {
			fmt.Println(line)
			position := [3]float32{}
			records := strings.Split(line, " ")
			if len(records) != 4 {
				return fmt.Errorf("line %d split v expected 4, got %d", lineNumber, len(records))
			}
			val, err := strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("line %d parse x %s failed: %w", lineNumber, records[1], err)
			}
			position[0] = float32(val)
			val, err = strconv.ParseFloat(records[2], 32)
			if err != nil {
				return fmt.Errorf("line %d parse y %s failed: %w", lineNumber, records[2], err)
			}
			position[1] = float32(val)
			val, err = strconv.ParseFloat(records[3], 32)
			if err != nil {
				return fmt.Errorf("line %d parse z %s failed: %w", lineNumber, records[3], err)
			}
			position[2] = float32(val)
			objCache.vertices = append(objCache.vertices, position)
		}
		if strings.HasPrefix(line, "vn ") {
			normal := [3]float32{}
			records := strings.Split(line, " ")
			if len(records) != 4 {
				return fmt.Errorf("line %d split v expected 4, got %d", lineNumber, len(records))
			}
			val, err := strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("line %d parse normal x %s failed: %w", lineNumber, records[1], err)
			}
			normal[0] = float32(val)
			val, err = strconv.ParseFloat(records[2], 32)
			if err != nil {
				return fmt.Errorf("line %d parse normal y %s failed: %w", lineNumber, records[2], err)
			}
			normal[1] = float32(val)
			val, err = strconv.ParseFloat(records[3], 32)
			if err != nil {
				return fmt.Errorf("line %d parse normal z %s failed: %w", lineNumber, records[3], err)
			}
			normal[2] = float32(val)
			objCache.normals = append(objCache.normals, normal)
		}
		if strings.HasPrefix(line, "vt ") {
			uv := [2]float32{}
			records := strings.Split(line, " ")
			if len(records) != 3 {
				return fmt.Errorf("line %d split vt expected 3, got %d", lineNumber, len(records))
			}
			val, err := strconv.ParseFloat(records[1], 32)
			if err != nil {
				return fmt.Errorf("line %d parse x %s failed: %w", lineNumber, records[1], err)
			}
			uv[0] = float32(val)
			val, err = strconv.ParseFloat(records[2], 32)
			if err != nil {
				return fmt.Errorf("line %d parse y %s failed: %w", lineNumber, records[2], err)
			}
			uv[1] = float32(val)
			objCache.uvs = append(objCache.uvs, uv)
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read obj %s: %w", req.ObjPath, err)
	}

	//for i := 0; i < len(positions); i++ {
	//obj.Vertices = append(obj.Vertices, )
	//err = e.AddVertex(positions[i], rotations[i], uvs[i])

	//}
	return nil
}

func triangle(v, t, n int, objCache *objCache, obj *ObjData) (int, error) {

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

	uv := [2]float32{}
	if t != 0 && len(objCache.uvs) >= t {
		uv = objCache.uvs[t]
	}

	obj.Vertices = append(obj.Vertices, &common.Vertex{
		Position: vert,
		Normal:   norm,
		Uv:       uv,
	})

	index = len(obj.Vertices) - 1
	objCache.vertexLookup[fmt.Sprintf("%d/%d/%d", v, t, n)] = index
	return index, nil
}
