package photo

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"testing"
)

func makeJPEG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("makeJPEG: %v", err)
	}
	return buf.Bytes()
}

func makePNG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("makePNG: %v", err)
	}
	return buf.Bytes()
}

func imageSize(data []byte) (int, int) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}

func TestImageCompressor_InvalidInput(t *testing.T) {
	c := NewImageCompressor()
	input := []byte("not an image")
	out, err := c.Compress(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(out, input) {
		t.Error("expected original data to be returned for invalid input")
	}
}

func TestImageCompressor_UnsupportedFormat(t *testing.T) {
	// Minimal valid GIF89a (1x1 pixel)
	gif1x1 := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x2c, 0x00, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x02,
		0x00, 0x3b,
	}
	c := NewImageCompressor()
	out, err := c.Compress(context.Background(), gif1x1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(out, gif1x1) {
		t.Error("expected original GIF data to be returned unchanged")
	}
}

func TestImageCompressor_SmallJPEG(t *testing.T) {
	input := makeJPEG(t, 100, 100)
	c := NewImageCompressor()
	out, err := c.Compress(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, h := imageSize(out)
	if w == 0 {
		t.Fatal("output is not a valid image")
	}
	if w > maxLongEdge || h > maxLongEdge {
		t.Errorf("small image was upscaled: got %dx%d", w, h)
	}
}

func TestImageCompressor_LargeJPEG(t *testing.T) {
	input := makeJPEG(t, 3000, 2000)
	c := NewImageCompressor()
	out, err := c.Compress(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, h := imageSize(out)
	if w == 0 {
		t.Fatal("output is not a valid image")
	}
	if w > maxLongEdge || h > maxLongEdge {
		t.Errorf("large image was not resized: got %dx%d, want long edge ≤ %d", w, h, maxLongEdge)
	}
}

func TestImageCompressor_LargePNG(t *testing.T) {
	input := makePNG(t, 2500, 3000)
	c := NewImageCompressor()
	out, err := c.Compress(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, h := imageSize(out)
	if w == 0 {
		t.Fatal("output is not a valid image")
	}
	if w > maxLongEdge || h > maxLongEdge {
		t.Errorf("large PNG was not resized: got %dx%d, want long edge ≤ %d", w, h, maxLongEdge)
	}
}
