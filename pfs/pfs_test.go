// eqg is a pfs archive for EverQuest
package pfs

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/xackery/quail/common"
)

func TestPFS_Add(t *testing.T) {
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
			e := &PFS{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			if err := e.Add(tt.args.name, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PFS.Add() error = %v, wantErr %v", err, tt.wantErr)
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
		want    *PFS
		wantErr bool
	}{
		{name: "success", args: args{name: "test"}, want: &PFS{name: "test"}, wantErr: false},
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
		want    *PFS
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

func TestPFS_Remove(t *testing.T) {
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
			e := &PFS{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			if err := e.Remove(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("PFS.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPFS_Extract(t *testing.T) {
	type fields struct {
		name            string
		files           []*FileEntry
		ContentsSummary string
		fileCount       int
	}
	dirTest := common.DirTest(t)
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
			e := &PFS{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			got, err := e.Extract(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PFS.Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("PFS.Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPFS_File(t *testing.T) {
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
			e := &PFS{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			got, err := e.File(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("PFS.File() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PFS.File() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPFS_Close(t *testing.T) {
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
			e := &PFS{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			if err := e.Close(); (err != nil) != tt.wantErr {
				t.Errorf("PFS.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPFS_WriteFile(t *testing.T) {
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
			e := &PFS{
				name:            tt.fields.name,
				files:           tt.fields.files,
				ContentsSummary: tt.fields.ContentsSummary,
				fileCount:       tt.fields.fileCount,
			}
			if err := e.SetFile(tt.args.name, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PFS.SetFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
