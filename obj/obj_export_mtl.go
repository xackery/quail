package obj

import (
	"fmt"
	"os"
)

func exportMtl(obj *ObjData, mtlPath string) error {
	w, err := os.Create(mtlPath)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.WriteString(fmt.Sprintf("# exported by quail\n# Material Count: %d\n", len(obj.Materials)))
	if err != nil {
		return fmt.Errorf("export header: %w", err)
	}

	for _, m := range obj.Materials {
		_, err = w.WriteString(fmt.Sprintf("\nnewmtl %s\n", m.Name))
		if err != nil {
			return fmt.Errorf("newmtl: %w", err)
		}
		_, err = w.WriteString("Ka 1.000000 1.000000 1.000000\nKd 1.000000 1.000000 1.000000\nd 1.000000\nillum 2\n")
		if err != nil {
			return fmt.Errorf("ka: %w", err)
		}
		for _, p := range m.Properties {
			if p.Name == "e_TextureDiffuse0" && p.StrValue != "" {
				_, err = w.WriteString(fmt.Sprintf("map_Kd %s\n", p.StrValue))
				if err != nil {
					return fmt.Errorf("ka: %w", err)
				}
				continue
			}
			if p.Name == "e_TextureNormal0" && p.StrValue != "" {
				_, err = w.WriteString(fmt.Sprintf("map_Bump %s\n", p.StrValue))
				if err != nil {
					return fmt.Errorf("ka: %w", err)
				}
			}
		}
	}

	return nil
}
