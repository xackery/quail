package obj

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/model/geo"
)

func mtlImport(req *ObjRequest) error {
	rm, err := os.Open(req.MtlPath)
	if err != nil {
		return err
	}
	defer rm.Close()

	scanner := bufio.NewScanner(rm)
	lineNumber := 0
	var lastMaterial *geo.Material
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.HasPrefix(line, "newmtl") {
			if len(line) < 8 {
				return fmt.Errorf("%s line %d: newmtl is too short", req.MtlPath, lineNumber)
			}
			lastMaterial = materialByName(line[7:], req.Data)
			if lastMaterial == nil {
				lastMaterial = &geo.Material{Name: line[7:], ShaderName: "Opaque_MaxCB1.fx"}
				req.Data.Materials = append(req.Data.Materials, lastMaterial)
			}
			continue
		}
		if strings.HasPrefix(line, "map_Kd") {
			if lastMaterial == nil {
				return fmt.Errorf("map_kd line %d found before material definition", lineNumber)
			}
			lastMaterial.Properties = append(lastMaterial.Properties, &geo.Property{Name: "e_TextureDiffuse0", Value: line[7:]})
			continue
		}
		if strings.HasPrefix(line, "map_Bump") {
			if lastMaterial == nil {
				return fmt.Errorf("map_Bump line %d found before material definition", lineNumber)
			}
			lastMaterial.Properties = append(lastMaterial.Properties, &geo.Property{Name: "e_TextureNormal0", Value: line[9:]})
			continue
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read mtl %s: %w", req.MtlPath, err)
	}

	return nil
}
