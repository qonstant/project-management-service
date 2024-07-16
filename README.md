# Todo List App

## Goal

Creating a RESTful API for a Todo List microservice. The service is deployed on Render, containerized using Docker, and run with Docker Compose and Makefile. Input data validated, and the README.md file is comprehensive.

## Features
- **Database was created by using dbdiagram.**

![alt text](https://raw.githubusercontent.com/qonstant/todo-list/main/DBdiagram.png)

- **CRUD operations for managing tasks.**
  - CRUD was created using SQLC. For regeneration of CRUD:
    ```bash
    make sqlc
    ```

- **Swagger UI for API documentation.**
  - For regeneration of Swagger documentation:
    ```bash
    swag init -g internal/handlers/http/task.go
    ```

- **Docker support for containerization.**
  - For running up:
    ```bash
    make up
    ```
  - For shutting down:
    ```bash
    make down
    ```
  - For restart:
    ```bash
    make restart
    ```

- **Unit tests for validating the CRUD.**
  - For running unit tests with visualization:
    ```bash
    make test-html
    ```

## Publishing and Deployment
- Project published on GitHub.
- Project deployed on Render.
- Project containerized using Docker.
- Project runs using Docker Compose and Makefile.

## Prerequisites

- Go 1.22 or later
- Docker (for containerization)
- SQlC (for CRUD generation)
```bash
brew install sqlc
```
or
```bash
go get github.com/kyleconroy/sqlc/cmd/sqlc
```

## Getting Started

### Clone the Repository

```bash
https://github.com/qonstant/todo-list.git
cd todo-list
```
## Build and Run Locally

### Build the application:

```bash
make build
```

### Run the application:

```bash
make run
```
After running make run, the server will start and be accessible at http://localhost:8080.

Link to the deployment: https://todo-list-hl.onrender.com

### Health Check

Health can by checked by [LINK](https://todo-list-hl.onrender.com/status)

### Generate Swagger Documentation

```bash
make swagger
```
### Run Tests
```bash
make test
```
## Docker
### For docker compose up
```bash
make up
```
### For docker compose down
```bash
make down
```
### For restarting container
```bash
make restart
```

## API Endpoints

### Create a New Task
- URL: http://localhost:8080/tasks
- Method: POST
- Request Body:
```json
{
  "title": "Buy a book",
  "active_at": "2023-08-04"
}
```

### Update an Existing Task
- URL: http://localhost:8080/tasks/{ID}
- Method: PUT
- Request Body:
```json
{
  "title": "Buy a book - High Performance Applications",
  "activeAt": "2023-08-05"
}
```

### Delete a Task
- URL: http://localhost:8080/tasks/{ID}
- Method: DELETE

### Mark a Task as Done 
- URL: http://localhost:8080/tasks/{ID}/done
- Method: PUT

### List Tasks by Status 
- URL: /api/todo-list/tasks?status=active or /api/todo-list/tasks?status=done
- Method: GET

- Link: https://todo-list-hl.onrender.com/swagger/index.html

# Swagger: HTTP tutorial for beginners

1. Add comments to your API source code, See [Declarative Comments Format](#declarative-comments-format).

2. Download swag by using:
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```
To build from source you need [Go](https://golang.org/dl/) (1.17 or newer).

Or download a pre-compiled binary from the [release page](https://github.com/swaggo/swag/releases).

3. Run `swag init` in the project's root folder which contains the `main.go` file. This will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).
```sh
swag init
```

  Make sure to import the generated `docs/docs.go` so that your specific configuration gets `init`'ed. If your General API annotations do not live in `main.go`, you can let swag know with `-g` flag.
  ```sh
  swag init -g internal/handler/handler.go
  ```

4. (optional) Use `swag fmt` format the SWAG comment. (Please upgrade to the latest version)

  ```sh
  swag fmt
  ```

## Project Structure

- main.go: The main server implementation.
- Makefile: Makefile for building, running, testing, and Docker tasks.
- Dockerfile: Dockerfile for containerizing the application.
- internal/handlers/http: Contains the HTTP handlers for the API endpoints.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements.