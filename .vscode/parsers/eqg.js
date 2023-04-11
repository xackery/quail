registerFileType((fileExt, filePath, fileData) => {
	// Check for wav extension
	if (fileExt == 'eqg') {
		const headerArray = fileData.getBytesAt(4, 3);
		const header = String.fromCharCode(...headerArray)
		if (header == 'PFS')
			return true;
	}
	return false;
});

registerParser(() => {
	addStandardHeader();
	read(4);
	dirOffset = getNumberValue();
	addRow('dirOffset', dirOffset, 'location of directory chunk');
	read(4);
	addRow('header', getStringValue(), 'header (PFS)');
	read(4);
	addRow('version', getNumberValue(), 'version (131072)');
	setOffset(dirOffset);
	read(4);
	fileCount = getNumberValue();
	addRow('fileCount', fileCount, 'number of files in pfs');
});