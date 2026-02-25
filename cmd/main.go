package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/Aquila-f/photo-slider/internal/handler"
	"github.com/Aquila-f/photo-slider/internal/photo"
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

	r := photo.NewDirResolver(absDir)
	rd := photo.NewFileReader()
	lh := handler.NewListHandler(r)
	rh := handler.NewReadHandler(r, rd)

	router := setupRouter(lh, rh)

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
