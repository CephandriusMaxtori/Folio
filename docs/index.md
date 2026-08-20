---
layout: default
title: Folio
---

# Folio

A self-hosted digital library server for manga, comics, and ebooks. Optimized for Raspberry Pi.

Built with Go. No CGo dependencies. Single binary deployment.

## Features

- **Format support**: CBZ, CBR, ZIP, RAR, EPUB
- **OPDS 1.2**: Compatible with KOReader, Calibre, and other OPDS readers
- **Multi-user**: Role-based access control (admin/user)
- **Reading progress**: Tracks page, scroll position, and on-deck queue
- **Collections & reading lists**: Organize your library
- **Annotations**: Highlight and take notes
- **Full-text search**: Find anything across your library
- **Web UI**: Dark-themed SPA with book upload
- **ARM optimized**: Low memory footprint, designed for Raspberry Pi 4/5

## Quick Start

1. Download the latest binary for your platform from [Releases](https://github.com/CephandriusMaxtori/Folio/releases)
2. Run it:

```bash
./folio
```

3. Open `http://localhost:8080` in your browser
4. Register an account (first user becomes admin)
5. Create a library and point it at your book folder
6. Click **Scan** to index your library

## Links

- [Getting Started](getting-started.html)
- [Configuration](configuration.html)
- [API Reference](api.html)
- [KOReader Setup](koreader.html)
- [Build from Source](building.html)
