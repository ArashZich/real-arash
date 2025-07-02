package database

import (
	"gitag.ir/armogroup/armo/services/reality/database/seeders"
	"gorm.io/gorm"
)

func SeedAllTables(db *gorm.DB) {
	seeders.RoleSeeder(db)
	seeders.UserSeeder(db)
	seeders.CategorySeeder(db)
	seeders.PlanSeeder(db)
}
