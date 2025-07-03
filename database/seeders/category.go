package seeders

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/dtp"
	"gorm.io/gorm"
)

func CategorySeeder(db *gorm.DB) {
	categorySeeders := []models.Category{
		{
			//ID: 1
			ParentID: dtp.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Title:            "home_kitchen",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/placeholder.svg",
			AcceptedFileType: "",
			ARPlacement: dtp.NullString{
				String: "",
				Valid:  false,
			},
			URL: dtp.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			//ID: 2
			ParentID: dtp.NullInt64{
				Int64: 1,
				Valid: true,
			},
			Title:            "furniture",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/furniture.svg",
			AcceptedFileType: "complex",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 3
			ParentID: dtp.NullInt64{
				Int64: 1,
				Valid: true,
			},
			Title:            "home_appliances",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/home_appliances.svg",
			AcceptedFileType: "complex",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 4
			ParentID: dtp.NullInt64{
				Int64: 1,
				Valid: true,
			},
			Title:            "wall_mounted",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/wall_mounted.svg",
			AcceptedFileType: "complex",
			ARPlacement: dtp.NullString{
				String: "wall",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 5
			ParentID: dtp.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Title:            "accessories",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/placeholder.svg",
			AcceptedFileType: "",
			ARPlacement: dtp.NullString{
				String: "",
				Valid:  false,
			},
			URL: dtp.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			//ID: 6
			ParentID: dtp.NullInt64{
				Int64: 5,
				Valid: true,
			},
			Title:            "watch",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/watch.svg",
			AcceptedFileType: "glb",
			ARPlacement: dtp.NullString{
				String: "hand",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://tryon.armogroup.tech/wrist",
				Valid:  true,
			},
		},
		{
			//ID: 7
			ParentID: dtp.NullInt64{
				Int64: 5,
				Valid: true,
			},
			Title:            "bracelet",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/wristband.svg",
			AcceptedFileType: "glb",
			ARPlacement: dtp.NullString{
				String: "hand",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://tryon.armogroup.tech/wrist",
				Valid:  true,
			},
		},
		{
			//ID: 8
			ParentID: dtp.NullInt64{
				Int64: 5,
				Valid: true,
			},
			Title:            "earrings",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/earrings.svg",
			AcceptedFileType: "image",
			ARPlacement: dtp.NullString{
				String: "face",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://reality.armogroup.tech/tryon/serve",
				Valid:  true,
			},
		},
		{
			//ID: 9
			ParentID: dtp.NullInt64{
				Int64: 5,
				Valid: true,
			},
			Title:            "necklace",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/necklace.svg",
			AcceptedFileType: "image",
			ARPlacement: dtp.NullString{
				String: "face",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://reality.armogroup.tech/tryon/serve",
				Valid:  true,
			},
		},
		{
			//ID: 10
			ParentID: dtp.NullInt64{
				Int64: 5,
				Valid: true,
			},
			Title:            "ring",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/ring.svg",
			AcceptedFileType: "glb",
			ARPlacement: dtp.NullString{
				String: "hand",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://tryon.armogroup.tech/ring",
				Valid:  true,
			},
		},
		{
			//ID: 11
			ParentID: dtp.NullInt64{
				Int64: 5,
				Valid: true,
			},
			Title:            "glasses",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/glasses.svg",
			AcceptedFileType: "glb",
			ARPlacement: dtp.NullString{
				String: "face",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://tryon-glasses.armogroup.tech",
				// String: "https://tryon.armogroup.tech/glasses",

				Valid: true,
			},
		},
		{
			//ID: 12
			ParentID: dtp.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Title:            "printed_products",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/placeholder.svg",
			AcceptedFileType: "",
			ARPlacement: dtp.NullString{
				String: "",
				Valid:  false,
			},
			URL: dtp.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			//ID: 13
			ParentID: dtp.NullInt64{
				Int64: 12,
				Valid: true,
			},
			Title:            "video",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/video.svg",
			AcceptedFileType: "video",
			ARPlacement: dtp.NullString{
				String: "media",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://media.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 14
			ParentID: dtp.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Title:            "carpet_wallpaper_flooring",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/placeholder.svg",
			AcceptedFileType: "",
			ARPlacement: dtp.NullString{
				String: "",
				Valid:  false,
			},
			URL: dtp.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			//ID: 15
			ParentID: dtp.NullInt64{
				Int64: 14,
				Valid: true,
			},
			Title:            "carpet",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/carpet.svg",
			AcceptedFileType: "multi-images",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},

		{
			//ID: 16
			ParentID: dtp.NullInt64{
				Int64: 14,
				Valid: true,
			},
			Title:            "stone_ceramic_flooring",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/stone_ceramic_flooring.svg",
			AcceptedFileType: "multi-images",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 17
			ParentID: dtp.NullInt64{
				Int64: 14,
				Valid: true,
			},
			Title:            "stone_ceramic_wall_covering",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/stone_ceramic_wall_overing.svg",
			AcceptedFileType: "multi-images",
			ARPlacement: dtp.NullString{
				String: "wall",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 18
			ParentID: dtp.NullInt64{
				Int64: 14,
				Valid: true,
			},
			Title:            "frame",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/frame.svg",
			AcceptedFileType: "multi-images",
			ARPlacement: dtp.NullString{
				String: "wall",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 19
			ParentID: dtp.NullInt64{
				Int64: 14,
				Valid: true,
			},
			Title:            "curtain",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/curtain.svg",
			AcceptedFileType: "multi-images",
			ARPlacement: dtp.NullString{
				String: "wall",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 20
			ParentID: dtp.NullInt64{
				Int64: 14,
				Valid: true,
			},
			Title:            "wallpaper",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/stone_ceramic_wall_overing.svg",
			AcceptedFileType: "multi-images",
			ARPlacement: dtp.NullString{
				String: "wall",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		// ⭐️ کتگوری جدید اضافه شده - رگال (برای batch processing چندتا عکس)
		{
			//ID: 34 (این ID ممکنه متفاوت باشه)
			ParentID: dtp.NullInt64{
				Int64: 14, // زیر carpet_wallpaper_flooring
				Valid: true,
			},
			Title:            "regal", // نام درست برای رگال/قفسه
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/regal.svg", // فعلاً همین، بعداً آیکون رگال
			AcceptedFileType: "multi-images",                                               // برای batch processing
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech/showroom", // URL صحیح برای showroom
				Valid:  true,
			},
		},
		{
			//ID: 21
			ParentID: dtp.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Title:            "industrial_products",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/placeholder.svg",
			AcceptedFileType: "",
			ARPlacement: dtp.NullString{
				String: "",
				Valid:  false,
			},
			URL: dtp.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			//ID: 22
			ParentID: dtp.NullInt64{
				Int64: 21,
				Valid: true,
			},
			Title:            "industrial_products",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/industrial_products.svg",
			AcceptedFileType: "complex",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 23
			ParentID: dtp.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Title:            "food_industry",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/placeholder.svg",
			AcceptedFileType: "",
			ARPlacement: dtp.NullString{
				String: "",
				Valid:  false,
			},
			URL: dtp.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			//ID: 24
			ParentID: dtp.NullInt64{
				Int64: 23,
				Valid: true,
			},
			Title:            "food_industry",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/food_industry.svg",
			AcceptedFileType: "complex",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 25
			ParentID: dtp.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Title:            "clothing_fashion",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/placeholder.svg",
			AcceptedFileType: "",
			ARPlacement: dtp.NullString{
				String: "",
				Valid:  false,
			},
			URL: dtp.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			//ID: 26
			ParentID: dtp.NullInt64{
				Int64: 25,
				Valid: true,
			},
			Title:            "shoes",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/shoes.svg",
			AcceptedFileType: "glb",
			ARPlacement: dtp.NullString{
				String: "foot",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://reality.armogroup.tech/tryon/serve",
				Valid:  true,
			},
		},
		{
			//ID: 27
			ParentID: dtp.NullInt64{
				Int64: 25,
				Valid: true,
			},
			Title:            "fashion",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/outfit.svg",
			AcceptedFileType: "complex",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 28
			ParentID: dtp.NullInt64{
				Int64: 25,
				Valid: true,
			},
			Title:            "hat",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/hat.svg",
			AcceptedFileType: "glb",
			ARPlacement: dtp.NullString{
				String: "face",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://tryon.armogroup.tech/hat",
				Valid:  true,
			},
		},
		{
			//ID: 29
			ParentID: dtp.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Title:            "showroom",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/placeholder.svg",
			AcceptedFileType: "",
			ARPlacement: dtp.NullString{
				String: "",
				Valid:  false,
			},
			URL: dtp.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			//ID: 30
			ParentID: dtp.NullInt64{
				Int64: 29,
				Valid: true,
			},
			Title:            "gamification",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/gamification.svg",
			AcceptedFileType: "complex",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech",
				Valid:  true,
			},
		},
		{
			//ID: 31
			ParentID: dtp.NullInt64{
				Int64: 5,
				Valid: true,
			},
			Title:            "bow_tie",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/bow_tie.svg",
			AcceptedFileType: "image",
			ARPlacement: dtp.NullString{
				String: "face",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://reality.armogroup.tech/tryon/serve",
				Valid:  true,
			},
		},
		{
			//ID: 32
			ParentID: dtp.NullInt64{
				Int64: 29,
				Valid: true,
			},
			Title:            "ar_3d",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/AR%2B3D.svg",
			AcceptedFileType: "complex",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech/showroom",
				Valid:  true,
			},
		},
		{
			//ID: 33
			ParentID: dtp.NullInt64{
				Int64: 29,
				Valid: true,
			},
			Title:            "3d",
			UserID:           2,
			IconUrl:          "https://armogroup.storage.iran.liara.space/icons/3D.svg",
			AcceptedFileType: "glb",
			ARPlacement: dtp.NullString{
				String: "floor",
				Valid:  true,
			},
			URL: dtp.NullString{
				String: "https://webar.armogroup.tech/showroom",
				Valid:  true,
			},
		},
	}

	err := db.Create(&categorySeeders).Error
	if err != nil {
		log.Fatal(err)
	}
}
