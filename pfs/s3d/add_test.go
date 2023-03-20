package s3d

import (
	"testing"

	"github.com/xackery/quail/common"
)

func TestS3D_Add(t *testing.T) {
	type fields struct {
		name      string
		ShortName string
		files     []common.Filer
		fileCount int
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
		{name: string("Add"), fields: fields{name: "test"}, args: args{name: "test", data: []byte("test")}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &S3D{
				name:      tt.fields.name,
				ShortName: tt.fields.ShortName,
				files:     tt.fields.files,
				fileCount: tt.fields.fileCount,
			}
			if err := e.Add(tt.args.name, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("S3D.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
