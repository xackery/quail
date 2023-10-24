package wld

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
)

func Test_decodeFirst(t *testing.T) {
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
		//{"gequip.s3d", "gequip.wld", 0, &First{NameRef: 1414544642}, false},
		{"gequip4.s3d", "gequip4.wld", 0, &First{NameRef: 1414544642}, false},
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
			world := common.NewWld("Test")
			err = Decode(world, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to decode wld %s: %s", tt.file, err.Error())
			}

			for _, frag := range world.Fragments {
				if frag.FragCode() != tt.fragIndex {
					continue
				}
				t.Fatalf("frag %d: %+v", frag.FragCode(), frag)
			}
			got, err := decodeFirst(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeFirst() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirst_Encode(t *testing.T) {
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
		//{"gequip.s3d", "gequip.wld", &First{NameRef: 1414544642}, false},
		{"gequip4.s3d", "gequip4.wld", &First{NameRef: 1414544642}, false},
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
			got, err := decodeFirst(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeFirst() error = %v, wantErr %v", err, tt.wantErr)
			}

			buf := &bytes.Buffer{}
			err = got.Encode(buf)
			if err != nil {
				t.Errorf("encodeFirst() error = %v", err)
			}
			// create a buf with readseeker
			got, err = decodeFirst(bytes.NewReader(buf.Bytes()))
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeFirst() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}
