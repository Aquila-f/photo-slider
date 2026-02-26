package photo

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
)

const (
	maxLongEdge = 1920
	jpegQuality = 80
)

type Compressor interface {
	Compress(ctx context.Context, data []byte) ([]byte, error)
}

type ImageCompressor struct{}

func NewImageCompressor() *ImageCompressor {
	return &ImageCompressor{}
}

func (c *ImageCompressor) Compress(_ context.Context, data []byte) ([]byte, error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return data, nil
	}

	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return data, nil
	}

	img = imaging.Fit(img, maxLongEdge, maxLongEdge, imaging.CatmullRom)

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(jpegQuality))
	case "png":
		err = imaging.Encode(&buf, img, imaging.PNG)
	case "gif":
		err = imaging.Encode(&buf, img, imaging.GIF)
	default:
		return data, nil
	}
	if err != nil {
		return data, nil
	}

	return buf.Bytes(), nil
}
