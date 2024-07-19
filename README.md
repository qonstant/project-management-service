# Project Management Service

## Goal

Creating a RESTful API for a Project Management service. The service is deployed on Render, containerized using Docker, and run with Docker Compose and Makefile. 

## Features
- **Database was created by using dbdiagram.**

![alt text](https://raw.githubusercontent.com/qonstant/project-management-service/main/DBdiagram.png)

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
https://github.com/qonstant/project-management-service.git
cd project-management-service
```
## Build and Run Locally

### Build and Run the application:

```bash
make up
```

After running it, the server will start and be accessible at 

```bash
http://localhost:8080/swagger/index.html
```

Link to the deployment: 
```bash
https://project-management-service-gjpy.onrender.com/swagger/index.html
```

### Health Check

Health can by checked by [LINK](https://project-management-service-gjpy.onrender.com/status)

### Generate Swagger Documentation

```bash
make swagger
```
### Run Tests
```bash
make tests
```

## Docker
### For database container running
```bash
make postgres
```
### For database creation
```bash
make createdb
```
### For database drop
```bash
make dropdb
```
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
### All API endpoints can be accessed through swagger, but here is data for post requests

### Create a New User
- URL: http://localhost:8080/users
- URL: https://project-management-service-gjpy.onrender.com/users
- Method: POST
- Description: Create a new user(We can'c create a project without user ID)
- Request Body:
```json
{
  "email": "nazarbayev@nu.edu.kz",
  "full_name": "Nazarbayev N",
  "role": "ex P"
}
```

### Create a New Project
- URL: http://localhost:8080/projects
- URL: https://project-management-service-gjpy.onrender.com/projects
- Method: POST
- Description: Create a new project(We can'c create a task without project ID)
- Request Body:
```json
{
  "description": "Kazakhstan is the greatest country in the world, all other countries are run by little girls.",
  "end_date": "2024-12-12",
  "manager_id": 2,
  "name": "Jana Kazakhstan",
  "start_date": "2024-07-19"
}
```

### Create a New Task
- URL: http://localhost:8080/tasks
- URL: https://project-management-service-gjpy.onrender.com/tasks
- Method: POST
- Description: Create a new task. The request can be sent with or without a completion date. Depending on whether a completion date is provided, you will need to set the valid field to true or false respectively.
- Request Body:
```json
{
  "title": "Jana Kazakhstan",
  "description": "Do something",
  "priority": "high",
  "status": "new",
  "assignee_id": 14,
  "project_id": 1,
  "completion_date": {
    "time": "2024-12-12",
    "valid": true
  }
}
```
OR
```json
{
  "title": "LRT"
  "description": "Construct LRT",
  "priority": "low",
  "status": "in_progress",
  "assignee_id": 1,
  "project_id": 4,
  "completion_date": {
    "time": null,
    "valid": false
  },
}
```

- [SWAGGER](https://project-management-service-gjpy.onrender.com/swagger/index.html)

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

- Makefile: Makefile for building, running, testing, and Docker tasks.
- Dockerfile: Dockerfile for containerizing the application.
- internal/handlers: Contains the HTTP handlers for the API endpoints.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements.