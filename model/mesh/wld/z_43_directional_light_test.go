package wld

import (
	"testing"
)

// TODO: no refs
func TestWLD_directionalLightRead(t *testing.T) {
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
		43,                     //fragCode
		-1,                     //fragIndex
		e.directionalLightRead) //callback
}
