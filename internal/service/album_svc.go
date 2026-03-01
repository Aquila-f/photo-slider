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
	albumMapper  domain.Mapper
	maxDepth     int
}

func NewAlbumService(sources map[string]*domain.Source, albums map[string]*domain.Album, strategy domain.AlbumStrategy, mapper domain.Mapper, maxDepth int) *AlbumService {
	return &AlbumService{albumSources: sources, albums: albums, strategy: strategy, albumMapper: mapper, maxDepth: maxDepth}
}

func (s *AlbumService) SyncAlbums(ctx context.Context) error {
	for _, src := range s.albumSources {
		snaps, err := src.Provider.Walk(ctx, "", s.maxDepth)
		if err != nil {
			continue
		}
		albums, err := s.strategy.GenerateAlbums(ctx, snaps, src.ID)
		if err != nil {
			continue
		}
		for i := range albums {
			s.albums[albums[i].UID] = &albums[i]
		}
	}
	return nil
}

func (s *AlbumService) ListAlbums(_ context.Context) []domain.AlbumItem {
	items := make([]domain.AlbumItem, 0, len(s.albums))
	for path, album := range s.albums {
		items = append(items, domain.AlbumItem{Name: album.Name, Key: s.albumMapper.Encode(path)})
	}
	return items
}

func (s *AlbumService) ListPhoto(ctx context.Context, albumKey string) ([]string, error) {
	albumUID, err := s.albumMapper.Decode(albumKey)
	if err != nil {
		return nil, domain.ErrAlbumNotFound
	}
	album, ok := s.albums[albumUID]
	if !ok {
		return nil, domain.ErrAlbumNotFound
	}

	tokens := make([]string, 0, len(album.Photos))
	for _, p := range album.Photos {
		tokens = append(tokens, p.FilePath)
	}
	return tokens, nil
}

func (s *AlbumService) ReadPhoto(ctx context.Context, albumKey, photoToken string) ([]byte, error) {
	albumUID, err := s.albumMapper.Decode(albumKey)
	if err != nil {
		return nil, domain.ErrAlbumNotFound
	}
	album, ok := s.albums[albumUID]
	if !ok {
		return nil, domain.ErrAlbumNotFound
	}
	src, ok := s.albumSources[album.SourceID]
	if !ok {
		return nil, domain.ErrSourceNotFound
	}

	return src.Provider.ReadFile(ctx, path.Join(album.Dir, photoToken))
}
