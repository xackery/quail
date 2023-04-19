package wld

import (
	"testing"
)

func TestWLD_textureRead(t *testing.T) {
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
		4,             //fragCode
		-1,            //fragIndex
		e.textureRead) //callback
}
