package ani

import (
	"fmt"
	"os"
)

// BlenderExport exports a ANI file to a directory for use in blender
func (e *ANI) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	bw, err := os.Create(fmt.Sprintf("%s/animation.txt", path))
	if err != nil {
		return fmt.Errorf("create animation.txt: %w", err)
	}
	defer bw.Close()
	bw.WriteString("name|frame_count\n")

	for _, b := range e.bones {
		bw.WriteString(fmt.Sprintf("%s|%d\n", b.Name, b.FrameCount))

		bfw, err := os.Create(fmt.Sprintf("%s/frame_%s.txt", path, b.Name))
		if err != nil {
			return fmt.Errorf("create frame_%s.txt: %w", b.Name, err)
		}
		bfw.WriteString("bone_name|frame|translation|rotation|scale\n")
		for _, bf := range b.Frames {
			bfw.WriteString(fmt.Sprintf("%s|%d|%s|%s|%s\n", b.Name, bf.Milliseconds, bf.Translation, bf.Rotation, bf.Scale))
		}
		bfw.Close()
	}

	return nil
}
