package database

import (
	"errors"

	"github.com/deFarro/letsdoit_backend/app/data"
	"github.com/go-pg/pg"
	"log"
	"github.com/go-pg/pg/orm"
	"fmt"
	"crypto/md5"
)

type Session struct {
	ID string
	UserID string
}

var mockUser = data.User{
	ID:           "34b7da764b21d298ef307d04d8152dc5",
	Username:     "tom",
	PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
}

var mockUser2 = data.User{
	ID:           "4ff9fc6e4e5d5f590c4f2134a8cc96d1",
	Username:     "jack",
	PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
}

var dbUsers = []data.User{
	mockUser,
	{
		ID:           "321",
		Username:     "jack",
		PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
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
func PrepopulateDatabase(db *pg.DB) error {
	// populate users
	err := db.CreateTable(&data.User{}, &orm.CreateTableOptions{})
	if err == nil {
		err := db.Insert(&dbUsers)
		if err != nil {
			return err
		}
		log.Println("Database is populated with users")
	} else {
		log.Println(err)
	}

	// populate todos
	err = db.CreateTable(&data.Todo{}, &orm.CreateTableOptions{})
	if err == nil {
		err := db.Insert(&dbTodosSlice)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Database is populated with todos")
	} else {
		log.Println(err)
	}

	return nil
}

// DropTables deletes all tables
func DropTables(db *pg.DB) error {
	for _, model := range []interface{}{&data.User{}, &data.Todos{}} {
		err := db.DropTable(model, &orm.DropTableOptions{
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
func FetchTodos(db *pg.DB) (data.SortedTodos, error) {
	var todos data.Todos
	err := db.Model(&todos).Select()
	if err != nil {
		return data.SortedTodos{}, err
	}

	return todos.Sort(), nil
}

// FetchUser mocks database request
func FetchUser(db *pg.DB, user data.User) (data.User, error) {
	userID := fmt.Sprintf("%x", md5.Sum([]byte(user.Username)))
	currentUser := data.User{ ID: userID }
	err := db.Select(&currentUser)
	if err != nil {
		return data.User{}, errors.New("user not found")
	}

	if currentUser.PasswordHash != user.PasswordHash {
		return data.User{}, errors.New("wrong password")
	}

	return currentUser, nil
}

// DropSession function to clear sessionID
func DropSession(sessionID string) error {
	for _, dbUser := range dbUsers {
		if dbUser.SessionID == sessionID {
			dbUser.SessionID = ""

			return nil
		}
	}

	return errors.New("sessionID was not found")
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
