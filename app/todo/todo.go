package todo

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
	"github.com/deFarro/letsdoit_backend/app/user"
	"errors"
)

// Todo is the type for a single todo
type Todo struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ID          string     `json:"id"`
	Author      user.PublicUser `json:"author"`
}

// SortedTodos is the type for all todos
type SortedTodos struct {
	Upcoming   Todos `json:"upcoming"`
	InProgress Todos `json:"inprogress"`
	Completed  Todos `json:"completed"`
}

// Todos is the type for a list of todos
type Todos []Todo

//
type TodoTransporter interface {
	GetUserByID(string) (user.User, error)
	GetUserBySessionID(string) (user.User, error)
	GetTodoByID(id string) (Todo, error)
	SelectAllTodos() (Todos, error)
	InsertTodo(Todo) error
	UpdateTodo(Todo) error
	DeleteTodo(Todo) error
}

// Sort method distribute todos to buckets based on status
func (tds Todos) Sort() SortedTodos {
	result := SortedTodos{
		Upcoming: []Todo{},
		InProgress: []Todo{},
		Completed: []Todo{},
	}

	for _, todo := range tds {
		switch todo.Status {
		case "upcoming":
			result.Upcoming = append(result.Upcoming, todo)

		case "inprogress":
			result.InProgress = append(result.InProgress, todo)

		case "completed":
			result.Completed = append(result.Completed, todo)
		}
	}

	return result
}

// GenerateTodoID generates unique ID for a todo
func (t Todo) GenerateTodoID() string {
	salt := strconv.FormatInt(time.Now().Unix(), 10)

	return fmt.Sprintf("%x", md5.Sum([]byte(t.Title+t.Description+salt)))
}

// IsAllowedToEdit checks if content is kept. If so anyone can change todo's status
func (todo1 Todo) IsAllowedToEdit(todo2 Todo, user user.User) bool {
	contendKept := todo1.Title == todo2.Title && todo1.Description == todo2.Description

	return contendKept || todo1.Author.ID == user.ID
}

// Fetch fetches all todos from db
func (todo Todo) Fetch(tr TodoTransporter) (SortedTodos, error) {
	todos, err := tr.SelectAllTodos()
	if err != nil {
		return SortedTodos{}, err
	}

	return todos.Sort(), nil
}

// AddModifyTodo adds/updates todo to database
func (todo Todo) AddModify(tr TodoTransporter, sessionID string) (Todo, error) {
	// If todo is a new one, generate ID and save it in db
	if todo.ID == "" {
		todo.ID = todo.GenerateTodoID()

		err := tr.InsertTodo(todo)
		if err != nil {
			return Todo{}, err
		}

		return todo, nil
	}

	// If todo exists, check if current user have rights to edit and update it
	user, err := tr.GetUserBySessionID(sessionID)
	if err != nil {
		return Todo{}, err
	}

	fetchedTodo, err := tr.GetTodoByID(todo.ID)
	if err != nil {
		return Todo{}, err
	}

	if !fetchedTodo.IsAllowedToEdit(todo, user) {
		return Todo{}, errors.New("editing is forbidden")
	}

	err = tr.UpdateTodo(todo)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

// FlushTodo deletes todo from db
func (todo *Todo) Flush(tr TodoTransporter, sessionID string) error {
	user, err := tr.GetUserBySessionID(sessionID)
	if err != nil {
		return err
	}

	fetchedTodo, err := tr.GetTodoByID(todo.ID)
	if err != nil {
		return err
	}

	if user.ID != fetchedTodo.Author.ID {
		return errors.New("delete is forbidden")
	}

	err = tr.DeleteTodo(fetchedTodo)
	if err != nil {
		return err
	}

	return nil
}
