package wld

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"gopkg.in/yaml.v3"
)

func TestDecode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	tests := []struct {
		name string
	}{
		{"crushbone"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.name, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(fmt.Sprintf("%s.wld", tt.name))
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.name, err.Error())
			}

			world := common.NewWld("test")
			err = Decode(world, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
			}

		})
	}
}

func TestDumpFragOffset(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)

	tests := []struct {
		name string
	}{
		{"gequip4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s.s3d", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.name, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(fmt.Sprintf("%s.wld", tt.name))
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.name, err.Error())
			}

			world := common.NewWld("test")
			err = Decode(world, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
			}

			path := fmt.Sprintf("%s/%s.wld.yaml", dirTest, tt.name)
			w, err := os.Create(path)
			if err != nil {
				t.Fatalf("failed to create %s: %s", tt.name, err.Error())
			}
			enc := yaml.NewEncoder(w)
			enc.SetIndent(2)
			err = enc.Encode(world.Fragments)
			if err != nil {
				t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
			}
			w.Close()
			fmt.Println("wrote", path)
		})
	}
}
