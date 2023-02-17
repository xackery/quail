package blend

import (
	"fmt"
	"os"
)

func blendPath() string {
	path := fmt.Sprintf("/usr/local/blender/blender")
	_, err := os.Stat(path)
	if err == nil {
		return path
	}
	return ""
}
