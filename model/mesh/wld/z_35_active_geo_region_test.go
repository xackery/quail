package wld

import (
	"testing"
)

// TODO: no refs used in any s3d - Tacc
func TestWLD_activeGeoRegionRead(t *testing.T) {
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
		35,                    //fragCode
		-1,                    //fragIndex
		e.activeGeoRegionRead) //callback
}
