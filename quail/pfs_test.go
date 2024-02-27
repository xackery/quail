package quail

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
)

func TestQuail_PfsRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	type fields struct {
		Models []*common.Model
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//{name: "invalid", args: args{path: "invalid.txt"}, wantErr: true},
		//{name: "valid", args: args{path: "it13900.eqg"}, wantErr: false},
		//{name: "valid", args: args{path: "broodlands.eqg"}, wantErr: false},
		//{name: "valid", args: args{path: "freportn_chr.s3d"}, wantErr: false},
		//{name: "valid", args: args{path: "crushbone.s3d"}, wantErr: false},
		//{name: "valid", args: args{path: "it13968.eqg"}, wantErr: false},
		//{name: "valid", args: args{path: "dbx.eqg"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Quail{
				Models: tt.fields.Models,
			}
			if err := e.PfsRead(eqPath + "/" + tt.args.path); (err != nil) != tt.wantErr {
				t.Fatalf("Quail.ImportPfs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuail_PfsWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()
	type fields struct {
		Models []*common.Model
	}
	type args struct {
		fileVersion uint32
		pfsVersion  int
		srcPath     string
	}
	tests := []struct {
		fields  fields
		args    args
		wantErr bool
	}{
		//{args: args{srcPath: "it13900.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{args: args{srcPath: "dbx.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{args: args{srcPath: "freportn_chr.s3d", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{args: args{srcPath: "bloodfields.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{args: args{srcPath: "qeynos2.s3d", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{args: args{srcPath: "wrm.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.args.srcPath, func(t *testing.T) {
			e := &Quail{
				Models: tt.fields.Models,
			}

			if err := e.PfsRead(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Fatalf("pfs import %s error = %v", tt.args.srcPath, err)
			}

			//e.Models[0].Bones = []def.Bone{}
			//e.Models[0].Animations = []def.BoneAnimation{}

			if err := e.PfsWrite(tt.args.fileVersion, tt.args.pfsVersion, dirTest+"/"+filepath.Base(tt.args.srcPath)); (err != nil) != tt.wantErr {
				t.Fatalf("pfs export %s error = %v, wantErr %v", tt.args.srcPath, err, tt.wantErr)
			}
		})
	}
}

