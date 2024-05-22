package helper

func WriteStringHash(hash string) []byte {
	hashKey := []byte{0x95, 0x3A, 0xC5, 0x2A, 0x95, 0x7A, 0x95, 0x6A}
	out := make([]byte, len(hash))
	for i := 0; i < len(hash); i++ {
		out[i] = hash[i] ^ hashKey[i%8]
	}
	return out
}

func ReadStringHash(hash []byte) string {
	hashKey := []byte{0x95, 0x3A, 0xC5, 0x2A, 0x95, 0x7A, 0x95, 0x6A}
	out := make([]byte, len(hash))
	for i := 0; i < len(hash); i++ {
		out[i] = hash[i] ^ hashKey[i%8]
	}
	return string(out)
}
