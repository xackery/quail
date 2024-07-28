package rawfrag

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/tag"
)

func TestFragment(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest()
	tests := []struct {
		path      string
		file      string
		fragIndex int
		isDump    bool
	}{
		// 0x00 Default
		// 0x01 PaletteFile
		// 0x02 UserData
		// OK 0x03 TextureList bminfo
		{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 1},
		// OK 0x04 Texture SimpleSpriteDef
		{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 2},
		// OK 0x05 TextureRef SimpleSprite
		{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 3},
		// 0x06 TwoDSprite

		// 0x07 TwoDSpriteRef
		// 0x08 ThreeDSprite
		// 0x09 ThreeDSpriteRef
		// 0x0A FourDSprite
		// 0x0B FourDSpriteRef
		// 0x0C ParticleSprite
		// 0x0D ParticleSpriteRef
		// 0x0E CompositeSprite
		// 0x0F CompositeSpriteRef
		// 0x10 SkeletonTrack
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 324},
		// 0x11 SkeletonTrackRef
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1367},
		// 0x12 Track
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1437},
		// 0x13 TrackRef
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1436},
		// 0x14 Model
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1368},
		// 0x15 ModelRef
		// 0x16 Sphere
		// 0x17 Polyhedron
		// 0x18 PolyhedronRef
		// 0x19 SphereList
		// 0x1A SphereListRef
		// 0x1B Light
		// 0x1C LightRef
		// 0x1D PointLightOld
		// 0x1E PointLightOldRef
		// 0x1F Sound
		// 0x20 SoundRef
		// 0x21 WorldTree
		// 0x22 Region
		// 0x23 ActiveGeoRegion
		// 0x24 SkyRegion
		// 0x25 DirectionalLightOld
		// 0x26 BlitSprite
		// 0x27 BlitSpriteRef
		// 0x28 PointLight
		// 0x29 Zone
		// 0x2A AmbientLight
		// 0x2B DirectionalLight
		// 0x2C DMSprite
		// 0x2D DMSpriteRef
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 321},
		// 0x2E DMTrack
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 1369},
		// 0x2F DMTrackRef
		{path: "globalfroglok_chr.s3d", file: "globalfroglok_chr.wld", fragIndex: 77},
		// OK 0x30 Material MaterialDef
		{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 4},
		// 0x31 MaterialList
		// 0x32 DMRGBTrack
		// 0x33 DMRGBTrackRef
		// 0x34 ParticleCloud
		// 0x35 First
		// 0x36 Mesh
		{path: "globalelm_chr.s3d", file: "globalelm_chr.wld", fragIndex: 110},
		// 0x37 MeshAnimated
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
	if !common.IsTestExtensive() {
		tests = []struct {
			path      string
			file      string
			fragIndex int
			isDump    bool
		}{
			//{path: "global_chr.s3d", file: "global_chr.wld", fragIndex: 4},
			{path: "load2.s3d", file: "load2.wld", fragIndex: 0},
		}
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

				err = reader.Read(r)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) read: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}

				if tt.isDump {
					os.WriteFile(fmt.Sprintf("%s/%s.src.hex", dirTest, tt.file), srcData, 0644)
					tag.Write(fmt.Sprintf("%s/%s.src.hex.tags", dirTest, tt.file))
				}

				buf := common.NewByteSeekerTest()
				buf.Write(srcData[:4])

				err = reader.Write(buf)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) write: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}

				dstData := buf.Bytes()

				if tt.isDump {
					os.WriteFile(fmt.Sprintf("%s/%s.dst.hex", dirTest, tt.file), dstData, 0644)
					tag.Write(fmt.Sprintf("%s/%s.dst.hex.tags", dirTest, tt.file))
				}

				err := common.ByteCompareTest(srcData, dstData)
				if err != nil {
					t.Fatalf("frag %d 0x%x (%s) mismatch: %s", i, reader.FragCode(), FragName(int(reader.FragCode())), err.Error())
				}
			}
			if tt.fragIndex != 0 {
				log.Debugf("Processed 1 fragment @ %d", tt.fragIndex)
				return
			}
			log.Debugf("Processed %d fragments", len(fragments))
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
