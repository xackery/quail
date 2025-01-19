![quail-logo](https://github.com/user-attachments/assets/5f816cc1-b7c0-412a-b885-e5e607e04f97)
[![quail](quail.png)](https://github.com/xackery/quail/releases/latest) [![GoDoc](https://godoc.org/github.com/xackery/quail?status.svg)](https://godoc.org/github.com/xackery/quail) [![Go Report Card](https://goreportcard.com/badge/github.com/xackery/quail)](https://goreportcard.com/report/github.com/xackery/quail)


Quail is a command line Ever**Qu**est asset manager.

You can [find the latest download of quail in releases](https://github.com/xackery/quail/releases/latest)


Features include:

- zip/unzip pfs archives such as eqg, s3d, pak, and pfs
- convert between pfs archives and .quail folder format (also known as [wce (WorldCom Emu)](https://docs.eqemu.io/client/wcemu/))
- inspect binary files like wld, mod, mds
- tree to visualize binary files

## Status

Quail currently in an early preview status. While many functionality goals have even reached, they are buggy and not supporting every use case.


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


# Running Tests and Setting Up a Development Environment

At the root of the repo, you'll see a file called `.env_default`. Copy it to the file `.env`, and edit the EQ_PATH variable inside to point to your EQ directory. NOTE that this EQ directory will get folders made like _gequip.s3d/, or _test_data/, so you may want to consider pointing to a unused EQ path if you don't like the clutter of folders on top.

# External resources

[gltf writer for wld fragments in lantern](https://github.com/vermadas/LanternExtractor/blob/vermadas/multi_inject/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs)

[fragment overview ref](https://github.com/cjab/libeq/blob/0aff154702fe122fa726fb7fbb43a079d8f3a138/crates/libeq_wld/docs/README.md)

[eq-sage ref](https://gitlab.com/knervous/eq-sage/-/tree/master/src/lib?ref_type=heads)
