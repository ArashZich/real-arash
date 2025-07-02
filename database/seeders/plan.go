package seeders

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"
)

func PlanSeeder(db *gorm.DB) {
	planSeeders := []models.Plan{

		{
			// ID: 1
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_50",
			DayLength:      90,
			ProductLimit:   50,
			StorageLimitMB: 50,
			IconUrl:        "www.armogroup.com/assets/images/2",
			Price:          1200000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 600000,
			UserID:          2,
		},

		{
			// ID: 2
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_50",
			DayLength:      180,
			ProductLimit:   50,
			StorageLimitMB: 50,
			IconUrl:        "",
			Price:          2000000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 1000000,
			UserID:          2,
		},

		{
			// ID: 3
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_50",
			DayLength:      365,
			ProductLimit:   50,
			StorageLimitMB: 50,
			IconUrl:        "",
			Price:          3000000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 1500000,
			UserID:          2,
		},

		{
			// ID: 4
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_500,onboarding",
			DayLength:      90,
			ProductLimit:   500,
			StorageLimitMB: 500,
			IconUrl:        "",
			Price:          3000000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 1500000,
			UserID:          2,
		},

		{
			// ID: 5
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_500,onboarding",
			DayLength:      180,
			ProductLimit:   500,
			StorageLimitMB: 500,
			IconUrl:        "",
			Price:          5000000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 2500000,
			UserID:          2,
		},

		{
			// ID: 6
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_500,onboarding",
			DayLength:      365,
			ProductLimit:   500,
			StorageLimitMB: 500,
			IconUrl:        "",
			Price:          8000000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 4000000,
			UserID:          2,
		},

		{
			// ID: 7
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_1024,onboarding,domain_customization,logo_customization",
			DayLength:      90,
			ProductLimit:   500,
			StorageLimitMB: 1024,
			IconUrl:        "",
			Price:          8000000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 4000000,
			UserID:          2,
		},

		{
			// ID: 8
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_1024,onboarding,domain_customization,logo_customization",
			DayLength:      180,
			ProductLimit:   500,
			StorageLimitMB: 1024,
			IconUrl:        "",
			Price:          11000000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 5500000,
			UserID:          2,
		},

		{
			// ID: 9
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_1024,onboarding,domain_customization,logo_customization",
			DayLength:      365,
			ProductLimit:   500,
			StorageLimitMB: 1024,
			IconUrl:        "",
			Price:          18000000,
			Categories: []*models.Category{
				{ID: 14},
			},
			DiscountedPrice: 9000000,
			UserID:          2,
		},

		{
			// ID: 10
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_1_5g",
			DayLength:      90,
			ProductLimit:   50,
			StorageLimitMB: 750,
			IconUrl:        "www.armogroup.com/assets/images/2",
			Price:          1200000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 600000,
			UserID:          2,
		},

		{
			// ID: 11
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_1_5g",
			DayLength:      180,
			ProductLimit:   50,
			StorageLimitMB: 750,
			IconUrl:        "",
			Price:          2000000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 1000000,
			UserID:          2,
		},

		{
			// ID: 12
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_1_5g",
			DayLength:      365,
			ProductLimit:   50,
			StorageLimitMB: 750,
			IconUrl:        "",
			Price:          3000000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 1500000,
			UserID:          2,
		},

		{
			// ID: 13
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_7_5g,onboarding",
			DayLength:      90,
			ProductLimit:   250,
			StorageLimitMB: 3750,
			IconUrl:        "",
			Price:          3000000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 1500000,
			UserID:          2,
		},

		{
			// ID: 14
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_7_5g,onboarding",
			DayLength:      180,
			ProductLimit:   250,
			StorageLimitMB: 3750,
			IconUrl:        "",
			Price:          5000000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 2500000,
			UserID:          2,
		},

		{
			// ID: 15
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_7_5g,onboarding",
			DayLength:      365,
			ProductLimit:   250,
			StorageLimitMB: 3750,
			IconUrl:        "",
			Price:          8000000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 4000000,
			UserID:          2,
		},

		{
			// ID: 16
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_10g,onboarding,domain_customization,logo_customization",
			DayLength:      90,
			ProductLimit:   250,
			StorageLimitMB: 5000,
			IconUrl:        "",
			Price:          5000000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 2500000,
			UserID:          2,
		},

		{
			// ID: 17
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_10g,onboarding,domain_customization,logo_customization",
			DayLength:      180,
			ProductLimit:   250,
			StorageLimitMB: 5000,
			IconUrl:        "",
			Price:          8000000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 4000000,
			UserID:          2,
		},

		{
			// ID: 18
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_10g,onboarding,domain_customization,logo_customization",
			DayLength:      365,
			ProductLimit:   250,
			StorageLimitMB: 5000,
			IconUrl:        "",
			Price:          12000000,
			Categories: []*models.Category{
				{ID: 1},
				{ID: 21},
				{ID: 23},
				{ID: 25},
				{ID: 29},
			},
			DiscountedPrice: 6000000,
			UserID:          2,
		},

		{
			// ID: 19
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_750m",
			DayLength:      90,
			ProductLimit:   50,
			StorageLimitMB: 500,
			IconUrl:        "www.armogroup.com/assets/images/2",
			Price:          1200000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 600000,
			UserID:          2,
		},

		{
			// ID: 20
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_750m",
			DayLength:      180,
			ProductLimit:   50,
			StorageLimitMB: 500,
			IconUrl:        "",
			Price:          2000000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 1000000,
			UserID:          2,
		},

		{
			// ID: 21
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_750m",
			DayLength:      365,
			ProductLimit:   50,
			StorageLimitMB: 500,
			IconUrl:        "",
			Price:          3000000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 1500000,
			UserID:          2,
		},

		{
			// ID: 22
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_3_67g,onboarding",
			DayLength:      90,
			ProductLimit:   250,
			StorageLimitMB: 2500,
			IconUrl:        "",
			Price:          3000000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 1500000,
			UserID:          2,
		},

		{
			// ID: 23
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_3_67g,onboarding",
			DayLength:      180,
			ProductLimit:   250,
			StorageLimitMB: 2500,
			IconUrl:        "",
			Price:          5000000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 2500000,
			UserID:          2,
		},

		{
			// ID: 24
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_3_67g,onboarding",
			DayLength:      365,
			ProductLimit:   250,
			StorageLimitMB: 2500,
			IconUrl:        "",
			Price:          8000000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 4000000,
			UserID:          2,
		},

		{
			// ID: 25
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_5g,onboarding,domain_customization,logo_customization",
			DayLength:      90,
			ProductLimit:   250,
			StorageLimitMB: 3000,
			IconUrl:        "",
			Price:          5000000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 2500000,
			UserID:          2,
		},

		{
			// ID: 26
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_5g,onboarding,domain_customization,logo_customization",
			DayLength:      180,
			ProductLimit:   250,
			StorageLimitMB: 3000,
			IconUrl:        "",
			Price:          8000000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 4000000,
			UserID:          2,
		},

		{
			// ID: 27
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_5g,onboarding,domain_customization,logo_customization",
			DayLength:      365,
			ProductLimit:   250,
			StorageLimitMB: 3000,
			IconUrl:        "",
			Price:          12000000,
			Categories: []*models.Category{
				{ID: 5},
			},
			DiscountedPrice: 6000000,
			UserID:          2,
		},

		{
			// ID: 28
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_700m,upload_file,insert_location,insert_website_contact_social",
			DayLength:      90,
			ProductLimit:   50,
			StorageLimitMB: 500,
			IconUrl:        "www.armogroup.com/assets/images/2",
			Price:          1200000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 600000,
			UserID:          2,
		},

		{
			// ID: 29
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_700m,upload_file,insert_location,insert_website_contact_social",
			DayLength:      180,
			ProductLimit:   50,
			StorageLimitMB: 500,
			IconUrl:        "",
			Price:          2000000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 1000000,
			UserID:          2,
		},

		{
			// ID: 30
			Title:          "starter",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_700m,upload_file,insert_location,insert_website_contact_social",
			DayLength:      365,
			ProductLimit:   50,
			StorageLimitMB: 500,
			IconUrl:        "",
			Price:          3000000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 1500000,
			UserID:          2,
		},

		{
			// ID: 31
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_3_42g,onboarding,upload_file,insert_location,insert_website_contact_social",
			DayLength:      90,
			ProductLimit:   250,
			StorageLimitMB: 2500,
			IconUrl:        "",
			Price:          3000000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 1500000,
			UserID:          2,
		},

		{
			// ID: 32
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_3_42g,onboarding,upload_file,insert_location,insert_website_contact_social",
			DayLength:      180,
			ProductLimit:   250,
			StorageLimitMB: 2500,
			IconUrl:        "",
			Price:          5000000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 2500000,
			UserID:          2,
		},

		{
			// ID: 33
			Title:          "pro",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_3_42g,onboarding,upload_file,insert_location,insert_website_contact_social",
			DayLength:      365,
			ProductLimit:   250,
			StorageLimitMB: 2500,
			IconUrl:        "",
			Price:          8000000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 4000000,
			UserID:          2,
		},

		{
			// ID: 34
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_5g,onboarding,domain_customization,logo_customization,upload_file,insert_location,insert_website_contact_social",
			DayLength:      90,
			ProductLimit:   250,
			StorageLimitMB: 3000,
			IconUrl:        "",
			Price:          5000000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 2500000,
			UserID:          2,
		},

		{
			// ID: 35
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_5g,onboarding,domain_customization,logo_customization,upload_file,insert_location,insert_website_contact_social",
			DayLength:      180,
			ProductLimit:   250,
			StorageLimitMB: 3000,
			IconUrl:        "",
			Price:          8000000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 4000000,
			UserID:          2,
		},

		{
			// ID: 36
			Title:          "premium",
			Description:    "unlimited_views,link_qrcode,android_ios_support,user_analytics,hosting_5g,onboarding,domain_customization,logo_customization,upload_file,insert_location,insert_website_contact_social",
			DayLength:      365,
			ProductLimit:   250,
			StorageLimitMB: 3000,
			IconUrl:        "",
			Price:          12000000,
			Categories: []*models.Category{
				{ID: 12},
			},
			DiscountedPrice: 6000000,
			UserID:          2,
		},
	}

	err := db.Create(&planSeeders).Error
	if err != nil {
		log.Fatal(err)
	}
}
