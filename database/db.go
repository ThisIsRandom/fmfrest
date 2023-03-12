package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connection struct {
	Instance *gorm.DB
}

func NewDatabaseConnection(dsn string) (*Connection, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var seeder Seeder

	db.Migrator().DropTable(
		&User{},
		&ContactInformation{},
		&Role{},
		&Profile{},
		&Image{},
		&Advertisement{},
		&Task{},
		&TaskImage{},
		&MessageStream{},
		&Message{},
	)

	db.AutoMigrate(
		&User{},
		&ContactInformation{},
		&Role{},
		&Profile{},
		&Image{},
		&Advertisement{},
		&Task{},
		&TaskImage{},
		&MessageStream{},
		&Message{},
	)

	seeder.Seed(db)

	if err != nil {
		return nil, err
	}

	return &Connection{
		Instance: db,
	}, nil
}
