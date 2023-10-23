package lay

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

func TestEncode(t *testing.T) {
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
		// .lay|2|crs.lay|crs.eqg crs.eqg pfs import: decodePrt crs.lay: 1 names materialID 0x41400000 not found
		{name: "crs.eqg"},
		// .lay|2|ddv.lay|ddv.eqg ddv.eqg pfs import: decodePrt ddv.lay: 1 names materialID 0x42800000 not found
		{name: "ddv.eqg"},
		// .lay|2|prt.lay|prt.eqg prt.eqg pfs import: decodePrt prt.lay: 1 names materialID 0x41400000 not found
		{name: "prt.eqg"},
		// .lay|2|rkp.lay|rkp.eqg rkp.eqg pfs import: decodePrt rkp.lay: 1 names materialID 0x41400000 not found
		{name: "rkp.eqg"},
		// .lay|3|rat.lay|rat.eqg rat.eqg pfs import: decodePrt rat.lay: 1 names materialID 0x420000 not found
		{name: "rat.eqg"},
		// .lay|4|aam.lay|aam.eqg
		{name: "aam.eqg"},
		// .lay|4|ahf.lay|ahf.eqg
		{name: "ahf.eqg"},
		// .lay|4|ahm.lay|ahm.eqg
		{name: "ahm.eqg"},
		// .lay|4|ala.lay|ala.eqg
		{name: "ala.eqg"},
		// .lay|4|alg.lay|alg.eqg
		{name: "alg.eqg"},
		// .lay|4|amy.lay|amy.eqg
		{name: "amy.eqg"},
		// .lay|4|cwc.lay|cwc.eqg cwc.eqg pfs import: decodePrt cwc.lay: 0 names colorTexture 0xffffffff not found
		{name: "cwc.eqg"},
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

				model := common.NewModel("")
				err = Decode(model, bytes.NewReader(file.Data()))
				if err != nil {
					err = os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					if err != nil {
						t.Fatalf("failed to write %s: %s", tt.name, err.Error())
					}
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
				}
				fmt.Println("decoded", tt.name)
				fmt.Printf("%+v\n", model)

				//encode
				buf := bytes.NewBuffer(nil)
				err = Encode(model, buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				//decode
				model2 := common.NewModel("")
				err = Decode(model2, bytes.NewReader(buf.Bytes()))
				if err != nil {
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
				}

				if len(model.Layers) != len(model2.Layers) {
					t.Fatalf("layers mismatch: %d != %d", len(model.Layers), len(model2.Layers))
				}

				for i := range model.Layers {
					if model.Layers[i].Material != model2.Layers[i].Material {
						t.Fatalf("material mismatch: %s != %s", model.Layers[i].Material, model2.Layers[i].Material)
					}
					if model.Layers[i].Diffuse != model2.Layers[i].Diffuse {
						t.Fatalf("diffuse mismatch: %s != %s", model.Layers[i].Diffuse, model2.Layers[i].Diffuse)
					}
					if model.Layers[i].Normal != model2.Layers[i].Normal {
						t.Fatalf("normal mismatch: %s != %s", model.Layers[i].Normal, model2.Layers[i].Normal)
					}
				}

			}
		})
	}
}
