package wld

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
)

func Test_decodePaletteFile(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		path      string
		file      string
		fragIndex int
		want      common.FragmentReader
		wantErr   bool
	}{
		{"gequip4.s3d", "gequip4.wld", 0, &PaletteFile{NameRef: 1414544642}, false},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.path))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.file, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(tt.file)
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.file, err.Error())
			}
			world := common.NewWld("test")
			err = Decode(world, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to decode wld %s: %s", tt.file, err.Error())
			}
			/*
				for _, frag := range world.Fragments {
					if frag.FragCode() != tt.fragIndex {
						continue
					}
					t.Fatalf("frag %d: %+v", frag.FragCode(), frag)
				}

				got, err := decodePaletteFile(bytes.NewReader(data))
				if (err != nil) != tt.wantErr {
					t.Errorf("decodePaletteFile() error = %v, wantErr %v", err, tt.wantErr)
				}

				d, ok := got.(*PaletteFile)
				if !ok {
					t.Errorf("decodePaletteFile() got = %T, want %T", got, tt.want)
				}
				if d.NameRef != tt.want.(*PaletteFile).NameRef {
					t.Errorf("decodePaletteFile() got = %v, want %v", d.NameRef, tt.want.(*PaletteFile).NameRef)
				}

				//if !reflect.DeepEqual(got, tt.want) {
				//	t.Errorf("decodePaletteFile() = %v, want %v", got, tt.want)
				//}
			*/
		})
	}
}

func TestPaletteFile_Encode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		path    string
		file    string
		want    common.FragmentReader
		wantErr bool
	}{
		{"gequip4.s3d", "gequip4.wld", &PaletteFile{NameRef: 1414544642}, false},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {

		})
	}
}
