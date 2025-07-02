package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M281Token(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Token{})
	if err != nil {
		log.Fatal(err)
	}
}
