package main

import (
	"context"
	"embed"
	"flag"
	"log"

	"github.com/Aquila-f/photo-slider/internal/config"
	"github.com/Aquila-f/photo-slider/internal/domain"
	"github.com/Aquila-f/photo-slider/internal/handler"
	"github.com/Aquila-f/photo-slider/internal/mapper"
	"github.com/Aquila-f/photo-slider/internal/photo"
	"github.com/Aquila-f/photo-slider/internal/service"
	"github.com/Aquila-f/photo-slider/internal/storage"
	"github.com/Aquila-f/photo-slider/internal/strategy"
)

//go:embed static/*
var staticFS embed.FS

func main() {
	// Parse CLI flags: server port and config file path.
	port := flag.String("port", "8080", "server port")
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	// Load application configuration from the specified YAML file.
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Build a Source map from configured directories, each backed by the local filesystem.
	sources := make(map[string]*domain.Source, len(cfg.Sources))
	for _, dir := range cfg.Sources {
		src := &domain.Source{
			ID:       dir,
			Provider: storage.NewLocalFSProvider(dir),
		}
		sources[src.ID] = src
	}
	albums := make(map[string]*domain.Album)

	// Initialize the album service and sync albums from all sources.
	svc := service.NewAlbumService(sources, albums, strategy.NewFolderAlbumStrategy(), mapper.NewBase64Mapper(), 3)
	if err := svc.SyncAlbums(context.Background()); err != nil {
		log.Fatalf("failed to sync albums: %v", err)
	}

	// Wire up the HTTP API and router with image compression and a 256-entry LRU cache.
	api := handler.NewAlbumAPI(svc, photo.NewImageCompressor(), photo.NewFixedSizeMapCacher(256), photo.NewEXIFExtractor())
	router := handler.SetupRouter(staticFS, api)

	// Log registered sources and albums before starting the server.
	log.Printf("Serving %d source(s), %d album(s)", len(cfg.Sources), len(albums))
	for _, a := range albums {
		log.Printf("  album: %s (%d photos)", a.Name, len(a.Photos))
	}
	log.Printf("Open http://localhost:%s", *port)
	if err := router.Run(":" + *port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
