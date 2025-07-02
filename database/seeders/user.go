package seeders

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/faker"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) {
	users := []models.User{
		{
			Name:            "آرزو",
			Phone:           "09125793220",
			Grade:           1000,
			LastName:        "زمانی",
			Nickname:        "arezoo",
			AvatarUrl:       "",
			Password:        "$2a$12$qFAj39gAkQCIj9wLnG4qDe15G9QLD03fqYfjucLjizr584LL3FGDq", // 9ol.)P:?
			CountryCode:     "IR",
			PhoneVerifiedAt: faker.SQLNow(),
			UID:             uuid.New(),
			IDCode: dtp.NullString{
				String: "09125793220",
				Valid:  true,
			},
			Username: dtp.NullString{
				String: "09125793220",
				Valid:  true,
			},
			Email: dtp.NullString{
				String: "zamanyarezoo@gmail.com",
				Valid:  true,
			},
			EmailVerifiedAt:    faker.SQLNow(),
			ProfileCompletedAt: faker.SQLNow(),
			MadeOfficialAt:     faker.SQLNow(),
			SuspendedAt:        dtp.NullTime{},
			CompanyName:        dtp.NullString{Valid: true, String: "ArezooSoft"},
			Invite: &models.Invite{
				Code:  dtp.NullString{Valid: true, String: "zich"},
				Limit: -1,
			},
			Roles: []*models.Role{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
		},
		{
			UID:             uuid.New(),
			Name:            "آرش",
			Phone:           "09354219008",
			LastName:        "زمانی",
			CountryCode:     "IR",
			Password:        "$2a$12$OZUUKc.QoPjGyITeU76m6uoCAzrrv2uy1ZiJlGwj.uXgVT0yU3UbG", // 9ol.)P:?
			PhoneVerifiedAt: faker.SQLNow(),
			Username: dtp.NullString{
				String: "09354219008",
				Valid:  true,
			},
			Email: dtp.NullString{
				String: "arashzich1992@gmail.com",
				Valid:  true,
			},
			EmailVerifiedAt: faker.SQLNow(),
			CompanyName:     dtp.NullString{Valid: true, String: "ZichSoft"},
			Roles: []*models.Role{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
		},
	}

	err := db.Create(&users).Error
	if err != nil {
		log.Fatal(err)
	}
}
