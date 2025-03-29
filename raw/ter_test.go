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

func TestTerRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()
	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .ter|1|ter_temple01.ter|fhalls.eqg
		{name: "fhalls.eqg"},
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
	//dirTest := helper.DirTest()

	tests := []struct {
		eqg  string
		name string
	}{
		// .ter|1|ter_temple01.ter|fhalls.eqg
		//{eqg: "fhalls.eqg", name: "ter_temple01.ter"}, // PASS
		// .ter|2|ter_abyss01.ter|thenest.eqg
		//{eqg: "thenest.eqg", name: "ter_abyss01.ter"}, // PASS
		// .ter|2|ter_bazaar.ter|bazaar.eqg
		//{eqg: "bazaar.eqg", name: "ter_bazaar.ter"}, // PASS
		// .ter|2|ter_upper.ter|riftseekers.eqg
		//{eqg: "riftseekers.eqg", name: "ter_upper.ter"}, // PASS
		// .ter|2|ter_volcano.ter|delvea.eqg
		//{eqg: "delvea.eqg", name: "ter_volcano.ter"}, // PASS
		// .ter|2|ter_volcano.ter|delveb.eqg
		//{eqg: "delveb.eqg", name: "ter_volcano.ter"}, // PASS
		// .ter|3|ter_aalishai.ter|aalishai.eqg
		//{eqg: "aalishai.eqg", name: "ter_aalishai.ter"}, // PASS
		// .ter|3|ter_akhevatwo.ter|akhevatwo.eqg
		//{eqg: "akhevatwo.eqg", name: "ter_akhevatwo.ter"}, // PASS
		// .ter|3|ter_am_main.ter|arxmentis.eqg
		//{eqg: "arxmentis.eqg", name: "ter_am_main.ter"}, // PASS
		// .ter|3|ter_arena.ter|arena.eqg
		//{eqg: "alkabormare.eqg", name: "ter_faydark.ter"},
		{eqg: "arena.eqg", name: "ter_arena.ter"}, // PASS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.eqg))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.eqg, err.Error())
			}

			data, err := archive.File(tt.name)
			if err != nil {
				t.Fatalf("failed to open file %s: %s", tt.name, err.Error())
			}

			ter := &Ter{}
			err = ter.Read(bytes.NewReader(data))
			if err != nil {

				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = ter.Write(buf)
			if err != nil {
				t.Fatalf("failed to write %s: %s", tt.name, err.Error())
			}

			fmt.Println("src size:", len(data))
			fmt.Println("dst size:", buf.Len())

			os.WriteFile(fmt.Sprintf("%s/%s", helper.DirTest(), tt.name), buf.Bytes(), 0644)

			ter2 := &Ter{}
			err = ter2.Read(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			buf2 := bytes.NewBuffer(nil)
			err = ter2.Write(buf2)
			if err != nil {
				t.Fatalf("failed to write %s: %s", tt.name, err.Error())
			}

			if len(ter2.Materials) != len(ter.Materials) {
				t.Fatalf("%s write: material count mismatch %d vs %d", tt.name, len(ter2.Materials), len(ter.Materials))
			}

			if len(ter2.Vertices) != len(ter.Vertices) {
				t.Fatalf("%s write: vertex count mismatch %d vs %d", tt.name, len(ter2.Vertices), len(ter.Vertices))
			}

			if len(ter2.Faces) != len(ter.Faces) {
				t.Fatalf("%s write: triangle count mismatch %d vs %d", tt.name, len(ter2.Faces), len(ter.Faces))
			}

			for i := 0; i < len(ter.Faces); i++ {
				if ter2.Faces[i].MaterialName != ter.Faces[i].MaterialName {
					t.Fatalf("%s write: face %d material name mismatch %s vs %s", tt.name, i, ter2.Faces[i].MaterialName, ter.Faces[i].MaterialName)
				}
			}

			err = helper.ByteCompareTest(ter.name.data(), ter.name.data())
			if err != nil {
				t.Fatalf("Name data mismatch: %s", err.Error())
			}

			err = helper.ByteCompareTest(buf.Bytes(), buf2.Bytes())
			if err != nil {
				t.Fatalf("%s write: byte compare failed: %s", tt.name, err.Error())
			}
			// err = helper.ByteCompareTest(data, buf.Bytes())
			// if err != nil {
			// 	t.Fatalf("%s byteCompare: %s", tt.name, err)
			// }

		})
	}
}
