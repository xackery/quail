package ani

import (
	"io"
	"testing"

	"github.com/xackery/quail/common"
)

func TestDecode(t *testing.T) {
	type args struct {
		animation *common.Animation
		r         io.ReadSeeker
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .ani|1|sidl_ba_1_tln.ani|tln.eqg
		// .ani|2|stnd_ba_1_exo.ani|exo.eqg eye_chr.s3d pfs import: s3d load: decode: dirName for crc 655939147 not found
		// .ani|2|walk_ba_1_vaf.ani|vaf.eqg valdeholm.eqg pfs import: eqg load: decode: read nameData unexpected EOF
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.animation, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
