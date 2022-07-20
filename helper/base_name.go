package helper

import (
	"regexp"
	"strings"
)

var (
	numEndingRegex = regexp.MustCompile("[0-9]+$")
)

// BaseName simplifies a name to a base one
func BaseName(in string) string {
	if strings.Contains(in, ".") {
		in = in[0:strings.Index(in, ".")]
	}

	in = numEndingRegex.ReplaceAllString(in, "")
	in = strings.TrimSuffix(in, "_")
	return in
}
