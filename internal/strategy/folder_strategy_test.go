package strategy

import (
	"context"
	"testing"

	"github.com/Aquila-f/photo-slider/internal/domain"
)

func newStrategy() *FolderAlbumStrategy {
	return NewFolderAlbumStrategy()
}

func TestFolderAlbumStrategy_BasicAlbum(t *testing.T) {
	snaps := []domain.DirSnapshot{
		{
			Path: "2024/summer",
			Files: []domain.FileInfo{
				{Name: "a.jpg", IsDir: false},
				{Name: "b.png", IsDir: false},
			},
		},
	}

	albums, err := newStrategy().GenerateAlbums(context.Background(), snaps, "src1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums) != 1 {
		t.Fatalf("expected 1 album, got %d", len(albums))
	}
	a := albums[0]
	if a.Name != "2024/summer" {
		t.Errorf("Name = %q, want %q", a.Name, "2024/summer")
	}
	if a.SourceID != "src1" {
		t.Errorf("SourceID = %q, want %q", a.SourceID, "src1")
	}
	if a.UID != "src1/2024/summer" {
		t.Errorf("UID = %q, want %q", a.UID, "src1/2024/summer")
	}
	if a.Dir != "2024/summer" {
		t.Errorf("Dir = %q, want %q", a.Dir, "2024/summer")
	}
	if len(a.Photos) != 2 {
		t.Errorf("Photos count = %d, want 2", len(a.Photos))
	}
}

func TestFolderAlbumStrategy_RootPathBecomesDefault(t *testing.T) {
	snaps := []domain.DirSnapshot{
		{
			Path:  "",
			Files: []domain.FileInfo{{Name: "cover.jpg", IsDir: false}},
		},
	}

	albums, err := newStrategy().GenerateAlbums(context.Background(), snaps, "src1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums) != 1 {
		t.Fatalf("expected 1 album, got %d", len(albums))
	}
	if albums[0].Name != "default" {
		t.Errorf("Name = %q, want %q", albums[0].Name, "default")
	}
}

func TestFolderAlbumStrategy_SkipsNonImageFiles(t *testing.T) {
	snaps := []domain.DirSnapshot{
		{
			Path: "docs",
			Files: []domain.FileInfo{
				{Name: "readme.txt", IsDir: false},
				{Name: "photo.jpg", IsDir: false},
				{Name: "video.mp4", IsDir: false},
			},
		},
	}

	albums, err := newStrategy().GenerateAlbums(context.Background(), snaps, "src1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums[0].Photos) != 1 {
		t.Errorf("Photos count = %d, want 1", len(albums[0].Photos))
	}
	if albums[0].Photos[0].FilePath != "photo.jpg" {
		t.Errorf("FilePath = %q, want %q", albums[0].Photos[0].FilePath, "photo.jpg")
	}
}

func TestFolderAlbumStrategy_SkipsDirectoryEntries(t *testing.T) {
	snaps := []domain.DirSnapshot{
		{
			Path: "gallery",
			Files: []domain.FileInfo{
				{Name: "subdir", IsDir: true},
				{Name: "photo.jpg", IsDir: false},
			},
		},
	}

	albums, err := newStrategy().GenerateAlbums(context.Background(), snaps, "src1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums[0].Photos) != 1 {
		t.Errorf("Photos count = %d, want 1 (subdir should be excluded)", len(albums[0].Photos))
	}
}

func TestFolderAlbumStrategy_SkipsSnapshotWithNoImages(t *testing.T) {
	snaps := []domain.DirSnapshot{
		{
			Path:  "empty",
			Files: []domain.FileInfo{{Name: "note.txt", IsDir: false}},
		},
	}

	albums, err := newStrategy().GenerateAlbums(context.Background(), snaps, "src1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums) != 0 {
		t.Errorf("expected 0 albums, got %d", len(albums))
	}
}

func TestFolderAlbumStrategy_EmptySnapshots(t *testing.T) {
	albums, err := newStrategy().GenerateAlbums(context.Background(), nil, "src1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums) != 0 {
		t.Errorf("expected 0 albums, got %d", len(albums))
	}
}

func TestFolderAlbumStrategy_MultipleSnapshots(t *testing.T) {
	snaps := []domain.DirSnapshot{
		{
			Path:  "a",
			Files: []domain.FileInfo{{Name: "1.jpg", IsDir: false}},
		},
		{
			Path:  "b",
			Files: []domain.FileInfo{{Name: "note.txt", IsDir: false}}, // no images → skipped
		},
		{
			Path:  "c",
			Files: []domain.FileInfo{{Name: "2.png", IsDir: false}},
		},
	}

	albums, err := newStrategy().GenerateAlbums(context.Background(), snaps, "src1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums) != 2 {
		t.Errorf("expected 2 albums, got %d", len(albums))
	}
}

func TestFolderAlbumStrategy_PhotoInfoFields(t *testing.T) {
	snaps := []domain.DirSnapshot{
		{
			Path:  "trip",
			Files: []domain.FileInfo{{Name: "sunset.jpg", IsDir: false}},
		},
	}

	albums, err := newStrategy().GenerateAlbums(context.Background(), snaps, "src1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := albums[0].Photos[0]
	if p.AlbumName != "trip" {
		t.Errorf("PhotoInfo.AlbumName = %q, want %q", p.AlbumName, "trip")
	}
	if p.FilePath != "sunset.jpg" {
		t.Errorf("PhotoInfo.FilePath = %q, want %q", p.FilePath, "sunset.jpg")
	}
}
