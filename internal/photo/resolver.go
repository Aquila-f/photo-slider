package photo

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Resolver interface {
	List(ctx context.Context) ([]string, error)
	Resolve(ctx context.Context, token string) (string, error)
}

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".webp": true,
}

type DirResolver struct {
	dir string
}

func NewDirResolver(dir string) *DirResolver {
	return &DirResolver{dir: dir}
}

func (r *DirResolver) List(_ context.Context) ([]string, error) {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if imageExts[strings.ToLower(filepath.Ext(e.Name()))] {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func (r *DirResolver) Resolve(_ context.Context, token string) (string, error) {
	full := filepath.Join(r.dir, filepath.Clean("/"+token))
	if !strings.HasPrefix(full, r.dir+string(filepath.Separator)) && full != r.dir {
		return "", os.ErrNotExist
	}
	if _, err := os.Stat(full); err != nil {
		return "", err
	}
	return full, nil
}
