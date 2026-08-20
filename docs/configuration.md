---
layout: default
title: Configuration - Folio
---

# Configuration

Folio uses a `config.yaml` file. It looks for it in the working directory.

## config.yaml

```yaml
port: 8080

database:
  dsn: "postgres://folio:folio@localhost:5432/folio?sslmode=disable"

jwt_secret: change-me-in-production

library:
  paths:
    - /mnt/usb/manga
    - /mnt/usb/comics
    - /mnt/usb/books
```

## Environment Variables

All config values can be set via environment variables:

| Variable | Config Key | Default |
|----------|-----------|---------|
| `FOLIO_PORT` | `port` | `8080` |
| `FOLIO_DB_DSN` | `database.dsn` | `postgres://folio:folio@localhost:5432/folio?sslmode=disable` |
| `FOLIO_JWT_SECRET` | `jwt_secret` | `change-me-in-production` |

## Database

Folio uses PostgreSQL. You need a running PostgreSQL instance.

### Quick setup with Docker

```bash
docker run -d --name folio-pg \
  -e POSTGRES_USER=folio \
  -e POSTGRES_PASSWORD=folio \
  -e POSTGRES_DB=folio \
  -p 5432:5432 \
  postgres:16-alpine
```

Then start Folio with:

```bash
FOLIO_DB_DSN="postgres://folio:folio@localhost:5432/folio?sslmode=disable" ./folio
```

Tables are created automatically on first run via GORM AutoMigrate.

## Supported Formats

| Format | Extension | Notes |
|--------|-----------|-------|
| Comic Book ZIP | `.cbz`, `.zip` | Pages are images inside a ZIP archive |
| Comic Book RAR | `.cbr`, `.rar` | Pages are images inside a RAR archive |
| EPUB | `.epub` | eBooks with chapters and metadata |
| PDF | planned | Not yet supported |

## File Naming

Folio parses series/volume/chapter information from file paths:

```
/mnt/books/manga/One Piece/Volume 01.cbz
/mnt/books/manga/One Piece/Volume 01/Chapter 001.cbz
```

The folder structure is: `Library Root > Series Name > Volume/Chapter files`

For flat structures (all files in one folder), each file becomes a chapter in a single volume.
