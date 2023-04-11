package blender

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func Test_import(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "it13926.eqg", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if err := import_blender(&cobra.Command{}, []string{fmt.Sprintf("%s/%s", eqPath, tt.name), "test/out/"}); (err != nil) != tt.wantErr {
			//	t.Errorf("import() error = %v, wantErr %v", err, tt.wantErr)
			//}
			err := import_blender(&cobra.Command{}, []string{fmt.Sprintf("%s/_%s", "test", tt.name), fmt.Sprintf("test/%s", tt.name)})
			if err != nil {
				t.Errorf("import() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
