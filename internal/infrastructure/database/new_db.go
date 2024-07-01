package database

import (
	"emailn/internal/domain/campaign"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewBD() *gorm.DB {
	dsn := "host=localhost user=admin password=admin dbname=emailn port=5434 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("fail to connect to database")
	}

	db.AutoMigrate(&campaign.Campaign{}, &campaign.Contact{})

	return db
}
