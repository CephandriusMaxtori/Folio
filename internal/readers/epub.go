package readers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/simp-lee/epub"
)

type EpubBook struct {
	*epub.Book
	filePath string
}

func OpenEpub(path string) (*EpubBook, error) {
	b, err := epub.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open epub: %w", err)
	}
	return &EpubBook{Book: b, filePath: path}, nil
}

func (b *EpubBook) Title() string {
	md := b.Metadata()
	if len(md.Titles) > 0 {
		return md.Titles[0]
	}
	return filepath.Base(b.filePath)
}

func (b *EpubBook) Author() string {
	md := b.Metadata()
	if len(md.Authors) > 0 {
		return md.Authors[0].Name
	}
	return "Unknown"
}

func (b *EpubBook) CoverImage() ([]byte, string, error) {
	c, err := b.Cover()
	if err != nil {
		return nil, "", err
	}
	return c.Data, c.MediaType, nil
}

func (b *EpubBook) ChapterCount() int {
	return len(b.Chapters())
}

func (b *EpubBook) ChapterText(index int) (string, error) {
	chs := b.Chapters()
	if index < 0 || index >= len(chs) {
		return "", fmt.Errorf("chapter index out of range")
	}
	return chs[index].TextContent()
}

func (b *EpubBook) FileSize() int64 {
	info, err := os.Stat(b.filePath)
	if err != nil {
		return 0
	}
	return info.Size()
}
