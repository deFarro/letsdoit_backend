package database

type Todos []Todo

type Todo struct {
	Title string
	Description string
	ID string
	Author User
}

type User struct {
	ID string
	Username string
}

func FetchTodos() Todos {
	user := User{ ID: "123", Username: "John"}

	return Todos{
		Todo{
			Title: "Todo 1",
			Description: "Do something",
			ID: "1",
			Author: user,
		},
		Todo{
			Title: "Todo 1",
			Description: "Do something",
			ID: "1",
			Author: user,
		},
		Todo{
			Title: "Todo 1",
			Description: "Do something",
			ID: "1",
			Author: user,
		},
	}
}