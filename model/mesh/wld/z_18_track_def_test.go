package wld

import (
	"testing"
)

// TODO: all skeletonCounts 0 on test
func TestWLD_trackDefRead(t *testing.T) {
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
		18,             //fragCode
		-1,             //fragIndex
		e.trackDefRead) //callback
}
