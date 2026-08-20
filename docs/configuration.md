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
  path: ./data/folio.db

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
| `FOLIO_DB_PATH` | `database.path` | `./data/folio.db` |
| `FOLIO_JWT_SECRET` | `jwt_secret` | `change-me-in-production` |

## Database

Folio uses SQLite with WAL mode. The database file is created automatically.

**Backup**: Copy `data/folio.db` to back up your library index, reading progress, and user accounts. The actual book files are not modified.

**WAL mode** allows concurrent reads while writing, which is important for the background scanner.

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
