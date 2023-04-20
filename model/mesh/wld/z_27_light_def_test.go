package wld

import (
	"testing"
)

func TestWLD_lightDefRead(t *testing.T) {
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
		27,             //fragCode
		-1,             //fragIndex
		e.lightDefRead) //callback
}
