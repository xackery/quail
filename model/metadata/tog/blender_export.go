package tog

import (
	"fmt"
	"os"

	"github.com/xackery/quail/dump"
)

// BlenderExport exports a TOG file to a directory for use in blender
func (e *TOG) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	if len(e.objects) > 0 {
		lw, err := os.Create(fmt.Sprintf("%s/object.txt", path))
		if err != nil {
			return fmt.Errorf("create object.txt: %w", err)
		}

		defer lw.Close()
		lw.WriteString("name|position|rotation|scale|file_type|file_name\n")
		for _, obj := range e.objects {
			lw.WriteString(dump.Str(obj.Name) + "|")
			lw.WriteString(dump.Str(obj.Position) + "|")
			lw.WriteString(dump.Str(obj.Rotation) + "|")
			lw.WriteString(dump.Str(obj.Scale) + "|")
			lw.WriteString(dump.Str(obj.FileType) + "|")
			lw.WriteString(dump.Str(obj.FileName) + "\n")
		}
	}

	return nil
}
