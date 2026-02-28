package strategy

import (
	"context"

	"github.com/Aquila-f/photo-slider/internal/domain"
)

type FolderAlbumStrategy struct {}

func NewFolderAlbumStrategy() *FolderAlbumStrategy {
	return &FolderAlbumStrategy{}
}

func (s *FolderAlbumStrategy) GenerateAlbum(ctx context.Context, files []domain.FileInfo, name, sourceId string) (domain.Album, error) {
	var Photos []domain.PhotoInfo
	for _, f := range files {
		if !f.IsDir && domain.IsImage(f.Name) {
			Photos = append(Photos, domain.PhotoInfo{AlbumName: name, FilePath: f.Name})
		}
	}
	if len(Photos) == 0 {
		return domain.Album{}, nil
	}

	return domain.Album{
		Name:     name,
		SourceID: sourceId,
		Photos:   Photos,
	}, nil
}