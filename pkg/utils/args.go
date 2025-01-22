package utils

import (
	"fmt"
	"strings"
)

// ReplaceFileWithBase64 replaces the file path in `-f` or `--filename` flags with its Base64-encoded content.
// It returns the modified arguments and any error encountered.
func ReplaceFileWithBase64(args []string, encodeFunc func(string) (string, error)) ([]string, error) {
	var filePath string

	// Iterate through the arguments to find `-f` or `--filename`
	for i, arg := range args {
		if strings.HasPrefix(arg, "-f") || strings.HasPrefix(arg, "--filename") {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") { // Handle `-f filename`
				filePath = args[i+1]

				// Encode the file content to Base64
				encodedFile, err := encodeFunc(filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to encode file: %w", err)
				}

				// Replace the file path with the encoded Base64 value
				args[i+1] = encodedFile
			} else if strings.Contains(arg, "=") { // Handle `-f=filename`
				filePath = strings.SplitN(arg, "=", 2)[1]

				// Encode the file content to Base64
				encodedFile, err := encodeFunc(filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to encode file: %w", err)
				}

				// Replace the `-f=filename` argument with `--filename=[encoded]`
				args[i] = fmt.Sprintf("--filename=%s", encodedFile)
			}
		}
	}

	return args, nil
}

// ReplaceBase64WithFile decodes Base64-encoded content in `-f` or `--file` flags
// and writes the decoded content to a temporary file. It returns the modified
// arguments and any error encountered.
func ReplaceBase64WithFile(args []string, decodeFunc func(string) (string, error)) ([]string, error) {
	var decodedFilePath string

	// Iterate through the arguments to find `-f` or `--file`
	for i, arg := range args {
		if strings.HasPrefix(arg, "-f") || strings.HasPrefix(arg, "--filename") {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") { // Handle `-f [Base64-encoded]`
				encodedContent := args[i+1]
				// Replace the Base64 value with the temporary file path
				decodedFilePath, _ = decodeFunc(encodedContent)
				args[i+1] = decodedFilePath
			} else if strings.Contains(arg, "=") { // Handle `-f=[Base64-encoded]`
				encodedContent := strings.SplitN(arg, "=", 2)[1]
				// Replace the Base64 value with the temporary file path
				decodedFilePath, _ = decodeFunc(encodedContent)
				args[i] = fmt.Sprintf("--filename=%s", decodedFilePath)
			}
		}
	}

	return args, nil
}
