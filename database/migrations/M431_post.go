package migrations

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"
)

func M431Post(db *gorm.DB) {
	// AutoMigrate to update the schema to match the struct
	err := db.Migrator().AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatal(err)
	}

	// Check if the column "Views" exists
	hasViews := db.Migrator().HasColumn(&models.Post{}, "Views")
	if !hasViews {
		// Add the "Views" column with default value of 0
		err = db.Exec("ALTER TABLE posts ADD COLUMN views INT DEFAULT 0").Error
		if err != nil {
			log.Fatal(err)
		}

		// Update existing records to set the default value of Views to 0 if Views is NULL
		err = db.Model(&models.Post{}).Where("views IS NULL").Update("views", 0).Error
		if err != nil {
			log.Fatal(err)
		}
	}
}
