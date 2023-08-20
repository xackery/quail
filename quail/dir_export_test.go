package quail

import (
	"os"
	"reflect"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
)

func TestQuail_DirExport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	type fields struct {
		Meshes []*common.Model
	}
	type args struct {
		srcPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//{name: "invalid", args: args{srcPath: "invalid.txt"}}, wantErr: true},
		//{name: "valid", args: args{srcPath: "it13900.eqg"},  wantErr: false},
		//{name: "valid", args: args{srcPath: "it12095.eqg"},  wantErr: false},
		//{name: "valid", args: args{srcPath: "box.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "broodlands.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "thuledream.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "feerrott2.eqg"}, wantErr: false}, //type 4 zone
		//{name: "valid", args: args{srcPath: "arena2.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "steamfontmts.eqg"}, wantErr: false},
		{name: "valid", args: args{srcPath: "bazaar.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "freportn_chr.s3d"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "gequip4.s3d"}, wantErr: false},
	}

	for _, tt := range tests {
		os.RemoveAll("test")
		os.MkdirAll("test", 0755)
		t.Run(tt.name, func(t *testing.T) {
			quail := &Quail{
				Models: tt.fields.Meshes,
			}

			if err := quail.PFSImport(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Errorf("Quail.ImportPFS() error = %+v", err)
			}

			if err := quail.DirExport("test" + "/" + tt.args.srcPath); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ExportDir() error = %+v, wantErr %+v", err, tt.wantErr)
			}
		})
	}
}

func TestQuail_DirExportPFS(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	type fields struct {
		Meshes []*common.Model
	}
	type args struct {
		srcPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//{name: "invalid", args: args{srcPath: "invalid.txt"}}, wantErr: true},
		//{name: "valid", args: args{srcPath: "it13900.eqg"},  wantErr: false},
		//{name: "valid", args: args{srcPath: "it12095.eqg"},  wantErr: false},
		{name: "valid", args: args{srcPath: "it13968.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "broodlands.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "freportn_chr.s3d"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "gequip4.s3d"}, wantErr: false},
	}

	for _, tt := range tests {
		os.RemoveAll("test")
		os.MkdirAll("test", 0755)
		t.Run(tt.name, func(t *testing.T) {
			quail := &Quail{
				Models: tt.fields.Meshes,
			}

			if err := quail.PFSImport(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Errorf("Quail.ImportPFS() error = %+v", err)
			}

			if err := quail.DirExport("test" + "/" + tt.args.srcPath); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ExportDir() error = %+v, wantErr %+v", err, tt.wantErr)
			}

			newPath := "test" + "/" + tt.args.srcPath[0:len(tt.args.srcPath)-4] + ".quail"
			if err := quail.DirImport(newPath); err != nil {
				t.Errorf("Quail.ImportDir() error = %+v", err)
			}

			if err := quail.PFSExport(1, 1, "test"+"/"+tt.args.srcPath); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ExportPFS() error = %+v, wantErr %+v", err, tt.wantErr)
			}

		})
	}
}

