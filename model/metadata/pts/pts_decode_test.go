package pts

import (
	"io"
	"testing"

	"github.com/xackery/quail/common"
)

func TestDecode(t *testing.T) {
	type args struct {
		point *common.ParticlePoint
		r     io.ReadSeeker
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .pts|1|aam.pts|aam.eqg
		// .pts|1|ae3.pts|ae3.eqg
		// .pts|1|ahf.pts|ahf.eqg
		// .pts|1|ahm.pts|ahm.eqg
		// .pts|1|aie.pts|aie.eqg
		// .pts|1|ala.pts|ala.eqg
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.point, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
