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
}

// NewDatabase creates new database
func NewDatabase(settings config.Config) (Database, error) {
	db := pg.Connect(&pg.Options{
		Database: settings.DatabaseName,
		User: settings.DatabaseUser,
	})

	database := Database{
		DB: *db,
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

var mockUser = data.User{
	ID:           "34b7da764b21d298ef307d04d8152dc5",
	Username:     "tom",
	PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
}

var dbUsers = []data.User{
	mockUser,
	{
		ID:           "4ff9fc6e4e5d5f590c4f2134a8cc96d1",
		Username:     "jack",
		PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
	},
	{
		ID:           "098f6bcd4621d373cade4e832627b4f6",
		Username:     "test",
		PasswordHash: "098f6bcd4621d373cade4e832627b4f6",
	},
}

var dbTodos = data.Todos{
	data.Todo{
		Title:       "Todo 1",
		Description: "Do something",
		Status:      "upcoming",
		ID:          "0",
		Author:      mockUser.Public(),
	},
	data.Todo{
		Title:       "Todo 2",
		Description: "Do another something",
		Status:      "upcoming",
		ID:          "1",
		Author:      mockUser.Public(),
	},
	data.Todo{
		Title:       "Todo 3",
		Description: "Do something more",
		Status:      "completed",
		ID:          "2",
		Author:      mockUser.Public(),
	},
	data.Todo{
		Title:       "Todo 4",
		Description: "Do something then",
		Status:      "inprogress",
		ID:          "3",
		Author:      mockUser.Public(),
	},
}

var dbTodosSlice = []data.Todo{
	{
		Title:       "Todo 1",
		Description: "Do something",
		Status:      "upcoming",
		ID:          "0",
		Author:      mockUser.Public(),
	},
	{
		Title:       "Todo 2",
		Description: "Do another something",
		Status:      "upcoming",
		ID:          "1",
		Author:      mockUser.Public(),
	},
	{
		Title:       "Todo 3",
		Description: "Do something more",
		Status:      "completed",
		ID:          "2",
		Author:      mockUser.Public(),
	},
	{
		Title:       "Todo 4",
		Description: "Do something then",
		Status:      "inprogress",
		ID:          "3",
		Author:      mockUser.Public(),
	},
}

// PrepopulateDatabase populates database wwith users and todos if it's empty
func (db *Database) PrepopulateDatabase() error {
	// populate users
	err := db.DB.CreateTable(&data.User{}, &orm.CreateTableOptions{})
	if err == nil {
		err := db.DB.Insert(&dbUsers)
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
		err := db.DB.Insert(&dbTodosSlice)
		if err != nil {
			return err
		}
		fmt.Println("Database is populated with todos")
	} else {
		log.Println(err)
	}

	// create table for sessions
	err = db.DB.CreateTable(&Session{}, &orm.CreateTableOptions{
		Temp: true,
	})
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

// FetchTodos mocks database request
func (db *Database) FetchTodos() (data.SortedTodos, error) {
	var todos data.Todos
	err := db.DB.Model(&todos).Select()
	if err != nil {
		return data.SortedTodos{}, err
	}

	return todos.Sort(), nil
}

// FetchUser mocks database request
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

	err = db.DB.Insert(&session)
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
func FlushTodo(sessionID, todoID string) error {
	user, err := GetUserBySessionID(sessionID)
	if err != nil {
		return err
	}

	todo, i, err := GetTodoByID(todoID)
	if err != nil {
		return err
	}

	if user.ID == todo.Author.ID {
		DeleteTodo(i)
	}

	return errors.New("delete is forbidden")
}

// UpdateTodo replaces/adds todo to database
func UpdateTodo(sessionID string, todo data.Todo) (data.Todo, error) {
	user, err := GetUserBySessionID(sessionID)
	if err != nil {
		return data.Todo{}, err
	}

	if todo.ID == "" {
		todo.ID = todo.GenerateTodoID()
		dbTodos = append(dbTodos, todo)

		return todo, nil
	}

	for i, dbTodo := range dbTodos {
		if dbTodo.ID == todo.ID && IsAllowedToEdit(dbTodo, todo, user) {
			dbTodos[i] = todo

			return todo, nil
		}
	}

	return data.Todo{}, errors.New("todo not found or editing is forbidden")
}

// GetUserBySessionID searches for user with session ID
func GetUserBySessionID(sessionID string) (data.User, error) {
	for _, dbUser := range dbUsers {
		if dbUser.SessionID == sessionID {
			return dbUser, nil
		}
	}

	return data.User{}, errors.New("user not found")
}

// GetUserBySessionID searches for user with session ID
func (db *Database) GetUserBySessionID(sessionID string) (data.User, error) {
	session := Session{ ID: sessionID }
	err := db.DB.Select(&session)
	if err != nil {
		return data.User{}, err
	}

	user := data.User{ ID: session.UserID }
	err = db.DB.Select(&user)
	if err != nil {
		return data.User{}, err
	}

	return user, nil
}

// GetUserByID searches for user by ID
func (db *Database) GetUserByID(id string) (data.User, error) {
	user := data.User{ ID: id }
	err := db.DB.Select(&user)
	if err != nil {
		return data.User{}, err
	}

	return user, nil
}

// GetTodoByID searches for todo by ID
func (db *Database) GetTodoByID(id string) (data.Todo, error) {
	todo := data.Todo{ ID: id }
	err := db.DB.Select(&todo)
	if err != nil {
		return data.Todo{}, err
	}

	return todo, nil
}

// GetTodoByID searches for todo with ID
func GetTodoByID(id string) (data.Todo, int, error) {
	for i, dbTodo := range dbTodos {
		if dbTodo.ID == id {
			return dbTodo, i, nil
		}
	}

	return data.Todo{}, -1, errors.New("todo not found")
}

// DeleteTodo removes todo from db
func DeleteTodo(i int) {
	dbTodos = append(dbTodos[:i], dbTodos[i+1:]...)
}

// IsAllowedToEdit checks if content is kept. If so anyone can change todo's status
func IsAllowedToEdit(todo1, todo2 data.Todo, user data.User) bool {
	contendKept := todo1.Title == todo2.Title && todo1.Description == todo2.Description

	return contendKept || todo1.Author.ID == user.ID
}
