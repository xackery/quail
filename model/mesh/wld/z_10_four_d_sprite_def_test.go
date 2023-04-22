package wld

import (
	"testing"
)

// TODO: no refs used in any s3d - Tacc
func TestWLD_fourDSpriteDefRead(t *testing.T) {
	e, err := New("test", nil)
	if err != nil {
		t.Fatalf("new: %v", err)
		return
	}
	fragmentTests(t,
		true, //single run stop
		[]string{
			"yxtta_obj.s3d",
		},
		10,                   //fragCode
		-1,                   //fragIndex
		e.fourDSpriteDefRead) //callback
}
