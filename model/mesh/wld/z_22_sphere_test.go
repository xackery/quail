package wld

import (
	"testing"
)

func TestWLD_sphereRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		true, //single run stop
		[]string{
			"qeynos2.s3d",
		},
		22,           //fragCode
		-1,           //fragIndex
		e.sphereRead) //callback
}
