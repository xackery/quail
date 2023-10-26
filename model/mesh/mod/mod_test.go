package mod

import (
	"io"
	"os"
	"testing"

	"github.com/xackery/quail/common"
)

func TestDecode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	type args struct {
		model *common.Model
		r     io.ReadSeeker
	}
	tests := []struct {
		name    string
		args    args
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

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.model, tt.args.r); (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
