package helper

func ReadStringHash(hash []byte) string {
	hashKey := []byte{0x95, 0x3A, 0xC5, 0x2A, 0x95, 0x7A, 0x95, 0x6A}
	out := ""
	for i := 0; i < len(hash); i++ {
		out += string(hash[i] ^ hashKey[i%8])
	}
	if len(out) == 0 {
		return out
	}
	if out[len(out)-1:] == "\x00" {
		out = out[:len(out)-1]
	}
	return out
}
