package helper

import "math"

func EulerToQuaternion(in [3]float32) [4]float32 {
	out := [4]float32{}
	c1 := math.Cos(float64(in[0] / 2))
	c2 := math.Cos(float64(in[1] / 2))
	c3 := math.Cos(float64(in[2] / 2))
	s1 := math.Sin(float64(in[0] / 2))
	s2 := math.Sin(float64(in[1] / 2))
	s3 := math.Sin(float64(in[2] / 2))

	out[0] = float32(s1*c2*c3 - c1*s2*s3)
	out[1] = float32(c1*s2*c3 + s1*c2*s3)
	out[2] = float32(c1*c2*s3 - s1*s2*c3)
	out[3] = float32(c1*c2*c3 + s1*s2*s3)

	return out
}
