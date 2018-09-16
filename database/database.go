package database

import (
	"errors"

	"github.com/letsdoit_back/app/data"
)

// FetchTodos mocks database request
func FetchTodos() data.Todos {
	user := data.User{ID: "123", Username: "John"}

	return data.Todos{
		data.Todo{
			Title:       "Todo 1",
			Description: "Do something",
			ID:          "1",
			Author:      user,
		},
		data.Todo{
			Title:       "Todo 2",
			Description: "Do another something",
			ID:          "2",
			Author:      user,
		},
		data.Todo{
			Title:       "Todo 3",
			Description: "Do something more",
			ID:          "3",
			Author:      user,
		},
	}
}

// FetchUser mocks database request
func FetchUser(user data.User) (data.User, error) {
	dbUsers := []data.User{
		data.User{
			ID:       "123",
			Username: "John",
			Password: "qwerty",
		},
	}

	for _, dbUser := range dbUsers {
		if user.Username == dbUser.Username && user.Password == dbUser.Password {
			dbUser.SessionID = dbUser.GenerateSessionID()

			return dbUser, nil
		}
	}

	return data.User{}, errors.New("User not found")
}