func TestQuail_DirExportImport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	type fields struct {
		Meshes []*common.Model
	}
	type args struct {
		srcPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//{name: "invalid", args: args{srcPath: "invalid.txt"}}, wantErr: true},
		//{name: "valid", args: args{srcPath: "it13900.eqg"},  wantErr: false},
		//{name: "valid", args: args{srcPath: "it12095.eqg"},  wantErr: false},
		{name: "valid", args: args{srcPath: "dbx.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "broodlands.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "freportn_chr.s3d"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "gequip4.s3d"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "dbx.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "it13968.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "i27.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "plane.eqg"}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Quail{
				Models: tt.fields.Meshes,
			}

			if err := e.PFSImport(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Errorf("Quail.ImportPFS() error = %+v", err)
			}

			if err := e.DirExport("test" + "/" + tt.args.srcPath); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ExportDir() error = %+v, wantErr %+v", err, tt.wantErr)
			}

			e2 := &Quail{}

			if err := e2.DirImport("test" + "/" + tt.args.srcPath[0:len(tt.args.srcPath)-4] + ".quail"); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ImportDir() error = %+v, wantErr %+v", err, tt.wantErr)
			}

			if len(e.Models) != len(e2.Models) {
				t.Fatalf("mesh count mismatch, %d != %d", len(e.Models), len(e2.Models))
			}

			for i, mesh := range e.Models {
				if mesh.Name != e2.Models[i].Name {
					t.Fatalf("mesh name mismatch, %s != %s", mesh.Name, e2.Models[i].Name)
				}

				if len(mesh.Vertices) != len(e2.Models[i].Vertices) {
					t.Fatalf("mesh vertex count mismatch, %d != %d", len(mesh.Vertices), len(e2.Models[i].Vertices))
				}

				/*for j, vert := range mesh.Vertices {
					if !reflect.DeepEqual(vert, e2.Meshes[i].Vertices[j]) {
						t.Fatalf("mesh vertex %d mismatch, %+v != %+v", j, vert, e2.Meshes[i].Vertices[j])
					}
				}*/

				if len(mesh.Triangles) != len(e2.Models[i].Triangles) {
					t.Fatalf("mesh triangle count mismatch, %d != %d", len(mesh.Triangles), len(e2.Models[i].Triangles))
				}

				for j, tri := range mesh.Triangles {
					if !reflect.DeepEqual(tri, e2.Models[i].Triangles[j]) {
						t.Fatalf("mesh triangle mismatch, %+v != %+v", tri, e2.Models[i].Triangles[j])
					}
				}

				if len(mesh.Materials) != len(e2.Models[i].Materials) {
					t.Fatalf("mesh material count mismatch, %d != %d", len(mesh.Materials), len(e2.Models[i].Materials))
				}

				/*for j, mat := range mesh.Materials {
					if !reflect.DeepEqual(mat, e2.Meshes[i].Materials[j]) {
						t.Fatalf("mesh material mismatch, %+v != %+v", mat, e2.Meshes[i].Materials[j])
					}
				}*/

				if len(mesh.Bones) != len(e2.Models[i].Bones) {
					t.Fatalf("mesh bone count mismatch, %d != %d", len(mesh.Bones), len(e2.Models[i].Bones))
				}

				/*for j, bone := range mesh.Bones {
					if !reflect.DeepEqual(bone, e2.Meshes[i].Bones[j]) {
						t.Fatalf("mesh bone mismatch, %+v != %+v", bone, e2.Meshes[i].Bones[j])
					}
				}*/

				if len(mesh.ParticlePoints) != len(e2.Models[i].ParticlePoints) {
					t.Fatalf("mesh particle point count mismatch, %d != %d", len(mesh.ParticlePoints), len(e2.Models[i].ParticlePoints))
				}

				/* for j, pp := range mesh.ParticlePoints {
					if !reflect.DeepEqual(pp, e2.Meshes[i].ParticlePoints[j]) {
						t.Fatalf("mesh particle point mismatch, %+v != %+v", pp, e2.Meshes[i].ParticlePoints[j])
					}
				} */

				if len(mesh.ParticleRenders) != len(e2.Models[i].ParticleRenders) {
					t.Fatalf("mesh particle render count mismatch, %d != %d", len(mesh.ParticleRenders), len(e2.Models[i].ParticleRenders))
				}

				for j, pr := range mesh.ParticleRenders {

					for k, entry := range pr.Entries {

						cmp := e2.Models[i].ParticleRenders[j].Entries[k]

						log.Debugf("%+v vs %+v", entry, cmp)
						if entry.Duration != cmp.Duration {
							t.Fatalf("mesh particle render entry duration mismatch, %+v != %+v", entry.Duration, cmp.Duration)
						}

						if entry.ID != cmp.ID {
							t.Fatalf("mesh particle render entry id mismatch, %+v != %+v", entry.ID, cmp.ID)
						}

						if entry.ID2 != cmp.ID2 {
							t.Fatalf("mesh particle render entry id2 mismatch, %+v != %+v", entry.ID2, cmp.ID2)
						}

						if entry.ParticlePoint != cmp.ParticlePoint {
							t.Fatalf("mesh particle render entry particle point mismatch, %+v != %+v", entry.ParticlePoint, cmp.ParticlePoint)
						}

						/*if reflect.DeepEqual(entry.ParticlePointSuffix, cmp.ParticlePointSuffix) {
							t.Fatalf("mesh particle render entry particle point suffix mismatch, %+v != %+v", entry.ParticlePointSuffix, cmp.ParticlePointSuffix)
						}*/

						if entry.UnknownA1 != cmp.UnknownA1 {
							t.Fatalf("mesh particle render entry unknown a1 mismatch, %+v != %+v", entry.UnknownA1, cmp.UnknownA1)
						}

						if entry.UnknownA2 != cmp.UnknownA2 {
							t.Fatalf("mesh particle render entry unknown a2 mismatch, %+v != %+v", entry.UnknownA2, cmp.UnknownA2)
						}

						if entry.UnknownA3 != cmp.UnknownA3 {
							t.Fatalf("mesh particle render entry unknown a3 mismatch, %+v != %+v", entry.UnknownA3, cmp.UnknownA3)
						}

						if entry.UnknownA4 != cmp.UnknownA4 {
							t.Fatalf("mesh particle render entry unknown a4 mismatch, %+v != %+v", entry.UnknownA4, cmp.UnknownA4)
						}

						if entry.UnknownA5 != cmp.UnknownA5 {
							t.Fatalf("mesh particle render entry unknown a5 mismatch, %+v != %+v", entry.UnknownA5, cmp.UnknownA5)
						}

						if entry.UnknownB != cmp.UnknownB {
							t.Fatalf("mesh particle render entry unknown b mismatch, %+v != %+v", entry.UnknownB, cmp.UnknownB)
						}

						if entry.UnknownC != cmp.UnknownC {
							t.Fatalf("mesh particle render entry unknown c mismatch, %+v != %+v", entry.UnknownC, cmp.UnknownC)
						}

						if entry.UnknownFFFFFFFF != cmp.UnknownFFFFFFFF {
							t.Fatalf("mesh particle render entry unknown ffffffff mismatch, %+v != %+v", entry.UnknownFFFFFFFF, cmp.UnknownFFFFFFFF)
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
							t.Fatalf("animation bone mismatch, %+v != %+v", bone, anim2.Bones[j])
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
