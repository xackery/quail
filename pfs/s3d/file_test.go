package s3d

import (
	"reflect"
	"testing"

	"github.com/xackery/quail/pfs/archive"
)

func TestS3D_File(t *testing.T) {
	type fields struct {
		name      string
		ShortName string
		files     []archive.Filer
		fileCount int
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
		errStr  string
	}{
		{name: string("File"), fields: fields{name: "test"}, args: args{name: "test"}, errStr: "test not found", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &S3D{
				name:      tt.fields.name,
				ShortName: tt.fields.ShortName,
				files:     tt.fields.files,
				fileCount: tt.fields.fileCount,
			}
			got, err := e.File(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("S3D.File() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errStr {
				t.Errorf("S3D.File() error = %v, want %v", err, tt.errStr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S3D.File() = %v, want %v", got, tt.want)
			}
		})
	}
}
