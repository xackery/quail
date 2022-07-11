# quail - Ever[Q]uest [U]niversal [A]rchive, [I]mport and [L]oader tool

Quail parses EverQuest files found inside on pfs compressed archives (*.eqg and *.s3d files). The overall goal is to design conversion commands to and from these specialized formats to more common, popular formats.

Extension|Notes
---|---
ani|30% - Load functionality prototyped
blend|10% - Needs a lot of work, python dependency to script
eqg|80% - Load/Save working, EQ client fails to support saved data
lit|10% - Load prototyped
mod|60% - Load/Save prototypd, GLTF birectional support prototyped
mds|0% - Not yet implemented
obj|80% - Load/Save working, bugs need to be sorted out (lightwave obj)
s3d|50% - Load/Save prototyped, EQ client fails to support saved data, some fragments unsupported
ter|60% - Load/Save prototyped, GLTF bidirectional support prototyped
tog|10% - Save template prototyped
wld|50% - Load/Save prototyped, needs attention
zon|50% - Load/Save prototyped, needs attention

## External resources

[gltf writer for wld fragments in lantern](https://github.com/vermadas/LanternExtractor/blob/vermadas/multi_inject/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs)