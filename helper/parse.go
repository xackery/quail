package helper

import "strconv"

// ParseInt parses an int from a string, returning fallback if it fails
func ParseInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ParseFloat parses a float from a string, returning fallback if it fails
func ParseFloat(s string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ParseBool parses a bool from a string, returning fallback if it fails
func ParseBool(s string) (bool, error) {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return val, nil
}

// ParseUint parses a uint from a string, returning fallback if it fails
func ParseUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// ParseInt16 parses a int16 from a string, returning fallback if it fails
func ParseInt16(s string) (int16, error) {
	val, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(val), nil
}

// ParseInt32 parses a int32 from a string, returning fallback if it fails
func ParseInt32(s string) (int32, error) {
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}

// ParseInt64 parses a int64 from a string, returning fallback if it fails
func ParseInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ParseUint16 parses a uint16 from a string, returning fallback if it fails
func ParseUint16(s string) (uint16, error) {
	val, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(val), nil
}

// ParseUint32 parses a uint32 from a string, returning fallback if it fails
func ParseUint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(val), nil
}

// ParseUint64 parses a uint64 from a string, returning fallback if it fails
func ParseUint64(s string) (uint64, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ParseFloat32 parses a float32 from a string, returning fallback if it fails
func ParseFloat32(s string) (float32, error) {
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(val), nil
}

// ParseFloat64 parses a float64 from a string, returning fallback if it fails
func ParseFloat64(s string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ParseInt8 parses a int8 from a string, returning fallback if it fails
func ParseInt8(s string) (int8, error) {
	val, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(val), nil
}

// ParseUint8 parses a uint8 from a string, returning fallback if it fails
func ParseUint8(s string) (uint8, error) {
	val, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(val), nil
}

func ParseFloat64Slice3(s []string) ([3]float64, error) {
	var arr [3]float64
	for i := 0; i < 3; i++ {
		val, err := ParseFloat(s[i])
		if err != nil {
			return arr, err
		}
		arr[i] = val
	}
	return arr, nil
}

func ParseFloat32Slice3(s []string) ([3]float32, error) {
	var arr [3]float32
	for i := 0; i < 3; i++ {
		val, err := ParseFloat32(s[i])
		if err != nil {
			return arr, err
		}
		arr[i] = val
	}
	return arr, nil
}

func ParseFloat64Slice2(s []string) ([2]float64, error) {
	var arr [2]float64
	for i := 0; i < 2; i++ {
		val, err := ParseFloat(s[i])
		if err != nil {
			return arr, err
		}
		arr[i] = val
	}
	return arr, nil
}

func ParseFloat32Slice2(s []string) ([2]float32, error) {
	var arr [2]float32
	for i := 0; i < 2; i++ {
		val, err := ParseFloat32(s[i])
		if err != nil {
			return arr, err
		}
		arr[i] = val
	}
	return arr, nil
}

func ParseUint16Slice3(s []string) ([3]uint16, error) {
	var arr [3]uint16
	for i := 0; i < 3; i++ {
		val, err := ParseUint16(s[i])
		if err != nil {
			return arr, err
		}
		arr[i] = val
	}
	return arr, nil
}

func ParseUint8Slice4(s []string) ([4]uint8, error) {
	var arr [4]uint8
	for i := 0; i < 4; i++ {
		val, err := ParseUint8(s[i])
		if err != nil {
			return arr, err
		}
		arr[i] = val
	}
	return arr, nil
}

func ParseFloat32Slice6(s []string) ([6]float32, error) {
	var arr [6]float32
	for i := 0; i < 6; i++ {
		val, err := ParseFloat32(s[i])
		if err != nil {
			return arr, err
		}
		arr[i] = val
	}
	return arr, nil
}

func ParseFloat32Slice4(s []string) ([4]float32, error) {
	var arr [4]float32
	for i := 0; i < 4; i++ {
		val, err := ParseFloat32(s[i])
		if err != nil {
			return arr, err
		}
		arr[i] = val
	}
	return arr, nil
}
