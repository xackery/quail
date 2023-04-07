package dump

import (
	"fmt"

	"github.com/xackery/quail/model/geo"
)

func Str(in interface{}) string {
	out := ""
	switch src := in.(type) {
	case nil:
		return "nil"
	case []uint32:
		for _, v := range src {
			out += fmt.Sprintf("%d,", v)
		}
	case []int32:
		for _, v := range src {
			out += fmt.Sprintf("%d,", v)
		}
	case []uint:
		for _, v := range src {
			out += fmt.Sprintf("%d,", v)
		}
	case []int16:
		for _, v := range src {
			out += fmt.Sprintf("%d,", v)
		}
	case string:
		return src
	case uint32:
		out = fmt.Sprintf("%d", src)
	case int:
		out = fmt.Sprintf("%d", src)
	case int8:
		out = fmt.Sprintf("%d", src)
	case int16:
		out = fmt.Sprintf("%d", src)
	case int32:
		out = fmt.Sprintf("%d", src)
	case float32:
		out = fmt.Sprintf("%0.2f", src)
	case float64:
		out = fmt.Sprintf("%0.2f", src)
	case *geo.RGBA:
		out = fmt.Sprintf("%d,%d,%d,%d", src.R, src.G, src.B, src.A)
	case *geo.Vector3:
		out = fmt.Sprintf("%0.2f,%0.2f,%0.2f", src.X, src.Y, src.Z)
	case *geo.Vector2:
		out = fmt.Sprintf("%0.2f,%0.2f", src.X, src.Y)
	case *geo.Quad4:
		out = fmt.Sprintf("%0.2f,%0.2f,%0.2f,%0.2f", src.X, src.Y, src.Z, src.W)
	default:
		return fmt.Sprintf("%s", src)
	}
	return out
}
