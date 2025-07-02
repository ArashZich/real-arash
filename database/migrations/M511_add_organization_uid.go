// M371_add_organization_uid.go
package migrations

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func M511AddOrganizationUid(db *gorm.DB) {
	// 1. اول چک می‌کنیم آیا ستون‌ها وجود دارند، اگر نه اضافه می‌کنیم
	var columnExists bool
	err := db.Raw(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = 'organizations' 
            AND column_name = 'organization_uid'
        );
    `).Scan(&columnExists).Error
	if err != nil {
		log.Fatal("Failed to check if column exists: ", err)
	}

	if !columnExists {
		err = db.Exec(`
            ALTER TABLE organizations 
            ADD COLUMN organization_uid uuid;
        `).Error
		if err != nil {
			log.Fatal("Failed to add organization_uid to organizations: ", err)
		}
	}

	// چک برای products
	err = db.Raw(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = 'products' 
            AND column_name = 'organization_uid'
        );
    `).Scan(&columnExists).Error
	if err != nil {
		log.Fatal("Failed to check if column exists: ", err)
	}

	if !columnExists {
		err = db.Exec(`
            ALTER TABLE products 
            ADD COLUMN organization_uid uuid;
        `).Error
		if err != nil {
			log.Fatal("Failed to add organization_uid to products: ", err)
		}
	}

	// چک برای documents
	err = db.Raw(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = 'documents' 
            AND column_name = 'organization_uid'
        );
    `).Scan(&columnExists).Error
	if err != nil {
		log.Fatal("Failed to check if column exists: ", err)
	}

	if !columnExists {
		err = db.Exec(`
            ALTER TABLE documents 
            ADD COLUMN organization_uid uuid;
        `).Error
		if err != nil {
			log.Fatal("Failed to add organization_uid to documents: ", err)
		}
	}

	// 2. برای organization های موجود که uuid ندارند uuid تولید می‌کنیم
	var organizations []models.Organization
	err = db.Where("organization_uid IS NULL").Find(&organizations).Error
	if err != nil {
		log.Fatal("Failed to fetch organizations: ", err)
	}

	for _, org := range organizations {
		newUUID := uuid.New()

		// UUID را برای organization به‌روز می‌کنیم
		err = db.Exec(`
            UPDATE organizations 
            SET organization_uid = ? 
            WHERE id = ? AND organization_uid IS NULL
        `, newUUID, org.ID).Error
		if err != nil {
			log.Fatal("Failed to update organization: ", err)
		}

		// UUID را برای products مرتبط به‌روز می‌کنیم
		err = db.Exec(`
            UPDATE products 
            SET organization_uid = ? 
            WHERE organization_id = ? AND organization_uid IS NULL
        `, newUUID, org.ID).Error
		if err != nil {
			log.Fatal("Failed to update products: ", err)
		}

		// UUID را برای documents مرتبط به‌روز می‌کنیم
		err = db.Exec(`
            UPDATE documents 
            SET organization_uid = ? 
            WHERE organization_id = ? AND organization_uid IS NULL
        `, newUUID, org.ID).Error
		if err != nil {
			log.Fatal("Failed to update documents: ", err)
		}
	}

	// 3. اضافه کردن محدودیت NOT NULL (اگر هنوز اضافه نشده)
	err = db.Exec(`
        ALTER TABLE organizations 
        ALTER COLUMN organization_uid SET NOT NULL;
    `).Error
	if err != nil {
		log.Printf("Warning: Failed to set not null constraint (might already exist): %v", err)
	}

	// 4. اضافه کردن محدودیت UNIQUE (اگر هنوز اضافه نشده)
	err = db.Exec(`
        ALTER TABLE organizations 
        ADD CONSTRAINT IF NOT EXISTS organizations_uid_unique UNIQUE (organization_uid);
    `).Error
	if err != nil {
		log.Printf("Warning: Failed to add unique constraint (might already exist): %v", err)
	}
}
