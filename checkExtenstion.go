package kafekoding_api

import "strings"

func checkExtension(ext string) bool {
	if !(strings.Contains(ext, "jpg") || strings.Contains(ext, ".png")) {
		return false
	} else {
		return true
	}
}
