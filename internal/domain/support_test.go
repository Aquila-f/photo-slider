package domain

import "testing"

func TestIsImage(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		// supported extensions
		{"jpg", "photo.jpg", true},
		{"jpeg", "photo.jpeg", true},
		{"png", "photo.png", true},
		{"webp", "photo.webp", true},

		// case insensitive
		{"JPG uppercase", "photo.JPG", true},
		{"JPEG uppercase", "photo.JPEG", true},
		{"PNG uppercase", "photo.PNG", true},
		{"WEBP uppercase", "photo.WEBP", true},
		{"mixed case", "photo.JpG", true},

		// unsupported extensions
		{"gif", "photo.gif", false},
		{"mp4", "video.mp4", false},
		{"txt", "readme.txt", false},

		// edge cases
		{"no extension", "photofile", false},
		{"empty string", "", false},
		{"dot only", ".", false},
		{"double extension", "photo.jpg.bak", false},
		{"hidden jpg file", ".jpg", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsImage(tt.filename)
			if got != tt.want {
				t.Errorf("IsImage(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}
