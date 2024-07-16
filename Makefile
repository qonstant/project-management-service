postgres:
	docker run --name management-db -p 5432:5432 -e POSTGRES_USER=management-user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=management-db -d postgres:latest

createdb:
	docker exec -it management-db createdb --username=management-user --owner=management-user management-db

dropdb:
	docker exec -it management-db dropdb --username=management-user management-db

migrateup:
	migrate -path db/migrations -database "postgresql://management-user:password@localhost:5432/management-db?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://management-user:password@localhost:5432/management-db?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://management-user:password@localhost:5432/management-db?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://management-user:password@localhost:5432/management-db?sslmode=disable" -verbose down 1

up:
	docker compose up -d

down:
	docker compose down
	docker rmi management-service

restart: down up

sqlc:
	sqlc generate

test-html:
	@echo "Creation of UI for tests..."
	@cd db/sqlc && go test -coverprofile=cover.txt ./...
	@cd db/sqlc && go tool cover -html=cover.txt

test:
	@cd db/sqlc && go test

server:
	go run main.go

coverfile:
	go test -coverprofile=c.out
	go tool cover -html="c.out"

swagger:
	swag init -g internal/handlers/http/task.go

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc tests server mock storetest coverfile