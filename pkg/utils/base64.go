package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

// EncodeFileToBase64String reads the content of a file, encodes it in Base64, and returns the encoded string.
func EncodeFileToBase64String(filePath string) (string, error) {
	// Read the file content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Encode the content to Base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

// DecodeBase64StringToFile decodes a Base64 string and writes the output to a specified file.
func DecodeBase64StringToFile(encodedContent string) (string, error) {
	// Decode the Base64 content
	decodedContent, err := base64.StdEncoding.DecodeString(encodedContent)
	if err != nil {
		return "decodedContent", fmt.Errorf("failed to decode Base64 string: %w", err)
	}

	// Write the decoded content to a temporary file
	tmpFile, err := os.CreateTemp("", "decoded-file-temp.yaml")
	if err != nil {
		return "decodedContent", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(decodedContent)
	if err != nil {
		return "decodedContent", fmt.Errorf("failed to write to temporary file: %w", err)
	}

	return tmpFile.Name(), nil
}
