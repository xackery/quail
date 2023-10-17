// Package mat repesents metadata related to material decoding
package mat

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
)

func TestDecodeMaterial(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name  string
		index int
	}{
		{"crushbone", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			material := &common.Material{}
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
			err = DecodeMaterialDef(material, &nameRef, r)
			if err != nil {
				t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
			}
			fmt.Printf("%s has loaded material\n", tt.name)
		})
	}
}
