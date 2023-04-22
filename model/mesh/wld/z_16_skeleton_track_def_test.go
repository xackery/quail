package wld

import (
	"testing"
)

// TODO: broken at skinToAnimRefs
func TestWLD_skeletonTrackDefRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		false, //single run stop
		[]string{
			"gequip.s3d",
		},
		16,                     //fragCode
		-1,                     //fragIndex
		e.skeletonTrackDefRead) //callback
}
