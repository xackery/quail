package blender

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/pfs/eqg"
)

func Test_export_import_cmp(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{name: "it13926.eqg", fileName: "it13926.mod", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("newFile: %s", err.Error())
				return
			}
			c1, err := mod.NewFile(tt.fileName, pfs, tt.fileName)
			if err != nil {
				t.Fatalf("mod newfile: %s", err.Error())
				return
			}

			buf1 := &bytes.Buffer{}
			err = c1.Encode(buf1)
			if err != nil {
				t.Fatalf("c1 Encode: %s", err.Error())
				return
			}

			err = export(&cobra.Command{}, []string{fmt.Sprintf("%s/%s", eqPath, tt.name), "test/"})
			if err != nil {
				t.Fatalf("export: %s", err.Error())
			}
			err = import_blender(&cobra.Command{}, []string{fmt.Sprintf("%s/_%s", "test", tt.name), fmt.Sprintf("test/%s", tt.name)})
			if err != nil {
				t.Fatalf("import_blender: %s", err.Error())
			}

			pfs, err = eqg.NewFile(fmt.Sprintf("test/%s", tt.name))
			if err != nil {
				t.Fatalf("newFile: %s", err.Error())
				return
			}

			c2, err := mod.NewFile(tt.fileName, pfs, tt.fileName)
			if err != nil {
				t.Fatalf("mod newfile: %s", err.Error())
				return
			}

			buf2 := &bytes.Buffer{}
			err = c2.Encode(buf2)
			if err != nil {
				t.Fatalf("c2 Encode: %s", err.Error())
				return
			}
			w1, err := os.Create("test/c1.hex")
			if err != nil {
				t.Fatalf("os.Create: %s", err.Error())
				return
			}
			_, err = w1.Write(buf1.Bytes())
			if err != nil {
				t.Fatalf("w1.Write: %s", err.Error())
				return
			}
			w1.Close()

			w2, err := os.Create("test/c2.hex")
			if err != nil {
				t.Fatalf("os.Create: %s", err.Error())
				return
			}
			_, err = w2.Write(buf2.Bytes())
			if err != nil {
				t.Fatalf("w2.Write: %s", err.Error())
				return
			}
			w2.Close()

			fmt.Println(len(buf1.Bytes()), len(buf2.Bytes()))
			if len(buf1.Bytes()) != len(buf2.Bytes()) {
				t.Fatalf("buf1 and buf2 are not the same length")
				return
			}

			for offset, v := range buf1.Bytes() {
				if buf2.Bytes()[offset] != v {
					t.Fatalf("buf2 does not match buf1 at offset %d", offset)
					return
				}
			}

			eqgData, err := os.ReadFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("eqg1.ReadFile: %s", err.Error())
				return
			}
			eqgData2, err := os.ReadFile(fmt.Sprintf("test/%s", tt.name))
			if err != nil {
				t.Fatalf("eqg2.ReadFile: %s", err.Error())
				return
			}
			if len(eqgData) != len(eqgData2) {
				t.Fatalf("eqg1 and eqg2 are not the same length")
				return
			}
			for i, v := range eqgData {
				if eqgData2[i] != v {
					t.Fatalf("eqg2 does not match eqg1 at offset %d", i)
					return
				}
			}

		})
	}
}
