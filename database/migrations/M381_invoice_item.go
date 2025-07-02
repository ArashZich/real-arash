package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M381InvoiceItem(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.InvoiceItem{})
	if err != nil {
		log.Fatal(err)
	}
}
