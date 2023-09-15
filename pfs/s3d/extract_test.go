package s3d

import (
	"testing"

	"github.com/xackery/quail/pfs/archive"
)

func TestS3D_Extract(t *testing.T) {
	type fields struct {
		name      string
		ShortName string
		files     []archive.Filer
		fileCount int
	}
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
		{name: string("Extract"), fields: fields{name: "test"}, args: args{path: "test"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &S3D{
				name:      tt.fields.name,
				ShortName: tt.fields.ShortName,
				files:     tt.fields.files,
				fileCount: tt.fields.fileCount,
			}
			got, err := e.Extract(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("S3D.Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Fatalf("S3D.Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}
