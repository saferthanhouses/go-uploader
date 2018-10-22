package utils

import (
	"mime/multipart"
	"net/http"
)

// Brazenly Stolen from here: https://golangcode.com/get-the-content-type-of-file/
func GetFileContentType(file multipart.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)

	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
