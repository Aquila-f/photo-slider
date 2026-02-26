package photo

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDirResolver_List(t *testing.T) {
	dir := t.TempDir()

	files := []string{"b.jpg", "a.PNG", "c.webp", "note.txt", "d.gif", "e.jpeg"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte{}, 0644); err != nil {
			t.Fatal(err)
		}
	}

	r := NewDirResolver(dir)
	got, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// only image files, sorted alphabetically (ext is lowercased for matching)
	want := []string{"a.PNG", "b.jpg", "c.webp", "d.gif", "e.jpeg"}
	if len(got) != len(want) {
		t.Fatalf("List() len = %d, want %d; got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestDirResolver_List_Empty(t *testing.T) {
	r := NewDirResolver(t.TempDir())
	got, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 0 {
		t.Errorf("List() = %v, want empty", got)
	}
}

func TestDirResolver_List_InvalidDir(t *testing.T) {
	r := NewDirResolver("/nonexistent/path")
	_, err := r.List(context.Background())
	if err == nil {
		t.Error("List() expected error for invalid dir, got nil")
	}
}
