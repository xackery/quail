package wld

import (
	"fmt"
	"os"
	"testing"
)

func TestWLD_textureListRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		true, //single run stop
		[]string{
			"gequip.s3d",
		},
		3,                 //fragCode
		-1,                //fragIndex
		e.textureListRead) //callback
}

func TestWLD_textureListReadWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	tests := []struct {
		name       string
		fragOffset int
		wantErr    bool
	}{
		{name: "gequip.s3d", fragOffset: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := New(tt.name, nil)
			if err != nil {
				t.Fatalf("%s new error = %v", tt.name, err)
			}

			err = compareReadAndWrite(t, fmt.Sprintf("%s/_test_data/%s/", eqPath, tt.name), 3, tt.fragOffset, e)
			if err != nil && !tt.wantErr {
				t.Errorf("%s compareReadAndWrite error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

		})
	}
}
