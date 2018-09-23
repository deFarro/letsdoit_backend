package database

import (
	"errors"

	"github.com/deFarro/letsdoit_backend/app/data"
)

var mockUser = data.User{ID: "123", Username: "tom"}

var dbUsers = []data.User{
	data.User{
		ID:           "123",
		Username:     "tom",
		PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
	},
}

var dbTodos = data.Todos{
	data.Todo{
		Title:       "Todo 1",
		Description: "Do something",
		Status:      "notCompleted",
		ID:          "1",
		Author:      mockUser,
	},
	data.Todo{
		Title:       "Todo 2",
		Description: "Do another something",
		Status:      "notCompleted",
		ID:          "2",
		Author:      mockUser,
	},
	data.Todo{
		Title:       "Todo 3",
		Description: "Do something more",
		Status:      "completed",
		ID:          "3",
		Author:      mockUser,
	},
}

// FetchTodos mocks database request
func FetchTodos() data.Todos {
	return dbTodos
}

// FetchUser mocks database request
func FetchUser(user data.User) (data.User, error) {
	for i, dbUser := range dbUsers {
		if user.Username == dbUser.Username && user.PasswordHash == dbUser.PasswordHash {
			dbUsers[i].SessionID = dbUser.GenerateSessionID()

			return dbUsers[i], nil
		}
	}

	return data.User{}, errors.New("user not found")
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
