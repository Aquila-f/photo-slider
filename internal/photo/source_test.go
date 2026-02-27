package photo

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDirSource_List(t *testing.T) {
	dir := t.TempDir()

	files := []string{"b.jpg", "a.PNG", "c.webp", "note.txt", "d.gif", "e.jpeg"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte{}, 0644); err != nil {
			t.Fatal(err)
		}
	}

	s := NewDirSource(dir)
	got, err := s.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := []string{"a.PNG", "b.jpg", "c.webp", "e.jpeg"}
	if len(got) != len(want) {
		t.Fatalf("List() len = %d, want %d; got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestDirSource_List_Empty(t *testing.T) {
	s := NewDirSource(t.TempDir())
	got, err := s.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 0 {
		t.Errorf("List() = %v, want empty", got)
	}
}

func TestDirSource_List_InvalidDir(t *testing.T) {
	s := NewDirSource("/nonexistent/path")
	_, err := s.List(context.Background())
	if err == nil {
		t.Error("List() expected error for invalid dir, got nil")
	}
}

func TestDirSource_Read(t *testing.T) {
	dir := t.TempDir()
	content := []byte("fake image data")
	if err := os.WriteFile(filepath.Join(dir, "photo.jpg"), content, 0644); err != nil {
		t.Fatal(err)
	}

	s := NewDirSource(dir)
	got, err := s.Read(context.Background(), "photo.jpg")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("Read() = %q, want %q", got, content)
	}
}

func TestDirSource_Read_Traversal(t *testing.T) {
	dir := t.TempDir()
	s := NewDirSource(dir)

	_, err := s.Read(context.Background(), "../secret.txt")
	if err == nil {
		t.Error("Read() expected error for path traversal, got nil")
	}
}

func TestDirSource_Read_NotFound(t *testing.T) {
	s := NewDirSource(t.TempDir())
	_, err := s.Read(context.Background(), "missing.jpg")
	if err == nil {
		t.Error("Read() expected error for missing file, got nil")
	}
}
