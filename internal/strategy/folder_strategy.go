package strategy

import (
	"context"

	"github.com/Aquila-f/photo-slider/internal/domain"
)

type FolderAlbumStrategy struct{}

func NewFolderAlbumStrategy() *FolderAlbumStrategy {
	return &FolderAlbumStrategy{}
}

func (s *FolderAlbumStrategy) GenerateAlbums(ctx context.Context, snaps []domain.DirSnapshot, sourceId string) ([]domain.Album, error) {
	var albums []domain.Album
	for _, snap := range snaps {
		name := snap.Path
		if name == "" {
			name = "default"
		}
		var photos []domain.PhotoInfo
		for _, f := range snap.Files {
			if !f.IsDir && domain.IsImage(f.Name) {
				photos = append(photos, domain.PhotoInfo{AlbumName: name, FilePath: f.Name})
			}
		}
		if len(photos) == 0 {
			continue
		}
		albums = append(albums, domain.Album{
			Name:     name,
			SourceID: sourceId,
			UID:      sourceId + "/" + snap.Path,
			Dir:      snap.Path,
			Photos:   photos,
		})
	}
	return albums, nil
}
