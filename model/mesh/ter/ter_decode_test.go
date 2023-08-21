package ter

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
		// .ter|1|ter_temple01.ter|fhalls.eqg
		// .ter|2|ter_abyss01.ter|thenest.eqg
		// .ter|2|ter_bazaar.ter|bazaar.eqg
		//.ter|2|ter_upper.ter|riftseekers.eqg
		//.ter|2|ter_volcano.ter|delvea.eqg
		//.ter|2|ter_volcano.ter|delveb.eqg
		//.ter|3|ter_aalishai.ter|aalishai.eqg
		//.ter|3|ter_akhevatwo.ter|akhevatwo.eqg
		//.ter|3|ter_am_main.ter|arxmentis.eqg
		//.ter|3|ter_arena.ter|arena.eqg
		//.ter|3|ter_arena.ter|arena2.eqg
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.model, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
