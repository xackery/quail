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

func TestTerRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)
	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .ter|1|ter_temple01.ter|fhalls.eqg
		//{name: "fhalls.eqg"},
		// .ter|2|ter_abyss01.ter|thenest.eqg
		// .ter|2|ter_bazaar.ter|bazaar.eqg
		// .ter|2|ter_upper.ter|riftseekers.eqg
		// .ter|2|ter_volcano.ter|delvea.eqg
		// .ter|2|ter_volcano.ter|delveb.eqg
		// .ter|3|ter_aalishai.ter|aalishai.eqg
		// .ter|3|ter_akhevatwo.ter|akhevatwo.eqg
		// .ter|3|ter_am_main.ter|arxmentis.eqg
		// .ter|3|ter_arena.ter|arena.eqg
		// .ter|3|ter_arena.ter|arena2.eqg
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".ter" {
					continue
				}
				ter := &Ter{}

				err = ter.Read(bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}

func TestTerWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)
	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .ter|1|ter_temple01.ter|fhalls.eqg
		// .ter|2|ter_abyss01.ter|thenest.eqg
		// .ter|2|ter_bazaar.ter|bazaar.eqg
		// .ter|2|ter_upper.ter|riftseekers.eqg
		// .ter|2|ter_volcano.ter|delvea.eqg
		// .ter|2|ter_volcano.ter|delveb.eqg
		// .ter|3|ter_aalishai.ter|aalishai.eqg
		// .ter|3|ter_akhevatwo.ter|akhevatwo.eqg
		// .ter|3|ter_am_main.ter|arxmentis.eqg
		// .ter|3|ter_arena.ter|arena.eqg
		// .ter|3|ter_arena.ter|arena2.eqg
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".ter" {
					continue
				}
				ter := &Ter{}
				err = ter.Read(bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = ter.Write(buf)
				if err != nil {
					t.Fatalf("failed to write %s: %s", tt.name, err.Error())
				}

				//srcData := file.Data()
				//dstData := buf.Bytes()
				/*for i := 0; i < len(srcData); i++ {
					if len(dstData) <= i {
						min := 0
						max := len(srcData)
						fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
						max = len(dstData)
						fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))

						t.Fatalf("%s src eof at offset %d (dst is too large by %d bytes)", tt.name, i, len(dstData)-len(srcData))
					}
					if len(dstData) <= i {
						t.Fatalf("%s dst eof at offset %d (dst is too small by %d bytes)", tt.name, i, len(srcData)-len(dstData))
					}
					if srcData[i] == dstData[i] {
						continue
					}

					fmt.Printf("%s mismatch at offset %d (src: 0x%x vs dst: 0x%x aka %d)\n", tt.name, i, srcData[i], dstData[i], dstData[i])
					max := i + 16
					if max > len(srcData) {
						max = len(srcData)
					}

					min := i - 16
					if min < 0 {
						min = 0
					}
					fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
					if max > len(dstData) {
						max = len(dstData)
					}

					fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))
					//os.WriteFile(fmt.Sprintf("%s/_src_%s", dirTest, file.Name()), file.Data(), 0644)
					//os.WriteFile(fmt.Sprintf("%s/_dst_%s", dirTest, file.Name()), buf.Bytes(), 0644)
					t.Fatalf("%s write: data mismatch", tt.name)
				}*/
			}
		})
	}
}
