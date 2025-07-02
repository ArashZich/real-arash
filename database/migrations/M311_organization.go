package migrations

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"

	"log"
)

func M311Organization(db *gorm.DB) {
	// فقط AutoMigrate ساده بدون اضافه کردن محدودیت‌های organization_uid
	err := db.Migrator().AutoMigrate(&models.Organization{})
	if err != nil {
		log.Fatal(err)
	}

	// حذف محدودیت NOT NULL و UNIQUE اگر از قبل وجود داشته
	err = db.Exec(`
		 ALTER TABLE organizations 
		 DROP CONSTRAINT IF EXISTS organizations_uid_unique;
	 `).Error
	if err != nil {
		log.Fatal("Failed to drop unique constraint: ", err)
	}

	err = db.Exec(`
		 ALTER TABLE organizations 
		 ALTER COLUMN organization_uid DROP NOT NULL;
	 `).Error
	if err != nil {
		log.Fatal("Failed to drop not null constraint: ", err)
	}
}
