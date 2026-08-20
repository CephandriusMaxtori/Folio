---
layout: default
title: KOReader Setup - Folio
---

# KOReader Setup

Folio is fully compatible with KOReader's built-in OPDS browser. This means you can browse and download books directly from your Folio server on your e-reader.

## Prerequisites

- A running Folio server (e.g., on your Raspberry Pi)
- KOReader installed on your e-reader
- Both devices on the same network (or Folio accessible via internet)

## Step 1: Create an API Key

1. Log in to the Folio web UI
2. Go to **Settings** or click your username
3. Navigate to **API Keys**
4. Click **Create** and give it a name (e.g., "KOReader")
5. Copy the generated key

## Step 2: Add Folio to KOReader

1. Open KOReader
2. Go to **Tools > OPDS catalogs**
3. Tap **+** to add a new catalog
4. Enter the catalog details:
   - **Name**: `Folio` (or whatever you like)
   - **URL**: `http://your-server-ip:8080/opds/YOUR_API_KEY`
5. Save

## Step 3: Browse and Download

1. Go to **Tools > OPDS catalogs**
2. Select **Folio**
3. Browse: Libraries > Series > Volume > Chapter
4. Tap a chapter to download it
5. The book will be saved to KOReader's downloads folder

## fetcher.koplugin

If you use the [fetcher.koplugin](https://github.com/nicoboss/fetcher.koplugin) plugin for automatic syncing:

1. Install fetcher.koplugin in KOReader
2. In fetcher settings, add the Folio OPDS URL
3. The plugin will use KOReader's built-in OPDS browser to download books

## Troubleshooting

**"Connection refused"**: Make sure Folio is running and the port is open.

**"Invalid API key"**: Double-check the API key in the URL. Create a new one if needed.

**Books not showing**: Make sure you've scanned your library in the Folio web UI.

**Slow browsing**: The first browse triggers metadata loading. Subsequent visits are faster since Folio caches in SQLite.
