package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a model or image",
	Long: `View an EverQuest asset

Supported extensions: gltf, mod, ter
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("parse path: %w", err)
		}
		if path == "" {
			if len(args) < 1 {
				return cmd.Usage()
			}
			path = args[0]
		}

		file, err := cmd.Flags().GetString("file")
		defer func() {
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}()
		if file == "" {
			if len(args) >= 2 {
				file = args[1]
			}
		}

		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path check: %w", err)
		}
		if fi.IsDir() {
			return fmt.Errorf("view requires a target file, directory provided")
		}

		buf := bytes.NewBuffer(nil)
		err = viewLoad(buf, path, file)
		if err != nil {
			return fmt.Errorf("viewLoad: %w", err)
		}

		view := &viewEngine{
			width:         796,
			height:        448,
			drawDebugText: true,
		}

		err = view.load(buf.Bytes())
		if err != nil {
			return fmt.Errorf("decode gltf to viewer: %w", err)
		}

		err = ebiten.RunGame(view)
		if err != nil {
			return fmt.Errorf("run viewer: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.PersistentFlags().String("path", "", "path to view")
	viewCmd.PersistentFlags().String("file", "", "file to view (in eqg/s3d)")
}
