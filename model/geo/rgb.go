package geo

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/helper"
)

// RGB represents R,G,B as uint8
type RGB struct {
	R uint8
	G uint8
	B uint8
}

// AToRGB converts a string to a RGB
func AToRGB(s string) *RGB {
	parts := strings.Split(s, ",")
	if len(parts) < 3 {
		return nil
	}
	return &RGB{
		R: helper.AtoU8(parts[0]),
		G: helper.AtoU8(parts[1]),
		B: helper.AtoU8(parts[2]),
	}
}

// String returns a string representation of the RGB
func (e *RGB) String() string {
	return fmt.Sprintf("%d,%d,%d", e.R, e.G, e.B)
}

// RGBA represents R,G,B,A as uint8
type RGBA struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// AtoRGBA converts a string to a RGBA
func AtoRGBA(s string) RGBA {
	parts := strings.Split(s, ",")
	if len(parts) < 4 {
		return RGBA{}
	}
	return RGBA{
		R: helper.AtoU8(parts[0]),
		G: helper.AtoU8(parts[1]),
		B: helper.AtoU8(parts[2]),
		A: helper.AtoU8(parts[3]),
	}
}

// String returns a string representation of the RGBA
func (e RGBA) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", e.R, e.G, e.B, e.A)
}
