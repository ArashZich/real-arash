package migrations

import (
	"gorm.io/gorm"

	"log"
)

func M412AddViewIndexes(db *gorm.DB) {
	// فقط ایندکس product_uid را اضافه می‌کنیم چون visit_uid قبلاً اضافه شده
	err := db.Exec("CREATE INDEX IF NOT EXISTS idx_views_product_uid ON views(product_uid)").Error
	if err != nil {
		log.Fatal(err)
	}
}
