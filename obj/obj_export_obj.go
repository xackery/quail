package obj

import (
	"fmt"
	"os"
)

func exportObjFile(obj *ObjData, objPath string) error {

	if obj.Name == "" {
		obj.Name = "unnamed"
	}
	w, err := os.Create(objPath)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.WriteString("# exported by quail\n\n")
	if err != nil {
		return fmt.Errorf("export header: %w", err)
	}

	_, err = w.WriteString(fmt.Sprintf("mtllib %s.mtl\no %s\n", obj.Name, obj.Name))
	if err != nil {
		return fmt.Errorf("mtllib: %w", err)
	}

	for _, e := range obj.Vertices {
		fmt.Println(e)
		_, err = w.WriteString(fmt.Sprintf("v %0.6f %0.6f %0.6f\n", e.Position.X, e.Position.Y, e.Position.Z))
		if err != nil {
			return fmt.Errorf("export pos: %w", err)
		}
	}
	for _, e := range obj.Vertices {
		_, err = w.WriteString(fmt.Sprintf("vt %0.6f %0.6f\n", e.Uv.X, e.Uv.Y))
		if err != nil {
			return fmt.Errorf("export uv: %w", err)
		}
	}
	for _, e := range obj.Vertices {
		_, err = w.WriteString(fmt.Sprintf("vn %0.6f %0.6f %0.6f\n", e.Normal.X, e.Normal.Y, e.Normal.Z))
		if err != nil {
			return fmt.Errorf("export normal: %w", err)
		}
	}

	lastMaterial := ""
	group := 0
	for _, e := range obj.Triangles {
		if lastMaterial != e.MaterialName {
			lastMaterial = e.MaterialName
			_, err = w.WriteString(fmt.Sprintf("usemtl %s\ns off\ng piece%d\n", e.MaterialName, group))
			if err != nil {
				return fmt.Errorf("usemtl: %w", err)
			}
			group++
		}
		_, err = w.WriteString(fmt.Sprintf("f %d/%d/%d %d/%d/%d %d/%d/%d\n", int(e.Index.X+1), int(e.Index.X+1), int(e.Index.X+1), int(e.Index.Y+1), int(e.Index.Y+1), int(e.Index.Y+1), int(e.Index.Z+1), int(e.Index.Z+1), int(e.Index.Z+1)))
		if err != nil {
			return fmt.Errorf("f: %w", err)
		}
	}

	return nil
}
