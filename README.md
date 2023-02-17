# quail - Ever[Q]uest [U]niversal [A]rchive, [I]mport and [L]oader tool

[![GoDoc](https://godoc.org/github.com/xackery/quail?status.svg)](https://godoc.org/github.com/xackery/quail) [![Go Report Card](https://goreportcard.com/badge/github.com/xackery/quail)](https://goreportcard.com/report/github.com/xackery/quail) [![Platform Tests & Build](https://github.com/xackery/quail/actions/workflows/build_workflow.yml/badge.svg)](https://github.com/xackery/quail/actions/workflows/build_workflow.yml)

![quail](quail.png)

Quail parses EverQuest files found inside on pfs compressed archives (*.eqg and *.s3d files). The overall goal is to design conversion commands to and from these specialized formats to more common, popular formats.

File extensions are broken into the following categories:
## Pfs

Pfs represents packaged files

Extension|%|Ver|Notes
---|---|---|---
eqg|80%|EQG|Version 1 - Contains most assets of modern EQ files
eqg|80%|EQG|Version 2 - Contains most assets of modern EQ files
eqg|80%|EQG|Version 3 - Contains most assets of modern EQ files
eqg|30%|EQG|Version 4 - Contains most assets of modern EQ files
s3d|40%|S3D|Version 1 - Decode/Encode prototyped, EQ client fails to support Encoded data, some fragments unsupported
s3d|40%|S3D|Version 2 - Decode/Encode prototyped, EQ client fails to support Encoded data, some fragments unsupported

## Model/Mesh

Model meshes represents geometry data, and some times metadata

Extension|%|Ver|Notes
---|---|---|---
ter|40%|EQG|terrain data, Decode/Encode prototyped, GLTF bidirectional support prototyped
mds|20%|EQG|model data, early prototyping
mod|20%|EQG|model data, Decode/Encode prototypd, GLTF birectional support prototyped
wld|20%|S3D|terrain/model megapack data, Decode/Encode prototyped, needs attention

## Model/Metadata

Model metadata tells additional details or variations of a mesh

Extension|%|Ver|Notes
---|---|---|---
tog|10%|EQG|object meta data - Encode template prototyped
zon|50%|EQG|zone metadata - Decode/Encode prototyped, needs attention
pts|30%|EQG|partical meta data, early prototype
ani|10%|EQG|animation nodes and positions, missing features
lit|10%|EQG|light data, decode prototyped
lay|40%|EQG|layered material metadata - prototype in
prt|30%|EQG|particle rendering - early prototype

## Model/Plugin

Model plugins are 3rd party export/import file types

Extension|%|Ver|Notes
---|---|---|---
blend|10%|EQG|Blender 3d modeling, python required still
obj|80%|EQG|lightform model OBJ export, - Decode/Encode working, bugs need to be sorted out (lightwave obj)
gltf|10%|EQG|GLTF 3d modeling - needs a lot of work

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
