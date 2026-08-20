package readers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nwaples/rardecode/v2"
)

type ArchiveEntry struct {
	Name    string
	Size    int64
	IsImage bool
}

type ArchiveReader struct {
	entries  []ArchiveEntry
	filePath string
	fileType string
}

func Open(path string) (*ArchiveReader, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".cbz", ".zip":
		return openZip(path)
	case ".cbr", ".rar":
		return openRar(path)
	default:
		return nil, fmt.Errorf("unsupported archive format: %s", ext)
	}
}

func openZip(path string) (*ArchiveReader, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var entries []ArchiveEntry
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := f.FileInfo().Name()
		ext := strings.ToLower(filepath.Ext(name))
		isImage := isImageExt(ext)
		entries = append(entries, ArchiveEntry{
			Name:    name,
			Size:    int64(f.UncompressedSize64),
			IsImage: isImage,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return naturalSort(entries[i].Name, entries[j].Name)
	})

	return &ArchiveReader{
		entries:  entries,
		filePath: path,
		fileType: strings.TrimPrefix(filepath.Ext(path), "."),
	}, nil
}

func openRar(path string) (*ArchiveReader, error) {
	r, err := rardecode.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var entries []ArchiveEntry
	for {
		h, err := r.Next()
		if err != nil {
			break
		}
		if h.IsDir {
			continue
		}
		ext := strings.ToLower(filepath.Ext(h.Name))
		isImage := isImageExt(ext)
		entries = append(entries, ArchiveEntry{
			Name:    filepath.Base(h.Name),
			Size:    h.UnPackedSize,
			IsImage: isImage,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return naturalSort(entries[i].Name, entries[j].Name)
	})

	return &ArchiveReader{
		entries:  entries,
		filePath: path,
		fileType: filepath.Ext(path)[1:],
	}, nil
}

func (a *ArchiveReader) Pages() int {
	count := 0
	for _, e := range a.entries {
		if e.IsImage {
			count++
		}
	}
	return count
}

func (a *ArchiveReader) GetPage(num int) ([]byte, string, error) {
	imgIdx := 0
	for _, e := range a.entries {
		if !e.IsImage {
			continue
		}
		if imgIdx == num {
			ext := strings.ToLower(filepath.Ext(e.Name))
			mime := extToMime(ext)
			data, err := a.readFile(e.Name)
			return data, mime, err
		}
		imgIdx++
	}
	return nil, "", fmt.Errorf("page %d not found", num)
}

func (a *ArchiveReader) GetPageInfo(num int) (string, int64, error) {
	imgIdx := 0
	for _, e := range a.entries {
		if !e.IsImage {
			continue
		}
		if imgIdx == num {
			return e.Name, e.Size, nil
		}
		imgIdx++
	}
	return "", 0, fmt.Errorf("page %d not found", num)
}

func (a *ArchiveReader) readFile(name string) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(a.filePath))
	switch ext {
	case ".cbz", ".zip":
		return readZipFile(a.filePath, name)
	case ".cbr", ".rar":
		return readRarFile(a.filePath, name)
	}
	return nil, fmt.Errorf("unsupported format")
}

func readZipFile(archivePath, entryName string) ([]byte, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if filepath.Base(f.Name) == entryName || f.Name == entryName {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("entry not found: %s", entryName)
}

func readRarFile(archivePath, entryName string) ([]byte, error) {
	r, err := rardecode.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for {
		h, err := r.Next()
		if err != nil {
			break
		}
		if filepath.Base(h.Name) == entryName || h.Name == entryName {
			return io.ReadAll(r)
		}
	}
	return nil, fmt.Errorf("entry not found: %s", entryName)
}

func (a *ArchiveReader) DetectCover() ([]byte, string, error) {
	if len(a.entries) == 0 {
		return nil, "", fmt.Errorf("no entries")
	}
	for _, e := range a.entries {
		if e.IsImage {
			ext := strings.ToLower(filepath.Ext(e.Name))
			mime := extToMime(ext)
			data, err := a.readFile(e.Name)
			if err != nil {
				return nil, "", err
			}
			if _, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
				return data, mime, nil
			}
		}
	}
	return nil, "", fmt.Errorf("no valid cover found")
}

func (a *ArchiveReader) Close() error {
	return nil
}

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff":
		return true
	}
	return false
}

func extToMime(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	default:
		return "application/octet-stream"
	}
}

func naturalSort(a, b string) bool {
	return a < b
}
