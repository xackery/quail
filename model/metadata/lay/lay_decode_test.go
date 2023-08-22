package lay

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/tag"
)

func TestDecode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
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

	os.RemoveAll("test")
	os.MkdirAll("test", 0755)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".lay" {
					continue
				}
				model := &common.Model{}
				err = Decode(model, bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("test/%s", file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("test/%s.tags", file.Name()))
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}
