package helper

import (
	"reflect"
	"testing"
)

func TestWriteStringHash(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"success", args{"test"}, []byte{225, 95, 182, 94}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteStringHash(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WriteStringHash() = %v, want %v", got, tt.want)
			}

			if got := ReadStringHash(tt.want); got != tt.args.hash {
				t.Errorf("ReadStringHash() = %v, want %v", got, tt.args.hash)
			}
		})
	}
}
