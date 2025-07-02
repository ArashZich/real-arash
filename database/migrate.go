// database/migrate.go
package database

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/database/migrations"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if config.AppConfig.Environment != "development" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are you sure you want to migrate the database in not development environment? (y/n): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.TrimSpace(response)
		if response != "y" {
			log.Fatal("Process terminated by the user.")
		}
	}
	M000Extensions(db)

	// migrate the dependent many2many before the main table
	migrations.M231Role(db)
	migrations.M222User(db)
	migrations.M228Verification(db)
	migrations.M229Invite(db)
	migrations.M230RegisterInvite(db)
	migrations.M281Token(db)
	migrations.M291Category(db)
	migrations.M311Organization(db)
	migrations.M321Plan(db)
	migrations.M331PackageItem(db)
	migrations.M341Package(db)
	migrations.M351Product(db)
	migrations.M361Document(db)
	migrations.M371Coupon(db)
	migrations.M381InvoiceItem(db)
	migrations.M391Invoice(db)
	migrations.M411View(db)
	migrations.M421Notification(db)
	migrations.M431Post(db)                   // اضافه کردن مدل پست جدید
	migrations.M511AddOrganizationUid(db)     // اضافه کردن ستون organization_uid
	migrations.M512AddShowroomUrl(db)         // اضافه کردن ستون showroom_url
	migrations.M513AddOrganizationType(db)    // اضافه کردن ستون organization_type
	migrations.M514AddShopLinkToDocuments(db) // اضافه کردن ستون shop_link به documents
	migrations.M412AddViewIndexes(db)         // اضافه کردن ایندکس‌های جدول views
}

func M000Extensions(db *gorm.DB) {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}
}
