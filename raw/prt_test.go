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

		// .prt|3|bat.prt|bat.eqg bat.eqg pfs import: decodePrt bat.prt: invalid version 3, wanted 4+
		//{name: "bat.eqg"},
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
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
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
		// .prt|3|bat.prt|bat.eqg bat.eqg pfs import: decodePrt bat.prt: invalid version 3, wanted 4+
		//{name: "bat.eqg"}, // FIXME: v3 or below anim support
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
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = prt.Write(buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				/* srcData := file.Data()
				dstData := buf.Bytes()

				for i := 0; i < len(srcData); i++ {
					if len(dstData) <= i {
						min := 0
						max := len(srcData)
						fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
						max = len(dstData)
						fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))

						t.Fatalf("%s src eof at offset %d (dst is too large by %d bytes)", tt.name, i, len(dstData)-len(srcData))
					}
					if len(dstData) <= i {
						t.Fatalf("%s dst eof at offset %d (dst is too small by %d bytes)", tt.name, i, len(srcData)-len(dstData))
					}
					if srcData[i] == dstData[i] {
						continue
					}

					fmt.Printf("%s mismatch at offset %d (src: 0x%x vs dst: 0x%x aka %d)\n", tt.name, i, srcData[i], dstData[i], dstData[i])
					max := i + 16
					if max > len(srcData) {
						max = len(srcData)
					}

					min := i - 16
					if min < 0 {
						min = 0
					}
					fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
					if max > len(dstData) {
						max = len(dstData)
					}

					fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))
					t.Fatalf("%s encode: data mismatch", tt.name)
				} */
			}
		})
	}
}
