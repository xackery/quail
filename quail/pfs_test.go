package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/helper"
)

func TestQuail_PfsRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
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
			e := &Quail{}
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
	dirTest := helper.DirTest()
	type args struct {
		fileVersion uint32
		pfsVersion  int
		srcPath     string
	}
	tests := []struct {
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
			e := &Quail{}

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
	dirTest := helper.DirTest()
	type args struct {
		fileVersion uint32
		pfsVersion  int
		srcPath     string
	}
	tests := []struct {
		name    string
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
			e := &Quail{}

			if err := e.PfsRead(eqPath + "/" + tt.args.srcPath); err != nil {
				t.Fatalf("Quail.ImportPfs() error = %v", err)
			}

			//e.Models[0].Bones = []def.Bone{}
			//e.Models[0].Animations = []def.BoneAnimation{}

			if err := e.PfsWrite(tt.args.fileVersion, tt.args.pfsVersion, dirTest+"/"+filepath.Base(tt.args.srcPath)); (err != nil) != tt.wantErr {
				t.Fatalf("Quail.ExportPfs() error = %v, wantErr %v", err, tt.wantErr)
			}

			e2 := &Quail{}

			if err := e2.PfsRead(dirTest + "/" + filepath.Base(tt.args.srcPath)); err != nil {
				t.Fatalf("Quail.ImportPfs() error = %v", err)
			}

			//e2.Models[0].Bones = []def.Bone{}
			//e2.Models[0].Animations = []def.BoneAnimation{}

			fmt.Printf("seems like a clean roundtrip for %s\n", filepath.Base(tt.args.srcPath))
		})
	}
}
