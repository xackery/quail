package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/tag"
)

func TestLodRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)

	// TODO: add lod
	tests := []struct {
		name    string
		lodName string
	}{}

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
				lod := &Lod{}
				err = lod.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}

func TestLodWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)

	tests := []struct {
		name    string
		lodName string
	}{}

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
				lod := &Lod{}

				err = lod.Read(bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = lod.Write(buf)
				if err != nil {
					t.Fatalf("failed to write %s: %s", tt.name, err.Error())
				}

				lod2 := &Lod{}
				err = lod2.Read(bytes.NewReader(buf.Bytes()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				if len(lod.Entries) != len(lod2.Entries) {
					t.Fatalf("lod mismatch: %d != %d", len(lod.Entries), len(lod2.Entries))
				}

				for i := range lod.Entries {
					if lod.Entries[i].Category != lod2.Entries[i].Category {
						t.Fatalf("category mismatch: %s != %s", lod.Entries[i].Category, lod2.Entries[i].Category)
					}
					if lod.Entries[i].Distance != lod2.Entries[i].Distance {
						t.Fatalf("distance mismatch: %d != %d", lod.Entries[i].Distance, lod2.Entries[i].Distance)
					}

					if lod.Entries[i].ObjectName != lod2.Entries[i].ObjectName {
						t.Fatalf("object name mismatch: %s != %s", lod.Entries[i].ObjectName, lod2.Entries[i].ObjectName)
					}
				}
			}
		})
	}
}
