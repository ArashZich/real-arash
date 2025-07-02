package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M421Notification(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Notification{})
	if err != nil {
		log.Fatal(err)
	}
}
