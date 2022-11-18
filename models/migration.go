package models

import "gorm.io/gorm"

type Migration interface {
	Migrate() error
}

type gormMigration struct {
	DB *gorm.DB
}

func (g gormMigration) Migrate() error {
	// err = g.DB.Migrator().AutoMigrate(&PasswordRestore{})
	// if err != nil {
	// 	return err
	// }

	return g.DB.Migrator().AutoMigrate(&User{})
}

func NewGormMigration(DB *gorm.DB) Migration {
	return &gormMigration{
		DB: DB,
	}
}
