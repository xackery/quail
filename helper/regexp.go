package helper

import (
	"fmt"
	"regexp"
)

var (
	regexpCache = make(map[string]*regexpEntry)
)

type regexpEntry struct {
	pattern string
	reg     *regexp.Regexp
}

// QuickParse parses a key using a pattern, removing first element (the matched line)
func RegexpMatch(key string, pattern string, search string) ([]string, error) {
	cache, ok := regexpCache[key]
	if ok {
		if cache.pattern != pattern {
			return []string{}, fmt.Errorf("pattern mismatch: %s != %s", cache.pattern, pattern)
		}

		matches := cache.reg.FindStringSubmatch(search)
		if len(matches) < 2 {
			return []string{}, nil
		}
		return matches[1:], nil
	}
	cache = &regexpEntry{pattern: pattern}
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return []string{}, err
	}
	cache.reg = reg
	regexpCache[key] = cache
	matches := reg.FindStringSubmatch(search)
	if len(matches) < 2 {
		return []string{}, nil
	}

	return matches[1:], nil
}
