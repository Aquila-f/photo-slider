package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Aquila-f/photo-slider/internal/domain"
	"github.com/Aquila-f/photo-slider/internal/mapper"
	"github.com/Aquila-f/photo-slider/internal/strategy"
)

// mockProvider implements domain.StorageProvider for testing.
type mockProvider struct {
	walkResult []domain.DirSnapshot
	walkErr    error
	files      map[string][]byte // filePath -> content
}

func (m *mockProvider) ListDir(_ context.Context, _ string) ([]domain.FileInfo, error) {
	return nil, nil
}

func (m *mockProvider) Walk(_ context.Context, _ string, _ int) ([]domain.DirSnapshot, error) {
	return m.walkResult, m.walkErr
}

func (m *mockProvider) ReadFile(_ context.Context, filePath string) ([]byte, error) {
	data, ok := m.files[filePath]
	if !ok {
		return nil, errors.New("file not found: " + filePath)
	}
	return data, nil
}

// newTestService wires real strategy + mapper with a mock provider.
// Returns the service plus the underlying maps so tests can inspect/inject state.
func newTestService(provider domain.StorageProvider, sourceID string) (*AlbumService, map[string]*domain.Source, map[string]*domain.Album) {
	sources := map[string]*domain.Source{
		sourceID: {ID: sourceID, Provider: provider},
	}
	albums := map[string]*domain.Album{}
	svc := NewAlbumService(
		sources,
		albums,
		strategy.NewFolderAlbumStrategy(),
		mapper.NewSlashMapper(),
		3,
	)
	return svc, sources, albums
}

// --- SyncAlbums ---

func TestAlbumService_SyncAlbums_PopulatesAlbums(t *testing.T) {
	provider := &mockProvider{
		walkResult: []domain.DirSnapshot{
			{
				Path:  "2024/summer",
				Files: []domain.FileInfo{{Name: "beach.jpg"}},
			},
		},
	}
	svc, _, _ := newTestService(provider, "src1")

	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items := svc.ListAlbums(context.Background())
	if len(items) != 1 {
		t.Fatalf("expected 1 album, got %d", len(items))
	}
	if items[0].Name != "2024/summer" {
		t.Errorf("Name = %q, want %q", items[0].Name, "2024/summer")
	}
}

func TestAlbumService_SyncAlbums_WalkErrorContinues(t *testing.T) {
	provider := &mockProvider{walkErr: errors.New("disk error")}
	svc, _, _ := newTestService(provider, "src1")

	// must not propagate the walk error
	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
	if items := svc.ListAlbums(context.Background()); len(items) != 0 {
		t.Errorf("expected 0 albums, got %d", len(items))
	}
}

func TestAlbumService_SyncAlbums_MultipleSourcesIndependent(t *testing.T) {
	snap := func(path, file string) domain.DirSnapshot {
		return domain.DirSnapshot{Path: path, Files: []domain.FileInfo{{Name: file}}}
	}

	sources := map[string]*domain.Source{
		"srcA": {ID: "srcA", Provider: &mockProvider{walkResult: []domain.DirSnapshot{snap("a", "1.jpg")}}},
		"srcB": {ID: "srcB", Provider: &mockProvider{walkErr: errors.New("fail")}},
		"srcC": {ID: "srcC", Provider: &mockProvider{walkResult: []domain.DirSnapshot{snap("c", "2.jpg")}}},
	}
	albums := map[string]*domain.Album{}
	svc := NewAlbumService(sources, albums, strategy.NewFolderAlbumStrategy(), mapper.NewSlashMapper(), 3)

	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// srcB failed; only srcA and srcC should produce albums
	if len(albums) != 2 {
		t.Errorf("expected 2 albums, got %d", len(albums))
	}
}

// --- ListAlbums ---

func TestAlbumService_ListAlbums_ReturnsEncodedKey(t *testing.T) {
	provider := &mockProvider{
		walkResult: []domain.DirSnapshot{
			{Path: "a/b", Files: []domain.FileInfo{{Name: "x.jpg"}}},
		},
	}
	svc, _, _ := newTestService(provider, "src1")
	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items := svc.ListAlbums(context.Background())
	if len(items) != 1 {
		t.Fatalf("expected 1 item")
	}
	// UID = "src1/a/b"  →  SlashMapper.Encode  →  "src1_a_b"
	if items[0].Key != "src1_a_b" {
		t.Errorf("Key = %q, want %q", items[0].Key, "src1_a_b")
	}
}

