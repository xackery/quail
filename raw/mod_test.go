package raw

import (
	"bytes"
	"os"
	"testing"

	"github.com/xackery/quail/pfs"
)

func TestModRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		eqg     string
		file    string
		wantErr bool
	}{

		// .mod|0|obp_fob_tree.mod|oldfieldofbone.eqg oldfieldofbone.eqg pfs import: decodeMod obp_fob_tree.mod: invalid header EQLO, wanted EQGM
		// .mod|0|obp_fob_tree.mod|oldfieldofboneb.eqg oldfieldofboneb.eqg pfs import: decodeMod obp_fob_tree.mod: invalid header EQLO, wanted EQGM
		// .mod|1|arch.mod|dranik.eqg
		// .mod|1|aro.mod|aro.eqg
		// .mod|1|col_b04.mod|b04.eqg b04.eqg pfs import: decodeMod col_b04.mod: material shader not found
		// .mod|2|boulder_lg.mod|broodlands.eqg
		// .mod|2|et_door01.mod|stillmoona.eqg
		// .mod|3|.mod|paperbaghat.eqg
		// .mod|3|it11409.mod|undequip.eqg
		{eqg: "undequip.eqg", file: "it11409.mod"},
	}
	for _, tt := range tests {
		t.Run(tt.eqg, func(t *testing.T) {
			archive, err := pfs.NewFile(eqPath + "/" + tt.eqg)
			if err != nil {
				t.Fatalf("pfs.NewFile() error = %v", err)
			}

			data, err := archive.File(tt.file)
			if err != nil {
				t.Fatalf("archive.Open() error = %v", err)
			}

			mod := &Mod{}
			err = mod.Read(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
