package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M411View(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.View{})
	if err != nil {
		log.Fatal(err)
	}
}
