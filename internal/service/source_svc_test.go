package service

import (
	"context"
	"testing"

	"github.com/Aquila-f/photo-slider/internal/domain"
	"github.com/Aquila-f/photo-slider/internal/mapper"
	"github.com/Aquila-f/photo-slider/internal/strategy"
)

func newTestSourceService(providers map[string]*mockProvider) (*SourceService, *AlbumService, map[string]*domain.Album) {
	sources := make(map[string]*domain.Source, len(providers))
	for id, p := range providers {
		sources[id] = &domain.Source{ID: id, Provider: p}
	}
	albums := map[string]*domain.Album{}

	sourceSvc := NewSourceService(sources, func(id string) domain.StorageProvider {
		if p, ok := providers[id]; ok {
			return p
		}
		return &mockProvider{}
	})
	albumSvc := NewAlbumService(sourceSvc, albums, strategy.NewFolderAlbumStrategy(), mapper.NewBase64Mapper(), 3)
	sourceSvc.SetRegistrar(albumSvc)
	return sourceSvc, albumSvc, albums
}

func TestSourceService_ListSources(t *testing.T) {
	sourceSvc, _, _ := newTestSourceService(map[string]*mockProvider{
		"src1": {},
		"src2": {},
	})

	ids, err := sourceSvc.ListSources(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 sources, got %d", len(ids))
	}
}

func TestSourceService_AddSource_RegistersAlbums(t *testing.T) {
	dir := t.TempDir()

	providers := map[string]*mockProvider{}
	sourceSvc, albumSvc, _ := newTestSourceService(providers)

	// Add a new source using a real temp directory
	providers[dir] = &mockProvider{
		walkResult: []domain.DirSnapshot{
			{Path: "photos", Files: []domain.FileInfo{{Name: "img.jpg"}}},
		},
	}

	if err := sourceSvc.AddSource(context.Background(), dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify source was added
	ids, _ := sourceSvc.ListSources(context.Background())
	if len(ids) != 1 {
		t.Errorf("expected 1 source, got %d", len(ids))
	}

	// Verify albums were registered
	items := albumSvc.ListAlbums(context.Background())
	if len(items) != 1 {
		t.Errorf("expected 1 album after AddSource, got %d", len(items))
	}
}

func TestSourceService_AddSource_NonExistentPath(t *testing.T) {
	sourceSvc, _, _ := newTestSourceService(map[string]*mockProvider{})

	err := sourceSvc.AddSource(context.Background(), "/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for non-existent path, got nil")
	}

	// Verify source was NOT added
	ids, _ := sourceSvc.ListSources(context.Background())
	if len(ids) != 0 {
		t.Errorf("expected 0 sources, got %d", len(ids))
	}
}

func TestSourceService_AddSource_DuplicateIsNoop(t *testing.T) {
	sourceSvc, _, _ := newTestSourceService(map[string]*mockProvider{
		"existing": {},
	})

	if err := sourceSvc.AddSource(context.Background(), "existing"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ids, _ := sourceSvc.ListSources(context.Background())
	if len(ids) != 1 {
		t.Errorf("expected 1 source (no duplicate), got %d", len(ids))
	}
}

func TestSourceService_DeleteSource_RemovesAlbums(t *testing.T) {
	sourceSvc, albumSvc, _ := newTestSourceService(map[string]*mockProvider{
		"src1": {
			walkResult: []domain.DirSnapshot{
				{Path: "gallery", Files: []domain.FileInfo{{Name: "a.jpg"}}},
			},
		},
	})

	// Sync albums first
	if err := albumSvc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if items := albumSvc.ListAlbums(context.Background()); len(items) != 1 {
		t.Fatalf("expected 1 album before delete, got %d", len(items))
	}

	// Delete the source
	if err := sourceSvc.DeleteSource(context.Background(), "src1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify source removed
	ids, _ := sourceSvc.ListSources(context.Background())
	if len(ids) != 0 {
		t.Errorf("expected 0 sources after delete, got %d", len(ids))
	}

	// Verify albums also removed
	if items := albumSvc.ListAlbums(context.Background()); len(items) != 0 {
		t.Errorf("expected 0 albums after delete, got %d", len(items))
	}
}
