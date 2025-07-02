package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M229Invite(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Invite{})
	if err != nil {
		log.Fatal(err)
	}
}
