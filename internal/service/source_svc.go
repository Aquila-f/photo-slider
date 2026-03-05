package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Aquila-f/photo-slider/internal/domain"
)

type SourceService struct {
	sources         map[string]*domain.Source
	providerFactory domain.ProviderFactory
	registrar       domain.AlbumRegistrar
}

func NewSourceService(sources map[string]*domain.Source, factory domain.ProviderFactory) *SourceService {
	return &SourceService{sources: sources, providerFactory: factory}
}

// SetRegistrar sets the AlbumRegistrar used to sync albums when sources change.
// This breaks the circular dependency between SourceService and AlbumService.
func (s *SourceService) SetRegistrar(r domain.AlbumRegistrar) {
	s.registrar = r
}

func (s *SourceService) GetSource(id string) (*domain.Source, bool) {
	src, ok := s.sources[id]
	return src, ok
}

func (s *SourceService) AllSources() map[string]*domain.Source {
	return s.sources
}

func (s *SourceService) ListSources(_ context.Context) ([]string, error) {
	ids := make([]string, 0, len(s.sources))
	for id := range s.sources {
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *SourceService) AddSource(ctx context.Context, id string) error {
	if _, exists := s.sources[id]; exists {
		return nil
	}
	info, err := os.Stat(id)
	if err != nil {
		return fmt.Errorf("path not found: %s", id)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", id)
	}
	src := &domain.Source{
		ID:       id,
		Provider: s.providerFactory(id),
	}
	if s.registrar != nil {
		if err := s.registrar.RegisterAlbumsForSource(ctx, src); err != nil {
			return err
		}
	}
	s.sources[id] = src
	return nil
}

func (s *SourceService) DeleteSource(_ context.Context, id string) error {
	delete(s.sources, id)
	if s.registrar != nil {
		s.registrar.RemoveAlbumsBySource(id)
	}
	return nil
}
