package common

import (
	"regexp"
	"strconv"
)

var (
	numEnding         = regexp.MustCompile("_([0-9]{2})$")
	spriteModelEnding = regexp.MustCompile("_s([0-9]{2})_m([0-9]{2})$")
)

// NumberEnding returns a numeric ending to a string with ending pattern _00
func NumberEnding(in string) int {
	rawEnding := numEnding.FindStringSubmatch(in)
	if len(rawEnding) == 0 {
		return -1
	}
	if len(rawEnding[1]) <= 0 {
		return -1
	}
	spriteIndex, err := strconv.Atoi(rawEnding[1])
	if err != nil {
		return -1
	}
	return spriteIndex
}

func SpriteModelEnding(in string) (int, int) {
	spriteModelEnding := spriteModelEnding.FindAllStringSubmatch(in, -1)
	if len(spriteModelEnding) == 0 {
		return 0, 0
	}
	if len(spriteModelEnding) <= 0 {
		return 0, 0
	}
	record := spriteModelEnding[0]

	if len(record) <= 2 {
		return 0, 0
	}
	spriteIndex, err := strconv.Atoi(record[1])
	if err != nil {
		return 0, 0
	}
	modelIndex, err := strconv.Atoi(record[2])
	if err != nil {
		return 0, 0
	}
	return spriteIndex, modelIndex
}
