package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M291Category(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Category{})
	if err != nil {
		log.Fatal(err)
	}
}
