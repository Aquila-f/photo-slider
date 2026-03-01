package main

import (
	"context"
	_ "embed"
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

//go:embed static/index.html
var indexHTML string

func main() {
	port := flag.String("port", "8080", "server port")
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	sources := make(map[string]*domain.Source, len(cfg.Sources))
	for _, dir := range cfg.Sources {
		src := &domain.Source{
			ID:       dir,
			Provider: storage.NewLocalFSProvider(dir),
		}
		sources[src.ID] = src
	}
	albums := make(map[string]*domain.Album)

	svc := service.NewAlbumService(sources, albums, strategy.NewFolderAlbumStrategy(), mapper.NewSlashMapper(), 3)
	if err := svc.SyncAlbums(context.Background()); err != nil {
		log.Fatalf("failed to sync albums: %v", err)
	}

	api := handler.NewAlbumAPI(svc, photo.NewImageCompressor(), photo.NewFixedSizeMapCacher(256))
	router := handler.SetupRouter(indexHTML, api)

	log.Printf("Serving %d source(s)", len(cfg.Sources))
	log.Printf("Open http://localhost:%s", *port)
	if err := router.Run(":" + *port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
