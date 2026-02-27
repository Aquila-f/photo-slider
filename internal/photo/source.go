package photo

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type Source interface {
	List(ctx context.Context) ([]string, error)
	Read(ctx context.Context, token string) ([]byte, error)
}

var imageExts = map[string]struct{}{
	".jpg": {}, ".jpeg": {}, ".png": {}, ".webp": {},
}

type DirSource struct {
	dir string
}

func NewDirSource(dir string) *DirSource {
	return &DirSource{dir: dir}
}

func (s *DirSource) List(_ context.Context) ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if _, ok := imageExts[strings.ToLower(filepath.Ext(e.Name()))]; ok {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func (s *DirSource) Read(_ context.Context, token string) ([]byte, error) {
	full := filepath.Join(s.dir, filepath.Clean("/"+token))
	if !strings.HasPrefix(full, s.dir+string(filepath.Separator)) && full != s.dir {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(full)
}
