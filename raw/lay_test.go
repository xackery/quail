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

func TestLayRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()
	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .lay|2|crs.lay|crs.eqg crs.eqg pfs import: readPrt crs.lay: 1 names materialID 0x41400000 not found
		{name: "crs.eqg"}, // PASS
		// .lay|2|ddv.lay|ddv.eqg ddv.eqg pfs import: readPrt ddv.lay: 1 names materialID 0x42800000 not found
		//{name: "ddv.eqg"}, // PASS
		// .lay|2|prt.lay|prt.eqg prt.eqg pfs import: readPrt prt.lay: 1 names materialID 0x41400000 not found
		//{name: "prt.eqg"}, // PASS
		// .lay|2|rkp.lay|rkp.eqg rkp.eqg pfs import: readPrt rkp.lay: 1 names materialID 0x41400000 not found
		//{name: "rkp.eqg"}, // PASS
		// .lay|3|rat.lay|rat.eqg rat.eqg pfs import: readPrt rat.lay: 1 names materialID 0x420000 not found
		//{name: "rat.eqg"}, // PASS
		// .lay|4|aam.lay|aam.eqg
		//{name: "aam.eqg"}, // PASS
		// .lay|4|ahf.lay|ahf.eqg
		//{name: "ahf.eqg"}, // PASS
		// .lay|4|ahm.lay|ahm.eqg
		//{name: "ahm.eqg"}, // PASS
		// .lay|4|ala.lay|ala.eqg
		//{name: "ala.eqg"}, // PASS
		// .lay|4|alg.lay|alg.eqg
		//{name: "alg.eqg"}, // PASS
		// .lay|4|amy.lay|amy.eqg
		//{name: "amy.eqg"}, // PASS
		// .lay|4|cwc.lay|cwc.eqg cwc.eqg pfs import: readPrt cwc.lay: 0 names colorTexture 0xffffffff not found
		//{name: "cwc.eqg"}, // PASS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open pfs %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".lay" {
					continue
				}

				lay := &Lay{}
				err = lay.Read(bytes.NewReader(file.Data()))
				if err != nil {
					err = os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					if err != nil {
						t.Fatalf("failed to write %s: %s", tt.name, err.Error())
					}
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}
			}
		})
	}
}

func TestLayWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()
	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .lay|2|crs.lay|crs.eqg crs.eqg pfs import: readPrt crs.lay: 1 names materialID 0x41400000 not found
		//{name: "crs.eqg"}, // PASS
		// .lay|2|ddv.lay|ddv.eqg ddv.eqg pfs import: readPrt ddv.lay: 1 names materialID 0x42800000 not found
		//{name: "ddv.eqg"}, // PASS
		// .lay|2|prt.lay|prt.eqg prt.eqg pfs import: readPrt prt.lay: 1 names materialID 0x41400000 not found
		//{name: "prt.eqg"}, // PASS
		// .lay|2|rkp.lay|rkp.eqg rkp.eqg pfs import: readPrt rkp.lay: 1 names materialID 0x41400000 not found
		//{name: "rkp.eqg"}, // PASS
		// .lay|3|rat.lay|rat.eqg rat.eqg pfs import: readPrt rat.lay: 1 names materialID 0x420000 not found
		//{name: "rat.eqg"}, // PASS
		// .lay|4|aam.lay|aam.eqg
		//{name: "aam.eqg"}, // PASS
		// .lay|4|ahf.lay|ahf.eqg
		//{name: "ahf.eqg"}, // PASS
		// .lay|4|ahm.lay|ahm.eqg
		//{name: "ahm.eqg"}, // PASS
		// .lay|4|ala.lay|ala.eqg
		//{name: "ala.eqg"}, // PASS
		// .lay|4|alg.lay|alg.eqg
		//{name: "alg.eqg"}, // PASS
		// .lay|4|amy.lay|amy.eqg
		//{name: "amy.eqg"}, // PASS
		// .lay|4|cwc.lay|cwc.eqg cwc.eqg pfs import: readPrt cwc.lay: 0 names colorTexture 0xffffffff not found
		{name: "cwc.eqg"}, // PASS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open pfs %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".lay" {
					continue
				}

				lay := &Lay{}
				err = lay.Read(bytes.NewReader(file.Data()))
				if err != nil {
					err = os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					if err != nil {
						t.Fatalf("failed to write %s: %s", tt.name, err.Error())
					}
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				//encode
				buf := bytes.NewBuffer(nil)
				err = lay.Write(buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				//read
				lay2 := &Lay{}
				err = lay2.Read(bytes.NewReader(buf.Bytes()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				if len(lay.Entries) != len(lay2.Entries) {
					t.Fatalf("layers mismatch: %d != %d", len(lay.Entries), len(lay2.Entries))
				}

				for i := range lay.Entries {
					if lay.Entries[i].Material != lay2.Entries[i].Material {
						t.Fatalf("material mismatch: %s != %s", lay.Entries[i].Material, lay2.Entries[i].Material)
					}
					if lay.Entries[i].Diffuse != lay2.Entries[i].Diffuse {
						t.Fatalf("diffuse mismatch: %s != %s", lay.Entries[i].Diffuse, lay2.Entries[i].Diffuse)
					}
					if lay.Entries[i].Normal != lay2.Entries[i].Normal {
						t.Fatalf("normal mismatch: %s != %s", lay.Entries[i].Normal, lay2.Entries[i].Normal)
					}
				}

			}
		})
	}
}
