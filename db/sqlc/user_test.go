package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "registration_date", "role"}).
		AddRow(1, "Test User", "test@example.com", now, "user")

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("Test User", "test@example.com", "user").
		WillReturnRows(rows)

	params := CreateUserParams{
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     "user",
	}

	user, err := queries.CreateUser(context.Background(), params)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "Test User", user.FullName)
	assert.Equal(t, "test@example.com", user.Email)
	assert.WithinDuration(t, now, user.RegistrationDate, time.Second)
	assert.Equal(t, "user", user.Role)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	mock.ExpectExec("DELETE FROM users").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = queries.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "registration_date", "role"}).
		AddRow(1, "Test User", "test@example.com", now, "user")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1 LIMIT 1").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := queries.GetUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "Test User", user.FullName)
	assert.Equal(t, "test@example.com", user.Email)
	assert.WithinDuration(t, now, user.RegistrationDate, time.Second)
	assert.Equal(t, "user", user.Role)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetUserTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	now := time.Now()
	completionDate := sql.NullTime{Time: now.AddDate(0, 0, 1), Valid: true}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee_id", "project_id", "creation_date", "completion_date"}).
		AddRow(1, "Test Task 1", "Description 1", TaskPriorityLow, TaskStatusNew, 1, 1, now, completionDate).
		AddRow(2, "Test Task 2", "Description 2", TaskPriorityMedium, TaskStatusNew, 1, 2, now, completionDate)

	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE assignee_id = \\$1 ORDER BY creation_date ASC").
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := queries.GetUserTasks(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, int64(1), tasks[0].ID)
	assert.Equal(t, "Test Task 1", tasks[0].Title)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestListUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "registration_date", "role"}).
		AddRow(1, "Alice Smith", "alice@example.com", now, "user").
		AddRow(2, "Bob Johnson", "bob@example.com", now, "admin")

	mock.ExpectQuery("SELECT (.+) FROM users ORDER BY full_name ASC").
		WillReturnRows(rows)

	users, err := queries.ListUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice Smith", users[0].FullName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchUsersByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	email := sql.NullString{String: "alice", Valid: true}
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "registration_date", "role"}).
		AddRow(1, "Alice Smith", "alice@example.com", now, "user")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE email ILIKE '%' || \\$1 || '%' ORDER BY email ASC").
		WithArgs(email).
		WillReturnRows(rows)

	users, err := queries.SearchUsersByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "Alice Smith", users[0].FullName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchUsersByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	name := sql.NullString{String: "Bob", Valid: true}
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "registration_date", "role"}).
		AddRow(2, "Bob Johnson", "bob@example.com", now, "admin")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE full_name ILIKE '%' || \\$1 || '%' ORDER BY full_name ASC").
		WithArgs(name).
		WillReturnRows(rows)

	users, err := queries.SearchUsersByName(context.Background(), name)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "Bob Johnson", users[0].FullName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	arg := UpdateUserParams{
		ID:               1,
		FullName:         "Boris Smith",
		Email:            "boris@example.com",
		Role:             "user",
		RegistrationDate: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "full_name", "email", "registration_date", "role"}).
		AddRow(arg.ID, arg.FullName, arg.Email, arg.RegistrationDate, arg.Role)

	mock.ExpectQuery("UPDATE users SET (.+) WHERE id = \\$1 RETURNING (.+)").
		WithArgs(arg.ID, arg.FullName, arg.Email, arg.Role, arg.RegistrationDate).
		WillReturnRows(rows)

	user, err := queries.UpdateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.Equal(t, arg.ID, user.ID)
	assert.Equal(t, arg.FullName, user.FullName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
