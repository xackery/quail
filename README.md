# quail - Ever[Q]uest [U]niversal [A]rchive, [I]mport and [L]oader tool

[![GoDoc](https://godoc.org/github.com/xackery/quail?status.svg)](https://godoc.org/github.com/xackery/quail) [![Go Report Card](https://goreportcard.com/badge/github.com/xackery/quail)](https://goreportcard.com/report/github.com/xackery/quail) [![Platform Tests & Build](https://github.com/xackery/quail/actions/workflows/build_workflow.yml/badge.svg)](https://github.com/xackery/quail/actions/workflows/build_workflow.yml)

[![quail](quail.png)](https://github.com/xackery/quail/releases/latest)

Quail manages EverQuest files. [Find downloads in releases](https://github.com/xackery/quail/releases/latest)

## Goals
- EQG weapon bidirectional support via blender with V2

## TODO

## Usage

Quail has a number of commands that are displayed when the program is ran on it's own with no arguments:
```
Available Commands:
  blender     Export/Import special blender-friendly quail-addon formats
  compress    Create an eqg archive by compressing a directory
  debug       Debug a file
  extract     Extract an eqg or s3d archive
  help        Help about any command
  inspect     Inspect a file
```


File extensions are broken into the following categories:

## Pfs

Pfs represents packaged files

Extension|Notes
---|---|---|---
eqg|EverQuest Game Asset Pfs Archive
s3d|EverQuest Game Asset Pfs Archive (Legacy)
## Model/Mesh

Model meshes represents geometry data, and some times metadata

Extension|Notes
---|---
dat|zone mesh for version 4 zones
mds|npc, object and item mesh information
mod|npc, object and item mesh information
ter|zone mesh for version 1 to 3
wld|megapack of virtually all data via fragments, s3d legacy route

## Model/Metadata

Model metadata tells additional details or variations of a mesh

Extension|Notes
---|---
ani|bone animation data
edd|unsure yet
lay|layered texture data (variations of materials to swap texture)
lit|baked light data, vertex colors
lod|level of detail related information
prt|particle rendering data
pts|particle point data
tog|toggle data, definnition information
zon|lists zone metadata (e.g. terrain mesh name)

tog|zone object location data (and refs vertex lighting data)
eco|ecology objects (flora, rocks), used for randomizing and blending maps, V4 Zones
rfd|radial floria definitions, more grass related sway, displacement, etc data, V4 Zones

## Special Files

Name|Notes
---|---
floraexclusion.dat|flora exclusion areas, V4 zones
grass.eco|
lake.eco|
*.prj|zone project information, 3ds max file (can be ignored) V4 zones
dbg.txt|generated typically when lit data is missing, (can be ignored)

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

## Simple examples
Frozen Crusader Ornament, item id 81466, is it13926

# GLTF Extensions

// TODO:
- [lights](https://github.com/KhronosGroup/glTF/tree/main/extensions/2.0/Khronos/KHR_lights_punctual)

## External resources

[gltf writer for wld fragments in lantern](https://github.com/vermadas/LanternExtractor/blob/vermadas/multi_inject/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs)
[fragment overview ref](https://github.com/cjab/libeq/blob/0aff154702fe122fa726fb7fbb43a079d8f3a138/crates/libeq_wld/docs/README.md)