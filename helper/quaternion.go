package helper

// ApplyQuaternion transforms this vector by multiplying it by
// the specified quaternion and then by the quaternion inverse.
// It basically applies the rotation encoded in the quaternion to this vector.
func ApplyQuaternion(v [3]float32, q [4]float32) [3]float32 {
	x := v[0]
	y := v[1]
	z := v[2]

	qx := q[0]
	qy := q[1]
	qz := q[2]
	qw := q[3]

	// calculate quat * vector
	ix := qw*x + qy*z - qz*y
	iy := qw*y + qz*x - qx*z
	iz := qw*z + qx*y - qy*x
	iw := -qx*x - qy*y - qz*z
	// calculate result * inverse quat
	v[0] = ix*qw + iw*-qx + iy*-qz - iz*-qy
	v[1] = iy*qw + iw*-qy + iz*-qx - ix*-qz
	v[2] = iz*qw + iw*-qz + ix*-qy - iy*-qx
	return v
}
