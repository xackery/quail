package obj

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xackery/quail/common"
)

var (
	rePat = regexp.MustCompile(`([a-z]) (.*) (.*) (.*)\n`)
)

func importMatTxt(obj *ObjData, matTxtPath string) error {

	r, err := os.Open(matTxtPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("readall: %w", err)
	}
	matches := rePat.FindAllStringSubmatch(string(data), -1)

	for lineNumber, records := range matches {
		if len(records) < 5 {
			return fmt.Errorf("line %d has an invalid number of records", lineNumber)
		}
		records = strings.Split(records[0], " ")

		switch records[0] {
		case "m":
			material := materialByName(records[1], obj)
			if material == nil {
				material = &common.Material{
					Name: records[1],
				}
				obj.Materials = append(obj.Materials, material)
			}
			material.ShaderName = records[3]
			val, err := strconv.Atoi(records[2])
			if err != nil {
				return fmt.Errorf("line %d parse flag %s: %w", lineNumber, records[2], err)
			}
			material.Flag = uint32(val)
		case "e":
			material := materialByName(records[1], obj)
			if material == nil {
				material = &common.Material{
					Name: records[1],
				}
				obj.Materials = append(obj.Materials, material)
			}

			val, err := strconv.Atoi(records[3])
			if err != nil {
				return fmt.Errorf("line %d parse material type %s: %w", lineNumber, records[3], err)
			}

			prop := &common.Property{
				Name:      records[2],
				TypeValue: uint32(val),
				StrValue:  records[4],
			}
			material.Properties = append(material.Properties, prop)
		default:
			return fmt.Errorf("line %d has an unsupported definition prefix: %s", lineNumber, records[0])
		}
	}
	return nil
}
