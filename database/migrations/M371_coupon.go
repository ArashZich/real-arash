package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M371Coupon(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Coupon{})
	if err != nil {
		log.Fatal(err)
	}
}
