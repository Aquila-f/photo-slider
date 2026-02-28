package main

import (
	_ "embed"
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/Aquila-f/photo-slider/internal/domain"
	"github.com/Aquila-f/photo-slider/internal/handler"
	"github.com/Aquila-f/photo-slider/internal/photo"
	"github.com/Aquila-f/photo-slider/internal/service"
	"github.com/Aquila-f/photo-slider/internal/storage"
	"github.com/Aquila-f/photo-slider/internal/strategy"
)

//go:embed static/index.html
var indexHTML string

func main() {
	dir := flag.String("dir", ".", "photo directory to serve")
	port := flag.String("port", "8080", "server port")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil || !isDir(absDir) {
		log.Fatalf("invalid directory: %s", *dir)
	}

	// TODO: allow users to select photo sources (e.g. local dir, remote) via config file
	src := &domain.Source{
		ID:       "local",
		Provider: storage.NewLocalFSProvider(absDir),
	}
	sources := map[string]*domain.Source{src.ID: src}
	albums := make(map[string]*domain.Album)

	svc := service.NewAlbumService(sources, albums, strategy.NewFolderAlbumStrategy(), nil, 3)
	if err := svc.SyncAlbums(context.Background()); err != nil {
		log.Fatalf("failed to sync albums: %v", err)
	}

	api := handler.NewAlbumAPI(svc, photo.NewImageCompressor(), photo.NewFixedSizeMapCacher(256))
	router := handler.SetupRouter(indexHTML, api)

	log.Printf("Serving photos from: %s", absDir)
	log.Printf("Open http://localhost:%s", *port)
	if err := router.Run(":" + *port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
