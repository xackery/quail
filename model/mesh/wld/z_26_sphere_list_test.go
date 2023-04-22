package wld

import (
	"testing"
)

// TODO: no refs
func TestWLD_sphereListRead(t *testing.T) {
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
		26,               //fragCode
		-1,               //fragIndex
		e.sphereListRead) //callback
}
