# photo-slider

A lightweight, self-contained local photo viewer. Point it at one or more directories and browse photos by album in your browser.

## Features

- Scan multiple local source directories for image files (JPEG, PNG, WebP, GIF)
- Organise photos into albums by folder structure
- Browser-based slideshow UI with album switching
- On-the-fly image compression (max 1920px, JPEG quality 80) with fixed-size in-memory ring-buffer cache (256 entries)
- EXIF metadata extraction (camera model, date taken)
- Keyboard, mouse, and touch/swipe navigation
- Fullscreen mode
- Shuffle and auto-play with configurable interval
- REST API for programmatic access

## Quick Start

```bash
go build -o photo-slider ./cmd
./photo-slider -config config.yaml -port 8080
```

Open http://localhost:8080

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `config.yaml` | Path to config file |
| `-port` | `8080` | Server port |

## Configuration

Create a `config.yaml` (see `config.example.yaml`):

```yaml
sources:
  - /path/to/your/photos
  - /another/photo/directory
```

Each entry must be an existing directory. Relative paths are resolved to absolute paths at startup.

## UI Controls

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `←` / `→` | Previous / next photo |
| `Space` | Toggle play/pause slideshow |
| `f` | Toggle fullscreen |
| `Escape` | Exit fullscreen or stop playback |

### Touch

Swipe left/right to navigate between photos.

### Slideshow

Click play to auto-advance photos. Use the interval slider (1–30 seconds) to control speed. Enable shuffle to randomise photo order.

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/albums` | List all albums |
| GET | `/api/albums/:album` | List photo keys in an album (`?shuffle=true` to randomise) |
| GET | `/photos/:album/:key` | Serve a compressed photo |

Album and photo keys are base64 URL-encoded identifiers. Photo responses include `X-Photo-Taken-At` and `X-Photo-Model` headers when EXIF data is available.

## Architecture

```
cmd/              – entrypoint and static assets (HTML/JS/CSS)
internal/
  config/         – YAML config loader
  domain/         – core types, interfaces, and error definitions
  handler/        – Gin HTTP handlers and router
  mapper/         – base64 key encoder/decoder
  photo/          – image compressor, ring-buffer cache, EXIF extractor
  service/        – album sync and photo retrieval
  storage/        – local filesystem provider
  strategy/       – album generation (folder-based) and photo list strategies
```
