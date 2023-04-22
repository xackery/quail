package wld

import (
	"testing"
)

func TestWLD_particleCloudDefRead(t *testing.T) {
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
		52,                     //fragCode
		-1,                     //fragIndex
		e.particleCloudDefRead) //callback
}
