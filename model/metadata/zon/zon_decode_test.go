package zon

import (
	"io"
	"testing"

	"github.com/xackery/quail/common"
)

func TestDecode(t *testing.T) {
	type args struct {
		zone *common.Zone
		r    io.ReadSeeker
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .zon|1|anguish.zon|anguish.eqg
		// .zon|1|bazaar.zon|bazaar.eqg
		// .zon|1|bloodfields.zon|bloodfields.eqg
		// .zon|1|broodlands.zon|broodlands.eqg
		// .zon|1|catacomba.zon|dranikcatacombsa.eqg
		// .zon|1|wallofslaughter.zon|wallofslaughter.eqg
		// .zon|2|arginhiz.zon|arginhiz.eqg
		// .zon|2|guardian.zon|guardian.eqg
		// .zon|4|arthicrex_te.zon|arthicrex.eqg
		// .zon|4|ascent.zon|direwind.eqg
		// .zon|4|atiiki.zon|atiiki.eqg
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.zone, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
