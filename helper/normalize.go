package helper

func Normalize(q [4]float32) [4]float32 {
	l := q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3]
	if l == 0 {
		q[0] = 0
		q[1] = 0
		q[2] = 0
		q[3] = 1
	} else {
		l = 1 / l
		q[0] *= l
		q[1] *= l
		q[2] *= l
		q[3] *= l
	}
	return q
}
