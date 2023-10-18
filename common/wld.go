package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/tag"
)

var (
	fragNames = map[int]string{
		0:  "Default",
		1:  "PaletteFile",
		2:  "UserData",
		3:  "TextureList",
		4:  "Texture",
		5:  "TextureRef",
		6:  "TwoDSpriteDef",
		7:  "TwoDSprite",
		8:  "ThreeDSpriteDef",
		9:  "ThreeDSprite",
		10: "FourDSpriteDef",
		11: "FourDSprite",
		12: "ParticleSpriteDef",
		13: "ParticleSprite",
		14: "CompositeSpriteDef",
		15: "CompositeSprite",
		16: "SkeletonTrackDef",
		17: "SkeletonTrack",
		18: "TrackDef",
		19: "Track",
		20: "Model",
		21: "ObjectLocation",
		22: "Sphere",
		23: "PolyhedronDef",
		24: "Polyhedron",
		25: "SphereListDef",
		26: "SphereList",
		27: "LightDef",
		28: "Light",
		29: "PointLightOld",
		31: "SoundDef",
		32: "Sound",
		33: "WorldTree",
		34: "Region",
		35: "ActiveGeoRegion",
		36: "SkyRegion",
		37: "DirectionalLightOld",
		38: "BlitSpriteDef",
		39: "BlitSprite",
		40: "PointLight",
		41: "Zone",
		42: "AmbientLight",
		43: "DirectionalLight",
		44: "DMSpriteDef",
		45: "DMSprite",
		46: "DMTrackDef",
		47: "DMTrack",
		48: "Material",
		49: "MaterialList",
		50: "DMRGBTrackDef",
		51: "DMRGBTrack",
		52: "ParticleCloudDef",
		53: "First",
		54: "Mesh",
		55: "DMTrackDef2",
	}

	/*
			----         8/24/2023   1:10 AM           3707 z_08_three_d_sprite_def.go
		-a----         4/23/2023   3:17 PM            360 z_08_three_d_sprite_def_test.go
		-a----         8/18/2023   1:42 PM            775 z_09_three_d_sprite.go
		-a----         4/23/2023   3:17 PM            348 z_09_three_d_sprite_test.go
		-a----         8/24/2023   1:10 AM           1667 z_10_four_d_sprite_def.go
		-a----         4/23/2023   3:17 PM            397 z_10_four_d_sprite_def_test.go
		-a----         8/18/2023   1:42 PM            759 z_11_four_d_sprite.go
		-a----         4/23/2023   3:17 PM            360 z_11_four_d_sprite_test.go
		-a----         8/24/2023   1:10 AM           2169 z_12_particle_sprite_def.go
		-a----         4/23/2023   3:17 PM            382 z_12_particle_sprite_def_test.go
		-a----         8/18/2023   1:42 PM            823 z_13_particle_sprite.go
		-a----         4/23/2023   3:17 PM            370 z_13_particle_sprite_test.go
		-a----         8/18/2023   1:42 PM            672 z_14_composite_sprite_def.go
		-a----         4/23/2023   3:17 PM            386 z_14_composite_sprite_def_test.go
		-a----         8/18/2023   1:42 PM            667 z_15_composite_sprite.go
		-a----         4/23/2023   3:17 PM            374 z_15_composite_sprite_test.go
		-a----         8/24/2023   1:10 AM           1979 z_16_skeleton_track_def.go
		-a----         4/23/2023   3:17 PM            397 z_16_skeleton_track_def_test.go
		-a----         8/18/2023   1:42 PM            824 z_17_skeleton_track.go
		-a----         4/23/2023   3:17 PM            415 z_17_skeleton_track_test.go
		-a----         8/24/2023   1:10 AM           1702 z_18_track_def.go
		-a----         4/23/2023   3:17 PM            368 z_18_track_def_test.go
		-a----         8/18/2023   1:42 PM            988 z_19_track.go
		-a----         4/30/2023  12:36 PM            318 z_19_track_test.go
		-a----         8/24/2023   1:10 AM           1974 z_20_model.go
		-a----         4/23/2023   3:17 PM            318 z_20_model_test.go
		-a----         8/24/2023   1:10 AM           1612 z_21_object_location.go
		-a----         4/23/2023   3:17 PM            357 z_21_object_location_test.go
		-a----         8/18/2023   1:42 PM            719 z_22_sphere.go
		-a----         4/23/2023   3:17 PM            323 z_22_sphere_test.go
		-a----         8/24/2023   1:10 AM           1489 z_23_polyhedron_def.go
		-a----         4/23/2023   3:17 PM            350 z_23_polyhedron_def_test.go
		-a----         8/18/2023   1:42 PM            811 z_24_polyhedron.go
		-a----         4/23/2023   3:17 PM            338 z_24_polyhedron_test.go
		-a----         8/24/2023   1:10 AM           1136 z_25_sphere_list_def.go
		-a----         8/18/2023   1:42 PM            393 z_25_sphere_list_def_test.go
		-a----         8/18/2023   1:42 PM            805 z_26_sphere_list.go
		-a----         4/23/2023   3:17 PM            355 z_26_sphere_list_test.go
		-a----         8/24/2023   1:10 AM           1333 z_27_light_def.go
		-a----         4/23/2023   3:17 PM            332 z_27_light_def_test.go
		-a----         8/18/2023   1:42 PM            748 z_28_light.go
		-a----         4/23/2023   3:17 PM            320 z_28_light_test.go
		-a----         8/18/2023   1:42 PM            745 z_29_point_light_old.go
		-a----         4/23/2023   3:17 PM            367 z_29_point_light_old_test.go
		-a----         8/18/2023   1:42 PM            710 z_31_sound_def.go
		-a----         4/23/2023   3:17 PM            347 z_31_sound_def_test.go
		-a----         8/18/2023   1:42 PM            689 z_32_sound.go
		-a----         4/23/2023   3:17 PM            361 z_32_sound_test.go
		-a----         8/24/2023   1:10 AM           1246 z_33_world_tree.go
		-a----         4/23/2023   3:17 PM            335 z_33_world_tree_test.go
		-a----         8/24/2023   1:10 AM           3407 z_34_region.go
		-a----         4/23/2023   3:17 PM            323 z_34_region_test.go
		-a----         8/18/2023   1:42 PM            652 z_35_active_geo_region.go
		-a----         4/23/2023   3:17 PM            401 z_35_active_geo_region_test.go
		-a----         8/18/2023   1:42 PM            616 z_36_sky_region.go
		-a----         4/23/2023   3:17 PM            351 z_36_sky_region_test.go
		-a----         8/18/2023   1:42 PM            676 z_37_directional_light_old.go
		-a----         4/23/2023   3:17 PM            391 z_37_directional_light_old_test.go
		-a----         8/18/2023   1:42 PM            640 z_38_blit_sprite_def.go
		-a----         4/23/2023   3:17 PM            350 z_38_blit_sprite_def_test.go
		-a----         8/18/2023   1:42 PM            622 z_39_blit_sprite.go
		-a----         4/23/2023   3:17 PM            338 z_39_blit_sprite_test.go
		-a----         8/18/2023   1:42 PM            622 z_40_point_light.go
		-a----         4/23/2023   3:17 PM            355 z_40_point_light_test.go
		-a----         8/18/2023   1:42 PM            556 z_41_zone.go
		-a----         8/18/2023   1:42 PM            316 z_41_zone_test.go
		-a----         8/18/2023   1:42 PM            634 z_42_ambient_light.go
		-a----         8/18/2023   1:42 PM            348 z_42_ambient_light_test.go
		-a----         8/18/2023   1:42 PM            658 z_43_directional_light.go
		-a----         4/23/2023   3:17 PM            379 z_43_directional_light_test.go
		-a----         8/24/2023   1:10 AM           4227 z_44_dm_sprite_def.go
		-a----         4/23/2023   3:17 PM            342 z_44_dm_sprite_def_test.go
		-a----         8/18/2023   1:42 PM            610 z_45_dm_sprite.go
		-a----         4/23/2023   3:17 PM            330 z_45_dm_sprite_test.go
		-a----         8/18/2023   1:42 PM            622 z_46_dm_track_def.go
		-a----         8/18/2023   1:42 PM            604 z_47_dm_track.go
		-a----         4/23/2023   3:17 PM            343 z_47_dm_track_test.go
		-a----         8/18/2023   1:42 PM           3177 z_48_material.go
		-a----         4/23/2023   3:17 PM            330 z_48_material_test.go
		-a----         8/18/2023   1:42 PM            978 z_49_material_list.go
		-a----         4/23/2023   3:17 PM            346 z_49_material_list_test.go
		-a----         8/18/2023   1:42 PM            640 z_50_dm_r_g_b_track_def.go
		-a----         4/23/2023   3:17 PM            367 z_50_dm_r_g_b_track_def_test.go
		-a----         8/18/2023   1:42 PM            622 z_51_dm_r_g_b_track.go
		-a----         4/23/2023   3:17 PM            355 z_51_dm_r_g_b_track_test.go
		-a----         8/18/2023   1:42 PM            658 z_52_particle_cloud_def.go
		-a----         4/23/2023   3:17 PM            362 z_52_particle_cloud_def_test.go
		-a----         8/18/2023   1:42 PM            607 z_53_first.go
		-a----         4/23/2023   3:17 PM            335 z_53_first_test.go
		-a----         8/24/2023   1:10 AM           8345 z_54_mesh.go
		-a----         8/24/2023   1:10 AM            316 z_54_mesh_test.go
		-a----         8/18/2023   1:42 PM            651 z_55_dm_track_def2.go
		-a----         4/23/2023   3:17 PM            359 z_55_dm_track_def2_test.go
		-a----        10/18/2023  10:46 AM           5788 z_test.go
	*/
)

