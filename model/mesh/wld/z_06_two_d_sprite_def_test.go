package wld

import (
	"testing"
)

// texture names misaligned
func TestWLD_twoDSpriteDefRead(t *testing.T) {
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
		6,                   //fragCode
		-1,                  //fragIndex
		e.twoDSpriteDefRead) //callback
}
