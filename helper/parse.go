package helper

import "strconv"

// ParseInt parses an int from a string, returning fallback if it fails
func ParseInt(s string, fallback int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return val
}

// ParseFloat parses a float from a string, returning fallback if it fails
func ParseFloat(s string, fallback float64) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fallback
	}
	return val
}

// ParseBool parses a bool from a string, returning fallback if it fails
func ParseBool(s string, fallback bool) bool {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return fallback
	}
	return val
}

// ParseUint parses a uint from a string, returning fallback if it fails
func ParseUint(s string, fallback uint) uint {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fallback
	}
	return uint(val)
}

// ParseInt32 parses a int32 from a string, returning fallback if it fails
func ParseInt32(s string, fallback int32) int32 {
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return fallback
	}
	return int32(val)
}

// ParseInt64 parses a int64 from a string, returning fallback if it fails
func ParseInt64(s string, fallback int64) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fallback
	}
	return val
}

// ParseUint32 parses a uint32 from a string, returning fallback if it fails
func ParseUint32(s string, fallback uint32) uint32 {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return fallback
	}
	return uint32(val)
}

// ParseUint64 parses a uint64 from a string, returning fallback if it fails
func ParseUint64(s string, fallback uint64) uint64 {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fallback
	}
	return val
}

// ParseFloat32 parses a float32 from a string, returning fallback if it fails
func ParseFloat32(s string, fallback float32) float32 {
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return fallback
	}
	return float32(val)
}

// ParseFloat64 parses a float64 from a string, returning fallback if it fails
func ParseFloat64(s string, fallback float64) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fallback
	}
	return val
}

// ParseInt8 parses a int8 from a string, returning fallback if it fails
func ParseInt8(s string, fallback int8) int8 {
	val, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return fallback
	}
	return int8(val)
}

// ParseUint8 parses a uint8 from a string, returning fallback if it fails
func ParseUint8(s string, fallback uint8) uint8 {
	val, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return fallback
	}
	return uint8(val)
}
