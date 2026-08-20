# Folio

A self-hosted digital library server for manga, comics, and ebooks. Optimized for Raspberry Pi.

Built with Go. No CGo dependencies. Single binary deployment.

## Features

- Serve manga/comics (CBZ, CBR, CB7, ZIP, RAR, 7Z) and books (EPUB, PDF)
- OPDS 1.2 support for KOReader, Calibre, and other OPDS readers
- KOReader fetcher.koplugin compatible
- Multi-user with role-based access
- Reading progress tracking
- Collections and reading lists
- Annotations and bookmarks
- Full-text search
- Embedded web UI
- Low memory footprint, optimized for ARM

## Install

### Binary

Download the latest release for your platform and run:

```bash
./folio
```

### Docker

```bash
docker run -d \
  -p 8080:8080 \
  -v /path/to/config:/folio/config \
  -v /path/to/data:/folio/data \
  -v /path/to/books:/books \
  folio:latest
```

### Build from source

```bash
git clone https://github.com/CephandriusMaxtori/Folio.git
cd Folio
make build
```

## Configuration

Folio looks for `config.yaml` in the working directory or `./config/`.

```yaml
port: 8080
database:
  path: ./data/folio.db
jwt_secret: change-me-in-production
library:
  paths:
    - /books/manga
    - /books/comics
    - /books/epub
```

## KOReader Setup

1. Install the OPDS Catalog plugin in KOReader
2. Add a new catalog with URL: `http://your-server:8080/opds/YOUR_API_KEY`
3. In fetcher.koplugin settings, select the Folio catalog for sync

## License

MIT
