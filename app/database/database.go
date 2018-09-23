package database

import (
	"errors"

	"github.com/deFarro/letsdoit_backend/app/data"
)

var dbUsers = []data.User{
	data.User{
		ID:           "123",
		Username:     "tom",
		PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
	},
}

// FetchTodos mocks database request
func FetchTodos() data.Todos {
	user := data.User{ID: "123", Username: "tom"}

	return data.Todos{
		data.Todo{
			Title:       "Todo 1",
			Description: "Do something",
			Status:      "notCompleted",
			ID:          "1",
			Author:      user,
		},
		data.Todo{
			Title:       "Todo 2",
			Description: "Do another something",
			Status:      "notCompleted",
			ID:          "2",
			Author:      user,
		},
		data.Todo{
			Title:       "Todo 3",
			Description: "Do something more",
			Status:      "completed",
			ID:          "3",
			Author:      user,
		},
	}
}

// FetchUser mocks database request
func FetchUser(user data.User) (data.User, error) {
	for i, dbUser := range dbUsers {
		if user.Username == dbUser.Username && user.PasswordHash == dbUser.PasswordHash {
			dbUsers[i].SessionID = dbUser.GenerateSessionID()

			return dbUsers[i], nil
		}
	}

	return data.User{}, errors.New("User not found")
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
