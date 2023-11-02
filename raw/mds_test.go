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
		{eqg: "aam.eqg", file: "aam.mds"},
		// .mds|1|ae3.mds|ae3.eqg
		{eqg: "ae3.eqg", file: "ae3.mds"},
		// .mds|1|bcn.mds|harbingers.eqg
		{eqg: "harbingers.eqg", file: "bcn.mds"},
		// .mds|1|stnd_gnome_wave.mds|it12095.eqg
		{eqg: "it12095.eqg", file: "stnd_gnome_wave.mds"},
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
		// .mds|1|aam.mds|aam.eqg
		{eqg: "aam.eqg", file: "aam.mds"},
		// .mds|1|ae3.mds|ae3.eqg
		{eqg: "ae3.eqg", file: "ae3.mds"},
		// .mds|1|bcn.mds|harbingers.eqg
		{eqg: "harbingers.eqg", file: "bcn.mds"},
		// .mds|1|stnd_gnome_wave.mds|it12095.eqg
		{eqg: "it12095.eqg", file: "stnd_gnome_wave.mds"},
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

			if len(mds.Materials) != len(mds2.Materials) {
				t.Errorf("Materials mismatch, got %d, expected %d", len(mds.Materials), len(mds2.Materials))
			}

			if len(mds.Vertices) != len(mds2.Vertices) {
				t.Errorf("Vertices mismatch, got %d, expected %d", len(mds.Vertices), len(mds2.Vertices))
			}

			if len(mds.Triangles) != len(mds2.Triangles) {
				t.Errorf("Triangles mismatch, got %d, expected %d", len(mds.Triangles), len(mds2.Triangles))
			}

			if len(mds.Bones) != len(mds2.Bones) {
				t.Errorf("Bones mismatch, got %d, expected %d", len(mds.Bones), len(mds2.Bones))
			}

		})
	}
}
