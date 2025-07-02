package migrations

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"
)

func M513AddOrganizationType(db *gorm.DB) {
	// چک کردن وجود ستون
	if !db.Migrator().HasColumn(&models.Organization{}, "organization_type") {
		// اضافه کردن ستون با مقدار پیش‌فرض basic
		err := db.Exec(`
            ALTER TABLE organizations 
            ADD COLUMN organization_type TEXT NOT NULL DEFAULT 'basic';
        `).Error
		if err != nil {
			log.Fatal("Failed to add organization_type column: ", err)
		}

		// آپدیت تمام رکوردهای موجود به basic
		err = db.Exec(`
            UPDATE organizations 
            SET organization_type = 'basic' 
            WHERE organization_type IS NULL;
        `).Error
		if err != nil {
			log.Fatal("Failed to update existing records: ", err)
		}
	}
}
