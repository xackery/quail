// eqg is a pfs archive for EverQuest
package pfs

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/xackery/quail/common"
)

func TestPfs_Add(t *testing.T) {
	type fields struct {
		name            string
		files           []*FileEntry
		ContentsSummary string
		fileCount       int
	}
	type args struct {
		name string
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{}, args: args{name: "test", data: []byte("test")}, wantErr: false},
		{name: "fail", fields: fields{}, args: args{name: "-", data: []byte("test")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Pfs{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			if err := e.Add(tt.args.name, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Pfs.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *Pfs
		wantErr bool
	}{
		{name: "success", args: args{name: "test"}, want: &Pfs{name: "test"}, wantErr: false},
		{name: "fail", args: args{name: ""}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Pfs
		wantErr bool
	}{
		{name: "success", args: args{path: ""}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPfs_Remove(t *testing.T) {
	type fields struct {
		name            string
		files           []*FileEntry
		ContentsSummary string
		fileCount       int
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{files: []*FileEntry{{name: "test"}}}, args: args{name: "test"}, wantErr: false},
		{name: "fail", args: args{name: "test"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Pfs{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			if err := e.Remove(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Pfs.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPfs_Extract(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	type fields struct {
		name            string
		files           []*FileEntry
		ContentsSummary string
		fileCount       int
	}
	dirTest := common.DirTest()
	fmt.Println(dirTest)
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{name: "success", fields: fields{files: []*FileEntry{{name: "test", data: []byte("test")}}}, args: args{path: dirTest + "/test.pfs"}, want: "extracted 1", wantErr: false},
		{name: "fail", args: args{path: dirTest + "/test.pfs"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Pfs{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			got, err := e.Extract(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pfs.Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("Pfs.Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPfs_File(t *testing.T) {
	type fields struct {
		name            string
		files           []*FileEntry
		ContentsSummary string
		fileCount       int
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{"success", fields{files: []*FileEntry{{name: "test", data: []byte("test")}}}, args{name: "test"}, []byte("test"), false},
		{"fail", fields{files: []*FileEntry{{name: "test", data: []byte("test")}}}, args{name: ""}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Pfs{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			got, err := e.File(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pfs.File() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pfs.File() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPfs_Close(t *testing.T) {
	type fields struct {
		name            string
		files           []*FileEntry
		ContentsSummary string
		fileCount       int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "success", fields: fields{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Pfs{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			if err := e.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Pfs.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPfs_WriteFile(t *testing.T) {
	type fields struct {
		name            string
		files           []*FileEntry
		ContentsSummary string
		fileCount       int
	}

	type args struct {
		name string
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{}, args: args{name: "test", data: []byte("test")}, wantErr: false},
		{name: "fail", fields: fields{}, args: args{name: "-", data: []byte("test")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Pfs{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			if err := e.SetFile(tt.args.name, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Pfs.SetFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkEQG(b *testing.B) {
	eqgPath := os.Getenv("EQ_PATH")
	if eqgPath == "" {
		b.Skip("EQ_PATH not set")
	}

	for i := 0; i < b.N; i++ {
		pfs, err := NewFile(fmt.Sprintf("%s/xhf.eqg", eqgPath))
		if err != nil {
			b.Fatalf("Failed newfile: %s", err.Error())
		}
		pfs.Close()
	}
}
