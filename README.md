# Go Simple CRUD Application with Gin Framework

This is a simple CRUD (Create, Read, Update, Delete) application built in Go using the [Gin Web Framework](https://github.com/gin-gonic/gin). The application manages a list of `todos`, allowing users to create, read, update, and delete `todo` items.

## Features

- **Create a Todo**: Add a new `todo` item with a ID, description, and completion status.
- **Read Todos**: Retrieve the list of all `todos` or get details for a specific `todo`.
- **Update a Todo**: Edit an existing `todo` by updating its description and/or completion status.
- **Delete a Todo**: Remove a `todo` item from the list.

## Endpoints

- `GET /todos` - Retrieves the full list of todos.
- `POST /todos` - Creates a new todo item.
- `GET /todos/:id` - Retrieves details of a specific todo by ID.
- `PATCH /todos/:id` - Toggles the completion status of a todo.
- `PUT /todos/:id` - Updates the description and/or completion status of a specific todo.
- `DELETE /todos/:id` - Deletes a specific todo by ID.

## Quick Start

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or higher)
- [Gin Web Framework](https://github.com/gin-gonic/gin)

### Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/aleksandr-slobodian/go-simple-crud
   cd go-simple-crud
   ```

2. **Install dependencies**:

   ```bash
   go mod download
   ```

3. **Run the application**:

   with go:

   ```bash
   go run main.go
   ```

   with [Air - Live reload](https://github.com/air-verse/air):

   ```bash
   air
   ```

### Usage

1. **Create a new todo**:

   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"item": "Buy groceries", "completed": false}' http://localhost:9090/todos
   ```

2. **Retrieve all todos**:

   ```bash
   curl http://localhost:9090/todos
   ```

3. **Retrieve a specific todo**:

   ```bash
   curl http://localhost:9090/todos/1
   ```

4. **Update a todo**:

   ```bash
   curl -X PUT -H "Content-Type: application/json" -d '{"item": "Buy groceries", "completed": true}' http://localhost:9090/todos/1
   ```

5. **Update todo's completed status**:

   ```bash
   curl -X PUTCH http://localhost:9090/todos/1
   ```

6. **Delete a todo**:

   ```bash
   curl -X DELETE http://localhost:9090/todos/1
   ```

## License

This project is licensed under the MIT License.
