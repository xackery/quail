package raw

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestModRead(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		eqg     string
		file    string
		wantErr bool
	}{
		// .mod|0|obp_fob_tree.mod|oldfieldofbone.eqg oldfieldofbone.eqg pfs import: readMod obp_fob_tree.mod: invalid header EQLO, wanted EQGM
		//{eqg: "oldfieldofbone.eqg", file: "obp_fob_tree.mod"}, // TODO: EQLO v4 .mod?
		// .mod|0|obp_fob_tree.mod|oldfieldofboneb.eqg oldfieldofboneb.eqg pfs import: readMod obp_fob_tree.mod: invalid header EQLO, wanted EQGM
		//{eqg: "oldfieldofboneb.eqg", file: "obp_fob_tree.mod"}, // TODO: EQLO v4 .mod
		// .mod|1|arch.mod|dranik.eqg
		//{eqg: "dranik.eqg", file: "arch.mod"}, // PASS
		// .mod|1|aro.mod|aro.eqg
		//{eqg: "aro.eqg", file: "aro.mod"}, // PASS
		// .mod|1|col_b04.mod|b04.eqg b04.eqg pfs import: readMod col_b04.mod: material shader not found
		//{eqg: "b04.eqg", file: "col_b04.mod"}, // PASS
		// .mod|2|boulder_lg.mod|broodlands.eqg
		//{eqg: "broodlands.eqg", file: "boulder_lg.mod"}, // PASS
		// .mod|2|et_door01.mod|stillmoona.eqg
		//{eqg: "stillmoona.eqg", file: "et_door01.mod"}, // PASS
		// .mod|3|.mod|paperbaghat.eqg
		//{eqg: "paperbaghat.eqg", file: ".mod"}, // PASS
		// .mod|3|it11409.mod|undequip.eqg
		//{eqg: "undequip.eqg", file: "it11409.mod"}, // PASS
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

func TestModWrite(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	testDir := helper.DirTest()
	tests := []struct {
		eqg     string
		file    string
		wantErr bool
	}{

		{eqg: "bazaar.eqg", file: "obj_sackc.mod"}, // TODO: 7813 bytes remaining
		{eqg: "dbx.eqg", file: "dbx.mod"},          // PASS
		//	{eqg: "it12043.eqg", file: "it12043.mod"}, // PASS
		// .mod|0|obp_fob_tree.mod|oldfieldofbone.eqg oldfieldofbone.eqg pfs import: readMod obp_fob_tree.mod: invalid header EQLO, wanted EQGM
		//{eqg: "oldfieldofbone.eqg", file: "obp_fob_tree.mod"}, // TODO: EQLO v4 .mod?
		// .mod|0|obp_fob_tree.mod|oldfieldofboneb.eqg oldfieldofboneb.eqg pfs import: readMod obp_fob_tree.mod: invalid header EQLO, wanted EQGM
		//{eqg: "oldfieldofboneb.eqg", file: "obp_fob_tree.mod"}, // TODO: EQLO v4 .mod
		// .mod|1|arch.mod|dranik.eqg
		//{eqg: "dranik.eqg", file: "arch.mod"}, // PASS
		// .mod|1|aro.mod|aro.eqg
		//{eqg: "aro.eqg", file: "aro.mod"}, // PASS
		// .mod|1|col_b04.mod|b04.eqg b04.eqg pfs import: readMod col_b04.mod: material shader not found
		//{eqg: "b04.eqg", file: "col_b04.mod"}, // PASS
		// .mod|2|boulder_lg.mod|broodlands.eqg
		//{eqg: "broodlands.eqg", file: "boulder_lg.mod"}, // PASS
		// .mod|2|et_door01.mod|stillmoona.eqg
		//{eqg: "stillmoona.eqg", file: "et_door01.mod"}, // PASS
		// .mod|3|.mod|paperbaghat.eqg
		//{eqg: "paperbaghat.eqg", file: ".mod"}, // PASS
		// .mod|3|it11409.mod|undequip.eqg
		//{eqg: "undequip.eqg", file: "it11409.mod"}, // PASS
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

			fmt.Println(tt.file)

			mod := &Mod{}
			err = mod.Read(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			err = os.WriteFile(testDir+"/"+tt.file+".src.mod", data, 0644)
			if err != nil {
				t.Fatalf("os.WriteFile() error = %v", err)
			}

			buf := bytes.NewBuffer(nil)
			err = mod.Write(buf)
			if err != nil {
				t.Fatalf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			}

			mod2 := &Mod{}
			err = mod2.Read(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			buf2 := bytes.NewBuffer(nil)
			err = mod2.Write(buf2)
			if err != nil {
				t.Fatalf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			}

			err = os.WriteFile(testDir+"/"+tt.file+".dst.mod", buf2.Bytes(), 0644)
			if err != nil {
				t.Fatalf("os.WriteFile() error = %v", err)
			}

			if len(mod.Materials) != len(mod2.Materials) {
				t.Fatalf("Materials mismatch, got %d, expected %d", len(mod.Materials), len(mod2.Materials))
			}

			if len(mod.Vertices) != len(mod2.Vertices) {
				t.Fatalf("Vertices mismatch, got %d, expected %d", len(mod.Vertices), len(mod2.Vertices))
			}

			if len(mod.Faces) != len(mod2.Faces) {
				t.Fatalf("Triangles mismatch, got %d, expected %d", len(mod.Faces), len(mod2.Faces))
			}

			for i := 0; i < len(mod.Faces); i++ {
				if mod.Faces[i].MaterialName != mod.Faces[i].MaterialName {
					t.Fatalf("face %d material name mismatch %s vs %s", i, mod.Faces[i].MaterialName, mod.Faces[i].MaterialName)
				}
			}

			if len(mod.Bones) != len(mod2.Bones) {
				t.Fatalf("Bones mismatch, got %d, expected %d", len(mod.Bones), len(mod2.Bones))
			}

			err = helper.ByteCompareTest(mod.name.data(), mod2.name.data())
			if err != nil {
				t.Fatalf("Name data mismatch: %s", err.Error())
			}

			err = helper.ByteCompareTest(buf.Bytes(), buf2.Bytes())
			if err != nil {
				t.Fatalf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
