package database

import (
	"github.com/deFarro/letsdoit_backend/app/todo"
	"github.com/deFarro/letsdoit_backend/app/user"
)

var user1 = user.User{
	ID:           "34b7da764b21d298ef307d04d8152dc5",
	Username:     "tom",
	PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
}

var user2 = user.User{
	ID:           "4ff9fc6e4e5d5f590c4f2134a8cc96d1",
	Username:     "jack",
	PasswordHash: "5f4dcc3b5aa765d61d8327deb882cf99",
}

var initialUsers = []user.User{
	user1,
	user2,
	{
		ID:           "098f6bcd4621d373cade4e832627b4f6",
		Username:     "test",
		PasswordHash: "098f6bcd4621d373cade4e832627b4f6",
	},
}

var initialTodos = []todo.Todo{
	{
		Title:       "Todo 1",
		Description: "Do something",
		Status:      "upcoming",
		ID:          "00f8a3d99253b2dc9916622fa94695081",
		Author:      user1.Public(),
	},
	{
		Title:       "Todo 2",
		Description: "Do another something",
		Status:      "upcoming",
		ID:          "3d10c8a7f51d9bf167d5ff88f2c16342",
		Author:      user1.Public(),
	},
	{
		Title:       "Todo 3",
		Description: "Do something more",
		Status:      "completed",
		ID:          "87650fde932ba43a878379b57644f006",
		Author:      user1.Public(),
	},
	{
		Title:       "Todo 4",
		Description: "Do something then",
		Status:      "inprogress",
		ID:          "f97544387f8e33da389029b9ca9f74c9",
		Author:      user2.Public(),
	},
}
