package mod

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/tag"
)

func TestDecodeMesh(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name       string
		isOldWorld bool
		index      int
	}{
		{"crushbone", true, 221},
		{"crushbone", true, 222},
		{"crushbone", true, 223},
	}
	os.RemoveAll("test")
	os.MkdirAll("test", 0755)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &common.Model{}
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.name, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(fmt.Sprintf("%s.wld", tt.name))
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.name, err.Error())
			}
			wld, err := common.WldOpen(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to load %s: %s", tt.name, err.Error())
			}
			defer wld.Close()

			data, err = wld.Fragment(tt.index)
			if err != nil {
				t.Fatalf("failed to find fragment %d %s: %s", tt.index, tt.name, err.Error())
			}

			nameRef := int32(0)
			r := bytes.NewReader(data)
			err = DecodeMesh(model, &nameRef, tt.isOldWorld, r)
			os.WriteFile(fmt.Sprintf("test/%s.hex", tt.name), data, 0644)
			tag.Write(fmt.Sprintf("test/%s.hex.tags", tt.name))
			if err != nil {
				t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
			}
			fmt.Printf("%s has loaded model %s\n", tt.name, model.Name)
		})
	}
}
