package services

import (
	"github.com/CephandriusMaxtori/Folio/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

// User
func (s *Service) CreateUser(user *models.User) error { return s.db.Create(user).Error }
func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := s.db.Where("email = ?", email).First(&u).Error
	return &u, err
}
func (s *Service) GetUserByID(id uint) (*models.User, error) {
	var u models.User
	err := s.db.First(&u, id).Error
	return &u, err
}
func (s *Service) ListUsers() ([]models.User, error) {
	var users []models.User
	err := s.db.Find(&users).Error
	return users, err
}

// Library
func (s *Service) CreateLibrary(lib *models.Library) error { return s.db.Create(lib).Error }
func (s *Service) GetLibrary(id uint) (*models.Library, error) {
	var l models.Library
	err := s.db.First(&l, id).Error
	return &l, err
}
func (s *Service) ListLibraries() ([]models.Library, error) {
	var libs []models.Library
	err := s.db.Find(&libs).Error
	return libs, err
}
func (s *Service) UpdateLibrary(lib *models.Library) error { return s.db.Save(lib).Error }
func (s *Service) DeleteLibrary(id uint) error { return s.db.Delete(&models.Library{}, id).Error }

// Series
func (s *Service) CreateSeries(series *models.Series) error { return s.db.Create(series).Error }
func (s *Service) GetSeries(id uint) (*models.Series, error) {
	var ser models.Series
	err := s.db.Preload("Volumes.Chapters").First(&ser, id).Error
	return &ser, err
}
func (s *Service) ListSeries(libraryID uint, sort string) ([]models.Series, error) {
	var series []models.Series
	q := s.db.Where("library_id = ?", libraryID)
	if sort != "" {
		q = q.Order(sort)
	} else {
		q = q.Order("sort_name")
	}
	err := q.Find(&series).Error
	return series, err
}
func (s *Service) UpdateSeries(series *models.Series) error { return s.db.Save(series).Error }
func (s *Service) DeleteSeries(id uint) error { return s.db.Delete(&models.Series{}, id).Error }

// Volume
func (s *Service) CreateVolume(vol *models.Volume) error { return s.db.Create(vol).Error }
func (s *Service) GetVolume(id uint) (*models.Volume, error) {
	var v models.Volume
	err := s.db.Preload("Chapters").First(&v, id).Error
	return &v, err
}

// Chapter
func (s *Service) CreateChapter(ch *models.Chapter) error { return s.db.Create(ch).Error }
func (s *Service) GetChapter(id uint) (*models.Chapter, error) {
	var c models.Chapter
	err := s.db.First(&c, id).Error
	return &c, err
}

// ReadingProgress
func (s *Service) GetProgress(userID, chapterID uint) (*models.ReadingProgress, error) {
	var p models.ReadingProgress
	err := s.db.Where("user_id = ? AND chapter_id = ?", userID, chapterID).First(&p).Error
	if err == gorm.ErrRecordNotFound {
		return &models.ReadingProgress{UserID: userID, ChapterID: chapterID}, nil
	}
	return &p, err
}
func (s *Service) SaveProgress(p *models.ReadingProgress) error {
	return s.db.Save(p).Error
}
func (s *Service) GetOnDeck(userID uint) ([]models.ReadingProgress, error) {
	var progress []models.ReadingProgress
	err := s.db.Where("user_id = ?", userID).Order("last_read DESC").Limit(20).Find(&progress).Error
	return progress, err
}

// Search
func (s *Service) Search(query string, libraryID uint) ([]models.Series, error) {
	var series []models.Series
	q := s.db.Where("name LIKE ?", "%"+query+"%")
	if libraryID > 0 {
		q = q.Where("library_id = ?", libraryID)
	}
	err := q.Find(&series).Error
	return series, err
}

// Collection
func (s *Service) CreateCollection(c *models.Collection) error { return s.db.Create(c).Error }
func (s *Service) ListCollections(userID uint) ([]models.Collection, error) {
	var cols []models.Collection
	err := s.db.Where("user_id = ?", userID).Preload("Items").Find(&cols).Error
	return cols, err
}
func (s *Service) AddToCollection(collectionID, seriesID uint, pos int) error {
	item := models.CollectionItem{
		CollectionID: collectionID,
		SeriesID:     seriesID,
		Position:     pos,
	}
	return s.db.Create(&item).Error
}

// ReadingList
func (s *Service) CreateReadingList(rl *models.ReadingList) error { return s.db.Create(rl).Error }
func (s *Service) ListReadingLists(userID uint) ([]models.ReadingList, error) {
	var lists []models.ReadingList
	err := s.db.Where("user_id = ?", userID).Preload("Items").Find(&lists).Error
	return lists, err
}
func (s *Service) AddToReadingList(readingListID, chapterID uint, pos int) error {
	item := models.ReadingListItem{
		ReadingListID: readingListID,
		ChapterID:     chapterID,
		Position:      pos,
	}
	return s.db.Create(&item).Error
}

// Annotation
func (s *Service) CreateAnnotation(a *models.Annotation) error { return s.db.Create(a).Error }
func (s *Service) ListAnnotations(userID, chapterID uint) ([]models.Annotation, error) {
	var anns []models.Annotation
	q := s.db.Where("user_id = ?", userID)
	if chapterID > 0 {
		q = q.Where("chapter_id = ?", chapterID)
	}
	err := q.Find(&anns).Error
	return anns, err
}
func (s *Service) UpdateAnnotation(a *models.Annotation) error { return s.db.Save(a).Error }
func (s *Service) DeleteAnnotation(id uint) error { return s.db.Delete(&models.Annotation{}, id).Error }

// Settings
func (s *Service) GetSettings(userID uint) (*models.Settings, error) {
	var st models.Settings
	err := s.db.Where("user_id = ?", userID).First(&st).Error
	if err == gorm.ErrRecordNotFound {
		st = models.Settings{UserID: userID}
		s.db.Create(&st)
		return &st, nil
	}
	return &st, err
}
func (s *Service) UpdateSettings(st *models.Settings) error { return s.db.Save(st).Error }

// APIKey
func (s *Service) CreateAPIKey(k *models.APIKey) error { return s.db.Create(k).Error }
func (s *Service) GetAPIKey(key string) (*models.APIKey, error) {
	var k models.APIKey
	err := s.db.Where("key = ?", key).First(&k).Error
	return &k, err
}
func (s *Service) DeleteAPIKey(id uint) error { return s.db.Delete(&models.APIKey{}, id).Error }

// Stats
func (s *Service) GetStats() (map[string]int64, error) {
	stats := make(map[string]int64)
	var count int64
	s.db.Model(&models.User{}).Count(&count); stats["users"] = count
	s.db.Model(&models.Library{}).Count(&count); stats["libraries"] = count
	s.db.Model(&models.Series{}).Count(&count); stats["series"] = count
	s.db.Model(&models.Volume{}).Count(&count); stats["volumes"] = count
	s.db.Model(&models.Chapter{}).Count(&count); stats["chapters"] = count
	return stats, nil
}
