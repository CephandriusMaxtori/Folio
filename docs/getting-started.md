---
layout: default
title: Getting Started - Folio
---

# Getting Started

## 1. Install

### Download a release

Grab the latest binary from [GitHub Releases](https://github.com/CephandriusMaxtori/Folio/releases).

| Platform | File |
|----------|------|
| Linux x86_64 | `folio-linux-amd64` |
| Linux ARM64 (Raspberry Pi 4/5) | `folio-linux-arm64` |
| Linux ARMv7 (Raspberry Pi 3) | `folio-linux-armv7` |
| Windows | `folio-windows-amd64.exe` |
| macOS | `folio-darwin-amd64` |

### Build from source

```bash
git clone https://github.com/CephandriusMaxtori/Folio.git
cd Folio
make build
```

## 2. Run

```bash
chmod +x folio-linux-arm64
./folio-linux-arm64
```

Folio will:
1. Create a `config.yaml` in the current directory (if none exists)
2. Create an SQLite database at `./data/folio.db`
3. Start the web server on port `8080`

## 3. First Login

1. Open `http://your-pi-ip:8080` in a browser
2. Register an account
3. The first user automatically gets the **admin** role

## 4. Add a Library

1. Click **+ New Library** on the home screen
2. Enter a name (e.g., "Manga")
3. Enter the folder path on your system (e.g., `/mnt/usb/manga`)
4. Click **Create**
5. Click **Scan** on the library card

Folio will scan the folder for supported files (CBZ, CBR, ZIP, RAR, EPUB) and index them.

## 5. Import Books

You can also import individual books via the web UI:

1. Click **Import** in the top bar
2. Select a library
3. Drag and drop or browse for a file
4. Click **Import**

## 6. Read

Click on a series, then a chapter to start reading. Use arrow keys or on-screen buttons to navigate pages.
