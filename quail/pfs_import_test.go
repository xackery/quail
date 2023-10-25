package quail

import (
	"os"
	"testing"

	"github.com/xackery/quail/common"
)

func TestQuail_PFSImport(t *testing.T) {
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
			if err := e.PFSImport(eqPath + "/" + tt.args.path); (err != nil) != tt.wantErr {
				t.Fatalf("Quail.ImportPFS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
