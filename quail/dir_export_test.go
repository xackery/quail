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
		Models []*common.Model
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

	os.RemoveAll("test")
	os.MkdirAll("test", 0755)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quail := &Quail{
				Models: tt.fields.Models,
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
		Models []*common.Model
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
		//{name: "valid", args: args{srcPath: "it13968.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "broodlands.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "freportn_chr.s3d"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "gequip4.s3d"}, wantErr: false},
		{name: "valid", args: args{srcPath: "crushbone.s3d"}, wantErr: false},
	}

	os.RemoveAll("test")
	os.MkdirAll("test", 0755)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quail := New()
			quail.Models = tt.fields.Models

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
		Models []*common.Model
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
		{name: "valid", args: args{srcPath: "it13900.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "it12095.eqg"},  wantErr: false},
		//{name: "valid", args: args{srcPath: "dbx.eqg"}, wantErr: false},
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
			e := New()
			e.Models = tt.fields.Models

			if err := e.PFSImport(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Errorf("Quail.ImportPFS() error = %+v", err)
			}

			if err := e.DirExport("test" + "/" + tt.args.srcPath); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ExportDir() error = %+v, wantErr %+v", err, tt.wantErr)
			}

			e2 := New()

			if err := e2.DirImport("test" + "/" + tt.args.srcPath[0:len(tt.args.srcPath)-4] + ".quail"); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ImportDir() error = %+v, wantErr %+v", err, tt.wantErr)
			}

			if len(e.Models) != len(e2.Models) {
				t.Fatalf("model count mismatch, %d != %d", len(e.Models), len(e2.Models))
			}

			for i, model := range e.Models {
				if model.Name != e2.Models[i].Name {
					t.Fatalf("model name mismatch, %s != %s", model.Name, e2.Models[i].Name)
				}

				if len(model.Vertices) != len(e2.Models[i].Vertices) {
					t.Fatalf("model vertex count mismatch, %d != %d", len(model.Vertices), len(e2.Models[i].Vertices))
				}

				/*for j, vert := range model.Vertices {
					if !reflect.DeepEqual(vert, e2.Models[i].Vertices[j]) {
						t.Fatalf("model vertex %d mismatch, %+v != %+v", j, vert, e2.Models[i].Vertices[j])
					}
				}*/

				if len(model.Triangles) != len(e2.Models[i].Triangles) {
					t.Fatalf("model triangle count mismatch, %d != %d", len(model.Triangles), len(e2.Models[i].Triangles))
				}

				for j, tri := range model.Triangles {
					if !reflect.DeepEqual(tri, e2.Models[i].Triangles[j]) {
						t.Fatalf("model triangle mismatch, %+v != %+v", tri, e2.Models[i].Triangles[j])
					}
				}

				if len(model.Materials) != len(e2.Models[i].Materials) {
					t.Fatalf("model material count mismatch, %d != %d", len(model.Materials), len(e2.Models[i].Materials))
				}

				/*for j, mat := range model.Materials {
					if !reflect.DeepEqual(mat, e2.Models[i].Materials[j]) {
						t.Fatalf("model material mismatch, %+v != %+v", mat, e2.Models[i].Materials[j])
					}
				}*/

				if len(model.Bones) != len(e2.Models[i].Bones) {
					t.Fatalf("model bone count mismatch, %d != %d", len(model.Bones), len(e2.Models[i].Bones))
				}

				/*for j, bone := range model.Bones {
					if !reflect.DeepEqual(bone, e2.Models[i].Bones[j]) {
						t.Fatalf("model bone mismatch, %+v != %+v", bone, e2.Models[i].Bones[j])
					}
				}*/

				if len(model.ParticlePoints) != len(e2.Models[i].ParticlePoints) {
					t.Fatalf("model particle point count mismatch, %d != %d", len(model.ParticlePoints), len(e2.Models[i].ParticlePoints))
				}

				/* for j, pp := range model.ParticlePoints {
					if !reflect.DeepEqual(pp, e2.Models[i].ParticlePoints[j]) {
						t.Fatalf("model particle point mismatch, %+v != %+v", pp, e2.Models[i].ParticlePoints[j])
					}
				} */

				if len(model.ParticleRenders) != len(e2.Models[i].ParticleRenders) {
					t.Fatalf("model particle render count mismatch, %d != %d", len(model.ParticleRenders), len(e2.Models[i].ParticleRenders))
				}

				for j, pr := range model.ParticleRenders {

					for k, entry := range pr.Entries {

						cmp := e2.Models[i].ParticleRenders[j].Entries[k]

						log.Debugf("%+v vs %+v", entry, cmp)
						if entry.Duration != cmp.Duration {
							t.Fatalf("model particle render entry duration mismatch, %+v != %+v", entry.Duration, cmp.Duration)
						}

						if entry.ID != cmp.ID {
							t.Fatalf("model particle render entry id mismatch, %+v != %+v", entry.ID, cmp.ID)
						}

						if entry.ID2 != cmp.ID2 {
							t.Fatalf("model particle render entry id2 mismatch, %+v != %+v", entry.ID2, cmp.ID2)
						}

						if entry.ParticlePoint != cmp.ParticlePoint {
							t.Fatalf("model particle render entry particle point mismatch, %+v != %+v", entry.ParticlePoint, cmp.ParticlePoint)
						}

						/*if reflect.DeepEqual(entry.ParticlePointSuffix, cmp.ParticlePointSuffix) {
							t.Fatalf("model particle render entry particle point suffix mismatch, %+v != %+v", entry.ParticlePointSuffix, cmp.ParticlePointSuffix)
						}*/

						if entry.UnknownA1 != cmp.UnknownA1 {
							t.Fatalf("model particle render entry unknown a1 mismatch, %+v != %+v", entry.UnknownA1, cmp.UnknownA1)
						}

						if entry.UnknownA2 != cmp.UnknownA2 {
							t.Fatalf("model particle render entry unknown a2 mismatch, %+v != %+v", entry.UnknownA2, cmp.UnknownA2)
						}

						if entry.UnknownA3 != cmp.UnknownA3 {
							t.Fatalf("model particle render entry unknown a3 mismatch, %+v != %+v", entry.UnknownA3, cmp.UnknownA3)
						}

						if entry.UnknownA4 != cmp.UnknownA4 {
							t.Fatalf("model particle render entry unknown a4 mismatch, %+v != %+v", entry.UnknownA4, cmp.UnknownA4)
						}

						if entry.UnknownA5 != cmp.UnknownA5 {
							t.Fatalf("model particle render entry unknown a5 mismatch, %+v != %+v", entry.UnknownA5, cmp.UnknownA5)
						}

						if entry.UnknownB != cmp.UnknownB {
							t.Fatalf("model particle render entry unknown b mismatch, %+v != %+v", entry.UnknownB, cmp.UnknownB)
						}

						if entry.UnknownC != cmp.UnknownC {
							t.Fatalf("model particle render entry unknown c mismatch, %+v != %+v", entry.UnknownC, cmp.UnknownC)
						}

						if entry.UnknownFFFFFFFF != cmp.UnknownFFFFFFFF {
							t.Fatalf("model particle render entry unknown ffffffff mismatch, %+v != %+v", entry.UnknownFFFFFFFF, cmp.UnknownFFFFFFFF)
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
