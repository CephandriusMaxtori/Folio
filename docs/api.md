---
layout: default
title: API Reference - Folio
---

# API Reference

All endpoints require a `Bearer` token (JWT) in the `Authorization` header, unless noted otherwise.

## Authentication

### Register

```
POST /api/auth/register
```

```json
{
  "email": "user@example.com",
  "password": "password",
  "username": "optional"
}
```

Response: `{ "token": "...", "user": { ... } }`

### Login

```
POST /api/auth/login
```

```json
{
  "email": "user@example.com",
  "password": "password"
}
```

Response: `{ "token": "...", "user": { ... } }`

### Get Current User

```
GET /api/auth/me
```

### API Keys

```
GET    /api/auth/apikeys        # List keys
POST   /api/auth/apikeys        # Create key { "name": "My Reader" }
DELETE /api/auth/apikeys/{id}   # Delete key
```

## Libraries

```
GET    /api/libraries           # List all
POST   /api/libraries           # Create { "name": "...", "type": "comic", "folder_paths": [...] }
PUT    /api/libraries/{id}      # Update
DELETE /api/libraries/{id}      # Delete
POST   /api/libraries/{id}/scan # Scan in background
```

## Series

```
GET /api/series?library_id=&sort=      # List series
GET /api/series/{id}                    # Series detail (includes volumes/chapters)
GET /api/series/{id}/volumes            # Volumes only
```

## Reader

```
GET /api/reader/chapter/{id}/pages      # Get page count
GET /api/reader/chapter/{id}/page/{num} # Get page content (image or text)
POST /api/reader/progress               # Save progress { chapterId, page, scrollPct }
GET /api/reader/on-deck                 # Continue reading queue
```

## Cover Images

```
GET /api/image/series-cover?seriesId={id}
GET /api/image/volume-cover?volumeId={id}
GET /api/image/chapter-cover?chapterId={id}
```

## Search

```
GET /api/search?q=one+piece&library_id=1
```

## Import

```
POST /api/import
Content-Type: multipart/form-data

Fields:
  file:        The book file (CBZ, CBR, ZIP, RAR, EPUB)
  library_id:  Target library ID
  series_name: Optional series name
```

## Collections

```
GET  /api/collections              # List
POST /api/collections              # Create { "name": "..." }
POST /api/collections/{id}/series  # Add series { "seriesId": 1, "position": 0 }
```

## Reading Lists

```
GET  /api/reading-lists                # List
POST /api/reading-lists                # Create { "name": "..." }
POST /api/reading-lists/{id}/chapters  # Add chapter { "chapterId": 1, "position": 0 }
```

## Annotations

```
GET    /api/annotations?chapter_id=1  # List for chapter
POST   /api/annotations               # Create
PUT    /api/annotations/{id}          # Update
DELETE /api/annotations/{id}          # Delete
```

## Settings

```
GET /api/settings    # Get user settings
PUT /api/settings    # Update settings
```

## Admin

```
GET /api/admin/stats  # { "users": 1, "libraries": 2, "series": 50, "chapters": 1200 }
```

## OPDS 1.2

OPDS endpoints use API key authentication via the URL path:

```
GET /opds/{apiKey}                                        # Root catalog
GET /opds/{apiKey}/libraries                              # Libraries
GET /opds/{apiKey}/library/{libraryId}                    # Series in library
GET /opds/{apiKey}/series/{seriesId}                      # Series detail
GET /opds/{apiKey}/series/{seriesId}/volume/{volumeId}    # Volume chapters
GET /opds/{apiKey}/series/{seriesId}/volume/{volumeId}/chapter/{chapterId}/download/{filename}
GET /opds/search/{apiKey}?searchTerms={query}             # Search
GET /opds/search/{apiKey}                                 # OpenSearch descriptor
```

Returns `application/atom+xml;profile=opds-catalog` XML feeds.
