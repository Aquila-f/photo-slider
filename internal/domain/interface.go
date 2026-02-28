package domain

import "context"

type Source struct {
	ID   string
	Provider StorageProvider
}

type FileInfo struct {
	Name string
	Path string
	IsDir bool
}

type StorageProvider interface {
	ListDir(ctx context.Context, path string) ([]FileInfo, error)
	ReadFile(ctx context.Context, filePath string) ([]byte, error)
}

type PhotoInfo struct {
	AlbumName string
	FilePath     string
}

type Album struct {
	SourceID string
	Name     string
	Path     string
	Photos   []PhotoInfo
}

type AlbumStrategy interface {
	GenerateAlbum(ctx context.Context, files []FileInfo, name, sourceId string) (Album, error)
}

type PathMapper interface {
	Encode(path string) string
	Decode(hash string) (string, error)
}