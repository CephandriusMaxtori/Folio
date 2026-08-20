package database

import (
	"github.com/CephandriusMaxtori/Folio/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Library{},
		&models.Series{},
		&models.Volume{},
		&models.Chapter{},
		&models.ReadingProgress{},
		&models.Bookmark{},
		&models.Collection{},
		&models.CollectionItem{},
		&models.ReadingList{},
		&models.ReadingListItem{},
		&models.Annotation{},
		&models.Settings{},
		&models.APIKey{},
	)
}
