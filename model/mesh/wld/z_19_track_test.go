package wld

import (
	"testing"
)

func TestWLD_trackRead(t *testing.T) {
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
		19,          //fragCode
		-1,          //fragIndex
		e.trackRead) //callback
}
