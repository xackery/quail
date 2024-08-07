package quail

import (
	"fmt"
	"os"
	"strings"
)

// DirWrite exports the quail target to a directory
func (q *Quail) DirWrite(path string) error {

	path = strings.TrimSuffix(path, ".eqg")
	path = strings.TrimSuffix(path, ".s3d")
	path = strings.TrimSuffix(path, ".quail")
	path += ".quail"

	_, err := os.Stat(path)
	if err == nil {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("path %s is not a directory", path)
	}

	if q.wld != nil {
		err = q.wld.WriteAscii(path, true)
		if err != nil {
			return err
		}
	}
	if q.wldObject != nil {
		err = q.wldObject.WriteAscii(path+"/_lights/", true)
		if err != nil {
			return err
		}
	}
	if q.wldLights != nil {
		err = q.wldLights.WriteAscii(path+"/_objects/", true)
		if err != nil {
			return err
		}
	}

	for name, texture := range q.Textures {
		err = os.WriteFile(path+"/"+name, texture, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
