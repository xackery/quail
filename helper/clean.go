package helper

import "github.com/xackery/quail/log"

// Clean will clean a string of invalid characters
func Clean(in string) string {
	out := in
	// ensure all are valid ascii
	for i := 0; i < len(out); i++ {
		log.Debugf("char %d: %d %s", i, out[i], string(out[i]))
		if out[i] == 0 {
			out = out[:i] + out[i+1:]
			continue
		}
		if out[i] > 127 {
			out = out[:i] + out[i+1:]
			continue
		}
	}
	return out
}
