package migrations

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"
)

func M512AddShowroomUrl(db *gorm.DB) {
	// اول چک میکنیم که ستون وجود نداره
	if !db.Migrator().HasColumn(&models.Organization{}, "showroom_url") {
		err := db.Exec(`
            ALTER TABLE organizations 
            ADD COLUMN showroom_url TEXT DEFAULT '';
        `).Error
		if err != nil {
			log.Fatal("Failed to add showroom_url column: ", err)
		}
	}
}
