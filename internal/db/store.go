package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/example/privai-vault/internal/models"
)

type GormStore struct {
	DB *gorm.DB
}

func NewGormStore(path string) (*GormStore, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &GormStore{DB: db}, nil
}

func (s *GormStore) AutoMigrate() error {
	return s.DB.AutoMigrate(&models.Query{})
}

func (s *GormStore) SaveQuery(q *models.Query) error {
	return s.DB.Create(q).Error
}

func (s *GormStore) ListQueries() ([]models.Query, error) {
	var out []models.Query
	err := s.DB.Order("created_at asc").Find(&out).Error
	return out, err
}
