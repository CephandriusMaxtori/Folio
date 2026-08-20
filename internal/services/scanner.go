package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/CephandriusMaxtori/Folio/internal/models"
	"github.com/CephandriusMaxtori/Folio/internal/readers"
	"gorm.io/gorm"
)

var archiveExts = map[string]bool{
	".cbz": true, ".zip": true,
	".cbr": true, ".rar": true,
}

var bookExts = map[string]bool{
	".epub": true,
}

type ScanResult struct {
	LibrariesScanned int
	SeriesFound      int
	ChaptersFound    int
	Errors           []string
}

func (s *Service) ScanLibrary(libID uint) (*ScanResult, error) {
	lib, err := s.GetLibrary(libID)
	if err != nil {
		return nil, err
	}

	result := &ScanResult{}
	folders := lib.GetFolders()

	for _, folder := range folders {
		s.scanFolder(folder, lib, result)
	}

	lib.LastScanned = time.Now()
	s.db.Save(lib)

	return result, nil
}

func (s *Service) scanFolder(folder string, lib *models.Library, result *ScanResult) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		seriesPath := filepath.Join(folder, entry.Name())
		series := s.findOrCreateSeries(lib.ID, entry.Name(), seriesPath)
		result.SeriesFound++

		s.scanSeries(series, result)
	}
}

func (s *Service) scanSeries(series *models.Series, result *ScanResult) {
	entries, err := os.ReadDir(series.LocalPath)
	if err != nil {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return naturalLess(entries[i].Name(), entries[j].Name())
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		fullPath := filepath.Join(series.LocalPath, name)

		if archiveExts[ext] {
			s.processArchiveFile(series, fullPath, name, ext, result)
		} else if bookExts[ext] {
			s.processBookFile(series, fullPath, name, ext, result)
		}
	}
}

func (s *Service) processArchiveFile(series *models.Series, path, name, ext string, result *ScanResult) {
	vol, ch := parseMangaName(name)

	volNum, _ := strconv.ParseFloat(vol, 64)
	volume := s.findOrCreateVolume(series.ID, volNum)
	chapter := s.findOrCreateChapter(volume.ID, series.ID, ch, name, path, ext[1:])

	archive, err := readers.Open(path)
	if err != nil {
		log.Printf("Failed to open %s: %v", path, err)
		return
	}
	defer archive.Close()

	if chapter.PageCount == 0 {
		chapter.PageCount = archive.Pages()
		s.db.Save(chapter)
	}

	if series.CoverPath == "" {
		data, _, err := archive.DetectCover()
		if err == nil {
			coverPath := filepath.Join(series.LocalPath, ".covers", series.Name+".jpg")
			os.MkdirAll(filepath.Dir(coverPath), 0755)
			os.WriteFile(coverPath, data, 0644)
			series.CoverPath = coverPath
			s.db.Save(series)
		}
	}

	result.ChaptersFound++
}

func (s *Service) processBookFile(series *models.Series, path, name, ext string, result *ScanResult) {
	ch := &models.Chapter{
		SeriesID:  series.ID,
		Number:    "1",
		Title:     strings.TrimSuffix(name, ext),
		FilePath:  path,
		FileType:  ext[1:],
		PageCount: 1,
	}

	var existing models.Chapter
	if err := s.db.Where("file_path = ?", path).First(&existing).Error; err == nil {
		return
	}

	s.db.Create(ch)
	result.ChaptersFound++
}

func (s *Service) findOrCreateSeries(libraryID uint, name, path string) *models.Series {
	var series models.Series
	err := s.db.Where("library_id = ? AND local_path = ?", libraryID, path).First(&series).Error
	if err == gorm.ErrRecordNotFound {
		series = models.Series{
			LibraryID: libraryID,
			Name:      name,
			SortName:  strings.ToLower(name),
			LocalPath: path,
		}
		s.db.Create(&series)
	}
	return &series
}

func (s *Service) findOrCreateVolume(seriesID uint, number float64) *models.Volume {
	var vol models.Volume
	err := s.db.Where("series_id = ? AND number = ?", seriesID, number).First(&vol).Error
	if err == gorm.ErrRecordNotFound {
		vol = models.Volume{
			SeriesID: seriesID,
			Number:   number,
			Name:     strconv.FormatFloat(number, 'f', -1, 64),
		}
		s.db.Create(&vol)
	}
	return &vol
}

func (s *Service) findOrCreateChapter(volumeID, seriesID uint, number, title, path, fileType string) *models.Chapter {
	var ch models.Chapter
	err := s.db.Where("file_path = ?", path).First(&ch).Error
	if err == gorm.ErrRecordNotFound {
		ch = models.Chapter{
			VolumeID: volumeID,
			SeriesID: seriesID,
			Number:   number,
			Title:    title,
			FilePath: path,
			FileType: fileType,
		}
		s.db.Create(&ch)
	}
	return &ch
}

var volRegex = regexp.MustCompile(`(?i)(?:vol|volume)[\s._-]*(\d+(?:\.\d+)?)`)
var chRegex = regexp.MustCompile(`(?i)(?:ch|chapter)[\s._-]*(\d+(?:\.\d+)?)`)

func parseMangaName(name string) (volume, chapter string) {
	vol := "1"
	ch := "1"

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	if m := volRegex.FindStringSubmatch(base); len(m) > 1 {
		vol = m[1]
	}
	if m := chRegex.FindStringSubmatch(base); len(m) > 1 {
		ch = m[1]
	}

	return vol, ch
}

func naturalLess(a, b string) bool {
	return a < b
}

// ImportFile saves an uploaded file into a library's folder and indexes it.
func (s *Service) ImportFile(libraryID uint, seriesName, filename string, data io.Reader) error {
	lib, err := s.GetLibrary(libraryID)
	if err != nil {
		return fmt.Errorf("library not found: %w", err)
	}

	folders := lib.GetFolders()
	if len(folders) == 0 {
		return fmt.Errorf("library has no folders configured")
	}

	destDir := filepath.Join(folders[0], seriesName)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	destPath := filepath.Join(destDir, filename)
	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	written, err := io.Copy(f, data)
	if err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed to write file: %w", err)
	}
	log.Printf("Imported %s (%d bytes) into library %d", filename, written, libraryID)

	series := s.findOrCreateSeries(libraryID, seriesName, destDir)
	ext := strings.ToLower(filepath.Ext(filename))

	if archiveExts[ext] {
		s.processArchiveFile(series, destPath, filename, ext, &ScanResult{})
	} else if bookExts[ext] {
		s.processBookFile(series, destPath, filename, ext, &ScanResult{})
	}

	return nil
}
