registerFileType((fileExt, filePath, fileData) => {
	// Check for wav extension
	if (fileExt == 'mod') {
		const headerArray = fileData.getBytesAt(0, 4);
		const header = String.fromCharCode(...headerArray)
		if (header == 'EQGM')
			return true;
	}
	return false;
});

registerParser(() => {
	addStandardHeader();
	read(4);
	addRow('header', getStringValue(), 'header (EQGM)');
	read(4);
	addRow('version', getNumberValue(), 'version (1, 2, 3)');
	read(4);
	const nameLength = getNumberValue();
	addRow('name', nameLength, 'name length');
	read(4);
	const materialCount = getNumberValue();
	addRow('material count', materialCount, 'material count');
	read(4);
	const verticesCount = getNumberValue();
	addRow('vertices count', verticesCount, 'vertices count');
	read(4);
	const trianglesCount = getNumberValue();
	addRow('triangles count', trianglesCount, 'triangles count');
	read(4);
	const bonesCount = getNumberValue();
	addRow('bones count', bonesCount, 'bones count');
	read(4);
	addRow('bones count', bonesCount, 'bones count');
	addRow('names', undefined, 'list of X names');
	addDetails(() => {
		for (let i = 0; i < nameLength; i++) {
			readUntil(0);
			addRow(i, getStringValue(), getStringValue());
			if (getOffset() >= nameLength) {
				break;
			}
			read(1);
		}
	}, true);
	setOffset(getOffset() + nameLength);
	addRow('materials', undefined, 'list of '+materialCount+' materials');
	addDetails(() => {
		for (let i = 0; i < materialCount; i++) {
			read(4);
			addRow('material id', getNumberValue(), 'material id');
			read(4);
			addRow('name offset', getNumberValue(), 'name offset');
			read(4);
			addRow('shader offset', getNumberValue(), 'shader offset');
			read(4);
			propertyCount = getNumberValue();
			addRow('properties', propertyCount, 'list of '+propertyCount+' properties');
			for (let j = 0; j < propertyCount; j++) {
				read(4);
				addRow('property '+i+' '+j+' name offset', getNumberValue(), 'property name offset');
				read(4);
				propertyType = getNumberValue();
				addRow('property type', propertyType, 'property type');
				if (propertyType == 0) {
					read(4);
					addRow('property value', getNumberValue(), 'property value');
				} else {
					read(4);
					addRow('property value', getNumberValue(), 'property value');
				}
				read(4);
				addRow('property value offset', getNumberValue(), 'property value offset');
			}
		}
	});

});