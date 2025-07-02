package migrations

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"
)

func M222User(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
}
