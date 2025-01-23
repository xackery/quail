package helper

import (
	"fmt"
	"regexp"
	"strconv"
)

// MaterialTagParse checks if the material tag represents a variation material and returns a prefix if applicable.
func MaterialTagParse(isChr bool, tag string) (string, error) {
	// Exit early if isChr is false
	if !isChr {
		return "", nil
	}

	// Check if the tag starts with "CLK" and the next 2 characters are numeric
	if len(tag) >= 5 && tag[:3] == "CLK" {
		if _, err := strconv.Atoi(tag[3:5]); err == nil {
			return "CLK04", nil
		}
	}

	// Regex pattern for other variation materials
	regex := regexp.MustCompile(`^[A-Z]{3}(CH|FA|FT|HE|HN|LG|TA|UA)\d{4}_MDF$`)

	if regex.MatchString(tag) && len(tag) == 13 {
		// Parse the 6th, 7th, 8th, and 9th characters
		sixthSeventh, err1 := strconv.Atoi(tag[5:7])
		eighthNinth, err2 := strconv.Atoi(tag[7:9])

		if err1 != nil || err2 != nil {
			return "", fmt.Errorf("failed to parse numeric parts of tag: %w", err1)
		}

		// Check conditions for being a variation material
		if sixthSeventh > 0 || eighthNinth > 10 {
			return tag[:3], nil
		}
	}

	// If no conditions are met, return an empty string
	return "", nil
}
