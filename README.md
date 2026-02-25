# photo-slider

A lightweight local photo viewer. Point it at a directory and browse photos in your browser.

## Planned Features

- Scan a local directory for image files
- Browse photos in a slider-style Web UI
- REST API to list and serve images

## Usage

```bash
go build -o photo-slider ./cmd
./photo-slider -dir ./photo -port 8080
```

Open http://localhost:8080

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/images` | List all images |
| GET | `/images/:key` | Serve a single image |
