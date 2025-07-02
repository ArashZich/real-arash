package migrations

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"
)

func M361Document(db *gorm.DB) {
	// انجام عملیات AutoMigrate برای ایجاد ستون‌های مورد نیاز
	err := db.Migrator().AutoMigrate(&models.Document{})
	if err != nil {
		log.Fatal(err)
	}

	// تغییر نوع داده‌ی product_uid به uuid
	err = db.Exec(`
        ALTER TABLE documents ALTER COLUMN product_uid TYPE uuid USING product_uid::uuid;
    `).Error
	if err != nil {
		log.Fatal("Failed to alter column type: ", err)
	}

	// اجرای کوئری به‌روزرسانی
	err = db.Exec(`
        UPDATE documents SET product_id = (
            SELECT id FROM products WHERE documents.product_uid = products.product_uid
        );
    `).Error
	if err != nil {
		log.Fatal("Failed to update product_id: ", err)
	}
}
