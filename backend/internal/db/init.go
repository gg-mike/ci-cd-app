package db

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(dbUrl string, cfg gorm.Config, migrate bool) error {
	var err error
	db, err = gorm.Open(postgres.Open(dbUrl), &cfg)
	if err != nil {
		return err
	}

	if !migrate {
		return nil
	}

	if db.Migrator().HasTable(&model.Secret{}) {
		db.Migrator().DropColumn(&model.Secret{}, "unique")
		db.Migrator().DropIndex(&model.Secret{}, "idx_secrets")
	}

	if db.Migrator().HasTable(&model.Variable{}) {
		db.Migrator().DropColumn(&model.Variable{}, "unique")
		db.Migrator().DropIndex(&model.Variable{}, "idx_variables")
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
