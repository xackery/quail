package quail

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/qfs"
)

// DirWrite exports the quail target to a directory
func (q *Quail) DirWrite(path string) error {

	if q.FileSystem == nil {
		q.FileSystem = &qfs.OSFS{}
	}

	path = strings.TrimSuffix(path, ".eqg")
	path = strings.TrimSuffix(path, ".s3d")
	path = strings.TrimSuffix(path, ".quail")
	path += ".quail"

	_, err := q.FileSystem.Stat(path)
	if err == nil {
		err = q.FileSystem.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	err = q.FileSystem.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	fi, err := q.FileSystem.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("path %s is not a directory", path)
	}

	if q.Wld != nil {
		err = q.Wld.WriteAscii(path)
		if err != nil {
			return err
		}
	}
	if q.WldObject != nil {
		err = q.WldObject.WriteAscii(path + "/_objects/")
		if err != nil {
			return err
		}
	}
	if q.WldLights != nil {
		err = q.WldLights.WriteAscii(path + "/_lights/")
		if err != nil {
			return err
		}
	}

	if len(q.Assets) > 0 {
		err = q.FileSystem.MkdirAll(path+"/assets/", 0755)
		if err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
	}

	for name, data := range q.Assets {

		/* data, err := fixWonkyDDS(name, texture)
		if err != nil {
			return err
		} */
		err = q.FileSystem.WriteFile(path+"/assets/"+name, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func fixWonkyDDS(name string, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	if string(data[0:3]) == "DDS" {
		//fmt.Println("Found DDS:", name)
		// change to png, blender doesn't like EQ dds
		img, err := dds.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Println("Failed to decode dds:", name, err, "fallback pink image")
			return data, nil
		}
		switch rgba := img.(type) {
		case *image.RGBA:
			buf := &bytes.Buffer{}
			err = png.Encode(buf, rgba)
			if err != nil {
				return nil, fmt.Errorf("png encode: %w", err)
			}
			return buf.Bytes(), nil
		case *image.NRGBA:
			newImg := image.NewRGBA(rgba.Rect)
			draw.Draw(newImg, newImg.Bounds(), rgba, rgba.Rect.Min, draw.Src)
			buf := &bytes.Buffer{}
			err = png.Encode(buf, newImg)
			if err != nil {
				return nil, fmt.Errorf("png encode: %w", err)
			}
			return buf.Bytes(), nil
		default:
			return nil, fmt.Errorf("unknown dds type %T", rgba)
		}
	}
	return data, nil
}
*/
