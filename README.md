# quail - Ever[Q]uest [U]niversal [A]rchive, [I]mport and [L]oader tool

[![GoDoc](https://godoc.org/github.com/xackery/quail?status.svg)](https://godoc.org/github.com/xackery/quail) [![Go Report Card](https://goreportcard.com/badge/github.com/xackery/quail)](https://goreportcard.com/report/github.com/xackery/quail) [![Platform Tests & Build](https://github.com/xackery/quail/actions/workflows/build_workflow.yml/badge.svg)](https://github.com/xackery/quail/actions/workflows/build_workflow.yml)

[![quail](quail.png)](https://github.com/xackery/quail/releases/latest)

Quail is a command line EverQuest pfs manager. The two primary ways to use it is with the extract/compress commands to turn a pfs archive into a folder to inspect and manipulate (then reverse), or do import/export to blender via [quail-addon](https://github.com/xackery/quail-addon).

You can [find the latest download of quail in releases](https://github.com/xackery/quail/releases/latest)


## Status

Quail currently in an early preview status. While many functionality goals have even reached, they are buggy and not supporting every use case.

## Usage

Quail has a number of commands that are displayed when the program is ran on it's own with no arguments:
```
Available Commands:
  completion  Generate the autocompletion script for the specified shell
  compress    Compress an eqg/s3d/pfs/pak folder named _file.ext/ to a pfs archive
  convert     Take one file, convert to another
  extract     Extract an pfs (eqg/s3d/pak/pfs) archive to a _file.ext/ folder
  extract-mod ExtractMod an pfs archive to a _file.ext/ folder
  help        Help about any command
  hybrid      Hybrid merge geometry with existing bone data
  inspect     Inspect an EverQuest asset
```

# EverQuest File Overview
## Pfs

Pfs represents packaged files, you can think of them as zip compressed archives but has a special format. They all use the same decoder and encoder.

Extension|Notes
---|---
eqg|EverQuest Game Asset Pfs Archive
s3d|EverQuest Game Asset Pfs Archive (Legacy)
pak|EverQuest Game Asset Pfs Archive (Legacy)
pfs|EverQuest Game Asset Pfs Archive (Legacy)

## Model/Mesh

Model meshes represents geometry data, and some times metadata

Extension|Notes
---|---
dat|zone model for version 4 zones
mds|npc, object and item model information
mod|npc, object and item model information
ter|zone model data for version 1 to 3, similar to mod/mds just without bone data
wld|megapack of model metadata and model information for s3d legacy pfs archives

## Model/Metadata

Model Metadata gives additional information about a mesh

Extension|Notes
---|---
ani|**Ani**mation data, frame by frame based on bone locations
edd|**E**mitter **d**efinition **d**ata, aligns with prt to give details about how particles and emitters work
lay|**Lay**ered texture data (Used for texture swaps/variations in a single model)
lit|**Li**gh**t** baking data, has same count as vertices
lod|**L**evel **o**f **D**etail related information, usually refs to additional meshes to render based on distance
prt|**P**article **R**endering **T**ransformations,
pts|**P**article **T**ransformation **S**tatements,
tog|**Tog**gle data
zon|**Zon**e placement data, gives information about the zone terrain and object placements. Version 3 and below this is a binary file, Version 4+ it is raw text in a 3DsMax format
eco|**Eco**logy metadata, used for randomizing and blending maps in Version 4 Zones
rfd|**R**adial **F**lora **D**ata, grass, rocks,and other placable greenery metadata

## Special Files

Name|Notes
---|---
floraexclusion.dat|Flora exclusion areas, Versin 4 zones use this to create ignores on RFD files
prj|3DS Max **Pr**o**j**ect files, this is used by internal team for opening a pfs mesh, doesn't appear to have any use for EverQuest.
dbg.txt|**D**e**b**u**g** log, shows the last export attempt internally, doesn't appear to have any use for EverQuest.


## Progress Checklist
- Import Support
  - Bugs EQG with V4 and zone files
  - Bugs S3D (not supported currently)
- Export Support
  - UV Fix bad alignment
- Bone/Animations
- Image Sequence support (animated textures)
- MOD support
- S3D/WLD support
- V4 Zone support

## World Fragments

Hex|Code|Name|Description
---|---|---|---
0x00|0|[Default](model/mesh/wld/z_00_default.go)|Default, unknown fragment
0x01|1|[PaletteFile](model/mesh/wld/z_01_palette_file.go)|Default palette file link
0x02|2|[UserData](model/mesh/wld/z_02_user_data.go)|
0x03|3|[TextureList](model/mesh/wld/z_03_texture_list.go)|TextureImages
0x04|4|[Texture](model/mesh/wld/z_04_texture.go)|
0x05|5|[TextureRef](model/mesh/wld/z_05_texture_ref.go)|
0x06|6|[TwoDSpriteDef](model/mesh/wld/z_06_two_d_sprite_def.go)|TwoDimensionalObject
0x07|7|[TwoDSprite](model/mesh/wld/z_07_two_d_sprite.go)|TwoDimensionalObjectReference
0x08|8|[ThreeDSpriteDef](model/mesh/wld/z_08_three_d_sprite_def.go)|Camera
0x09|9|[ThreeDSprite](model/mesh/wld/z_09_three_d_sprite.go)|CameraReference
0x0A|10|[FourDSpriteDef](model/mesh/wld/z_10_four_d_sprite_def.go)|
0x0B|11|[FourDSprite](model/mesh/wld/z_11_four_d_sprite.go)|
0x0C|12|[ParticleSpriteDef](model/mesh/wld/z_12_particle_sprite_def.go)|
0x0D|13|[ParticleSprite](model/mesh/wld/z_13_particle_sprite.go)|
0x0E|14|[CompositeSpriteDef](model/mesh/wld/z_14_composite_sprite_def.go)|
0x0F|15|[CompositeSprite](model/mesh/wld/z_15_composite_sprite.go)|
0x10|16|[skeletonTrackDef](model/mesh/wld/z_16_hierarchial_sprite_def.go)|SkeletonTrackSet
0x11|17|[skeletonTrack](model/mesh/wld/z_17_hierarchial_sprite.go)|SkeletonTrackSetReference
0x12|18|[TrackDef](model/mesh/wld/z_18_track_def.go)|MobSkeletonPieceTrack
0x13|19|[Track](model/mesh/wld/z_19_track.go)|MobSkeletonPieceTrackReference
0x14|20|[ActorDef](model/mesh/wld/z_20_actor_def.go)|Model
0x15|21|[Actor](model/mesh/wld/z_21_actor.go)|ObjectLocation
0x16|22|[Sphere](model/mesh/wld/z_22_sphere.go)|ZoneUnknown
0x17|23|[PolyhedronDef](model/mesh/wld/z_23_polyhedron_def.go)|PolygonAnimation
0x18|24|[Polyhedron](model/mesh/wld/z_24_polyhedron.go)|PolygonAnimationReference
0x19|25|[SphereListDef](model/mesh/wld/z_25_sphere_list_def.go)|
0x1A|26|[SphereList](model/mesh/wld/z_26_sphere_list.go)|
0x1B|27|[LightDef](model/mesh/wld/z_27_light_def.go)|LightSource
0x1C|28|[Light](model/mesh/wld/z_28_light.go)|LightSourceReference
0x1D|29|[PointLightOld](model/mesh/wld/z_29_point_light_old.go)|
0x1F|31|[SoundDef](model/mesh/wld/z_31_sound_def.go)|
0x20|32|[Sound](model/mesh/wld/z_32_sound.go)|
0x21|33|[WorldTree](model/mesh/wld/z_33_world_tree.go)|BspTree
0x22|34|[Region](model/mesh/wld/z_34_region.go)|BspRegion
0x23|35|[ActiveGeoRegion](model/mesh/wld/z_35_active_geo_region.go)|
0x24|36|[SkyRegion](model/mesh/wld/z_36_sky_region.go)|
0x25|37|[DirectionalLightOld](model/mesh/wld/z_37_directional_light_old.go)|
0x26|38|[BlitSpriteDef](model/mesh/wld/z_38_blit_sprite_def.go)|BlitSpriteDefinition
0x27|39|[BlitSprite](model/mesh/wld/z_39_blit_sprite.go)|BlitSpriteDeference
0x28|40|[PointLight](model/mesh/wld/z_40_point_light.go)|LightInfo
0x29|41|[Zone](model/mesh/wld/z_41_zone.go)|RegionFlag
0x2A|42|[AmbientLight](model/mesh/wld/z_42_ambient_light.go)|
0x2B|43|[DirectionalLight](model/mesh/wld/z_43_directional_light.go)|
0x2C|44|[DmSpriteDef](model/mesh/wld/z_44_dm_sprite_def.go)|AlternateMesh
0x2D|45|[DmSprite](model/mesh/wld/z_45_dm_sprite.go)|MeshReference
0x2E|46|[DmTrackDef](model/mesh/wld/z_46_dm_track_def.go)|
0x2F|47|[DmTrack](model/mesh/wld/z_47_dm_track.go)|MeshAnimatedVericesReference
0x30|48|[Material](model/mesh/wld/z_48_material.go)|
0x31|49|[MaterialList](model/mesh/wld/z_49_material_list.go)|
0x32|50|[DmRGBTrackDef](model/mesh/wld/z_50_dm_r_g_b_track_def.go)|VertexColor
0x33|51|[DmRGBTrack](model/mesh/wld/z_51_dm_r_g_b_track.go)|VertexColorReference
0x34|52|[ParticleCloudDef](model/mesh/wld/z_52_particle_cloud_def.go)|
0x35|53|[First](model/mesh/wld/z_53_first.go)|
0x36|54|[Mesh](model/mesh/wld/z_54_mesh.go)|
0x37|55|[DmTrackDef2](model/mesh/wld/z_55_dm_track_def2.go)|MeshAnimatedVertices


# Running Tests and Setting Up a Development Environment

At the root of the repo, you'll see a file called `.env_default`. Copy it to the file `.env`, and edit the EQ_PATH variable inside to point to your EQ directory. NOTE that this EQ directory will get folders made like _gequip.s3d/, or _test_data/, so you may want to consider pointing to a unused EQ path if you don't like the clutter of folders on top.

# References

- Frozen Crusader Ornament, item id 81466, is it13926


# External resources

[gltf writer for wld fragments in lantern](https://github.com/vermadas/LanternExtractor/blob/vermadas/multi_inject/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs)

[fragment overview ref](https://github.com/cjab/libeq/blob/0aff154702fe122fa726fb7fbb43a079d8f3a138/crates/libeq_wld/docs/README.md)