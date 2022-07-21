package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestView(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	path := "test/eq/arena.eqg"
	file := ""

	buf := bytes.NewBuffer(nil)
	err := viewLoad(buf, path, file)
	if err != nil {
		t.Fatalf("viewLoad: %s", err)
	}

	view := &viewEngine{
		width:         796,
		height:        448,
		drawDebugText: true,
	}

	err = view.load(buf.Bytes())
	if err != nil {
		t.Fatalf("decode gltf to viewer: %s", err)
	}
	err = ebiten.RunGame(view)
	if err != nil {
		t.Fatalf("run viewer: %s", err)
	}
}
