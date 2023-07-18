package quail

import (
	"os"
	"testing"

	"github.com/xackery/quail/quail/def"
)

func TestQuail_DirExport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	type fields struct {
		Meshes []*def.Mesh
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
		//{name: "valid", args: args{srcPath: "dbx.eqg"}, wantErr: false},
		{name: "valid", args: args{srcPath: "broodlands.eqg"}, wantErr: false},
		//{name: "valid", args: args{srcPath: "freportn_chr.s3d"},  wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quail := &Quail{
				Meshes: tt.fields.Meshes,
			}

			if err := quail.PFSImport(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Errorf("Quail.ImportPFS() error = %v", err)
			}

			if err := quail.DirExport("test" + "/" + tt.args.srcPath); (err != nil) != tt.wantErr {
				t.Errorf("Quail.ExportDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