func TestAlbumService_ListAlbums_EmptyWhenNoAlbums(t *testing.T) {
	provider := &mockProvider{}
	svc, _, _ := newTestService(provider, "src1")
	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if items := svc.ListAlbums(context.Background()); len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

// --- ListPhoto ---

func TestAlbumService_ListPhoto_ReturnsTokens(t *testing.T) {
	provider := &mockProvider{
		walkResult: []domain.DirSnapshot{
			{
				Path: "gallery",
				Files: []domain.FileInfo{
					{Name: "photo1.jpg"},
					{Name: "photo2.png"},
				},
			},
		},
	}
	svc, _, _ := newTestService(provider, "src1")
	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// UID = "src1/gallery" → encoded key = "src1_gallery"
	tokens, err := svc.ListPhoto(context.Background(), "src1_gallery")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
}

func TestAlbumService_ListPhoto_AlbumNotFound(t *testing.T) {
	provider := &mockProvider{}
	svc, _, _ := newTestService(provider, "src1")
	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := svc.ListPhoto(context.Background(), "no_such_album")
	if err != domain.ErrAlbumNotFound {
		t.Errorf("expected ErrAlbumNotFound, got: %v", err)
	}
}

// --- ReadPhoto ---

func TestAlbumService_ReadPhoto_ReturnsFileBytes(t *testing.T) {
	provider := &mockProvider{
		walkResult: []domain.DirSnapshot{
			{Path: "trips", Files: []domain.FileInfo{{Name: "sunset.jpg"}}},
		},
		files: map[string][]byte{
			"trips/sunset.jpg": []byte("fake-image-data"),
		},
	}
	svc, _, _ := newTestService(provider, "src1")
	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// key = "src1_trips", token = "sunset.jpg"
	data, err := svc.ReadPhoto(context.Background(), "src1_trips", "sunset.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "fake-image-data" {
		t.Errorf("data = %q, want %q", string(data), "fake-image-data")
	}
}

func TestAlbumService_ReadPhoto_JoinsAlbumDirWithToken(t *testing.T) {
	// Validate that ReadFile receives path.Join(album.Dir, photoToken).
	// If the path is assembled incorrectly, ReadFile returns an error.
	provider := &mockProvider{
		walkResult: []domain.DirSnapshot{
			{Path: "2024/summer", Files: []domain.FileInfo{{Name: "beach.jpg"}}},
		},
		files: map[string][]byte{
			"2024/summer/beach.jpg": []byte("img"),
		},
	}
	svc, _, _ := newTestService(provider, "src1")
	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.ReadPhoto(context.Background(), "src1_2024_summer", "beach.jpg"); err != nil {
		t.Errorf("unexpected error (wrong path?): %v", err)
	}
}

func TestAlbumService_ReadPhoto_AlbumNotFound(t *testing.T) {
	provider := &mockProvider{}
	svc, _, _ := newTestService(provider, "src1")
	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := svc.ReadPhoto(context.Background(), "bad_key", "photo.jpg")
	if err != domain.ErrAlbumNotFound {
		t.Errorf("expected ErrAlbumNotFound, got: %v", err)
	}
}

func TestAlbumService_ReadPhoto_SourceNotFound(t *testing.T) {
	provider := &mockProvider{}
	svc, _, albums := newTestService(provider, "src1")

	// Inject an album that references a source not in albumSources.
	// UID = "ghost/album" → encoded key = "ghost_album"
	albums["ghost/album"] = &domain.Album{
		UID:      "ghost/album",
		SourceID: "missing_src",
		Dir:      "some/dir",
	}

	_, err := svc.ReadPhoto(context.Background(), "ghost_album", "photo.jpg")
	if err != domain.ErrSourceNotFound {
		t.Errorf("expected ErrSourceNotFound, got: %v", err)
	}
}

// --- Idempotency ---

func TestAlbumService_SyncAlbums_Idempotent(t *testing.T) {
	provider := &mockProvider{
		walkResult: []domain.DirSnapshot{
			{Path: "gallery", Files: []domain.FileInfo{{Name: "photo.jpg"}}},
		},
	}
	svc, _, _ := newTestService(provider, "src1")

	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error on first sync: %v", err)
	}
	firstCount := len(svc.ListAlbums(context.Background()))

	if err := svc.SyncAlbums(context.Background()); err != nil {
		t.Fatalf("unexpected error on second sync: %v", err)
	}
	secondCount := len(svc.ListAlbums(context.Background()))

	if firstCount != secondCount {
		t.Errorf("SyncAlbums not idempotent: first=%d, second=%d albums", firstCount, secondCount)
	}
}
