package blender

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func Test_export(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "steamfontmts", args: args{cmd: &cobra.Command{}, args: []string{fmt.Sprintf("%s/steamfontmts.eqg", eqPath), "test/"}}, wantErr: false},
		{name: "it13900", args: args{cmd: &cobra.Command{}, args: []string{fmt.Sprintf("%s/it13900.eqg", eqPath), "test/"}}, wantErr: false},
		{name: "xhf", args: args{cmd: &cobra.Command{}, args: []string{fmt.Sprintf("%s/xhf.eqg", eqPath), "test/"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := export(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("export() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
