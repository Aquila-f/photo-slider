# Photo Slider

[中文文档](README_zh.md)

A lightweight, self-contained photo viewer. Point it at your local photo directories and browse slideshows in your browser — no database, no cloud, just a single binary.

<p align="center">
  <img src="https://img.shields.io/github/v/release/Aquila-f/photo-slider" alt="Release">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/license-MIT-blue" alt="License">
</p>

## Features

- Scan multiple directories for photos (JPEG, PNG, WebP, GIF)
- Auto-organize into albums by folder structure
- On-the-fly image compression (max 1920 px, JPEG quality 80) with in-memory LRU cache
- EXIF metadata display (camera model, date taken)
- Keyboard, mouse, and touch/swipe navigation
- Fullscreen mode with overlay info
- Shuffle and auto-play with adjustable interval (1–30 s)
- Manage source directories from the UI
- Single binary with embedded web assets — no external dependencies
- REST API for programmatic access

## Quick Start

### From source

```bash
go build -o photo-slider ./cmd
cp config.example.yaml config.yaml   # edit with your photo paths
./photo-slider
```

### From release

Download a binary from the [Releases](https://github.com/Aquila-f/photo-slider/releases) page, edit `config.example.yaml`, and run:

```bash
./photo-slider -config config.yaml
```

Then open **http://localhost:8080**.

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `config.yaml` | Path to configuration file |
| `-port` | `8080` | HTTP server port |

## Configuration

Create a `config.yaml` (see [config.example.yaml](config.example.yaml)):

```yaml
sources:
  - /path/to/your/photos
  - /another/photo/directory
```

Each entry must be an existing, readable directory. Relative paths are resolved to absolute paths at startup.

You can also add or remove sources at runtime through the web UI — click the **Sources** panel at the top of the page.

## Controls

### Keyboard

| Key | Action |
|-----|--------|
| `←` / `→` | Previous / next photo |
| `Space` | Toggle auto-play |
| `f` | Toggle fullscreen |
| `Esc` | Exit fullscreen or stop playback |

### Touch

Swipe left or right to navigate between photos.

### Slideshow

Click ▶ to auto-advance. Adjust the interval slider (1–30 s) to control speed. Enable **Shuffle** to randomize photo order.

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/sources` | List configured source directories |
| `POST` | `/api/sources` | Add a source directory |
| `DELETE` | `/api/sources` | Remove a source directory |
| `GET` | `/api/albums` | List all albums |
| `GET` | `/api/albums/:key` | List photo keys in an album (`?shuffle=true`) |
| `GET` | `/photos/:album/:key` | Serve a compressed photo |

Album and photo identifiers are Base64 URL-encoded. Photo responses include `X-Photo-Taken-At` (RFC 3339) and `X-Photo-Model` headers when EXIF data is available.

## Architecture

```
cmd/
  main.go             Entry point, dependency wiring
  static/             Embedded web UI (HTML, JS, CSS via Alpine.js)

internal/
  config/             YAML configuration loader
  domain/             Core types, interfaces, error definitions
  handler/            Gin HTTP handlers and router
  mapper/             Base64 key encoder/decoder
  photo/              Image compressor, ring-buffer LRU cache, EXIF extractor
  service/            Business logic (album sync, source management)
  storage/            Local filesystem provider
  strategy/           Album generation and photo list strategies
```

## Building

```bash
# Development
go build -o photo-slider ./cmd

# Run tests
go test ./...

# Cross-platform release (requires GoReleaser)
goreleaser release --snapshot --clean
```

## License

MIT
