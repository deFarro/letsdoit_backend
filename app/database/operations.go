package database

import (
	"github.com/deFarro/letsdoit_backend/app/user"
	"github.com/deFarro/letsdoit_backend/app/todo"
	"github.com/deFarro/letsdoit_backend/app/session"
)

// GetUserByID searches for user by ID
func (db Database) GetUserByID(id string) (user.User, error) {
	db.log("user", id, "select")

	fetchedUser := user.User{ ID: id }
	err := db.DB.Select(&fetchedUser)
	if err != nil {
		return user.User{}, err
	}

	return fetchedUser, nil
}

// GetUserBySessionID searches for user with session ID
func (db Database) GetUserBySessionID(id string) (user.User, error) {
	session, err := db.GetSessionByID(id)
	if err != nil {
		return user.User{}, err
	}

	fetchedUser, err := db.GetUserByID(session.UserID)
	if err != nil {
		return user.User{}, err
	}

	return fetchedUser, nil
}

// GetTodoByID searches for todo by ID
func (db Database) GetTodoByID(id string) (todo.Todo, error) {
	db.log("todo", id, "select")

	fetchedTodo := todo.Todo{ ID: id }
	err := db.DB.Select(&fetchedTodo)
	if err != nil {
		return todo.Todo{}, err
	}

	return fetchedTodo, nil
}

// SelectAllTodos selects all todos from db
func (db Database) SelectAllTodos() (todo.Todos, error) {
	db.log("todos", "all", "select")

	var todos todo.Todos
	err := db.DB.Model(&todos).Select()

	return todos, err
}

// InsertTodo inserts todo to db
func (db Database) InsertTodo(todo todo.Todo) error {
	db.log("todo", todo.ID, "insert")
	return db.DB.Insert(&todo)
}

// UpdateTodo updates todo to db
func (db Database) UpdateTodo(todo todo.Todo) error {
	db.log("todo", todo.ID, "update")
	return db.DB.Update(&todo)
}

// DeleteTodo delets todo from db
func (db Database) DeleteTodo(todo todo.Todo) error {
	db.log("todo", todo.ID, "delete")

	err := db.DB.Delete(&todo)
	if err != nil {
		return err
	}

	return nil
}

// GetSessionByID searches for session by ID
func (db Database) GetSessionByID(id string) (session.Session, error) {
	db.log("session", id, "select")

	fetchedSession := session.Session{ ID: id }
	err := db.DB.Select(&fetchedSession)
	if err != nil {
		return session.Session{}, err
	}

	return fetchedSession, nil
}

// InsertSession inserts session to db
func (db Database) InsertSession(session session.Session) error {
	db.log("session", session.ID, "insert")
	return db.DB.Insert(&session)
}

// DropSession function to clear sessionID
func (db Database) DropSession(s session.Session) error {
	return db.DB.Delete(&s)
}
