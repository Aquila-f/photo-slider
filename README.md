# photo-slider

A lightweight local photo viewer. Point it at one or more directories and browse photos by album in your browser.

## Features

- Scan multiple local source directories for image files
- Organise photos into albums by folder structure
- Browse albums and photos via a slider-style Web UI
- REST API to list albums, list photos, and serve images
- On-the-fly image compression with a fixed-size in-memory cache (ring-buffer, 256 entries)

## Configuration

Create a `config.yaml` (see `config.example.yaml`):

```yaml
sources:
  - /path/to/your/photos
  - /another/photo/directory
```

Each entry must be an existing directory. Relative paths are resolved to absolute paths at startup.

## Usage

```bash
go build -o photo-slider ./cmd
./photo-slider -config config.yaml -port 8080
```

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `config.yaml` | Path to config file |
| `-port` | `8080` | Server port |

Open http://localhost:8080

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/albums` | List all albums |
| GET | `/api/albums/:album` | List photo keys in an album |
| GET | `/photos/:album/:key` | Serve a single photo (compressed) |

Album keys are slash-encoded identifiers derived from the source path and folder name. Photo keys are unique tokens per file.

## Architecture

```
cmd/              – entrypoint, wires dependencies
internal/
  config/         – YAML config loader
  domain/         – core types and interfaces
  handler/        – Gin HTTP handlers and router
  mapper/         – slash-based key encoder/decoder
  photo/          – image compressor and fixed-size ring-buffer cache
  service/        – album sync and photo read logic
  storage/        – local filesystem provider
  strategy/       – folder-based album generation strategy
```
