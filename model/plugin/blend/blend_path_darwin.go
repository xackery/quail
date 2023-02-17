package blend

import (
	"fmt"
	"os"
)

func blendPath() string {
	path, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	path = fmt.Sprintf("%s/Library/Application Support/Steam/SteamApps/common/Blender/Blender.app/Contents/MacOS/blender", path)
	_, err = os.Stat(path)
	if err == nil {
		return path //strings.ReplaceAll(path, " ", "\\ ")
	}
	path = "Applications/Blender/blender.app/Contents/MacOS/blender"
	_, err = os.Stat(path)
	if err == nil {
		return path
	}
	return ""
}
