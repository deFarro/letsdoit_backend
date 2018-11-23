package database

import (
	"github.com/go-pg/pg"
	"log"
	"github.com/go-pg/pg/orm"
	"fmt"
	"github.com/deFarro/letsdoit_backend/app/config"
	"github.com/deFarro/letsdoit_backend/app/user"
	"github.com/deFarro/letsdoit_backend/app/todo"
	"github.com/deFarro/letsdoit_backend/app/session"
)

type Database struct {
	DB pg.DB
	Settings config.Config
}

// Method for logging database operetions
func (db Database) log(itemType, id, action string) {
	fmt.Printf("%s (id: %s) -> %s\n", itemType, id, action)
}

// NewDatabase creates new database and populates it
func NewDatabase(settings config.Config) (Database, error) {
	db := pg.Connect(&pg.Options{
		Addr: settings.DatabaseAddr,
		Database: settings.DatabaseName,
		User: settings.DatabaseUser,
		Password: settings.DatabasePassword,
	})

	database := Database{
		DB: *db,
		Settings: settings,
	}

	err := database.DropTables()
	if err != nil {
		return Database{}, err
	}

	err = database.PrepopulateDatabase()
	if err != nil {
		return Database{}, err
	}

	return database, nil
}

// PrepopulateDatabase populates database with users and todos if it's empty
func (db *Database) PrepopulateDatabase() error {
	// populate users
	err := db.DB.CreateTable(&user.User{}, &orm.CreateTableOptions{})
	if err == nil {
		err := db.DB.Insert(&initialUsers)
		if err != nil {
			return err
		}
		fmt.Println("Database is populated with users")
	} else {
		log.Println(err)
	}

	// populate todos
	err = db.DB.CreateTable(&todo.Todo{}, &orm.CreateTableOptions{})
	if err == nil {
		err := db.DB.Insert(&initialTodos)
		if err != nil {
			return err
		}
		fmt.Println("Database is populated with todos")
	} else {
		log.Println(err)
	}

	// create table for sessions
	err = db.DB.CreateTable(&session.Session{}, &orm.CreateTableOptions{})
	if err == nil {
		fmt.Println("Sessions table was created")
	} else {
		log.Println(err)
	}

	return nil
}

// DropTables deletes all tables
func (db *Database) DropTables() error {
	for _, model := range []interface{}{&user.User{}, &todo.Todos{}, &session.Session{}} {
		err := db.DB.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return err
		}
	}
	fmt.Println("Tables were dropped")

	return nil
}
