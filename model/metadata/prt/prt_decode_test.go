package prt

import (
	"io"
	"testing"

	"github.com/xackery/quail/common"
)

func TestDecode(t *testing.T) {
	type args struct {
		render *common.ParticleRender
		r      io.ReadSeeker
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .prt|3|bat.prt|bat.eqg bat.eqg pfs import: decodePrt bat.prt: invalid version 3, wanted 4+
		// .prt|3|btn.prt|btn.eqg btn.eqg pfs import: decodePrt btn.prt: invalid version 3, wanted 4+
		// .prt|3|chm.prt|chm.eqg chm.eqg pfs import: decodePrt chm.prt: invalid version 3, wanted 4+
		// .prt|3|clv.prt|clv.eqg clv.eqg pfs import: decodePrt clv.prt: invalid version 3, wanted 4+
		// .prt|3|ddm.prt|ddm.eqg ddm.eqg pfs import: decodePrt ddm.prt: invalid version 3, wanted 4+
		// .prt|3|dsf.prt|dsf.eqg dsf.eqg pfs import: decodePrt dsf.prt: invalid version 3, wanted 4+
		// .prt|3|dsg.prt|dsg.eqg dsg.eqg pfs import: decodePrt dsg.prt: invalid version 3, wanted 4+
		// .prt|3|fra.prt|fra.eqg fra.eqg pfs import: decodePrt fra.prt: invalid version 3, wanted 4+
		// .prt|3|mch.prt|mch.eqg mch.eqg pfs import: decodePrt mch.prt: invalid version 3, wanted 4+
		// .prt|3|mur.prt|mur.eqg mur.eqg pfs import: decodePrt mur.prt: invalid version 3, wanted 4+
		// .prt|3|rtn.prt|rtn.eqg rtn.eqg pfs import: decodePrt rtn.prt: invalid version 3, wanted 4+
		// .prt|3|scu.prt|scu.eqg scu.eqg pfs import: decodePrt scu.prt: invalid version 3, wanted 4+
		// .prt|3|tgo.prt|tgo.eqg tgo.eqg pfs import: decodePrt tgo.prt: invalid version 3, wanted 4+
		// .prt|3|tln.prt|tln.eqg tln.eqg pfs import: decodePrt tln.prt: invalid version 3, wanted 4+
		// .prt|4|cnp.prt|cnp.eqg
		// .prt|5|aam.prt|aam.eqg
		// .prt|5|ae3.prt|ae3.eqg
		// .prt|5|ahf.prt|ahf.eqg
		// .prt|5|ahm.prt|ahm.eqg

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.render, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
