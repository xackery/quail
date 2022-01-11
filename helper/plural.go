package helper

func Pluralize(in int) string {
	if in == 1 {
		return ""
	}
	return "s"
}
