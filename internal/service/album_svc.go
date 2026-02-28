package service

import (
	"context"
	"path"

	"github.com/Aquila-f/photo-slider/internal/domain"
)

type AlbumService struct {
	albumSources map[string]*domain.Source
	albums       map[string]*domain.Album
	strategy     domain.AlbumStrategy
	mapper       domain.PathMapper
	maxDepth     int
}

func NewAlbumService(sources map[string]*domain.Source, albums map[string]*domain.Album, strategy domain.AlbumStrategy, mapper domain.PathMapper, maxDepth int) *AlbumService {
	return &AlbumService{albumSources: sources, albums: albums, strategy: strategy, mapper: mapper, maxDepth: maxDepth}
}

func (s *AlbumService) helper(ctx context.Context, dirPath string, src *domain.Source, depth int) error {
	files, err := src.Provider.ListDir(ctx, dirPath)
	if err != nil {
		return err
	}

	album, err := s.strategy.GenerateAlbum(ctx, files, dirPath, src.ID)
	if err != nil {
		return err
	}
	if len(album.Photos) > 0 {
		s.albums[album.Name] = &album
	}

	depth++
	if depth > s.maxDepth {
		return nil
	}

	for _, file := range files {
		if file.IsDir {
			var subPath string
			if dirPath == "" {
				subPath = file.Name
			} else {
				subPath = dirPath + "/" + file.Name
			}
			if err := s.helper(ctx, subPath, src, depth); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *AlbumService) SyncAlbums(ctx context.Context) error {
	for _, src := range s.albumSources {
		if err := s.helper(ctx, "", src, 1); err != nil {
			continue
		}
	}
	return nil
}

func (s *AlbumService) ListAlbums(_ context.Context) []string {
	names := make([]string, 0, len(s.albums))
	for name := range s.albums {
		names = append(names, name)
	}
	return names
}

func (s *AlbumService) ListPhoto(ctx context.Context, albumName string) ([]string, error) {
	album, ok := s.albums[albumName]
	if !ok {
		return nil, domain.ErrAlbumNotFound
	}

	tokens := make([]string, 0, len(album.Photos))
	for _, p := range album.Photos {
		tokens = append(tokens, p.FilePath)
	}
	return tokens, nil
}

func (s *AlbumService) ReadPhoto(ctx context.Context, albumName, photoToken string) ([]byte, error) {
	album, ok := s.albums[albumName]
	if !ok {
		return nil, domain.ErrAlbumNotFound
	}
	src, ok := s.albumSources[album.SourceID]
	if !ok {
		return nil, domain.ErrSourceNotFound
	}

	return src.Provider.ReadFile(ctx, path.Join(albumName, photoToken))
}
