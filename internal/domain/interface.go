package domain

import "context"

type Source struct {
	ID   string
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
	UID     string
	SourceID string
	Name     string
	Dir      string
	Photos   []PhotoInfo
}

type AlbumStrategy interface {
	GenerateAlbums(ctx context.Context, snaps []DirSnapshot, sourceId string) ([]Album, error)
}

type Mapper interface {
	Encode(name string) string
	Decode(hash string) (string, error)
}