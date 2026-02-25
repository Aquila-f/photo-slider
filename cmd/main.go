package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/Aquila-f/photo-slider/internal/handler"
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

	// TODO: replace with real impls
	lh := handler.NewListHandler(nil)
	rh := handler.NewReadHandler(nil, nil)

	r := setupRouter(lh, rh)

	log.Printf("Serving photos from: %s", absDir)
	log.Printf("Open http://localhost:%s", *port)
	r.Run(":" + *port)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
