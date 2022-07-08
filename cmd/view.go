package cmd

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/tetra3d"
	"github.com/spf13/cobra"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
)

const screenWidth = 398
const screenHeight = 224

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a model or image",
	Long: `View an EverQuest asset

Supported extensions: mod
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
		defer func() {
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}()

		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path check: %w", err)
		}
		if fi.IsDir() {
			return fmt.Errorf("view requires a target file, directory provided")
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer f.Close()
		ext := strings.ToLower(filepath.Ext(path))

		//shortname := filepath.Base(path)
		//shortname = strings.TrimSuffix(shortname, filepath.Ext(shortname))
		type loader interface {
			Load(io.ReadSeeker) error
			common.GLTFExporter
			SetPath(string)
		}
		type loadTypes struct {
			instance  loader
			extension string
		}
		loads := []*loadTypes{
			//{instance: &ani.ANI{}, extension: ".ani"},
			//{instance: &eqg.EQG{}, extension: ".eqg"},
			{instance: &mod.MOD{}, extension: ".mod"},
			{instance: &ter.TER{}, extension: ".ter"},
			//{instance: &zon.ZON{}, extension: ".zon"},
		}

		for _, v := range loads {
			if ext != v.extension {
				continue
			}

			v.instance.SetPath(filepath.Dir(path))

			err = v.instance.Load(f)
			if err != nil {
				return fmt.Errorf("instance load %s: %w", v.extension, err)
			}

			buf := &bytes.Buffer{}
			err = v.instance.ExportGLTF(buf)
			if err != nil {
				return fmt.Errorf("export gltf %s: %w", v.extension, err)
			}

			v := &viewer{}
			err = v.load(buf.Bytes())
			if err != nil {
				return fmt.Errorf("load gltf to viewer: %w", err)
			}
			err = ebiten.RunGame(v)
			if err != nil {
				return fmt.Errorf("run viewer: %w", err)
			}
			return nil
		}

		return fmt.Errorf("open: unsupported extension %s on file %s", ext, filepath.Base(path))
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.PersistentFlags().String("path", "", "path to view")
}

type viewer struct {
	scene  *tetra3d.Scene
	camera *tetra3d.Camera
}

func (v *viewer) load(gltfData []byte) error {
	library, err := tetra3d.LoadGLTFData(gltfData, nil)
	if err != nil {
		return fmt.Errorf("loadGLTFData: %w", err)
	}
	v.scene = library.ExportedScene

	v.camera = tetra3d.NewCamera(screenWidth, screenHeight)
	v.camera.Move(0, 0, 12)
	v.scene.Root.AddChildren(v.camera)

	fmt.Println(v.scene.Root.HierarchyAsString())
	return nil
}

func (g *viewer) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		fmt.Println("q key was pressed, quitting")
		os.Exit(0)
	}
	return nil
}

func (g *viewer) Draw(screen *ebiten.Image) {

	// Call Camera.Clear() to clear its internal backing texture. This
	// should be called once per frame before drawing your Scene.
	g.camera.Clear()

	// Now we'll render the Scene from the camera. The Camera's ColorTexture will then
	// hold the result.

	// Below, we'll pass both the Scene and the scene root because 1) the Scene influences
	// how Models draw (fog, for example), and 2) we may not want to render
	// all Models.

	// Camera.RenderNodes() renders all Nodes in a tree, starting with the
	// Node specified. You can also use Camera.Render() to simply render a selection of
	// Models.
	g.camera.RenderNodes(g.scene, g.scene.Root)

	// Before drawing the result, clear the screen first.
	screen.Fill(color.RGBA{20, 30, 40, 255})

	// Draw the resulting texture to the screen, and you're done! You can
	// also visualize the depth texture with g.Camera.DepthTexture.
	screen.DrawImage(g.camera.ColorTexture(), nil)
}

func (g *viewer) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}
