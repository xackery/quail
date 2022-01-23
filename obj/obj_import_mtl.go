package obj

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/common"
)

func importMtl(obj *ObjData, mtlPath string) error {
	rm, err := os.Open(mtlPath)
	if err != nil {
		return err
	}
	defer rm.Close()

	scanner := bufio.NewScanner(rm)
	lineNumber := 0
	var lastMaterial *common.Material
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.HasPrefix(line, "newmtl") {
			if len(line) < 8 {
				return fmt.Errorf("%s line %d: newmtl is too short", mtlPath, lineNumber)
			}
			lastMaterial = materialByName(line[7:], obj)
			if lastMaterial == nil {
				lastMaterial = &common.Material{Name: line[7:], ShaderName: "Opaque_MaxCB1.fx"}
				obj.Materials = append(obj.Materials, lastMaterial)
			}
			continue
		}
		if strings.HasPrefix(line, "map_Kd") {
			if lastMaterial == nil {
				return fmt.Errorf("map_kd line %d found before material definition", lineNumber)
			}
			lastMaterial.Properties = append(lastMaterial.Properties, &common.Property{Name: "e_TextureDiffuse0", StrValue: line[7:]})
			continue
		}
		if strings.HasPrefix(line, "map_Bump") {
			if lastMaterial == nil {
				return fmt.Errorf("map_Bump line %d found before material definition", lineNumber)
			}
			lastMaterial.Properties = append(lastMaterial.Properties, &common.Property{Name: "e_TextureNormal0", StrValue: line[9:]})
			continue
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("read mtl %s: %w", mtlPath, err)
	}

	return nil
}
