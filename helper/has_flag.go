package helper

func HasFlag(flags uint32, flag uint32) bool {
	return flags&flag != 0
}
