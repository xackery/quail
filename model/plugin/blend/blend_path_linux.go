package blend

import (
	"os"
)

func blendPath() string {
	path := "/usr/local/blender/blender"
	_, err := os.Stat(path)
	if err == nil {
		return path
	}
	return ""
}
