package wld

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/log"
)

// TODO: no refs
func TestWLD_worldTreeRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "gequip.s3d", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e, err := New(tt.name, nil)
			if err != nil {
				t.Errorf("%s new error = %v", tt.name, err)
				return
			}

			count, err := parseFragments(fmt.Sprintf("%s/test_data/%s", eqPath, tt.name), 33, e.worldTreeRead)
			if err != nil && !tt.wantErr {
				t.Errorf("%s parseFragment error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			log.Debugf("%s total parsed: %d", tt.name, count)
		})
	}
}
