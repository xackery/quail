package cmd

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kvartborg/vector"
	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"
	"github.com/spf13/cobra"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
	"golang.org/x/image/font/basicfont"
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

		if ext == ".gltf" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			v := &viewer{
				width:         796,
				height:        448,
				drawDebugText: true,
			}
			err = v.load(data)
			if err != nil {
				return fmt.Errorf("load gltf to viewer: %w", err)
			}
			err = ebiten.RunGame(v)
			if err != nil {
				return fmt.Errorf("run viewer: %w", err)
			}
			return nil
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

			v := &viewer{
				width:         796,
				height:        448,
				drawDebugText: true,
			}

			//data, _ := ioutil.ReadFile("example.gltf")
			//err = v.load(data)
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
	width, height     int
	library           *tetra3d.Library
	scene             *tetra3d.Scene
	camera            *tetra3d.Camera
	cameraTilt        float64
	cameraRotate      float64
	prevMousePosition vector.Vector

	drawDebugText      bool
	drawDebugDepth     bool
	drawDebugNormals   bool
	drawDebugCenters   bool
	drawDebugWireframe bool
}

func (v *viewer) load(gltfData []byte) error {
	var err error
	v.library, err = tetra3d.LoadGLTFData(gltfData, nil)
	if err != nil {
		return fmt.Errorf("loadGLTFData: %w", err)
	}
	v.scene = v.library.ExportedScene

	v.camera = tetra3d.NewCamera(v.width, v.height)
	v.camera.Move(0, 0, 12)
	v.scene.Root.AddChildren(v.camera)
	v.scene.LightingOn = false

	fmt.Println(v.scene.Root.HierarchyAsString())
	return nil
}

func (v *viewer) Update() error {

	moveSpd := 0.05
	if ebiten.IsKeyPressed(ebiten.KeyQ) || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		fmt.Println("q or escape key was pressed, quitting")
		os.Exit(0)
	}
	// We use Camera.Rotation.Forward().Invert() because the camera looks down -Z (so its forward vector is inverted)
	forward := v.camera.LocalRotation().Forward().Invert()
	right := v.camera.LocalRotation().Right()

	pos := v.camera.LocalPosition()

	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		v.drawDebugWireframe = !v.drawDebugWireframe
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		v.drawDebugText = !v.drawDebugText
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		v.drawDebugNormals = !v.drawDebugNormals
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		pos = pos.Add(forward.Scale(moveSpd))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		pos = pos.Add(right.Scale(moveSpd))
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		pos = pos.Add(forward.Scale(-moveSpd))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		pos = pos.Add(right.Scale(-moveSpd))
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		pos[1] += moveSpd
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		pos[1] -= moveSpd
	}

	v.camera.SetLocalPosition(pos)

	// Rotating the camera with the mouse

	// Rotate and tilt the camera according to mouse movements
	/*mx, my := ebiten.CursorPosition()

	mv := vector.Vector{float64(mx), float64(my)}

	diff := mv.Sub(v.prevMousePosition)

	v.cameraTilt -= diff[1] * 0.005
	v.cameraRotate -= diff[0] * 0.005

	v.cameraTilt = math.Max(math.Min(v.cameraTilt, math.Pi/2-0.1), -math.Pi/2+0.1)

	tilt := tetra3d.NewMatrix4Rotate(1, 0, 0, v.cameraTilt)
	rotate := tetra3d.NewMatrix4Rotate(0, 1, 0, v.cameraRotate)

	// Order of this is important - tilt * rotate works, rotate * tilt does not, lol
	v.camera.SetLocalRotation(tilt.Mult(rotate))

	v.prevMousePosition = mv.Clone()
	*/
	/*	armature := v.scene.Root.ChildrenRecursive().ByName("Armature", true)[0].(*tetra3d.Node)
		armature.Rotate(0, 1, 0, 0.01)

		if inpututil.IsKeyJustPressed(ebiten.KeyF) {
			// armature.AnimationPlayer.FinishMode = tetra3d.FinishModeStop
			armature.AnimationPlayer().Play(v.library.Animations["ArmatureAction"])
		}

		armature.AnimationPlayer().Update(1.0 / 60)

		table := v.scene.Root.Get("Table").(*tetra3d.Model)
		table.AnimationPlayer().BlendTime = 0.1
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			table.AnimationPlayer().Play(v.library.Animations["SmoothRoll"])
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			table.AnimationPlayer().Play(v.library.Animations["StepRoll"])
		}

		table.AnimationPlayer().Update(1.0 / 60)
	*/
	return nil
}

func (v *viewer) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{60, 70, 80, 255})

	v.camera.Clear()

	v.camera.RenderNodes(v.scene, v.scene.Root)
	opt := &ebiten.DrawImageOptions{}
	w, h := v.camera.ColorTexture().Size()
	opt.GeoM.Scale(float64(v.width)/float64(w), float64(v.height)/float64(h))
	screen.DrawImage(v.camera.ColorTexture(), nil)

	if v.drawDebugWireframe {
		v.camera.DrawDebugWireframe(screen, v.scene.Root, colors.Red())
	}

	if v.drawDebugNormals {
		v.camera.DrawDebugNormals(screen, v.scene.Root, 0.5, colors.Blue())
	}

	if v.drawDebugText {
		v.camera.DrawDebugText(screen, 1, colors.White())
		txt := "T to toggle this text\nWASD: Move, Mouse: Look\n1 Key: Play [SmoothRoll] Animation On Table\n2 Key: Play [StepRoll] Animation on Table\nNote the animations can blend\nF Key: Play Animation on Skinned Mesh\nNote that the nodes move as well\nF4: Toggle fullscreen\nF6: Node Debug View\nQ: Quit"
		text.Draw(screen, txt, basicfont.Face7x13, 0, 140, color.RGBA{255, 0, 0, 255})
	}
}

func (v *viewer) Layout(w, h int) (int, int) {
	return v.width, v.height
}
