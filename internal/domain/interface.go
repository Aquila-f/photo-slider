package domain

import (
	"context"
	"time"
)

type Source struct {
	ID       string
	Provider StorageProvider
}

type FileInfo struct {
	Name  string
	Path  string
	IsDir bool
}

type DirSnapshot struct {
	Path  string
	Files []FileInfo
}

type StorageProvider interface {
	ListDir(ctx context.Context, path string) ([]FileInfo, error)
	Walk(ctx context.Context, root string, maxDepth int) ([]DirSnapshot, error)
	ReadFile(ctx context.Context, filePath string) ([]byte, error)
}

type AlbumItem struct {
	Name string
	Key  string
}

type PhotoInfo struct {
	AlbumName string
	FilePath  string
}

type Album struct {
	UID      string
	SourceID string
	Name     string
	Dir      string
	Photos   []PhotoInfo
}

type AlbumStrategy interface {
	GenerateAlbums(ctx context.Context, snaps []DirSnapshot, sourceId string) ([]Album, error)
}

type PhotoListStrategy interface {
	Arrange(ctx context.Context, tokens []string) ([]string, error)
}

type Mapper interface {
	Encode(name string) string
	Decode(hash string) (string, error)
}

// ProviderFactory creates a StorageProvider for a given source ID.
type ProviderFactory func(id string) StorageProvider

// AlbumRegistrar is implemented by AlbumService to allow SourceService
// to trigger album registration/removal without a circular dependency.
type AlbumRegistrar interface {
	RegisterAlbumsForSource(ctx context.Context, src *Source) error
	RemoveAlbumsBySource(sourceID string)
}

type PhotoMeta struct {
	TakenAt *time.Time
	Model   string
}

func (m *PhotoMeta) Headers() map[string]string {
	h := make(map[string]string)
	if m.TakenAt != nil {
		h["X-Photo-Taken-At"] = m.TakenAt.Format(time.RFC3339)
	}
	if m.Model != "" {
		h["X-Photo-Model"] = m.Model
	}
	return h
}

type MetaExtractor interface {
	Extract(ctx context.Context, data []byte) (*PhotoMeta, error)
}
