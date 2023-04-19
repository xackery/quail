package wld

import (
	"testing"
)

// TODO: all skeletontrackrefs are 0, yet flags seem to fluxuate
func TestWLD_skeletonTrackRead(t *testing.T) {
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
		17,                  //fragCode
		-1,                  //fragIndex
		e.skeletonTrackRead) //callback
}
