package service

import (
	"context"
	"log"
	"path"

	"github.com/Aquila-f/photo-slider/internal/domain"
)

type SourceReader interface {
	GetSource(id string) (*domain.Source, bool)
	AllSources() map[string]*domain.Source
}

type AlbumService struct {
	sourceReader SourceReader
	albums       map[string]*domain.Album
	strategy     domain.AlbumStrategy
	albumMapper  domain.Mapper
	maxDepth     int
}

func NewAlbumService(sourceReader SourceReader, albums map[string]*domain.Album, strategy domain.AlbumStrategy, mapper domain.Mapper, maxDepth int) *AlbumService {
	return &AlbumService{sourceReader: sourceReader, albums: albums, strategy: strategy, albumMapper: mapper, maxDepth: maxDepth}
}

func (s *AlbumService) SyncAlbums(ctx context.Context) error {
	for k := range s.albums {
		delete(s.albums, k)
	}
	for _, src := range s.sourceReader.AllSources() {
		if err := s.RegisterAlbumsForSource(ctx, src); err != nil {
			log.Printf("error syncing source %s: %v", src.ID, err)
		}
	}
	return nil
}

func (s *AlbumService) RegisterAlbumsForSource(ctx context.Context, src *domain.Source) error {
	snaps, err := src.Provider.Walk(ctx, "", s.maxDepth)
	if err != nil {
		return err
	}
	albums, err := s.strategy.GenerateAlbums(ctx, snaps, src.ID)
	if err != nil {
		return err
	}
	for i := range albums {
		s.albums[albums[i].UID] = &albums[i]
	}
	return nil
}

func (s *AlbumService) RemoveAlbumsBySource(sourceID string) {
	for k, album := range s.albums {
		if album.SourceID == sourceID {
			delete(s.albums, k)
		}
	}
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
	src, ok := s.sourceReader.GetSource(album.SourceID)
	if !ok {
		return nil, domain.ErrSourceNotFound
	}

	return src.Provider.ReadFile(ctx, path.Join(album.Dir, photoToken))
}
