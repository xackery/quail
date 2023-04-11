package helper

import "strconv"

// AtoU16 converts a string to uint16
func AtoU16(s string) uint16 {
	return uint16(AtoI32(s))
}

// AtoU32 converts a string to uint32
func AtoU32(s string) uint32 {
	return uint32(AtoI32(s))
}

// AtoU8 converts a string to uint8
func AtoU8(s string) uint8 {
	return uint8(AtoI32(s))
}

func AtoI16(s string) int16 {
	return int16(AtoI32(s))
}

func AtoI8(s string) int8 {
	return int8(AtoI32(s))
}

// AtoI32 converts a string to int
func AtoI32(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int32(i)
}

// AtoF32 converts a string to float32
func AtoF32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}

// AtoF64 converts a string to float64
func AtoF64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// AtoB converts a string to bool
func AtoB(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}