func TestQuail_PfsWriteImportExport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()
	type fields struct {
		Models []*common.Model
	}
	type args struct {
		fileVersion uint32
		pfsVersion  int
		srcPath     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//{name: "invalid", args: args{srcPath: "invalid.txt", fileVersion: 1, pfsVersion: 1}, wantErr: true},
		//{name: "load-save", args: args{srcPath: "it13900.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "dbx.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "mnt.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "mnt.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "mnt.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "bloodfields.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "bazaar.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "i9.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "it13968.eqg", fileVersion: 1, pfsVersion: 1}, wantErr: false},
		//{name: "load-save", args: args{srcPath: "freportn_chr.s3d", fileVersion: 1, pfsVersion: 1}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Quail{
				Models: tt.fields.Models,
			}

			if err := e.PfsRead(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Fatalf("Quail.ImportPfs() error = %v", err)
			}

			//e.Models[0].Bones = []def.Bone{}
			//e.Models[0].Animations = []def.BoneAnimation{}

			if err := e.PfsWrite(tt.args.fileVersion, tt.args.pfsVersion, dirTest+"/"+filepath.Base(tt.args.srcPath)); (err != nil) != tt.wantErr {
				t.Fatalf("Quail.ExportPfs() error = %v, wantErr %v", err, tt.wantErr)
			}

			e2 := &Quail{
				Models: tt.fields.Models,
			}

			if err := e2.PfsRead(dirTest + "/" + filepath.Base(tt.args.srcPath)); err != nil {
				t.Fatalf("Quail.ImportPfs() error = %v", err)
			}

			//e2.Models[0].Bones = []def.Bone{}
			//e2.Models[0].Animations = []def.BoneAnimation{}

			if len(e.Models) != len(e2.Models) {
				t.Fatalf("model count mismatch, %d != %d", len(e.Models), len(e2.Models))
			}

			for i, model := range e.Models {
				if model.Header.Name != e2.Models[i].Header.Name {
					t.Fatalf("model name mismatch, %s != %s", model.Header.Name, e2.Models[i].Header.Name)
				}

				if len(model.Vertices) != len(e2.Models[i].Vertices) {
					t.Fatalf("model vertex count mismatch, %d != %d", len(model.Vertices), len(e2.Models[i].Vertices))
				}

				for j, vert := range model.Vertices {
					if !reflect.DeepEqual(vert, e2.Models[i].Vertices[j]) {
						t.Fatalf("model vertex mismatch, %v != %v", vert, e2.Models[i].Vertices[j])
					}
				}

				if len(model.Triangles) != len(e2.Models[i].Triangles) {
					t.Fatalf("model triangle count mismatch, %d != %d", len(model.Triangles), len(e2.Models[i].Triangles))
				}

				for j, tri := range model.Triangles {
					if !reflect.DeepEqual(tri, e2.Models[i].Triangles[j]) {
						t.Fatalf("model triangle mismatch, %v != %v", tri, e2.Models[i].Triangles[j])
					}
				}

				if len(model.Materials) != len(e2.Models[i].Materials) {
					t.Fatalf("model material count mismatch, %d != %d", len(model.Materials), len(e2.Models[i].Materials))
				}

				for j, mat := range model.Materials {
					if !reflect.DeepEqual(mat, e2.Models[i].Materials[j]) {
						t.Fatalf("model material mismatch, %v != %v", mat, e2.Models[i].Materials[j])
					}
				}

				if len(model.Bones) != len(e2.Models[i].Bones) {
					t.Fatalf("model bone count mismatch, %d != %d", len(model.Bones), len(e2.Models[i].Bones))
				}

				for j, bone := range model.Bones {
					if !reflect.DeepEqual(bone, e2.Models[i].Bones[j]) {
						t.Fatalf("model bone mismatch, %v != %v", bone, e2.Models[i].Bones[j])
					}
				}

				if len(model.ParticlePoints) != len(e2.Models[i].ParticlePoints) {
					t.Fatalf("model particle point count mismatch, %d != %d", len(model.ParticlePoints), len(e2.Models[i].ParticlePoints))
				}

				for j, pp := range model.ParticlePoints {
					if !reflect.DeepEqual(pp, e2.Models[i].ParticlePoints[j]) {
						t.Fatalf("model particle point mismatch, %v != %v", pp, e2.Models[i].ParticlePoints[j])
					}
				}

				if len(model.ParticleRenders) != len(e2.Models[i].ParticleRenders) {
					t.Fatalf("model particle render count mismatch, %d != %d", len(model.ParticleRenders), len(e2.Models[i].ParticleRenders))
				}

				for j, pr := range model.ParticleRenders {

					for k, entry := range pr.Entries {

						cmp := e2.Models[i].ParticleRenders[j].Entries[k]

						log.Debugf("%v vs %v", entry, cmp)
						if entry.Duration != cmp.Duration {
							t.Fatalf("model particle render entry duration mismatch, %v != %v", entry.Duration, cmp.Duration)
						}

						if entry.ID != cmp.ID {
							t.Fatalf("model particle render entry id mismatch, %v != %v", entry.ID, cmp.ID)
						}

						if entry.ID2 != cmp.ID2 {
							t.Fatalf("model particle render entry id2 mismatch, %v != %v", entry.ID2, cmp.ID2)
						}

						if entry.ParticlePoint != cmp.ParticlePoint {
							t.Fatalf("model particle render entry particle point mismatch, %v != %v", entry.ParticlePoint, cmp.ParticlePoint)
						}

						/*if reflect.DeepEqual(entry.ParticlePointSuffix, cmp.ParticlePointSuffix) {
							t.Fatalf("model particle render entry particle point suffix mismatch, %v != %v", entry.ParticlePointSuffix, cmp.ParticlePointSuffix)
						}*/

						if entry.UnknownA1 != cmp.UnknownA1 {
							t.Fatalf("model particle render entry unknown a1 mismatch, %v != %v", entry.UnknownA1, cmp.UnknownA1)
						}

						if entry.UnknownA2 != cmp.UnknownA2 {
							t.Fatalf("model particle render entry unknown a2 mismatch, %v != %v", entry.UnknownA2, cmp.UnknownA2)
						}

						if entry.UnknownA3 != cmp.UnknownA3 {
							t.Fatalf("model particle render entry unknown a3 mismatch, %v != %v", entry.UnknownA3, cmp.UnknownA3)
						}

						if entry.UnknownA4 != cmp.UnknownA4 {
							t.Fatalf("model particle render entry unknown a4 mismatch, %v != %v", entry.UnknownA4, cmp.UnknownA4)
						}

						if entry.UnknownA5 != cmp.UnknownA5 {
							t.Fatalf("model particle render entry unknown a5 mismatch, %v != %v", entry.UnknownA5, cmp.UnknownA5)
						}

						if entry.UnknownB != cmp.UnknownB {
							t.Fatalf("model particle render entry unknown b mismatch, %v != %v", entry.UnknownB, cmp.UnknownB)
						}

						if entry.UnknownC != cmp.UnknownC {
							t.Fatalf("model particle render entry unknown c mismatch, %v != %v", entry.UnknownC, cmp.UnknownC)
						}

						if entry.UnknownFFFFFFFF != cmp.UnknownFFFFFFFF {
							t.Fatalf("model particle render entry unknown ffffffff mismatch, %v != %v", entry.UnknownFFFFFFFF, cmp.UnknownFFFFFFFF)
						}

					}
				}

			}

			if len(e.Animations) != len(e2.Animations) {
				t.Fatalf("animation count mismatch, %d != %d", len(e.Animations), len(e2.Animations))
			}

			for _, anim := range e.Animations {
				isFound := false
				for _, anim2 := range e2.Animations {
					if anim.Header.Name != anim2.Header.Name {
						continue
					}
					isFound = true

					if len(anim.Bones) != len(anim2.Bones) {
						t.Fatalf("animation bone count mismatch, %d != %d", len(anim.Bones), len(anim2.Bones))
					}

					for j, bone := range anim.Bones {
						if !reflect.DeepEqual(bone, anim2.Bones[j]) {
							t.Fatalf("animation bone mismatch, %v != %v", bone, anim2.Bones[j])
						}
					}

					break
				}
				if !isFound {
					t.Fatalf("animation name mismatch, %s not found", anim.Header.Name)
				}
			}

			log.Debugf("seems like a clean roundtrip for %s", filepath.Base(tt.args.srcPath))
		})
	}
}
