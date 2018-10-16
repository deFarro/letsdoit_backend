package router

import (
	"github.com/go-pg/pg"
	"github.com/deFarro/letsdoit_backend/app/config"
	"github.com/deFarro/letsdoit_backend/app/database"
)

type Router struct {
	Settings config.Config
	Database *pg.DB
}

// NewRouter creates new router instance
func NewRouter(settings config.Config) (Router, error) {
	db := pg.Connect(&pg.Options{
		Database: settings.DatabaseName,
		User: settings.DatabaseUser,
	})

	err := database.PrepopulateDatabase(db)
	//err := database.DropTables(db)
	if err != nil {
		return Router{}, err
	}

	return Router{
		Settings: settings,
		Database: db,
	}, nil
}