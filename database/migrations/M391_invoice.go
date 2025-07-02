package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M391Invoice(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Invoice{})
	if err != nil {
		log.Fatal(err)
	}
}
