package helper

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	currentAniCode      string
	currentAniModelCode string
	previousAnimations  = make(map[string]struct{})
)

// Dummy strings used in tag matching
var dummyStrings = []string{
	"10404P0", "2HNSWORD", "BARDING", "BELT", "BODY", "BONE",
	"BOW", "BOX", "DUMMY", "HUMEYE", "MESH", "POINT", "POLYSURF",
	"RIDER", "SHOULDER",
}

// Item patterns for non-character cases
var itemPatterns = []string{
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])IT\d+_TRACK$`,
	`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])_IT\d+_TRACK$`,
	`^([C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])){2}_IT\d+_TRACK$`,
}

// TrackAnimationParse parses the tag and returns the animation and model codes
func TrackAnimationParse(isChr bool, tag string) (animationTag string, modelTag string) {
	fmt.Printf("[TrackAnimationParse] Called with tag: %s, currentAniCode: %s, currentAniModelCode: %s\n", tag, currentAniCode, currentAniModelCode)
	// Check if the tag starts with currentAniCode + currentAniModelCode
	combinedCode := currentAniCode + currentAniModelCode
	if currentAniCode != "" && currentAniModelCode != "" && strings.HasPrefix(tag, combinedCode) {
		fmt.Printf("[TrackAnimationParse] Match found with combinedCode: %s\n", combinedCode)
		return currentAniCode, currentAniModelCode
	}

	// // there's an edge case where e.g. P01POINT01_TRACK is inside this
	// matches, err := RegexpMatch("pointPattern", `^([C,D,L,O,P,S,T]0[1-9]|[1-9][0-9]{2})(POINT)[0-9]{2}_TRACK`, tag)
	// if err != nil {
	// 	fmt.Println("pointPattern failed:", err)
	// }
	// if len(matches) == 2 {
	// 	return handleNewAniModelCode(matches[0], matches[1])
	// }

	// Check against previousAnimations
	for previous := range previousAnimations {
		if len(previous) == 0 {
			continue
		}
		if strings.HasPrefix(tag, previous) {
			parts := strings.Split(previous, ":")
			if len(parts) == 2 {
				fmt.Printf("[TrackAnimationParse] Match found in previousAnimations: %s\n", previous)
				return parts[0], parts[1]
			}
		}
	}

	// Check if the tag starts with the currentAniCode and contains a dummy string
	for _, dummy := range dummyStrings {
		if strings.HasPrefix(tag, currentAniCode) && strings.Contains(tag, dummy) {
			fmt.Printf("[TrackAnimationParse] Dummy match found: %s\n", dummy)
			return currentAniCode, currentAniModelCode
		}
	}

	// Handle special cases when isChr is true
	if isChr {
		fmt.Println("[TrackAnimationParse] Processing as character (isChr = true).")
		if strings.HasPrefix(tag, currentAniCode) {
			if currentAniModelCode == "SED" && len(tag) >= 6 && tag[3:6] == "FDD" {
				return currentAniCode, currentAniModelCode
			}
			if currentAniModelCode == "FMP" && len(tag) >= len(currentAniCode)+2 {
				suffixStartIndex := len(currentAniCode)
				for _, suffix := range []string{"PE", "CH", "NE", "HE", "BI", "FO", "TH", "CA", "BO"} {
					if strings.HasPrefix(tag[suffixStartIndex:], suffix) {
						return currentAniCode, currentAniModelCode
					}
				}
			}
			if currentAniModelCode == "SKE" && len(tag) >= len(currentAniCode)+2 {
				suffixStartIndex := len(currentAniCode)
				for _, suffix := range []string{"BI", "BO", "CA", "CH", "FA", "FI", "FO", "HA", "HE", "L_POINT", "NE", "PE", "R_POINT", "SH", "TH", "TO", "TU"} {
					if strings.HasPrefix(tag[suffixStartIndex:], suffix) {
						return currentAniCode, currentAniModelCode
					}
				}
			}
		}

		// // handle generic MIM_TRACK tags
		// matches, err = RegexpMatch("basePattern", `^([A-Z]{3})_TRACK$`, tag)
		// if err != nil {
		// 	fmt.Println("basePattern failed:", err)
		// }
		// if len(matches) == 1 {
		// 	return "", matches[0]
		// }

		// // handle generic bone MIMPE_TRACK tags
		// matches, err = RegexpMatch("bonePattern", `^([A-Z]{3})[A-Z]{2}_TRACK$`, tag)
		// if err != nil {
		// 	fmt.Println("bonePattern failed:", err)
		// }
		// if len(matches) == 1 {
		// 	return "", matches[0]
		// }

		patterns := []string{
			`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])([A-Z]){3}_TRACK$`,
			`^([C,D,L,O,P,S,T]0[1-9]|[1-9][0-9]{2})_([A-Z]{3})_TRACK$`,
			`^([C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])){2}_[A-Z]{3}_TRACK$`,
			`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A-Z]{3}[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])_[A-Z]{3}_TRACK$`,
			`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A,B,G][A-Z]{3}[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A,B,G]_[A-Z]{3}_TRACK$`,
			`^[C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])[A,B,G][C,D,L,O,P,S,T](0[1-9]|[1-9][0-9])_[A-Z]{3}_TRACK$`,
		}

		// for key, pattern := range patterns {
		// 	matches, err := RegexpMatch(key, pattern, tag)
		// 	if err != nil {
		// 		fmt.Println(key, "failed:", err)
		// 		continue
		// 	}
		// 	if len(matches) == 2 {
		// 		return handleNewAniModelCode(matches[0], matches[1])
		// 	}
		// }

		// Attempt to match root patterns
		for i, pattern := range patterns {
			matched, _ := regexp.MatchString(pattern, tag)
			if matched {
				switch i {
				case 0: // Pattern 1
					return handleNewAniModelCode(tag[:3], tag[3:6])
				case 1: // Pattern 2
					return handleNewAniModelCode(tag[:3], tag[7:10])
				case 2, 3: // Pattern 3 and 4
					return handleNewAniModelCode(tag[:3], tag[3:6])
				case 4: // Pattern 5
					return handleNewAniModelCode(tag[:4], tag[4:7])
				case 5: // Pattern 6
					return handleNewAniModelCode(tag[:4], tag[8:11])
				}
			}
		}

		// Fallback for isChr
		if len(tag) >= 6 {
			newAniCode := tag[:3]
			newModelCode := tag[3:6]
			fmt.Printf("[TrackAnimationParse] New animation code: %s, New model code: %s\n", newAniCode, newModelCode)

			return handleNewAniModelCode(newAniCode, newModelCode)
		}

		// If the tag is too short, return empty values
		return "", ""
	}

	// Special cases for isChr == false
	if strings.HasPrefix(tag, currentAniCode) {
		if currentAniModelCode == "IT157" && len(tag) >= 6 && tag[3:6] == "SNA" {
			return currentAniCode, currentAniModelCode
		}
		if currentAniModelCode == "IT61" && len(tag) >= 6 && tag[3:6] == "WIP" {
			return currentAniCode, currentAniModelCode
		}
	}

	// Handle item patterns if isChr is false
	for _, pattern := range itemPatterns {
		matched, _ := regexp.MatchString(pattern, tag)
		if matched {
			newAniCode := tag[:3]
			modelCodeStart := strings.Index(tag, "IT") + 2
			modelCodeEnd := modelCodeStart
			for modelCodeEnd < len(tag) && tag[modelCodeEnd] >= '0' && tag[modelCodeEnd] <= '9' {
				modelCodeEnd++
			}
			return handleNewAniModelCode(newAniCode, "IT"+tag[modelCodeStart:modelCodeEnd])
		}
	}

	// Default fallback for isChr == false
	if len(tag) >= 6 {
		aniCode := tag[:3]
		modelCode := "IT"
		for i := 3; i < len(tag); i++ {
			if tag[i] >= '0' && tag[i] <= '9' {
				modelCode += string(tag[i])
			} else {
				break
			}
		}
		return handleNewAniModelCode(aniCode, modelCode)
	}

	return "", ""
}

// Helper function to handle new animation and model codes
func handleNewAniModelCode(newAniCode, newModelCode string) (string, string) {
	previousAnimations[currentAniCode+currentAniModelCode] = struct{}{}
	currentAniCode = newAniCode
	currentAniModelCode = newModelCode
	return newAniCode, newModelCode
}
