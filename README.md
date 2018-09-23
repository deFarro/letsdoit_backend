# Let's do it back-end

Server side for Lets do it application (https://github.com/deFarro/letsdoit).

Allows user to log in/out, add new todos, edit/delete existing ones (only author can), change status for any todo.

Endpoints:

1. "/user/login" GET - authenticate user by checking username and password hash (provides session ID)
2. "/user/logout" GET - drop active user session
3. "/todos" GET - request all todos in database
4. "/todo" PUT - add new todo or edit existing one
5. "/todo" DELETE - delete existing todo

---

### Tech stack:
* Golang
* Postgress
* Docker
