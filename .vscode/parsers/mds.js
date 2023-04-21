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
	addStandardHeader();
	read(4);
	addRow('header', getStringValue(), 'header (EQGS)');
	read(4);
	addRow('version', getNumberValue(), 'version (1, 2, 3)');
	read(4);
	const nameLength = getNumberValue();
	addRow('nameLength', nameLength, 'name chunk length');
	read(4);
	const materialCount = getNumberValue();
	addRow('materialCount', materialCount, 'number of materials');
	read(4);
	const bonesCount = getNumberValue();
	read(4);
	addRow('bonesCount', bonesCount, 'number of bones');
	read(4);
	const subCount = getNumberValue();
	addRow('subCount', subCount, 'number of subs (unknown)');
	read(nameLength);

	addDetails(() => {
		for (let i = 0; i < materialCount; i++) {
			read(4);
			const materialID = getNumberValue();
			addRow('materialID', materialID, 'material ID');
			read(4);
			const nameOffset = getNumberValue();
			addRow('nameOffset', nameOffset, 'name offset');
			read(4);
			const shaderOffset = getNumberValue();
			addRow('shaderOffset', shaderOffset, 'shader offset');
			read(4);
			const propertyCount = getNumberValue();
			addRow('propertyCount', propertyCount, 'number of properties');
			read(propertyCount * 3*4);
		}
	}, true);

	read(boneCount * 14*4);
	read(2 * 4); // filler
	read(4);
	const verticesCount = getNumberValue();
	addRow('verticesCount', verticesCount, 'number of vertices');
	read(4);
	const triangleCount = getNumberValue();
	addRow('triangleCount', triangleCount, 'number of triangles');
	read(4);
	const boneAssignmentCount = getNumberValue();
	addRow('boneAssignmentCount', boneAssignmentCount, 'number of bone assignments');
	read(verticesCount * 8*4);
	if (version > 2) {
		read(2 * 4); //uv
	}
	read(triangleCount * 5*4);
});