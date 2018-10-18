package database

import (
	"errors"

	"github.com/deFarro/letsdoit_backend/app/data"
	"github.com/go-pg/pg"
	"log"
	"github.com/go-pg/pg/orm"
	"fmt"
	"crypto/md5"
	"github.com/deFarro/letsdoit_backend/app/config"
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

type Session struct {
	ID string
	UserID string
}

// PrepopulateDatabase populates database with users and todos if it's empty
func (db *Database) PrepopulateDatabase() error {
	// populate users
	err := db.DB.CreateTable(&data.User{}, &orm.CreateTableOptions{})
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
	err = db.DB.CreateTable(&data.Todo{}, &orm.CreateTableOptions{})
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
	err = db.DB.CreateTable(&Session{}, &orm.CreateTableOptions{})
	if err == nil {
		fmt.Println("Sessions table was created")
	} else {
		log.Println(err)
	}

	return nil
}

// DropTables deletes all tables
func (db *Database) DropTables() error {
	for _, model := range []interface{}{&data.User{}, &data.Todos{}, &Session{}} {
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

// FetchTodos fetches all todos from db
func (db *Database) FetchTodos() (data.SortedTodos, error) {
	todos, err := db.SelectAllTodos()
	if err != nil {
		return data.SortedTodos{}, err
	}

	return todos.Sort(), nil
}

// FetchUser fetches user from db, checks password, create user session and returns user
func (db *Database) FetchUser(user data.User) (data.User, error) {
	userID := fmt.Sprintf("%x", md5.Sum([]byte(user.Username)))

	currentUser, err := db.GetUserByID(userID)
	if err != nil {
		return data.User{}, err
	}

	if currentUser.PasswordHash != user.PasswordHash {
		return data.User{}, errors.New("wrong password")
	}

	session := Session{
		ID: currentUser.GenerateSessionID(),
		UserID: currentUser.ID,
	}

	err = db.InsertSession(session)
	if err != nil {
		return data.User{}, errors.New("cannot save session")
	}

	currentUser.SessionID = session.ID

	return currentUser, nil
}

// DropSession function to clear sessionID
func (db *Database) DropSession(sessionID string) error {
	return db.DB.Delete(&Session{ ID: sessionID })
}

// FlushTodo deletes todo from db
func (db *Database) FlushTodo(sessionID, todoID string) error {
	user, err := db.GetUserBySessionID(sessionID)
	if err != nil {
		return err
	}

	todo, err := db.GetTodoByID(todoID)
	if err != nil {
		return err
	}

	if user.ID != todo.Author.ID {
		return errors.New("delete is forbidden")
	}

	err = db.DeleteTodo(todo)
	if err != nil {
		return err
	}

	return nil
}

// AddModifyTodo adds/updates todo to database
func (db *Database) AddModifyTodo(sessionID string, todo data.Todo) (data.Todo, error) {
	// If todo is a new one, generate ID and save it in db
	if todo.ID == "" {
		todo.ID = todo.GenerateTodoID()

		err := db.InsertTodo(todo)
		if err != nil {
			return data.Todo{}, err
		}

		return todo, nil
	}

	// If todo exists, check if current user have rights to edit and update it
	user, err := db.GetUserBySessionID(sessionID)
	if err != nil {
		return data.Todo{}, err
	}

	dbTodo, err := db.GetTodoByID(todo.ID)
	if err != nil {
		return data.Todo{}, err
	}

	if !IsAllowedToEdit(dbTodo, todo, user) {
		return data.Todo{}, errors.New("editing is forbidden")
	}

	err = db.UpdateTodo(todo)
	if err != nil {
		return data.Todo{}, err
	}

	return todo, nil
}
