package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestLitRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()

	tests := []struct {
		name    string
		litName string
	}{
		//.lit|1|commons_inn_obj_lampc01.lit|commonlands.eqg
		//{name: "commonlands.eqg", litName: "commons_inn_obj_lampc01.lit"}, // PASS
		//.lit|1|communalhut_obj_treasureb01.lit|buriedsea.eqg
		//{name: "buriedsea.eqg", litName: "communalhut_obj_treasureb01.lit"}, // PASS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".lay" {
					continue
				}
				lit := &Lit{}
				err = lit.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}

func TestLitWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()

	tests := []struct {
		name    string
		litName string
	}{
		//.lit|1|commons_inn_obj_lampc01.lit|commonlands.eqg
		//{name: "commonlands.eqg", litName: "commons_inn_obj_lampc01.lit"}, // PASS
		//.lit|1|communalhut_obj_treasureb01.lit|buriedsea.eqg
		//{name: "buriedsea.eqg", litName: "communalhut_obj_treasureb01.lit"}, // PASS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".lay" {
					continue
				}
				lit := &Lit{}

				err = lit.Read(bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = lit.Write(buf)
				if err != nil {
					t.Fatalf("failed to write %s: %s", tt.name, err.Error())
				}

				lit2 := &Lit{}
				err = lit2.Read(bytes.NewReader(buf.Bytes()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				if len(lit.Entries) != len(lit2.Entries) {
					t.Fatalf("layers mismatch: %d != %d", len(lit.Entries), len(lit2.Entries))
				}

				for i := range lit.Entries {
					for j := 0; j < 5; j++ {
						if lit.Entries[i][j] != lit2.Entries[i][j] {
							t.Fatalf("layer mismatch: %d != %d", lit.Entries[i][j], lit2.Entries[i][j])
						}
					}
				}
			}
		})
	}
}
