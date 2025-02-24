package rawfrag

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestFragment(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()
	tests := []struct {
		path      string
		file      string
		fragIndex int
		isDump    bool
	}{
		// 0x00 Default
		// 0x01 DefaultPaletteFile
		// 0x02 UserData
		// 0x03 BMInfo
		{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 1},
		// 0x04 SimpleSpriteDef
		{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 2},
		// 0x05 SimpleSprite
		{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 3},
		// 0x06 Sprite2DDef
		//{path: "gequip.s3d", file: "gequip.wld", fragIndex: 1},
		// 0x07 Sprite2D
		// 0x08 Sprite3DDef
		// 0x09 Sprite3D
		// 0x0A Sprite4DDef
		// 0x0B Sprite4D
		// 0x0C ParticleSpriteDef
		// 0x0D ParticleSprite
		// 0x0E CompositeSpriteDef
		// 0x0F CompositeSprite
		// 0x10 HierarchicalSpriteDef
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 324},
		// 0x11 HierarchicalSprite
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1367},
		// 0x12 TrackDef
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1437},
		// 0x13 Track
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1436},
		// 0x14 ActorDef
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1368},
		// 0x15 Actor
		// 0x16 Sphere
		// 0x17 PolyhedronDef
		// 0x18 Polyhedron
		// 0x19 SphereListDef
		// 0x1A SphereList
		// 0x1B LightDef
		// 0x1C Light
		// 0x1D PointLightOld
		// 0x1E PointLightOldDef
		// 0x1F Sound
		// 0x20 SoundDef
		// 0x21 WorldTree
		{path: "crushbone.s3d", file: "crushbone.wld", fragIndex: 2917},
		// 0x22 Region
		{path: "crushbone.s3d", file: "crushbone.wld", fragIndex: 4051},
		// 0x23 ActiveGeoRegion
		// 0x24 SkyRegion
		// 0x25 DirectionalLightOld
		// 0x26 BlitSpriteDef
		{path: "gequip.s3d", file: "gequip.wld", fragIndex: 62},
		// 0x27 BlitSprite
		// 0x28 PointLight
		// 0x29 Zone
		{path: "crushbone.s3d", file: "crushbone.wld", fragIndex: 5720},
		// 0x2A AmbientLight
		// 0x2B DirectionalLight
		// 0x2C DMSpriteDef
		// 0x2D DMSprite
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 321},
		// 0x2E DMTrackDef
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1369},
		// 0x2F DMTrack
		{path: "globalfroglok_chr.s3d", file: "globalfroglok_chr.wld", fragIndex: 77},
		// 0x30 MaterialDef
		{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 4},
		// 0x31 MaterialPalette
		// 0x32 DmRGBTrackDef
		// 0x33 DmRGBTrack
		// 0x34 ParticleCloudDef
		// 0x35 GlobalAmbientLightDef
		// 0x36 DmSpriteDef2
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 110},
		// 0x37 DmTrackDef2
		{path: "qeynos_obj.s3d", file: "qeynos_obj.wld", fragIndex: 1001},
		//{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 557}, // tex coord count misaligned
		//{path: "gequip.s3d", file: "gequip.wld", fragIndex: 0}, // Mesh
		//{path: "gfaydark.s3d", file: "gfaydark.wld", fragIndex: 0}, // Mesh
		//{path: "frozenshadow.s3d", file: "frozenshadow.wld", fragIndex: 0}, // Mesh
		//{path: "crushbone.s3d", file: "crushbone.wld", fragIndex: 2916, isDump: true}, // PASS
		//{path: "crushbone.s3d", file: "crushbone.wld", fragIndex: 0}, // PASS
		//{path: "poknowledge.s3d", file: "poknowledge.wld", fragIndex: 0}, // PASS
		// {path: "gequip4.s3d", file: "gequip4.wld", fragIndex: 0}, // PASS
		// {path: "gequip3.s3d", file: "gequip3.wld", fragIndex: 0}, // PASS
		//{path: "gfaydark_obj.s3d", file: "gfaydark_obj.wld", fragIndex: 0}, // PASS
		//{path: "gequip2.s3d", file: "gequip2.wld", fragIndex: 22280}, // PASS
		//{path: "zel_v2_chr.s3d", file: "zel_v2_chr.wld", fragIndex: 0}, // PASS
		//{path: "wol_v3_chr.s3d", file: "wol_v3_chr.wld", fragIndex: 0}, // PASS
		//{path: "globalhuf_chr.s3d", file: "globalhuf_chr.wld", fragIndex: 0}, // PASS
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

			fragments, err := tmpFragments(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to read fragments: %s", err.Error())
			}

			fragName := "unknown"
			fragCode := 0

			for i := 1; i <= len(fragments); i++ {
				if tt.fragIndex != 0 && i != tt.fragIndex {
					continue
				}
				srcData := fragments[i]
				r := bytes.NewReader(srcData)
				reader := NewFrag(r)
				if reader == nil {
					t.Fatalf("frag %d read: unsupported fragment", i)
				}

				err = reader.Read(r, false)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) read: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}

				if tt.isDump {
					os.WriteFile(fmt.Sprintf("%s/%s.src.hex", dirTest, tt.file), srcData, 0644)
				}

				buf := helper.NewByteSeekerTest()
				_, err = buf.Write(srcData[:4])
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) write: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}

				err = reader.Write(buf, false)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) write: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}

				reader2 := NewFrag(buf)
				if reader2 == nil {
					t.Fatalf("frag %d read: unsupported fragment", i)
				}

				fragName = FragName(reader2.FragCode())
				fragCode = reader2.FragCode()

				err = reader2.Read(buf, false)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) read: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}

				diff := deep.Equal(reader, reader2)
				if diff != nil {
					t.Fatalf("frag %d 0x%x (%s) diff mismatch: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), diff)
				}

				dstData := buf.Bytes()

				if tt.isDump {
					os.WriteFile(fmt.Sprintf("%s/%s.dst.hex", dirTest, tt.file), dstData, 0644)
				}
				/*
					err := helper.ByteCompareTest(srcData, dstData)
					if err != nil {
						t.Fatalf("%s frag %d 0x%x (%s) mismatch: %s", tt.file, i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
					} */
			}
			if tt.fragIndex != 0 {
				fmt.Printf("0x%00X %s %s:%s:%d OK\n", fragCode, fragName, tt.path, tt.file, tt.fragIndex)
				return
			}
			fmt.Printf("Processed %d fragments\n", len(fragments))
		})
	}
}

func tmpFragments(r io.ReadSeeker) (fragments [][]byte, err error) {

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
		return nil, fmt.Errorf("read: %w", dec.Error())
	}
	return fragments, nil
}
