# quail - Ever[Q]uest [U]niversal [A]rchive, [I]mport and [L]oader tool

Quail is about EverQuest files, typically found inside eqg and s3d. The overall goal is to design conversion commands to and from these specialized formats to more common, popular formats.

Extension|CLI|Library
---|---|---
ani|inspect|load
blend|TBD|TBD
eqg|compress,extract,inspect|load,save
lit|TBD|load,save
mod|inspect|add,gtlfImport,objImport,load,save
mds|TODO!|
obj|TBD|mattxtExport,mtlExport,objExport,export,mattxtImport,mtlImport,objImport,import
s3d|TBD|load,save
ter|inspect|add,gltfExport,objExport,gltfImport,objImport,load,save
tog|TBD|save
wld|TBD|load,save
zon|inspect|load,save,import

## External resources

[gltf writer for wld fragments in lantern](https://github.com/vermadas/LanternExtractor/blob/vermadas/multi_inject/LanternExtractor/EQ/Wld/Exporters/GltfWriter.cs)