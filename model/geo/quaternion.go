package geo

// ApplyQuaternion transforms this vector by multiplying it by
// the specified quaternion and then by the quaternion inverse.
// It basically applies the rotation encoded in the quaternion to this vector.
func ApplyQuaternion(v Vector3, q Quad4) Vector3 {
	x := v.X
	y := v.Y
	z := v.Z

	qx := q.X
	qy := q.Y
	qz := q.Z
	qw := q.W

	// calculate quat * vector
	ix := qw*x + qy*z - qz*y
	iy := qw*y + qz*x - qx*z
	iz := qw*z + qx*y - qy*x
	iw := -qx*x - qy*y - qz*z
	// calculate result * inverse quat
	v.X = ix*qw + iw*-qx + iy*-qz - iz*-qy
	v.Y = iy*qw + iw*-qy + iz*-qx - ix*-qz
	v.Z = iz*qw + iw*-qz + ix*-qy - iy*-qx
	return v
}

// Normalize a quaternion
func Normalize(q Quad4) Quad4 {
	out := Quad4{}
	l := q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
	if l == 0 {
		out.X = 0
		out.Y = 0
		out.Z = 0
		out.W = 1
		return out
	}
	l = 1 / l
	out.X = q.X * l
	out.Y = q.Y * l
	out.Z = q.Z * l
	out.W = q.W * l
	return out
}
