# quail - Ever[Q]uest [U]niversal [A]rchive, [I]mport and [L]oader tool

Quail parses EverQuest files found inside on pfs compressed archives (*.eqg and *.s3d files). The overall goal is to design conversion commands to and from these specialized formats to more common, popular formats.

Extension|Notes
---|---
ani|animation (EQG), 30% - Load functionality prototyped
blend|Blender source file, 10% - Needs a lot of work, python dependency to script
eqg|pfs acrhive (EQG), 80% - Load/Save working, EQ client fails to support saved data
lay|layered material (EQG),
lit|light data (EQG), 10% - Load prototyped
mds|model data (EQG), 0% - Not yet implemented
mod|model data (EQG), 60% - Load/Save prototypd, GLTF birectional support prototyped
obj|lightform source file, 80% - Load/Save working, bugs need to be sorted out (lightwave obj)
prt|particle rendering (EQG), 0% -
pts|partical transform (EQG), 0% - 
s3d|pfs archive (S3D), 50% - Load/Save prototyped, EQ client fails to support saved data, some fragments unsupported
ter|terrain data (EQG), 60% - Load/Save prototyped, GLTF bidirectional support prototyped
tog|10% - Save template prototyped
wld|terrain/model megapack data (EQG),50% - Load/Save prototyped, needs attention
zon|zone index data (EQG), 50% - Load/Save prototyped, needs attention

# EQG Zone Versions

- 1 e.g.: .zon will define
- 4 e.g. shardslanding: .zon 


## External resources

[gltf writer for wld fragments in lantern](https://github.com/vermadas/LanternExtractor/blob/vermadas/multi_inject/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs)