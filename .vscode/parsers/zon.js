registerFileType((fileExt, filePath, fileData) => {
	// Check for wav extension
	if (fileExt == 'zon') {
		const headerArray = fileData.getBytesAt(0, 4);
		const header = String.fromCharCode(...headerArray)
		if (header == 'EQTZ')
			return true;
	}
	return false;
});

registerParser(() => {
	// Parse
	read(4);
	addRow('header', getStringValue(), 'header (EQTZ)');
	read(4);
	addRow('version', getNumberValue(), 'version (1)');
	read(4);
	const nameLength = getNumberValue();
	addRow('name', nameLength, 'name length');
	read(4);
	const modelCount = getNumberValue();
	addRow('model Count', modelCount, 'model count');
	read(4);
	const objectCount = getNumberValue();
	addRow('object Count', objectCount, 'object count');
	read(4);
	const regionCount = getNumberValue();
	addRow('region Count', regionCount, 'region count');
	read(4);
	const lightCount = getNumberValue();
	addRow('light Count', lightCount, 'light count');
});