package wld

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
)

func TestFragment(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		path      string
		file      string
		fragIndex int
	}{
		//{"gequip4.s3d", "gequip4.wld", 0},
		//{"gequip3.s3d", "gequip3.wld", 0},
		//{"gfaydark_obj.s3d", "gfaydark_obj.wld", 0},
		//{"crushbone.s3d", "crushbone.wld", 112}, // threedpsprite
		//{"poknowledge.s3d", "poknowledge.wld", 112}, // threedsprite
		//{"gequip2.s3d", "gequip2.wld", 22280}, // mesh
		//{"gequip.s3d", "gequip.wld", 972}, // mesh
		//{"gfaydark.s3d", "gfaydark.wld", 82}, // mesh
		//{"frozenshadow.s3d", "frozenshadow.wld", 82}, // mesh
		//{"zel_v2_chr.s3d", "zel_v2_chr.wld", 0}, // mesh
		//{"wol_v3_chr.s3d", "wol_v3_chr.wld", 0}, // mesh
		{"globalhuf_chr.s3d", "globalhuf_chr.wld", 0}, // mesh
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.path))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.file, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(tt.file)
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.file, err.Error())
			}

			fragments, err := tmpFragments(t, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read fragments: %s", err.Error())
			}

			for i := 1; i <= len(fragments); i++ {
				data := fragments[i-1]
				r := bytes.NewReader(data)

				dec := encdec.NewDecoder(r, binary.LittleEndian)

				fragCode := dec.Int32()

				decoder, ok := decoders[int(fragCode)]
				if !ok {
					t.Fatalf("frag %d 0x%x decode: unsupported fragment", i, fragCode)
				}

				frag, err := decoder(r)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) decode: %s", i, fragCode, common.FragName(int(fragCode)), err.Error())
				}

				buf := bytes.NewBuffer(nil)
				buf.Write(data[:4])

				err = frag.Encode(buf)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) encode: %s", i, fragCode, common.FragName(int(fragCode)), err.Error())
				}

				//if !reflect.DeepEqual(data, buf.Bytes()) {
				for i := 0; i < len(buf.Bytes()); i++ {
					if len(data) <= i {
						min := 0
						max := len(data)
						fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(data[min:max]))
						max = len(buf.Bytes())
						fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(buf.Bytes()[min:max]))

						t.Fatalf("frag %d 0x%x (%s) src eof at offset %d (dst is too large by %d bytes)", i, fragCode, common.FragName(int(fragCode)), i, len(buf.Bytes())-len(data))
					}
					if len(buf.Bytes()) <= i {
						t.Fatalf("frag %d 0x%x (%s) dst eof at offset %d (dst is too small by %d bytes)", i, fragCode, common.FragName(int(fragCode)), i, len(data)-len(buf.Bytes()))
					}
					if buf.Bytes()[i] == data[i] {
						continue
					}
					fmt.Printf("frag %d 0x%x (%s) mismatch at offset %d (src: 0x%x vs dst: 0x%x aka %d)\n", i, fragCode, common.FragName(int(fragCode)), i, data[i], buf.Bytes()[i], buf.Bytes()[i])
					max := i + 16
					if max > len(data) {
						max = len(data)
					}

					min := i - 16
					if min < 0 {
						min = 0
					}
					fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(data[min:max]))
					if max > len(buf.Bytes()) {
						max = len(buf.Bytes())
					}

					fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(buf.Bytes()[min:max]))
					t.Fatalf("frag %d 0x%x (%s) encode: data mismatch", i, fragCode, common.FragName(int(fragCode)))
				}
			}
		})
	}
}

func tmpFragments(t *testing.T, r io.ReadSeeker) (fragments [][]byte, err error) {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	_ = dec.Bytes(4)
	_ = int(dec.Uint32())

	fragmentCount := dec.Uint32()
	_ = dec.Uint32() //unk1
	_ = dec.Uint32() //unk2
	hashSize := dec.Uint32()
	_ = dec.Uint32() //unk3
	_ = dec.Bytes(int(hashSize))

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(fragmentCount); fragOffset++ {

		fragSize := dec.Uint32()
		totalFragSize += fragSize

		fragCode := dec.Bytes(4)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("frag position seek %d/%d: %w", fragOffset, fragmentCount, err)
		}
		data := make([]byte, fragSize)
		_, err = r.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read frag %d/%d: %w", fragOffset, fragmentCount, err)
		}

		data = append(fragCode, data...)

		fragments = append(fragments, data)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek end of frag %d/%d: %w", fragOffset, fragmentCount, err)
		}
	}

	if dec.Error() != nil {
		return nil, fmt.Errorf("decode: %w", dec.Error())
	}
	return fragments, nil
}
