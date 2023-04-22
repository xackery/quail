package wld

import (
	"testing"
)

// TODO: no refs
func TestWLD_pointLightOldRead(t *testing.T) {
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
		29,                  //fragCode
		-1,                  //fragIndex
		e.pointLightOldRead) //callback
}
