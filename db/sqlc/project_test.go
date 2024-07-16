package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	// Replace q.db with mock DB for testing purposes
	queries := New(db)

	// Capture the current time for consistent use
	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 1, 0)

	// Mock expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name", "description", "start_date", "end_date", "manager_id"}).
		AddRow(1, "Test Project", "Description", startDate, endDate, 123)

	// Expectation: QueryRowContext with expected arguments
	mock.ExpectQuery("INSERT INTO projects").
		WithArgs("Test Project", "Description", startDate, endDate, int64(123)).
		WillReturnRows(rows)

	// Prepare input params
	params := CreateProjectParams{
		Name:        "Test Project",
		Description: "Description",
		StartDate:   startDate,
		EndDate:     endDate,
		ManagerID:   123,
	}

	// Call the CreateProject method
	project, err := queries.CreateProject(context.Background(), params)

	// Verify results
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, int64(1), project.ID, "Expected project ID to match")
	assert.Equal(t, "Test Project", project.Name, "Expected project name to match")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDeleteProject(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	// Replace q.db with mock DB for testing purposes
	queries := New(db)

	// Mock expected result for ExecContext
	mock.ExpectExec("DELETE FROM projects").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the DeleteProject method
	err = queries.DeleteProject(context.Background(), 1)

	// Verify results
	assert.NoError(t, err, "Expected no error")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetProject(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	// Replace q.db with mock DB for testing purposes
	queries := New(db)

	// Capture the current time for consistent use
	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 1, 0)

	// Mock expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name", "description", "start_date", "end_date", "manager_id"}).
		AddRow(1, "Test Project", "Description", startDate, endDate, 123)

	// Expectation: QueryRowContext with expected arguments
	mock.ExpectQuery("SELECT id, name, description, start_date, end_date, manager_id FROM projects").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	// Call the GetProject method
	project, err := queries.GetProject(context.Background(), 1)

	// Verify results
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, int64(1), project.ID, "Expected project ID to match")
	assert.Equal(t, "Test Project", project.Name, "Expected project name to match")
	assert.Equal(t, "Description", project.Description, "Expected project description to match")
	assert.WithinDuration(t, startDate, project.StartDate, time.Second, "Expected project start date to match")
	assert.WithinDuration(t, endDate, project.EndDate, time.Second, "Expected project end date to match")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetProjectTasks(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	// Replace q.db with mock DB for testing purposes
	queries := New(db)

	// Mock expected rows to return
	rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee_id", "project_id", "creation_date", "completion_date"}).
		AddRow(1, "Task 1", "Description 1", "high", "Pending", 123, 1, time.Now(), nil).
		AddRow(2, "Task 2", "Description 2", "low", "InProgress", 456, 1, time.Now(), time.Now())

	// Expectation: QueryContext with expected arguments
	mock.ExpectQuery("SELECT id, title, description, priority, status, assignee_id, project_id, creation_date, completion_date FROM tasks").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	// Call the GetProjectTasks method
	tasks, err := queries.GetProjectTasks(context.Background(), 1)

	// Verify results
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 2, len(tasks), "Expected two tasks")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchProjectsByManager(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	// Replace q.db with mock DB for testing purposes
	queries := New(db)

	// Mock expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name", "description", "start_date", "end_date", "manager_id"}).
		AddRow(1, "Test Project", "Description", time.Now(), time.Now().AddDate(0, 1, 0), 123)

	// Expectation: QueryContext with expected arguments
	mock.ExpectQuery(`SELECT id, name, description, start_date, end_date, manager_id FROM projects WHERE manager_id = \$1 ORDER BY start_date ASC`).
		WithArgs(int64(123)).
		WillReturnRows(rows)

	// Call the SearchProjectsByManager method
	projects, err := queries.SearchProjectsByManager(context.Background(), 123)

	// Verify results
	assert.NoError(t, err, "Expected no error")
	assert.Len(t, projects, 1, "Expected one project")
	assert.Equal(t, int64(1), projects[0].ID, "Expected project ID to match")
	assert.Equal(t, "Test Project", projects[0].Name, "Expected project name to match")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchProjectsByTitle(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	// Replace q.db with mock DB for testing purposes
	queries := New(db)

	// Capture the current time for consistent use
	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 1, 0)

	// Mock expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name", "description", "start_date", "end_date", "manager_id"}).
		AddRow(1, "Test Project", "Description", startDate, endDate, 123)

	// Expectation: QueryContext with expected arguments
	mock.ExpectQuery(`SELECT id, name, description, start_date, end_date, manager_id FROM projects WHERE name ILIKE '%' || \$1 || '%' ORDER BY start_date ASC`).
		WithArgs("Test").
		WillReturnRows(rows)

	// Call the SearchProjectsByTitle method
	projects, err := queries.SearchProjectsByTitle(context.Background(), sql.NullString{String: "Test", Valid: true})

	// Verify results
	assert.NoError(t, err, "Expected no error")
	assert.Len(t, projects, 1, "Expected one project")
	assert.Equal(t, int64(1), projects[0].ID, "Expected project ID to match")
	assert.Equal(t, "Test Project", projects[0].Name, "Expected project name to match")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdateProject(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	// Replace q.db with mock DB for testing purposes
	queries := New(db)

	// Capture the current time for consistent use
	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 1, 0)

	// Mock expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name", "description", "start_date", "end_date", "manager_id"}).
		AddRow(1, "Updated Project", "Updated Description", startDate, endDate, 123)

	// Expectation: QueryRowContext with expected arguments
	mock.ExpectQuery("UPDATE projects SET name = \\$2, description = \\$3, start_date = \\$4, end_date = \\$5, manager_id = \\$6 WHERE id = \\$1 RETURNING id, name, description, start_date, end_date, manager_id").
		WithArgs(int64(1), "Updated Project", "Updated Description", startDate, endDate, int64(123)).
		WillReturnRows(rows)

	// Prepare input params
	params := UpdateProjectParams{
		ID:          1,
		Name:        "Updated Project",
		Description: "Updated Description",
		StartDate:   startDate,
		EndDate:     endDate,
		ManagerID:   123,
	}

	// Call the UpdateProject method
	project, err := queries.UpdateProject(context.Background(), params)

	// Verify results
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, int64(1), project.ID, "Expected project ID to match")
	assert.Equal(t, "Updated Project", project.Name, "Expected project name to match")
	assert.Equal(t, "Updated Description", project.Description, "Expected project description to match")
	assert.WithinDuration(t, startDate, project.StartDate, time.Second, "Expected project start date to match")
	assert.WithinDuration(t, endDate, project.EndDate, time.Second, "Expected project end date to match")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestListProjects(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	// Replace q.db with mock DB for testing purposes
	queries := New(db)

	// Capture the current time for consistent use
	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 1, 0)

	// Mock expected rows to return
	rows := sqlmock.NewRows([]string{"id", "name", "description", "start_date", "end_date", "manager_id"}).
		AddRow(1, "Project 1", "Description 1", startDate, endDate, 123).
		AddRow(2, "Project 2", "Description 2", startDate, endDate, 456)

	// Expectation: QueryContext
	mock.ExpectQuery("SELECT id, name, description, start_date, end_date, manager_id FROM projects").
		WillReturnRows(rows)

	// Call the ListProjects method
	projects, err := queries.ListProjects(context.Background())

	// Verify results
	assert.NoError(t, err, "Expected no error")
	assert.Len(t, projects, 2, "Expected two projects")
	assert.Equal(t, int64(1), projects[0].ID, "Expected project ID to match")
	assert.Equal(t, "Project 1", projects[0].Name, "Expected project name to match")
	assert.Equal(t, "Description 1", projects[0].Description, "Expected project description to match")
	assert.WithinDuration(t, startDate, projects[0].StartDate, time.Second, "Expected project start date to match")
	assert.WithinDuration(t, endDate, projects[0].EndDate, time.Second, "Expected project end date to match")
	assert.Equal(t, int64(2), projects[1].ID, "Expected project ID to match")
	assert.Equal(t, "Project 2", projects[1].Name, "Expected project name to match")
	assert.Equal(t, "Description 2", projects[1].Description, "Expected project description to match")
	assert.WithinDuration(t, startDate, projects[1].StartDate, time.Second, "Expected project start date to match")
	assert.WithinDuration(t, endDate, projects[1].EndDate, time.Second, "Expected project end date to match")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
