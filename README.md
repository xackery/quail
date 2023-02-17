# quail - Ever[Q]uest [U]niversal [A]rchive, [I]mport and [L]oader tool

[![GoDoc](https://godoc.org/github.com/xackery/quail?status.svg)](https://godoc.org/github.com/xackery/quail) [![Go Report Card](https://goreportcard.com/badge/github.com/xackery/quail)](https://goreportcard.com/report/github.com/xackery/quail) [![Total alerts](https://img.shields.io/lgtm/alerts/g/xackery/quail.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/xackery/quail/alerts/)

![quail](quail.png)

Quail parses EverQuest files found inside on pfs compressed archives (*.eqg and *.s3d files). The overall goal is to design conversion commands to and from these specialized formats to more common, popular formats.

File extensions are broken into the following categories:
## Pfs

Pfs represents packaged files
Extension|Notes
---|---
eqg|pfs acrhive (EQG), 95% - Decode/Encode working, EQ client fails to support Encoded data
s3d|pfs archive (S3D), 50% - Decode/Encode prototyped, EQ client fails to support Encoded data, some fragments unsupported

## Model/Mesh

Model meshes represents geometry data, and some times metadata

Extension|Notes
---|---
ter|terrain data (EQG), 60% - Decode/Encode prototyped, GLTF bidirectional support prototyped
mds|model data (EQG), 60% - Not yet implemented
mod|model data (EQG), 60% - Decode/Encode prototypd, GLTF birectional support prototyped
wld|terrain/model megapack data (S3D), 20% - Decode/Encode prototyped, needs attention

## Model/Metadata

Model metadata tells additional details or variations of a mesh

Extension|Notes
---|---
tog|object meta data (EQG), 10% - Encode template prototyped
zon|zone metadata (EQG), 50% - Decode/Encode prototyped, needs attention
pts|partical meta data (EQG), 30% - 
ani|animation nodes and positions (EQG), 0% - Decode functionality prototyped
lit|light data (EQG), 10% - Decode prototyped
lay|layered material metadata (EQG), 40%
prt|particle rendering (EQG), 30% -

## Model/Plugin

Model plugins are 3rd party export/import file types

Extension|Notes
---|---
blend|Blender 3d modeling, 10% - Needs a lot of work, python dependency to script
obj|lightform model OBJ export, 80% - Decode/Encode working, bugs need to be sorted out (lightwave obj)
gltf|GLTF 3d modeling, 10% - needs a lot of work

# EQG Zone Versions

- 1 e.g.: .zon will define
- 4 e.g. shardslanding: .zon 

# Problem Children
// TODO:
- b09.eqg mds
- aam.eqg mds
- ahf.eqg mds
- alg.eqg c_ala_bd_s24_c.dds not found
- alkabormare.eqg obp_td_pine_burnedb.mod td_pine_needles_c.dds not found
- anguish.eqg obj_walltorch447ch modelName obj_smallgate15 not found
- ans.eqg mds
- arelis.eqg obj_village_rubble_med_lod1.mod residence.dds not found
- arena.eqg obp_tower.mod grid_standard.dds

# GLTF Extensions

// TODO:
- [lights](https://github.com/KhronosGroup/glTF/tree/main/extensions/2.0/Khronos/KHR_lights_punctual)

## External resources

[gltf writer for wld fragments in lantern](https://github.com/vermadas/LanternExtractor/blob/vermadas/multi_inject/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs)
