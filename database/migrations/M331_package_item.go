package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M331PackageItem(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.PackageItem{})
	if err != nil {
		log.Fatal(err)
	}
}
