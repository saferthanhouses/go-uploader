package utils

import (
	"strings"
	"learnabout-filemanager/config"
)


func GetUploadUrl(id string, contentSuffix string) string {
	urlComponents := []string {
		config.Endpoint,
		config.SpaceName,
		id + "." + contentSuffix,
	}

	return "http://" + strings.Join(urlComponents, "/")
}

