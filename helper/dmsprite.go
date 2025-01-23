package helper

import (
	"regexp"
)

func DmSpriteDefTagParse(isChr bool, tag string) (string, error) {
	if !isChr {
		return "", nil // No processing needed if not a character model
	}

	regex1 := regexp.MustCompile(`^[A-Z]{3}0[1-9]_DMSPRITEDEF$`)
	regex2 := regexp.MustCompile(`^[A-Z]{3}HE0[1-9]_DMSPRITEDEF$`)

	if regex1.MatchString(tag) || regex2.MatchString(tag) {
		return tag[:3], nil // Return the 3-letter prefix
	}

	return "", nil
}
