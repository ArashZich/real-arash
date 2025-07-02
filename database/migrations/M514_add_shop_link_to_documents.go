package migrations

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"
)

func M514AddShopLinkToDocuments(db *gorm.DB) {
	// چک کردن وجود ستون
	if !db.Migrator().HasColumn(&models.Document{}, "shop_link") {
		// اضافه کردن ستون با مقدار پیش‌فرض NULL (اختیاری)
		err := db.Exec(`
            ALTER TABLE documents 
            ADD COLUMN shop_link TEXT DEFAULT NULL;
        `).Error
		if err != nil {
			log.Fatal("Failed to add shop_link column: ", err)
		}

		log.Println("shop_link column added to documents table successfully")
	} else {
		log.Println("shop_link column already exists in documents table")
	}
}
