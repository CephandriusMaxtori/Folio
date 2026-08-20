---
layout: default
title: Building from Source - Folio
---

# Building from Source

## Requirements

- Go 1.21+
- Make (optional)
- Git

No CGo or C compiler needed. All dependencies are pure Go.

## Quick Build

```bash
git clone https://github.com/CephandriusMaxtori/Folio.git
cd Folio
make build
```

This produces a `folio` binary for your current platform.

## Cross-Compile for Raspberry Pi

### Raspberry Pi 4/5 (ARM64)

```bash
make build-arm64
```

Produces `folio-linux-arm64`. Copy to your Pi and run:

```bash
scp folio-linux-arm64 pi@raspberrypi:~/
ssh pi@raspberrypi
chmod +x folio-linux-arm64
./folio-linux-arm64
```

### Raspberry Pi 3 (ARMv7)

```bash
make build-armv7
```

Produces `folio-linux-armv7`.

## Manual Cross-Compile

Without Make:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o folio-linux-arm64 ./cmd/folio
```

## Build for All Platforms

```bash
make build-all
```

Produces:
- `folio-linux-amd64`
- `folio-linux-arm64`
- `folio-linux-armv7`

## Optimized Build

The Makefile uses these flags by default:

```
-ldflags="-s -w"
```

- `-s`: Strip symbol table (smaller binary)
- `-w`: Strip DWARF debug info (smaller binary)

## Run Tests

```bash
make test
```

## Lint

```bash
make lint
```
