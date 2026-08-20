# Folio - Implementation Plan

A Go-based Kavita clone optimized for Raspberry Pi 4/5 (ARM64).
No CGo dependencies. Single binary deployment.

## Tech Stack

| Layer | Technology | Why |
|-------|-----------|-----|
| HTTP Router | `github.com/go-chi/chi/v5` | Lightweight, stdlib-compatible |
| Database | `modernc.org/sqlite` via `gorm.io/driver/sqlite` | Pure Go SQLite, no C toolchain |
| ORM | `gorm.io/gorm` | Auto-migrations, clean models |
| ZIP/CBZ | `archive/zip` (stdlib) | Zero dependencies |
| RAR/CBR | `github.com/nwaples/rardecode/v2` | Pure Go RAR decoder |
| 7Z/CB7 | `github.com/bodgit/sevenzip` | Pure Go 7z reader |
| EPUB | `github.com/simp-lee/epub` | Pure Go, ePub 2+3 |
| PDF | Client-side pdf.js | No server-side PDF processing |
| Auth | `github.com/golang-jwt/jwt/v5` + `golang.org/x/crypto` | JWT + bcrypt |
| Frontend | Vanilla HTML/CSS/JS (embedded `//go:embed`) | No Node.js build step |

## Project Structure

```
folio/
├── cmd/folio/main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── models/
│   ├── handlers/
│   ├── services/
│   ├── middleware/
│   ├── readers/
│   └── opds/
├── web/
│   ├── index.html
│   ├── css/style.css
│   └── js/
├── go.mod
├── Makefile
├── Dockerfile
├── Todo.md
└── README.md
```

## Database Schema (SQLite, WAL mode)

- **users**: id, email, password_hash, username, role, created_at
- **libraries**: id, name, type, folder_paths, last_scanned, folder_watching, created_at
- **series**: id, library_id, name, sort_name, local_path, cover_path, created_at
- **volumes**: id, series_id, number, name, cover_path, created_at
- **chapters**: id, volume_id, series_id, number, title, file_path, file_type, page_count, created_at
- **reading_progress**: id, user_id, chapter_id, page, scroll_pct, last_read, read_count, created_at
- **bookmarks**: id, user_id, chapter_id, page, created_at
- **collections**: id, user_id, name, cover_image, created_at
- **collection_items**: id, collection_id, series_id, position
- **reading_lists**: id, user_id, name, created_at
- **reading_list_items**: id, reading_list_id, chapter_id, position
- **annotations**: id, user_id, chapter_id, page, color, text, note, shared, created_at
- **settings**: id, user_id, theme, language, reader_mode, preferences (JSON)
- **api_keys**: id, user_id, key, name, expires_at, created_at

## API Endpoints

### REST (Web UI)
```
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/refresh
GET    /api/auth/me

GET    /api/libraries
POST   /api/libraries
PUT    /api/libraries/:id
DELETE /api/libraries/:id
POST   /api/libraries/:id/scan

GET    /api/series?library_id=&sort=&filter=
GET    /api/series/:id
GET    /api/series/:id/volumes
GET    /api/series/:id/metadata

GET    /api/reader/chapter/:id/pages
GET    /api/reader/chapter/:id/page/:num
POST   /api/reader/progress
GET    /api/reader/on-deck

GET    /api/search?q=&library_id=
GET    /api/collections
POST   /api/collections
POST   /api/collections/:id/series
GET    /api/reading-lists
POST   /api/reading-lists
POST   /api/reading-lists/:id/chapters

GET    /api/annotations?chapter_id=
POST   /api/annotations
PUT    /api/annotations/:id
DELETE /api/annotations/:id

GET    /api/settings
PUT    /api/settings
GET    /api/admin/stats
```

### OPDS (KOReader / fetcher.koplugin compatible)
```
GET  /opds/{apiKey}                                    # Root catalog
GET  /opds/{apiKey}/libraries                          # List libraries
GET  /opds/{apiKey}/library/{libraryId}                # Series in library
GET  /opds/{apiKey}/series/{seriesId}                  # Series detail
GET  /opds/{apiKey}/series/{seriesId}/volume/{volumeId}          # Volume chapters
GET  /opds/{apiKey}/series/{seriesId}/volume/{volumeId}/chapter/{chapterId}  # Chapter files
GET  /opds/{apiKey}/series/{seriesId}/volume/{volumeId}/chapter/{chapterId}/download/{filename}
GET  /opds/{apiKey}/collections
GET  /opds/{apiKey}/collection/{id}
GET  /opds/{apiKey}/reading-list
GET  /opds/{apiKey}/reading-list/{id}
GET  /opds/{apiKey}/want-to-read
GET  /opds/{apiKey}/search/{query}
GET  /opds/search/{apiKey}                             # OpenSearch descriptor
```

### Cover Images
```
GET /api/image/series-cover?seriesId={id}
GET /api/image/volume-cover?volumeId={id}
GET /api/image/chapter-cover?chapterId={id}
```

## RPi Optimizations

1. SQLite WAL mode - concurrent reads during writes
2. Connection pool: 2-4 max
3. Lazy thumbnail generation
4. Streaming archive reads (no full extraction to memory)
5. Bounded goroutine scanner (4 workers)
6. Gzip compression on all API responses
7. Long Cache-Control for static assets
8. ARM64 cross-compilation (CGO_ENABLED=0)
9. Binary size optimization (-ldflags="-s -w")
10. Low memory profile throughout

## KOReader / fetcher.koplugin Compatibility

The fetcher plugin uses KOReader's built-in OPDS browser to download books.
Users add Folio as an OPDS catalog in KOReader, then Fetcher auto-syncs.

Requirements:
- OPDS 1.2 XML/Atom feeds (primary)
- OPDS 2.0 JSON (optional)
- API key auth via URL path (like Kavita)
- Standard Atom link rel attributes (acquisition, image, thumbnail, subsection)
- OpenSearch descriptor for search
- Pagination support (prev/next links)
- Cover images and thumbnails
- Content-Disposition headers for downloads

## Implementation Steps

- [x] Step 1: Project init - go.mod, Makefile, README, .gitignore
- [x] Step 2: Config & DB - SQLite WAL, GORM, auto-migrate
- [x] Step 3: Models - all GORM model definitions
- [x] Step 4: Auth - JWT, bcrypt, login/register, middleware
- [x] Step 5: Library CRUD
- [x] Step 6: Scanner - library scanning engine
- [x] Step 7: Archive readers - CBZ/CBR
- [x] Step 8: EPUB reader
- [x] Step 9: Cover service
- [x] Step 10: Reader API - pages, progress, on-deck
- [x] Step 11: OPDS feeds - full OPDS 1.2 XML
- [x] Step 12: Search
- [x] Step 13: Collections & reading lists
- [x] Step 14: Annotations
- [x] Step 15: Admin
- [x] Step 16: Settings
- [ ] Step 17: Frontend SPA
- [ ] Step 18: Final polish - error handling, logging, Docker multi-arch
