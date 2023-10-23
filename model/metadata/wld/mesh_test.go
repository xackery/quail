package wld

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
)

func Test_decodeMesh(t *testing.T) {
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
		{"gequip3.s3d", "gequip3.wld", 0, &Mesh{NameRef: 1414544642}, false},
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
			world := &common.Wld{}
			err = Decode(world, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to decode wld %s: %s", tt.file, err.Error())
			}

			for _, frag := range world.Fragments {
				if frag.FragCode() != 0x36 {
					continue
				}
				//t.Fatalf("frag %d: %+v", frag.FragCode(), frag)
			}

			got, err := decodeMesh(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeMesh() error = %v, wantErr %v", err, tt.wantErr)
			}

			d, ok := got.(*Mesh)
			if !ok {
				t.Errorf("decodeMesh() got = %T, want %T", got, tt.want)
			}
			if d.NameRef != tt.want.(*Mesh).NameRef {
				t.Errorf("decodeMesh() got = %v, want %v", d.NameRef, tt.want.(*Mesh).NameRef)
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("decodeMesh() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestMesh_Encode(t *testing.T) {
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
		{"gequip4.s3d", "gequip4.wld", &Mesh{NameRef: 1414544642}, false},
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
			got, err := decodeMesh(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeMesh() error = %v, wantErr %v", err, tt.wantErr)
			}

			buf := &bytes.Buffer{}
			err = got.Encode(buf)
			if err != nil {
				t.Errorf("encodeMesh() error = %v", err)
			}
			// create a buf with readseeker
			got, err = decodeMesh(bytes.NewReader(buf.Bytes()))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeMesh() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got == nil {
				t.Errorf("decodeMesh() got = nil, want %v", tt.want)
			}

		})
	}
}
