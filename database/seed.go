package database

import "gorm.io/gorm"

type Seeder uint

var (
	role1    string
	role2    string
	role3    string
	mail     string
	mail2    string
	password string
	name     string
)

func (s *Seeder) Seed(db *gorm.DB) {
	role1 = "normal"
	role2 = "business"
	role3 = "admin"
	mail = "test@test.dk"
	password = "$2a$10$r5Dgfbo5TNX7nWurezTkwObY/8hvPJ5qfJDDIdujvcbbXRCeKh42y"
	name = "test"
	mail2 = "business@business.dk"

	roles := []Role{
		{Name: &role1},
		{Name: &role2},
		{Name: &role3},
	}

	users := []User{
		{
			Email:    &mail,
			Password: &password,
			Profile: Profile{
				Name: name,
			},
		},
		{
			Email:    &mail2,
			Password: &password,
			Profile: Profile{
				Name: name,
			},
		},
	}

	db.Create(&roles)
	db.Create(&users)
}
