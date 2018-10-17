package database

import "github.com/deFarro/letsdoit_backend/app/data"

// GetUserByID searches for user by ID
func (db *Database) GetUserByID(id string) (data.User, error) {
	db.log("user", id, "select")

	user := data.User{ ID: id }
	err := db.DB.Select(&user)
	if err != nil {
		return data.User{}, err
	}

	return user, nil
}

// GetUserBySessionID searches for user with session ID
func (db *Database) GetUserBySessionID(id string) (data.User, error) {
	session, err := db.GetSessionByID(id)
	if err != nil {
		return data.User{}, err
	}

	user, err := db.GetUserByID(session.UserID)
	if err != nil {
		return data.User{}, err
	}

	return user, nil
}

// GetTodoByID searches for todo by ID
func (db *Database) GetTodoByID(id string) (data.Todo, error) {
	db.log("todo", id, "select")

	todo := data.Todo{ ID: id }
	err := db.DB.Select(&todo)
	if err != nil {
		return data.Todo{}, err
	}

	return todo, nil
}

// SelectAllTodos selects all todos from db
func (db *Database) SelectAllTodos() (data.Todos, error) {
	db.log("todos", "all", "select")

	var todos data.Todos
	err := db.DB.Model(&todos).Select()

	return todos, err
}

// InsertTodo inserts todo to db
func (db *Database) InsertTodo(todo data.Todo) error {
	db.log("todo", todo.ID, "insert")
	return db.DB.Insert(&todo)
}

// UpdateTodo updates todo to db
func (db *Database) UpdateTodo(todo data.Todo) error {
	db.log("todo", todo.ID, "update")
	return db.DB.Update(&todo)
}

// DeleteTodo delets todo from db
func (db *Database) DeleteTodo(todo data.Todo) error {
	db.log("todo", todo.ID, "delete")

	err := db.DB.Delete(&todo)
	if err != nil {
		return err
	}

	return nil
}

// GetSessionByID searches for session by ID
func (db *Database) GetSessionByID(id string) (Session, error) {
	db.log("session", id, "select")

	session := Session{ ID: id }
	err := db.DB.Select(&session)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

// InsertSession inserts session to db
func (db *Database) InsertSession(session Session) error {
	db.log("session", session.ID, "insert")
	return db.DB.Insert(&session)
}

// IsAllowedToEdit checks if content is kept. If so anyone can change todo's status
func IsAllowedToEdit(todo1, todo2 data.Todo, user data.User) bool {
	contendKept := todo1.Title == todo2.Title && todo1.Description == todo2.Description

	return contendKept || todo1.Author.ID == user.ID
}
