package photo

import (
	"bytes"
	"context"

	"github.com/Aquila-f/photo-slider/internal/domain"
	"github.com/rwcarlsen/goexif/exif"
)

type EXIFExtractor struct{}

func NewEXIFExtractor() *EXIFExtractor {
	return &EXIFExtractor{}
}

const maxEXIFBytes = 64 * 1024

func (e *EXIFExtractor) Extract(_ context.Context, data []byte) (*domain.PhotoMeta, error) {
	meta := &domain.PhotoMeta{}

	if len(data) > maxEXIFBytes {
		data = data[:maxEXIFBytes]
	}
	x, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		return meta, nil
	}

	if tm, err := x.DateTime(); err == nil {
		meta.TakenAt = &tm
	}
	if tag, err := x.Get(exif.Model); err == nil {
		if s, err := tag.StringVal(); err == nil {
			meta.Model = s
		}
	}

	return meta, nil
}