type Wld struct {
	Version       int
	IsOldWorld    bool
	Names         map[int32]string // used temporarily while decoding a wld
	fragments     [][]byte         // used temporarily while decoding a wld
	Materials     []*Material
	FragmentCount uint32
	Models        []*Model
}

// WldOpen prepares a wld file, loading fragments. This is usually then called by Decode
func WldOpen(r io.ReadSeeker) (*Wld, error) {
	wld := &Wld{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	tag.New()

	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		return nil, fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	wld.Version = int(dec.Uint32())

	wld.IsOldWorld = false
	switch wld.Version {
	case 0x00015500:
		wld.IsOldWorld = true
	case 0x1000C800:
		wld.IsOldWorld = false
	default:
		return nil, fmt.Errorf("unknown wld identifier %d", wld.Version)
	}
	wld.FragmentCount = dec.Uint32()
	_ = dec.Uint32() //unk1
	_ = dec.Uint32() //unk2
	hashSize := dec.Uint32()
	_ = dec.Uint32() //unk3
	tag.Add(tag.LastPos(), dec.Pos(), "red", "header")
	hashRaw := dec.Bytes(int(hashSize))
	nameData := helper.ReadStringHash(hashRaw)

	wld.Names = make(map[int32]string)
	chunk := []rune{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			wld.Names[int32(lastOffset)] = string(chunk)
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "green", "namedata")

	if dec.Error() != nil {
		return nil, fmt.Errorf("decode: %w", dec.Error())
	}

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(wld.FragmentCount); fragOffset++ {

		fragSize := dec.Uint32()
		totalFragSize += fragSize

		fragCode := dec.Bytes(4)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("frag position seek %d/%d: %w", fragOffset, wld.FragmentCount, err)
		}
		data := make([]byte, fragSize)
		_, err = r.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read frag %d/%d: %w", fragOffset, wld.FragmentCount, err)
		}

		data = append(fragCode, data...)

		wld.fragments = append(wld.fragments, data)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek end of frag %d/%d: %w", fragOffset, wld.FragmentCount, err)
		}
	}
	return wld, nil
}

// Fragment returns data from a specific fragment, used primarily for tests
func (wld *Wld) Fragment(fragmentIndex int) ([]byte, error) {
	if fragmentIndex < 0 || fragmentIndex >= len(wld.fragments) {
		return nil, fmt.Errorf("fragment %d out of bounds", fragmentIndex)
	}
	return wld.fragments[fragmentIndex], nil
}

func (wld *Wld) Close() {
	wld.fragments = nil
	wld.Names = nil
	wld.Materials = nil
}

// FragName returns the name of a fragment
func FragName(fragCode int) string {
	if fragNames[fragCode] != "" {
		return fragNames[fragCode]
	}
	return fmt.Sprintf("unknownFrag%d", fragCode)
}
