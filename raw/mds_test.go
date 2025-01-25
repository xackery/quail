package raw

import (
	"bytes"
	"os"
	"testing"

	"github.com/xackery/quail/pfs"
)

func TestMdsRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		eqg     string
		file    string
		wantErr bool
	}{
		// .mds|1|aam.mds|aam.eqg
		//{eqg: "aam.eqg", file: "aam.mds"}, // PASS
		// .mds|1|ae3.mds|ae3.eqg
		//{eqg: "ae3.eqg", file: "ae3.mds"}, // PASS
		//{eqg: "djf.eqg", file: "djf.mds"}, // PASS
		{eqg: "mrd.eqg", file: "mrd.mds"}, // PASS
		// .mds|1|bcn.mds|harbingers.eqg
		//{eqg: "harbingers.eqg", file: "bcn.mds"}, // PASS
		// .mds|1|stnd_gnome_wave.mds|it12095.eqg
		//{eqg: "it12095.eqg", file: "stnd_gnome_wave.mds"}, // PASS
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

			mds := &Mds{}
			err = mds.Read(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMdsWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		eqg     string
		file    string
		wantErr bool
	}{
		{eqg: "mrd.eqg", file: "mrd.mds"}, // PASS
		// .mds|1|aam.mds|aam.eqg
		//{eqg: "aam.eqg", file: "aam.mds"}, // PASS
		// .mds|1|ae3.mds|ae3.eqg
		//{eqg: "ae3.eqg", file: "ae3.mds"}, // PASS
		// .mds|1|bcn.mds|harbingers.eqg
		//{eqg: "harbingers.eqg", file: "bcn.mds"}, // PASS
		// .mds|1|stnd_gnome_wave.mds|it12095.eqg
		//{eqg: "it12095.eqg", file: "stnd_gnome_wave.mds"}, // PASS
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

			mds := &Mds{}
			err = mds.Read(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			// srcNames := []string{}

			// chunk := []byte{}
			// for _, b := range mds.nameBuf {
			// 	if b != 0 {
			// 		chunk = append(chunk, b)
			// 		continue
			// 	}
			// 	srcNames = append(srcNames, string(chunk))
			// 	chunk = []byte{}
			// }

			buf := bytes.NewBuffer(nil)
			err = mds.Write(buf)
			if err != nil {
				t.Fatalf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			}

			mds2 := &Mds{}
			err = mds2.Read(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			// dstNames := []string{}
			// chunk = []byte{}
			// for _, b := range mds.nameBuf {
			// 	if b != 0 {
			// 		chunk = append(chunk, b)
			// 		continue
			// 	}
			// 	dstNames = append(dstNames, string(chunk))
			// 	chunk = []byte{}
			// }

			// sort.Strings(srcNames)
			// sort.Strings(dstNames)
			// for i := 0; i < len(srcNames); i++ {
			// 	if len(srcNames) > i && len(dstNames) > i {
			// 		fmt.Printf("%d src: %s, dst: %s\n", i, srcNames[i], dstNames[i])
			// 	}
			// }

			// if len(srcNames) != len(dstNames) {
			// 	t.Errorf("Name count mismatch, got %d, expected %d", len(srcNames), len(dstNames))
			// }

			if len(mds.Materials) != len(mds2.Materials) {
				t.Errorf("Materials mismatch, got %d, expected %d", len(mds.Materials), len(mds2.Materials))
			}

			if len(mds.Bones) != len(mds2.Bones) {
				t.Errorf("Bones mismatch, got %d, expected %d", len(mds.Bones), len(mds2.Bones))
			}

			if len(mds.Models) != len(mds2.Models) {
				t.Errorf("Models mismatch, got %d, expected %d", len(mds.Models), len(mds.Models))
			}

			// if len(mds.NameData()) != len(mds2.NameData()) {
			// 	t.Errorf("NameData mismatch, got %d, expected %d", len(mds2.NameData()), len(mds.NameData()))
			// }

		})
	}
}
