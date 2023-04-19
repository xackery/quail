package wld

import (
	"testing"
)

// TODO: no examples found
func TestWLD_userDataRead(t *testing.T) {
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
		1,              //fragCode
		-1,             //fragIndex
		e.userDataRead) //callback
}
