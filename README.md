# Folio

A self-hosted digital library server for manga, comics, and ebooks. Optimized for Raspberry Pi.

Built with Go. No CGo dependencies. Single binary deployment.

## Features

- Serve manga/comics (CBZ, CBR, ZIP, RAR) and books (EPUB)
- OPDS 1.2 support for KOReader, Calibre, and other OPDS readers
- KOReader fetcher.koplugin compatible
- Multi-user with role-based access
- Reading progress tracking
- Collections and reading lists
- Annotations and bookmarks
- Full-text search
- Web-based UI (dark theme)
- Import books via web upload
- Low memory footprint, optimized for ARM

## Install

### Binary

Download the latest release for your platform and run:

```bash
./folio
```

### Build from source

```bash
git clone https://github.com/CephandriusMaxtori/Folio.git
cd Folio
make build
```

Cross-compile for ARM64 (Raspberry Pi):

```bash
make build-arm64
```

## Configuration

Folio uses PostgreSQL for data storage.

```yaml
port: 8080
database:
  dsn: "postgres://folio:folio@localhost:5432/folio?sslmode=disable"
jwt_secret: change-me-in-production
library:
  paths:
    - /books/manga
    - /books/comics
```

Or use environment variables:

```bash
FOLIO_PORT=3000
FOLIO_DB_DSN="postgres://user:pass@localhost:5432/folio?sslmode=disable"
FOLIO_JWT_SECRET=your-secret
```

## KOReader Setup

1. Install the OPDS Catalog plugin in KOReader
2. Go to Settings > OPDS Catalogs
3. Add a new catalog with URL: `http://your-server:8080/opds/{API_KEY}`
4. To get your API key, log in to the web UI and go to Settings > API Keys

## License

MIT
