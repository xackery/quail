package wld

import (
	"testing"
)

func TestWLD_zoneRead(t *testing.T) {
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
		41,         //fragCode
		-1,         //fragIndex
		e.zoneRead) //callback
}
