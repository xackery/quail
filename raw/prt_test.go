package raw

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

func TestPrtRead(t *testing.T) {
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

		// .prt|3|bat.prt|bat.eqg bat.eqg pfs import: readPrt bat.prt: invalid version 3, wanted 4+
		//{name: "bat.eqg"},
		// .prt|3|btn.prt|btn.eqg btn.eqg pfs import: readPrt btn.prt: invalid version 3, wanted 4+
		// .prt|3|chm.prt|chm.eqg chm.eqg pfs import: readPrt chm.prt: invalid version 3, wanted 4+
		// .prt|3|clv.prt|clv.eqg clv.eqg pfs import: readPrt clv.prt: invalid version 3, wanted 4+
		// .prt|3|ddm.prt|ddm.eqg ddm.eqg pfs import: readPrt ddm.prt: invalid version 3, wanted 4+
		// .prt|3|dsf.prt|dsf.eqg dsf.eqg pfs import: readPrt dsf.prt: invalid version 3, wanted 4+
		// .prt|3|dsg.prt|dsg.eqg dsg.eqg pfs import: readPrt dsg.prt: invalid version 3, wanted 4+
		// .prt|3|fra.prt|fra.eqg fra.eqg pfs import: readPrt fra.prt: invalid version 3, wanted 4+
		// .prt|3|mch.prt|mch.eqg mch.eqg pfs import: readPrt mch.prt: invalid version 3, wanted 4+
		// .prt|3|mur.prt|mur.eqg mur.eqg pfs import: readPrt mur.prt: invalid version 3, wanted 4+
		// .prt|3|rtn.prt|rtn.eqg rtn.eqg pfs import: readPrt rtn.prt: invalid version 3, wanted 4+
		// .prt|3|scu.prt|scu.eqg scu.eqg pfs import: readPrt scu.prt: invalid version 3, wanted 4+
		// .prt|3|tgo.prt|tgo.eqg tgo.eqg pfs import: readPrt tgo.prt: invalid version 3, wanted 4+
		// .prt|3|tln.prt|tln.eqg tln.eqg pfs import: readPrt tln.prt: invalid version 3, wanted 4+
		// .prt|4|cnp.prt|cnp.eqg
		{name: "cnp.eqg"},
		// .prt|5|aam.prt|aam.eqg
		{name: "aam.eqg"},
		// .prt|5|ae3.prt|ae3.eqg
		{name: "ae3.eqg"},
		// .prt|5|ahf.prt|ahf.eqg
		{name: "ahf.eqg"},
		// .prt|5|ahm.prt|ahm.eqg
		{name: "ahm.eqg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".prt" {
					continue
				}
				prt := &Prt{}
				err = prt.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}

func TestPrtWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)

	tests := []struct {
		name    string
		wantErr bool
	}{
		// .prt|3|bat.prt|bat.eqg bat.eqg pfs import: readPrt bat.prt: invalid version 3, wanted 4+
		//{name: "bat.eqg"}, // FIXME: v3 or below anim support
		// .prt|3|btn.prt|btn.eqg btn.eqg pfs import: readPrt btn.prt: invalid version 3, wanted 4+
		// .prt|3|chm.prt|chm.eqg chm.eqg pfs import: readPrt chm.prt: invalid version 3, wanted 4+
		// .prt|3|clv.prt|clv.eqg clv.eqg pfs import: readPrt clv.prt: invalid version 3, wanted 4+
		// .prt|3|ddm.prt|ddm.eqg ddm.eqg pfs import: readPrt ddm.prt: invalid version 3, wanted 4+
		// .prt|3|dsf.prt|dsf.eqg dsf.eqg pfs import: readPrt dsf.prt: invalid version 3, wanted 4+
		// .prt|3|dsg.prt|dsg.eqg dsg.eqg pfs import: readPrt dsg.prt: invalid version 3, wanted 4+
		// .prt|3|fra.prt|fra.eqg fra.eqg pfs import: readPrt fra.prt: invalid version 3, wanted 4+
		// .prt|3|mch.prt|mch.eqg mch.eqg pfs import: readPrt mch.prt: invalid version 3, wanted 4+
		// .prt|3|mur.prt|mur.eqg mur.eqg pfs import: readPrt mur.prt: invalid version 3, wanted 4+
		// .prt|3|rtn.prt|rtn.eqg rtn.eqg pfs import: readPrt rtn.prt: invalid version 3, wanted 4+
		// .prt|3|scu.prt|scu.eqg scu.eqg pfs import: readPrt scu.prt: invalid version 3, wanted 4+
		// .prt|3|tgo.prt|tgo.eqg tgo.eqg pfs import: readPrt tgo.prt: invalid version 3, wanted 4+
		// .prt|3|tln.prt|tln.eqg tln.eqg pfs import: readPrt tln.prt: invalid version 3, wanted 4+
		// .prt|4|cnp.prt|cnp.eqg
		{name: "cnp.eqg"},
		// .prt|5|aam.prt|aam.eqg
		{name: "aam.eqg"},
		// .prt|5|ae3.prt|ae3.eqg
		{name: "ae3.eqg"},
		// .prt|5|ahf.prt|ahf.eqg
		{name: "ahf.eqg"},
		// .prt|5|ahm.prt|ahm.eqg
		{name: "ahm.eqg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".prt" {
					continue
				}
				prt := &Prt{}
				err = prt.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = prt.Write(buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				srcData := file.Data()
				dstData := buf.Bytes()

				err = common.ByteCompareTest(srcData, dstData)
				if err != nil {
					t.Fatalf("%s failed byteCompare: %s", tt.name, err)
				}
			}
		})
	}
}
