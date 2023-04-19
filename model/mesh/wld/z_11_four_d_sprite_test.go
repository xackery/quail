package wld

import (
	"testing"
)

// TOO: no refs
func TestWLD_fourDSpriteRead(t *testing.T) {
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
		11,                //fragCode
		-1,                //fragIndex
		e.fourDSpriteRead) //callback
}
