package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models/aggregation"
	"gorm.io/gorm"

	"log"
)

func M230RegisterInvite(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&aggregation.RegisterInvite{})
	if err != nil {
		log.Fatal(err)
	}
}
