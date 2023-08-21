package lay

import (
	"io"
	"testing"

	"github.com/xackery/quail/common"
)

func TestDecode(t *testing.T) {
	type args struct {
		model *common.Model
		r     io.ReadSeeker
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .lay|2|crs.lay|crs.eqg crs.eqg pfs import: decodePrt crs.lay: 1 names materialID 0x41400000 not found
		// .lay|2|ddv.lay|ddv.eqg ddv.eqg pfs import: decodePrt ddv.lay: 1 names materialID 0x42800000 not found
		// .lay|2|prt.lay|prt.eqg prt.eqg pfs import: decodePrt prt.lay: 1 names materialID 0x41400000 not found
		// .lay|2|rkp.lay|rkp.eqg rkp.eqg pfs import: decodePrt rkp.lay: 1 names materialID 0x41400000 not found
		// .lay|3|rat.lay|rat.eqg rat.eqg pfs import: decodePrt rat.lay: 1 names materialID 0x420000 not found
		// .lay|4|aam.lay|aam.eqg
		// .lay|4|ahf.lay|ahf.eqg
		// .lay|4|ahm.lay|ahm.eqg
		// .lay|4|ala.lay|ala.eqg
		// .lay|4|alg.lay|alg.eqg
		// .lay|4|amy.lay|amy.eqg
		// .lay|4|cwc.lay|cwc.eqg cwc.eqg pfs import: decodePrt cwc.lay: 0 names colorTexture 0xffffffff not found
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.model, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
