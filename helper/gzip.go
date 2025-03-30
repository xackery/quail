package helper

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

// GzipBase64Encode compresses data with gzip and encodes it as base64
func GzipBase64Encode(data []byte) (string, error) {
	// Create a buffer to write the compressed data to
	var compressed bytes.Buffer

	// Create a gzip writer that writes to the buffer
	gzipWriter := gzip.NewWriter(&compressed)

	// Write the data to the gzip writer
	_, err := gzipWriter.Write(data)
	if err != nil {
		return "", fmt.Errorf("failed to write to gzip writer: %w", err)
	}

	// Close the gzip writer to flush any pending data
	if err := gzipWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Encode the compressed data as base64 and return it as a string
	return base64.StdEncoding.EncodeToString(compressed.Bytes()), nil
}

// GzipBase64Decode decodes base64 data and decompresses it from gzip
func GzipBase64Decode(encodedData string) ([]byte, error) {
	// Decode the base64 string to bytes
	compressed, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create a reader for the compressed data
	gzipReader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Read all the decompressed data
	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	return decompressed, nil
}
