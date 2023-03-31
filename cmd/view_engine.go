package cmd

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kvartborg/vector"
	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"
	"golang.org/x/image/font/basicfont"
)

type viewEngine struct {
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

func (v *viewEngine) load(gltfData []byte) error {
	var err error

	v.library, err = tetra3d.LoadGLTFData(gltfData, &tetra3d.GLTFLoadOptions{
		DefaultToAutoTransparency: true,
	})
	if err != nil {
		return fmt.Errorf("loadGLTFData: %w", err)
	}
	v.scene = v.library.ExportedScene

	v.camera = tetra3d.NewCamera(v.width, v.height)
	v.camera.Move(0, 0, 12)
	v.scene.Root.AddChildren(v.camera)
	v.scene.World.LightingOn = false

	fmt.Println(v.scene.Root.HierarchyAsString())
	return nil
}

func (v *viewEngine) Update() error {

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

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		v.drawDebugWireframe = !v.drawDebugWireframe
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		v.drawDebugNormals = !v.drawDebugNormals
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		v.drawDebugCenters = !v.drawDebugCenters
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyY) {
		v.drawDebugDepth = !v.drawDebugDepth
	}

	boost := float64(1)
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		boost = 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
		(ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)) {
		pos = pos.Add(forward.Scale(moveSpd * boost))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		pos = pos.Add(right.Scale(moveSpd * boost))
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		pos = pos.Add(forward.Scale(-moveSpd * boost))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		pos = pos.Add(right.Scale(-moveSpd * boost))
	}

	v.camera.SetLocalPosition(pos.X, pos.Y, pos.Z)

	// Rotating the camera with the mouse

	// Rotate and tilt the camera according to mouse movements
	mx, my := ebiten.CursorPosition()

	mv := vector.Vector{float64(mx), float64(my)}

	diff := mv.Sub(v.prevMousePosition)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

		v.cameraTilt -= diff[1] * 0.005
		v.cameraRotate -= diff[0] * 0.005

		v.cameraTilt = math.Max(math.Min(v.cameraTilt, math.Pi/2-0.1), -math.Pi/2+0.1)

		tilt := tetra3d.NewMatrix4Rotate(1, 0, 0, v.cameraTilt)
		rotate := tetra3d.NewMatrix4Rotate(0, 1, 0, v.cameraRotate)
		// Order of this is important - tilt * rotate works, rotate * tilt does not, lol
		v.camera.SetLocalRotation(tilt.Mult(rotate))
	}
	v.prevMousePosition = mv.Clone()

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

func (v *viewEngine) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{60, 70, 80, 255})

	v.camera.Clear()

	v.camera.RenderNodes(v.scene, v.scene.Root)
	opt := &ebiten.DrawImageOptions{}
	w, h := v.camera.ColorTexture().Size()
	opt.GeoM.Scale(float64(v.width)/float64(w), float64(v.height)/float64(h))

	if v.drawDebugDepth {
		screen.DrawImage(v.camera.DepthTexture(), nil)
	} else {
		screen.DrawImage(v.camera.ColorTexture(), nil)
	}

	if v.drawDebugWireframe {
		v.camera.DrawDebugWireframe(screen, v.scene.Root, colors.Red())
	}

	if v.drawDebugNormals {
		v.camera.DrawDebugNormals(screen, v.scene.Root, 0.5, colors.Blue())
	}

	if v.drawDebugCenters {
		v.camera.DrawDebugCenters(screen, v.scene.Root, colors.SkyBlue())
	}

	if v.drawDebugText {

		v.camera.DrawDebugRenderInfo(screen, 1, colors.White())
		//1 Key: Play [SmoothRoll] Animation On Table
		//2 Key: Play [StepRoll] Animation on Table Note the animations can blend\n
		//F Key: Play Animation on Skinned Mesh\n Note that the nodes move as well

		txt := fmt.Sprintf(`WASD/Arrows: Move, Mouse: Look
C: Center
T: DebugText
V: Wireframe
Y: Depth
N: Normals
Q: Quit
Pos: %0.2f %0.2f %0.2f`, v.camera.LocalPosition().X, v.camera.LocalPosition().Y, v.camera.LocalPosition().Z)
		text.Draw(screen, txt, basicfont.Face7x13, 0, 140, color.RGBA{255, 0, 0, 255})
	}
}

func (v *viewEngine) Layout(w, h int) (int, int) {
	return v.width, v.height
}
