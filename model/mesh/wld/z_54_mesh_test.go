package wld

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/log"
)

func TestWLD_meshRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	tests := []struct {
		name       string
		isNewWorld bool
		wantErr    bool
	}{
		{name: "gequip.s3d", isNewWorld: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e, err := New(tt.name, nil)
			if err != nil {
				t.Errorf("%s new error = %v", tt.name, err)
				return
			}
			e.isOldWorld = !tt.isNewWorld

			count, err := parseFragments(fmt.Sprintf("%s/test_data/%s", eqPath, tt.name), 54, e.meshRead)
			if err != nil && !tt.wantErr {
				t.Errorf("%s parseFragment error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			log.Debugf("%s total parsed: %d", tt.name, count)
		})
	}
}
