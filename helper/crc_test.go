package helper

import (
	"testing"
)

func TestFilenameCRC32(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{name: string("FilenameCRC32"), args: args{name: "test"}, want: 1537663841},
		{name: "Test 2", args: args{name: "hello world"}, want: 2533725502},
		{name: "Test 3", args: args{name: "12345"}, want: 742322399},
		{name: "Test 4", args: args{name: "test.txt"}, want: 2138351979},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilenameCRC32(tt.args.name); got != tt.want {
				t.Fatalf("FilenameCRC32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: string("Validate"), args: args{data: []byte("test")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.data); (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateCRC16(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name         string
		args         args
		wantChecksum uint16
		wantErr      bool
	}{
		{name: string("GenerateCRC16"), args: args{data: []byte("test")}, wantChecksum: 32268, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChecksum, err := GenerateCRC16(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateCRC16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotChecksum != tt.wantChecksum {
				t.Fatalf("GenerateCRC16() = %v, want %v", gotChecksum, tt.wantChecksum)
			}
		})
	}
}

func TestGenerateCRC32(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name         string
		args         args
		wantChecksum uint32
		wantErr      bool
	}{
		{name: string("GenerateCRC32"), args: args{data: []byte("test")}, wantChecksum: 32268, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChecksum, err := GenerateCRC32(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateCRC32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotChecksum != tt.wantChecksum {
				t.Fatalf("GenerateCRC32() = %v, want %v", gotChecksum, tt.wantChecksum)
			}
		})
	}
}
