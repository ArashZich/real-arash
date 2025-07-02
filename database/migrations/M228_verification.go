package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M228Verification(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Verification{})
	if err != nil {
		log.Fatal(err)
	}
}
