package db

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(dbUrl string, migrate bool) error {
	var err error
	db, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{Logger: logger.Gorm()})
	if err != nil {
		return err
	}

	if !migrate {
		return nil
	}

	if err = db.AutoMigrate(
		&model.Worker{},
		&model.Project{},
		&model.Pipeline{},
		&model.Build{},
		&model.BuildLog{},
		&model.BuildStep{},
		&model.Secret{},
		&model.User{},
		&model.Variable{},
	); err != nil {
		return err
	}
	return nil
}

func Get() *gorm.DB {
	return db
}
