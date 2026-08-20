package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:255" json:"email"`
	PasswordHash string    `gorm:"size:255" json:"-"`
	Username     string    `gorm:"size:100" json:"username"`
	Role         string    `gorm:"size:20;default:user" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Library struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:255" json:"name"`
	Type         string    `gorm:"size:20" json:"type"`
	FolderPaths  string    `gorm:"type:text" json:"-"`
	LastScanned  time.Time `json:"last_scanned"`
	FolderWatch  bool      `gorm:"default:false" json:"folder_watching"`
	CreatedAt    time.Time `json:"created_at"`
	Folders      []string  `gorm:"-" json:"folder_paths"`
}

type Series struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LibraryID uint      `gorm:"index" json:"library_id"`
	Name      string    `gorm:"size:500" json:"name"`
	SortName  string    `gorm:"size:500" json:"sort_name"`
	LocalPath string    `gorm:"size:1000" json:"-"`
	CoverPath string    `gorm:"size:1000" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Library   Library   `gorm:"foreignKey:LibraryID" json:"-"`
	Volumes   []Volume  `gorm:"foreignKey:SeriesID" json:"volumes,omitempty"`
}

type Volume struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SeriesID  uint      `gorm:"index" json:"series_id"`
	Number    float64   `json:"number"`
	Name      string    `gorm:"size:255" json:"name"`
	CoverPath string    `gorm:"size:1000" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Series    Series    `gorm:"foreignKey:SeriesID" json:"-"`
	Chapters  []Chapter `gorm:"foreignKey:VolumeID" json:"chapters,omitempty"`
}

type Chapter struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VolumeID  uint      `gorm:"index" json:"volume_id"`
	SeriesID  uint      `gorm:"index" json:"series_id"`
	Number    string    `gorm:"size:50" json:"number"`
	Title     string    `gorm:"size:500" json:"title"`
	FilePath  string    `gorm:"size:1000" json:"-"`
	FileType  string    `gorm:"size:20" json:"file_type"`
	PageCount int       `json:"page_count"`
	CreatedAt time.Time `json:"created_at"`
	Volume    Volume    `gorm:"foreignKey:VolumeID" json:"-"`
}

type ReadingProgress struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ChapterID uint      `gorm:"index" json:"chapter_id"`
	Page      int       `json:"page"`
	ScrollPct float64   `json:"scroll_pct"`
	LastRead  time.Time `json:"last_read"`
	ReadCount int       `json:"read_count"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Chapter   Chapter   `gorm:"foreignKey:ChapterID" json:"-"`
}

type Bookmark struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ChapterID uint      `gorm:"index" json:"chapter_id"`
	Page      int       `json:"page"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Chapter   Chapter   `gorm:"foreignKey:ChapterID" json:"-"`
}

type Collection struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	UserID     uint             `gorm:"index" json:"user_id"`
	Name       string           `gorm:"size:255" json:"name"`
	CoverImage string           `gorm:"size:1000" json:"-"`
	CreatedAt  time.Time        `json:"created_at"`
	User       User             `gorm:"foreignKey:UserID" json:"-"`
	Items      []CollectionItem `gorm:"foreignKey:CollectionID" json:"items,omitempty"`
}

type CollectionItem struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	CollectionID uint       `gorm:"index" json:"collection_id"`
	SeriesID     uint       `gorm:"index" json:"series_id"`
	Position     int        `json:"position"`
	Collection   Collection `gorm:"foreignKey:CollectionID" json:"-"`
	Series       Series     `gorm:"foreignKey:SeriesID" json:"-"`
}

type ReadingList struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"index" json:"user_id"`
	Name      string           `gorm:"size:255" json:"name"`
	CreatedAt time.Time        `json:"created_at"`
	User      User             `gorm:"foreignKey:UserID" json:"-"`
	Items     []ReadingListItem `gorm:"foreignKey:ReadingListID" json:"items,omitempty"`
}

type ReadingListItem struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	ReadingListID uint      `gorm:"index" json:"reading_list_id"`
	ChapterID    uint        `gorm:"index" json:"chapter_id"`
	Position     int         `json:"position"`
	ReadingList  ReadingList `gorm:"foreignKey:ReadingListID" json:"-"`
	Chapter      Chapter     `gorm:"foreignKey:ChapterID" json:"-"`
}

type Annotation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ChapterID uint      `gorm:"index" json:"chapter_id"`
	Page      int       `json:"page"`
	Color     string    `gorm:"size:20" json:"color"`
	Text      string    `gorm:"type:text" json:"text"`
	Note      string    `gorm:"type:text" json:"note"`
	Shared    bool      `gorm:"default:false" json:"shared"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Chapter   Chapter   `gorm:"foreignKey:ChapterID" json:"-"`
}

type Settings struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UserID      uint   `gorm:"uniqueIndex" json:"user_id"`
	Theme       string `gorm:"size:50;default:dark" json:"theme"`
	Language    string `gorm:"size:10;default:en" json:"language"`
	ReaderMode  string `gorm:"size:20;default:single" json:"reader_mode"`
	Preferences string `gorm:"type:text" json:"preferences"`
	User        User   `gorm:"foreignKey:UserID" json:"-"`
}

type APIKey struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	Key       string     `gorm:"uniqueIndex;size:64" json:"key"`
	Name      string     `gorm:"size:100" json:"name"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
}
