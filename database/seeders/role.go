package seeders

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/services/role"
	"gorm.io/gorm"
)

func RoleSeeder(db *gorm.DB) {

	var roles []models.Role
	for _, v := range role.Roles {
		roles = append(
			roles, models.Role{
				Title: v,
			},
		)
	}
	err := db.Create(&roles).Error
	if err != nil {
		log.Fatal(err)
	}
}
