registerFileType((fileExt, filePath, fileData) => {
	if (fileExt == 'mds') {
		const headerArray = fileData.getBytesAt(0, 4);
		const header = String.fromCharCode(...headerArray)
		if (header == 'EQGS')
			return true;
	}
	return false;
});

registerParser(() => {
	read(4);
	addRow('header', getStringValue(), 'header (EQGS)');
	read(4);
	addRow('version', getNumberValue(), 'version (1, 2, 3)');
	read(4);
	const nameLength = getNumberValue();
	addRow('name', nameLength, 'name chunk length');
	read(4);
	const materialCount = getNumberValue();
	addRow('material count', materialCount, 'number of materials');
	read(4);
	const verticesCount = getNumberValue();
	addRow('vertices count', verticesCount, 'nunmber of vertices');
	read(4);
	const trianglesCount = getNumberValue();
	addRow('triangles count', trianglesCount, 'number of triangles');
	read(4);
	const bonesCount = getNumberValue();
	addRow('bones count', bonesCount, 'number of bones');
});