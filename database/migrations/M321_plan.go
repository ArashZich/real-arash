package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M321Plan(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Plan{})
	if err != nil {
		log.Fatal(err)
	}
}
