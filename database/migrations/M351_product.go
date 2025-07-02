package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M351Product(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal(err)
	}
}
