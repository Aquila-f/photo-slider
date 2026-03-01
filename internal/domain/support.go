package domain

import (
	"path/filepath"
	"strings"
)

var imageExts = map[string]struct{}{
	".jpg": {}, ".jpeg": {}, ".png": {}, ".webp": {},
}

func IsImage(name string) bool {
	_, ok := imageExts[strings.ToLower(filepath.Ext(name))]
	return ok
}
