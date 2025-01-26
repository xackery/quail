package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestAniRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()
	type args struct {
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// .ani|1|sidl_ba_1_tln.ani|tln.eqg
		{name: "tln.eqg"}, // PASS
		// .ani|2|stnd_ba_1_exo.ani|exo.eqg eye_chr.s3d pfs import: s3d load: decode: dirName for crc 655939147 not found
		{name: "exo.eqg"}, // PASS
		// .ani|2|walk_ba_1_vaf.ani|vaf.eqg valdeholm.eqg pfs import: eqg load: decode: read nameData unexpected EOF
		{name: "vaf.eqg"}, // PASS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open pfs %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".ani" {
					continue
				}
				ani := &Ani{}
				err = ani.Read(bytes.NewReader(file.Data()))
				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}
			}
		})
	}
}

func TestAniWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()

	tests := []struct {
		name    string
		wantErr bool
	}{
		// .ani|1|sidl_ba_1_tln.ani|tln.eqg
		//{name: "tln.eqg"}, // PASS
		// .ani|2|stnd_ba_1_exo.ani|exo.eqg
		//{name: "exo.eqg"}, // PASS
		// .ani|2|walk_ba_1_vaf.ani|vaf.eqg
		//{name: "vaf.eqg"}, // PASS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".ani" {
					continue
				}
				ani := &Ani{}
				err = ani.Read(bytes.NewReader(file.Data()))

				if err != nil {
					os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = ani.Write(buf)
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				ani2 := &Ani{}
				err = ani2.Read(bytes.NewReader(buf.Bytes()))
				if err != nil {
					t.Fatalf("failed to write %s: %s", tt.name, err.Error())
				}

				if ani.IsStrict != ani2.IsStrict {
					t.Fatalf("IsStrict mismatch: %t != %t", ani.IsStrict, ani2.IsStrict)
				}

				if ani.Version != ani2.Version {
					t.Fatalf("Version mismatch: %d != %d", ani.Version, ani2.Version)
				}

				if len(ani.Bones) != len(ani2.Bones) {
					t.Fatalf("Bone count mismatch: %d != %d", len(ani.Bones), len(ani2.Bones))
				}

				for i, bone := range ani.Bones {
					if len(bone.Frames) != len(ani2.Bones[i].Frames) {
						t.Fatalf("Bone frame count mismatch: %d != %d", len(bone.Frames), len(ani2.Bones[i].Frames))
					}

					if bone.Name != ani2.Bones[i].Name {
						t.Fatalf("Bone name mismatch: %s != %s", bone.Name, ani2.Bones[i].Name)
					}
				}

			}
		})
	}
}
