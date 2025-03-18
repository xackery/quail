package qfs

import (
	"io/fs"
	"reflect"
	"testing"
	"time"
)

func TestMFS_ReadFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		m       *MFS
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "file exists",
			m: func() *MFS {
				mfs := NewMemoryFS()
				mfs.WriteFile("test.txt", []byte("hello world"), 0644)
				return mfs
			}(),
			args:    args{name: "test.txt"},
			want:    []byte("hello world"),
			wantErr: false,
		},
		{
			name:    "file does not exist",
			m:       NewMemoryFS(),
			args:    args{name: "nonexistent.txt"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.ReadFile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MFS.ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MFS.ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMFS_Open(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		m       *MFS
		args    args
		want    fs.File
		wantErr bool
	}{
		{
			name: "file exists",
			m: func() *MFS {
				mfs := NewMemoryFS()
				mfs.WriteFile("test.txt", []byte("hello world"), 0644)
				return mfs
			}(),
			args:    args{name: "test.txt"},
			want:    &memReadWriter{data: []byte("hello world"), info: &fileInfo{name: "test.txt", size: 11, mode: 0644, modTime: time.Now()}},
			wantErr: false,
		},
		{
			name:    "file does not exist",
			m:       NewMemoryFS(),
			args:    args{name: "nonexistent.txt"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MFS.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MFS.Open() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMFS_Stat(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		m       *MFS
		args    args
		want    fs.FileInfo
		wantErr bool
	}{
		{
			name: "file exists",
			m: func() *MFS {
				mfs := NewMemoryFS()
				mfs.WriteFile("test.txt", []byte("hello world"), 0644)
				return mfs
			}(),
			args:    args{name: "test.txt"},
			want:    &fileInfo{name: "test.txt", size: 11, mode: 0644, modTime: time.Now()},
			wantErr: false,
		},
		{
			name:    "file does not exist",
			m:       NewMemoryFS(),
			args:    args{name: "nonexistent.txt"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Stat(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MFS.Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MFS.Stat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMFS_ReadDir(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		m       *MFS
		args    args
		want    []fs.DirEntry
		wantErr bool
	}{
		{
			name: "directory exists",
			m: func() *MFS {
				mfs := NewMemoryFS()
				mfs.WriteFile("dir/file1.txt", []byte("file1"), 0644)
				mfs.WriteFile("dir/file2.txt", []byte("file2"), 0644)
				return mfs
			}(),
			args: args{name: "dir"},
			want: []fs.DirEntry{
				&fileInfo{name: "dir/file1.txt", mode: 0644},
				&fileInfo{name: "dir/file2.txt", mode: 0644},
			},
			wantErr: false,
		},
		{
			name:    "directory does not exist",
			m:       NewMemoryFS(),
			args:    args{name: "nonexistent"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.ReadDir(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MFS.ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MFS.ReadDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMFS_RemoveAll(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		m       *MFS
		args    args
		wantErr bool
	}{
		{
			name: "file exists",
			m: func() *MFS {
				mfs := NewMemoryFS()
				mfs.WriteFile("test.txt", []byte("hello world"), 0644)
				return mfs
			}(),
			args:    args{name: "test.txt"},
			wantErr: false,
		},
		{
			name:    "file does not exist",
			m:       NewMemoryFS(),
			args:    args{name: "nonexistent.txt"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.RemoveAll(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("MFS.RemoveAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMFS_MkdirAll(t *testing.T) {
	type args struct {
		name string
		perm fs.FileMode
	}
	tests := []struct {
		name    string
		m       *MFS
		args    args
		wantErr bool
	}{
		{
			name:    "create directory",
			m:       NewMemoryFS(),
			args:    args{name: "dir", perm: 0755},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.MkdirAll(tt.args.name, tt.args.perm); (err != nil) != tt.wantErr {
				t.Errorf("MFS.MkdirAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMFS_WriteFile(t *testing.T) {
	type args struct {
		name string
		data []byte
		perm fs.FileMode
	}
	tests := []struct {
		name    string
		m       *MFS
		args    args
		wantErr bool
	}{
		{
			name:    "write file",
			m:       NewMemoryFS(),
			args:    args{name: "test.txt", data: []byte("hello world"), perm: 0644},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.WriteFile(tt.args.name, tt.args.data, tt.args.perm); (err != nil) != tt.wantErr {
				t.Errorf("MFS.WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memReader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		m       *memReadWriter
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "read data",
			m:       &memReadWriter{data: []byte("hello world")},
			args:    args{p: make([]byte, 5)},
			want:    5,
			wantErr: false,
		},
		{
			name:    "read past end",
			m:       &memReadWriter{data: []byte("hello")},
			args:    args{p: make([]byte, 10)},
			want:    5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("memReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("memReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memReader_Close(t *testing.T) {
	tests := []struct {
		name    string
		m       *memReadWriter
		wantErr bool
	}{
		{
			name:    "close reader",
			m:       &memReadWriter{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Close(); (err != nil) != tt.wantErr {
				t.Errorf("memReader.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memReader_Stat(t *testing.T) {
	tests := []struct {
		name    string
		m       *memReadWriter
		want    fs.FileInfo
		wantErr bool
	}{
		{
			name:    "get file info",
			m:       &memReadWriter{info: &fileInfo{name: "test.txt", size: 11, mode: 0644, modTime: time.Now()}},
			want:    &fileInfo{name: "test.txt", size: 11, mode: 0644, modTime: time.Now()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Stat()
			if (err != nil) != tt.wantErr {
				t.Errorf("memReader.Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memReader.Stat() = %v, want %v", got, tt.want)
			}
		})
	}
}
