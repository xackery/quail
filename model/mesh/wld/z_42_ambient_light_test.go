package wld

import (
	"testing"
)

// TODO: no refs
func TestWLD_ambientLightRead(t *testing.T) {
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
		42,                 //fragCode
		-1,                 //fragIndex
		e.ambientLightRead) //callback
}
