package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M231Role(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Role{})
	if err != nil {
		log.Fatal(err)
	}
}
