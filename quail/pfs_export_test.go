package quail

import (
	"os"
	"reflect"
	"testing"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/quail/def"
)

func TestQuail_PFSExport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	type fields struct {
		Meshes []*def.Mesh
	}
	type args struct {
		fileVersion uint32
		pfsVersion  int
		dstPath     string
		srcPath     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//{name: "invalid", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "invalid.txt"}, wantErr: true},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "it13900.eqg", dstPath: "test/it13900.eqg"}, wantErr: false},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "dbx.eqg", dstPath: "test/dbx.eqg"}, wantErr: false},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "freportn_chr.s3d", dstPath: "test/freport_n.eqg"}, wantErr: false},
		{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "bloodfields.eqg", dstPath: "test/bloodfields.eqg"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Quail{
				Meshes: tt.fields.Meshes,
			}

			if err := e.PFSImport(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Errorf("Quail.ImportPFS() error = %v", err)
			}

			//e.Meshes[0].Bones = []def.Bone{}
			//e.Meshes[0].Animations = []def.BoneAnimation{}

			if err := e.PFSExport(tt.args.fileVersion, tt.args.pfsVersion, tt.args.dstPath); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ExportPFS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuail_PFSExportImportExport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	type fields struct {
		Meshes []*def.Mesh
	}
	type args struct {
		fileVersion uint32
		pfsVersion  int
		dstPath     string
		srcPath     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//{name: "invalid", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "invalid.txt"}, wantErr: true},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "it13900.eqg", dstPath: "test/it13900.eqg"}, wantErr: false},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "dbx.eqg", dstPath: "test/dbx.eqg"}, wantErr: false},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "mnt.eqg", dstPath: "test/mnt.eqg"}, wantErr: false},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "mnt.eqg", dstPath: "test/mnt.eqg"}, wantErr: false},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "mnt.eqg", dstPath: "test/mnt.eqg"}, wantErr: false},
		{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "bloodfields.eqg", dstPath: "test/bloodfields.eqg"}, wantErr: false},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "i27.eqg", dstPath: "test/i27.eqg"}, wantErr: false},
		//{name: "load-save", args: args{fileVersion: 1, pfsVersion: 1, srcPath: "freportn_chr.s3d", dstPath: "test/freport_n.eqg"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Quail{
				Meshes: tt.fields.Meshes,
			}

			if err := e.PFSImport(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Errorf("Quail.ImportPFS() error = %v", err)
			}

			//e.Meshes[0].Bones = []def.Bone{}
			//e.Meshes[0].Animations = []def.BoneAnimation{}

			if err := e.PFSExport(tt.args.fileVersion, tt.args.pfsVersion, tt.args.dstPath); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ExportPFS() error = %v, wantErr %v", err, tt.wantErr)
			}

			e2 := &Quail{
				Meshes: tt.fields.Meshes,
			}

			if err := e2.PFSImport(tt.args.dstPath); err != nil {
				t.Errorf("Quail.ImportPFS() error = %v", err)
			}

			//e2.Meshes[0].Bones = []def.Bone{}
			//e2.Meshes[0].Animations = []def.BoneAnimation{}

			if len(e.Meshes) != len(e2.Meshes) {
				t.Fatalf("mesh count mismatch, %d != %d", len(e.Meshes), len(e2.Meshes))
			}

			for i, mesh := range e.Meshes {
				if mesh.Name != e2.Meshes[i].Name {
					t.Fatalf("mesh name mismatch, %s != %s", mesh.Name, e2.Meshes[i].Name)
				}

				if len(mesh.Vertices) != len(e2.Meshes[i].Vertices) {
					t.Fatalf("mesh vertex count mismatch, %d != %d", len(mesh.Vertices), len(e2.Meshes[i].Vertices))
				}

				for j, vert := range mesh.Vertices {
					if !reflect.DeepEqual(vert, e2.Meshes[i].Vertices[j]) {
						t.Fatalf("mesh vertex mismatch, %v != %v", vert, e2.Meshes[i].Vertices[j])
					}
				}

				if len(mesh.Triangles) != len(e2.Meshes[i].Triangles) {
					t.Fatalf("mesh triangle count mismatch, %d != %d", len(mesh.Triangles), len(e2.Meshes[i].Triangles))
				}

				for j, tri := range mesh.Triangles {
					if !reflect.DeepEqual(tri, e2.Meshes[i].Triangles[j]) {
						t.Fatalf("mesh triangle mismatch, %v != %v", tri, e2.Meshes[i].Triangles[j])
					}
				}

				if len(mesh.Materials) != len(e2.Meshes[i].Materials) {
					t.Fatalf("mesh material count mismatch, %d != %d", len(mesh.Materials), len(e2.Meshes[i].Materials))
				}

				for j, mat := range mesh.Materials {
					if !reflect.DeepEqual(mat, e2.Meshes[i].Materials[j]) {
						t.Fatalf("mesh material mismatch, %v != %v", mat, e2.Meshes[i].Materials[j])
					}
				}

				if len(mesh.Bones) != len(e2.Meshes[i].Bones) {
					t.Fatalf("mesh bone count mismatch, %d != %d", len(mesh.Bones), len(e2.Meshes[i].Bones))
				}

				for j, bone := range mesh.Bones {
					if !reflect.DeepEqual(bone, e2.Meshes[i].Bones[j]) {
						t.Fatalf("mesh bone mismatch, %v != %v", bone, e2.Meshes[i].Bones[j])
					}
				}

				if len(mesh.ParticlePoints) != len(e2.Meshes[i].ParticlePoints) {
					t.Fatalf("mesh particle point count mismatch, %d != %d", len(mesh.ParticlePoints), len(e2.Meshes[i].ParticlePoints))
				}

				for j, pp := range mesh.ParticlePoints {
					if !reflect.DeepEqual(pp, e2.Meshes[i].ParticlePoints[j]) {
						t.Fatalf("mesh particle point mismatch, %v != %v", pp, e2.Meshes[i].ParticlePoints[j])
					}
				}

				if len(mesh.ParticleRenders) != len(e2.Meshes[i].ParticleRenders) {
					t.Fatalf("mesh particle render count mismatch, %d != %d", len(mesh.ParticleRenders), len(e2.Meshes[i].ParticleRenders))
				}

				for j, pr := range mesh.ParticleRenders {

					for k, entry := range pr.Entries {

						cmp := e2.Meshes[i].ParticleRenders[j].Entries[k]

						log.Debugf("%v vs %v", entry, cmp)
						if entry.Duration != cmp.Duration {
							t.Fatalf("mesh particle render entry duration mismatch, %v != %v", entry.Duration, cmp.Duration)
						}

						if entry.ID != cmp.ID {
							t.Fatalf("mesh particle render entry id mismatch, %v != %v", entry.ID, cmp.ID)
						}

						if entry.ID2 != cmp.ID2 {
							t.Fatalf("mesh particle render entry id2 mismatch, %v != %v", entry.ID2, cmp.ID2)
						}

						if entry.ParticlePoint != cmp.ParticlePoint {
							t.Fatalf("mesh particle render entry particle point mismatch, %v != %v", entry.ParticlePoint, cmp.ParticlePoint)
						}

						/*if reflect.DeepEqual(entry.ParticlePointSuffix, cmp.ParticlePointSuffix) {
							t.Fatalf("mesh particle render entry particle point suffix mismatch, %v != %v", entry.ParticlePointSuffix, cmp.ParticlePointSuffix)
						}*/

						if entry.UnknownA != cmp.UnknownA {
							t.Fatalf("mesh particle render entry unknown a mismatch, %v != %v", entry.UnknownA, cmp.UnknownA)
						}

						if entry.UnknownB != cmp.UnknownB {
							t.Fatalf("mesh particle render entry unknown b mismatch, %v != %v", entry.UnknownB, cmp.UnknownB)
						}

						if entry.UnknownC != cmp.UnknownC {
							t.Fatalf("mesh particle render entry unknown c mismatch, %v != %v", entry.UnknownC, cmp.UnknownC)
						}

						if entry.UnknownFFFFFFFF != cmp.UnknownFFFFFFFF {
							t.Fatalf("mesh particle render entry unknown ffffffff mismatch, %v != %v", entry.UnknownFFFFFFFF, cmp.UnknownFFFFFFFF)
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
					if anim.Name != anim2.Name {
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
					t.Fatalf("animation name mismatch, %s not found", anim.Name)
				}
			}
		})
	}
}
