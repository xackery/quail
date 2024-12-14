![quail-logo](https://github.com/user-attachments/assets/5f816cc1-b7c0-412a-b885-e5e607e04f97)
[![quail](quail.png)](https://github.com/xackery/quail/releases/latest) [![GoDoc](https://godoc.org/github.com/xackery/quail?status.svg)](https://godoc.org/github.com/xackery/quail) [![Go Report Card](https://goreportcard.com/badge/github.com/xackery/quail)](https://goreportcard.com/report/github.com/xackery/quail)


Quail is a command line EverQuest archive manager.

You can use it to extract content from EQG, S3D, and other pfs-based files.

You can also use it with [WCEmu](https://docs.eqemu.io/client/wcemu/) to modify content.

For a more user-friendly approach, I suggest [quail-addon](https://github.com/xackery/quail-addon).

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

Pfs represents packaged files, you can think of them as zip compressed archives but has a special format

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

# Running Tests and Setting Up a Development Environment

At the root of the repo, you'll see a file called `.env_default`. Copy it to the file `.env`, and edit the EQ_PATH variable inside to point to your EQ directory. NOTE that this EQ directory will get folders made like _gequip.s3d/, or _test_data/, so you may want to consider pointing to a unused EQ path if you don't like the clutter of folders on top.

# References

- Frozen Crusader Ornament, item id 81466, is it13926

# External resources

[gltf writer for wld fragments in lantern](https://github.com/vermadas/LanternExtractor/blob/vermadas/multi_inject/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs)

[fragment overview ref](https://github.com/cjab/libeq/blob/0aff154702fe122fa726fb7fbb43a079d8f3a138/crates/libeq_wld/docs/README.md)

[eq-sage ref](https://gitlab.com/knervous/eq-sage/-/tree/master/src/lib?ref_type=heads)

polyhedron none
ambientlight none
blitspritedef none
dmrgbtrack none
====
--skipped
compositesprite
===
dmspritedef 0x01 0x02 0x800 0x1000
Flag 6146 found in global_chr.s3d/global_chr.wld/IVM_DMSPRITEDEF fragID 558 Flag 6146 (0x1802) -- 0x02 0x800 0x1000
Flag 6147 found in erudsxing_chr2.s3d/erudsxing_chr2.wld/HULL_DMSPRITEDEF fragID 74 Flag 6147 (0x1803) -- 0x02 0x800 0x1000
===
actordef 0x80 abysmal_obj.s3d/abysmal_obj.wld/ARMSHPSIGN301_ACTORDEF fragID 2046 Flag 128
===
polyhedrondef 0x01
===
dmspritedef2 0x02 0x4000 0x8000 0x10000
dmspritedef2 face 0x1000  https://gist.github.com/xackery/a587a84fd31f663d3a27c24e655ed0d1
Flag 3 found in abysmal_obj.s3d/abysmal_obj.wld/FWD301_DMSPRITEDEF fragID 170 Flag 3 (0x3) -- 0x02
Flag 16387 found in befallen_obj.s3d/befallen_obj.wld/SPIKEHEADHUM_DMSPRITEDEF fragID 382 Flag 16387 (0x4003) -- 0x02 0x4000
Flag 65539 found in abysmal.s3d/abysmal.wld/R63_DMSPRITEDEF fragID 220 Flag 65539 (0x10003) -- 0x02 0x10000
Flag 81923 found in abysmal_obj.s3d/abysmal_obj.wld/ARMSHPSIGN301_DMSPRITEDEF fragID 31 Flag 81923 (0x14003) -- 0x02 0x4000 0x10000
Flag 98307 found in acrylia.s3d/acrylia.wld/R1_DMSPRITEDEF fragID 299 Flag 98307 (0x18003) -- 0x02 0x8000 0x10000
Max flag found: 2 0x02
===
actor
Flag 46 found in abysmal.s3d/abysmal.wld/ fragID 11861 Flag 46 (0x2e) -- 0x02 0x04 0x08 0x20
Flag 558 found in timorous_lit.s3d/timorous_lit.wld/ fragID 13349 Flag 558 (0x22e) -- 0x02 0x04 0x08 0x20 0x200
Flag 814 found in timorous_lit.s3d/timorous_lit.wld/ fragID 13192 Flag 814 (0x32e) -- 0x02 0x04 0x08 0x20 0x100 0x200
Max flag found: 2 0x02
===


# Compiling WebAssembly API
Quail can be compiled to WebAssembly to be used in the browser. To compile to WebAssembly, first install `tinygo` and then run:

`GOOS=js GOARCH=wasm tinygo build --target wasm --gc leaking -o quail.wasm main_wasm.go`