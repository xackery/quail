package wld

import (
	"testing"
)

func TestWLD_threeDSpriteRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		true, //single run stop
		[]string{
			"gfaydark.s3d",
		},
		9,                  //fragCode
		-1,                 //fragIndex
		e.threeDSpriteRead) //callback
}
