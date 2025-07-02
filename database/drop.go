package database

import (
	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/models/aggregation"
	"gorm.io/gorm"

	"log"
)

func DropJoinTables(db *gorm.DB) {
	// TODO: handle errors
	err := db.Migrator().DropTable("user_craft")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable("user_permission")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable("user_role")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable("plan_categories")
	if err != nil {
		log.Fatal(err)
	}
}

func DropTables(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Verification{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable(&models.Invite{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable(&aggregation.RegisterInvite{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable(&models.Role{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Token{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Category{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Category{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Organization{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.PackageItem{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Package{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Plan{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Product{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Document{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.InvoiceItem{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Invoice{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Coupon{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.View{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Notification{})
	if err != nil {
		log.Fatal(err)
	}
}

func DropAll(db *gorm.DB) {
	DropJoinTables(db.Debug())
	DropTables(db.Debug())
}
