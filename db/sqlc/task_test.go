package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	now := time.Now()
	completionDate := sql.NullTime{Time: now.Add(48 * time.Hour), Valid: true}

	// Define expected rows
	rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee_id", "project_id", "creation_date", "completion_date"}).
		AddRow(1, "Sample Task", "This is a sample task", "medium", "Pending", 1, 1, now, completionDate.Time)

	// Mock the query
	mock.ExpectQuery("INSERT INTO tasks").
		WithArgs("Sample Task", "This is a sample task", "medium", "Pending", 1, 1, completionDate).
		WillReturnRows(rows)

	// Define the parameters for CreateTask
	params := CreateTaskParams{
		Title:          "Sample Task",
		Description:    "This is a sample task",
		Priority:       "medium", // Adjusted to match the string representation of TaskPriority
		Status:         "Pending",
		AssigneeID:     1,
		ProjectID:      1,
		CompletionDate: completionDate,
	}

	// Call the CreateTask method
	task, err := queries.CreateTask(context.Background(), params)

	// Perform assertions
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, int64(1), task.ID)
	assert.Equal(t, "Sample Task", task.Title)
	assert.Equal(t, "This is a sample task", task.Description)
	assert.Equal(t, TaskPriority("medium"), task.Priority) // Convert string to TaskPriority
	assert.Equal(t, TaskStatus("Pending"), task.Status)    // Convert string to TaskStatus
	assert.Equal(t, int64(1), task.AssigneeID)
	assert.Equal(t, int64(1), task.ProjectID)
	assert.WithinDuration(t, now, task.CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, task.CompletionDate.Time, time.Second)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}


func TestGetTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	now := time.Now()
	completionDate := sql.NullTime{Time: now.AddDate(0, 0, 1), Valid: true}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee_id", "project_id", "creation_date", "completion_date"}).
		AddRow(1, "Test Task", "Description", TaskPriorityLow, TaskStatusNew, 1, 1, now, completionDate)

	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = \\$1 LIMIT 1").
		WithArgs(1).
		WillReturnRows(rows)

	task, err := queries.GetTask(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "Description", task.Description)
	assert.Equal(t, TaskPriorityLow, task.Priority)
	assert.Equal(t, TaskStatusNew, task.Status)
	assert.Equal(t, int64(1), task.AssigneeID)
	assert.Equal(t, int64(1), task.ProjectID)
	assert.WithinDuration(t, now, task.CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, task.CompletionDate.Time, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestListTasks(t *testing.T) {
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
		AddRow(2, "Test Task 2", "Description 2", TaskPriorityMedium, TaskStatusInProgress, 2, 2, now, completionDate)

	mock.ExpectQuery("SELECT (.+) FROM tasks ORDER BY creation_date ASC").
		WillReturnRows(rows)

	tasks, err := queries.ListTasks(context.Background())

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, int64(1), tasks[0].ID)
	assert.Equal(t, "Test Task 1", tasks[0].Title)
	assert.Equal(t, "Description 1", tasks[0].Description)
	assert.Equal(t, TaskPriorityLow, tasks[0].Priority)
	assert.Equal(t, TaskStatusNew, tasks[0].Status)
	assert.Equal(t, int64(1), tasks[0].AssigneeID)
	assert.Equal(t, int64(1), tasks[0].ProjectID)
	assert.WithinDuration(t, now, tasks[0].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[0].CompletionDate.Time, time.Second)
	assert.Equal(t, int64(2), tasks[1].ID)
	assert.Equal(t, "Test Task 2", tasks[1].Title)
	assert.Equal(t, "Description 2", tasks[1].Description)
	assert.Equal(t, TaskPriorityMedium, tasks[1].Priority)
	assert.Equal(t, TaskStatusInProgress, tasks[1].Status)
	assert.Equal(t, int64(2), tasks[1].AssigneeID)
	assert.Equal(t, int64(2), tasks[1].ProjectID)
	assert.WithinDuration(t, now, tasks[1].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[1].CompletionDate.Time, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchTasksByAssignee(t *testing.T) {
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
		AddRow(2, "Test Task 2", "Description 2", TaskPriorityMedium, TaskStatusInProgress, 1, 2, now, completionDate)

	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE assignee_id = \\$1 ORDER BY creation_date ASC").
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := queries.SearchTasksByAssignee(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, int64(1), tasks[0].ID)
	assert.Equal(t, "Test Task 1", tasks[0].Title)
	assert.Equal(t, "Description 1", tasks[0].Description)
	assert.Equal(t, TaskPriorityLow, tasks[0].Priority)
	assert.Equal(t, TaskStatusNew, tasks[0].Status)
	assert.Equal(t, int64(1), tasks[0].AssigneeID)
	assert.Equal(t, int64(1), tasks[0].ProjectID)
	assert.WithinDuration(t, now, tasks[0].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[0].CompletionDate.Time, time.Second)
	assert.Equal(t, int64(2), tasks[1].ID)
	assert.Equal(t, "Test Task 2", tasks[1].Title)
	assert.Equal(t, "Description 2", tasks[1].Description)
	assert.Equal(t, TaskPriorityMedium, tasks[1].Priority)
	assert.Equal(t, TaskStatusInProgress, tasks[1].Status)
	assert.Equal(t, int64(1), tasks[1].AssigneeID)
	assert.Equal(t, int64(2), tasks[1].ProjectID)
	assert.WithinDuration(t, now, tasks[1].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[1].CompletionDate.Time, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchTasksByPriority(t *testing.T) {
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
		AddRow(2, "Test Task 2", "Description 2", TaskPriorityLow, TaskStatusInProgress, 2, 2, now, completionDate)

	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE priority = \\$1 ORDER BY creation_date ASC").
		WithArgs(TaskPriorityLow).
		WillReturnRows(rows)

	tasks, err := queries.SearchTasksByPriority(context.Background(), TaskPriorityLow)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, int64(1), tasks[0].ID)
	assert.Equal(t, "Test Task 1", tasks[0].Title)
	assert.Equal(t, "Description 1", tasks[0].Description)
	assert.Equal(t, TaskPriorityLow, tasks[0].Priority)
	assert.Equal(t, TaskStatusNew, tasks[0].Status)
	assert.Equal(t, int64(1), tasks[0].AssigneeID)
	assert.Equal(t, int64(1), tasks[0].ProjectID)
	assert.WithinDuration(t, now, tasks[0].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[0].CompletionDate.Time, time.Second)
	assert.Equal(t, int64(2), tasks[1].ID)
	assert.Equal(t, "Test Task 2", tasks[1].Title)
	assert.Equal(t, "Description 2", tasks[1].Description)
	assert.Equal(t, TaskPriorityLow, tasks[1].Priority)
	assert.Equal(t, TaskStatusInProgress, tasks[1].Status)
	assert.Equal(t, int64(2), tasks[1].AssigneeID)
	assert.Equal(t, int64(2), tasks[1].ProjectID)
	assert.WithinDuration(t, now, tasks[1].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[1].CompletionDate.Time, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchTasksByProject(t *testing.T) {
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
		AddRow(2, "Test Task 2", "Description 2", TaskPriorityMedium, TaskStatusInProgress, 2, 1, now, completionDate)

	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE project_id = \\$1 ORDER BY creation_date ASC").
		WithArgs(1).
		WillReturnRows(rows)

	tasks, err := queries.SearchTasksByProject(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, int64(1), tasks[0].ID)
	assert.Equal(t, "Test Task 1", tasks[0].Title)
	assert.Equal(t, "Description 1", tasks[0].Description)
	assert.Equal(t, TaskPriorityLow, tasks[0].Priority)
	assert.Equal(t, TaskStatusNew, tasks[0].Status)
	assert.Equal(t, int64(1), tasks[0].AssigneeID)
	assert.Equal(t, int64(1), tasks[0].ProjectID)
	assert.WithinDuration(t, now, tasks[0].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[0].CompletionDate.Time, time.Second)
	assert.Equal(t, int64(2), tasks[1].ID)
	assert.Equal(t, "Test Task 2", tasks[1].Title)
	assert.Equal(t, "Description 2", tasks[1].Description)
	assert.Equal(t, TaskPriorityMedium, tasks[1].Priority)
	assert.Equal(t, TaskStatusInProgress, tasks[1].Status)
	assert.Equal(t, int64(2), tasks[1].AssigneeID)
	assert.Equal(t, int64(1), tasks[1].ProjectID)
	assert.WithinDuration(t, now, tasks[1].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[1].CompletionDate.Time, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchTasksByStatus(t *testing.T) {
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
		AddRow(2, "Test Task 2", "Description 2", TaskPriorityMedium, TaskStatusNew, 2, 2, now, completionDate)

	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE status = \\$1 ORDER BY creation_date ASC").
		WithArgs(TaskStatusNew).
		WillReturnRows(rows)

	tasks, err := queries.SearchTasksByStatus(context.Background(), TaskStatusNew)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, int64(1), tasks[0].ID)
	assert.Equal(t, "Test Task 1", tasks[0].Title)
	assert.Equal(t, "Description 1", tasks[0].Description)
	assert.Equal(t, TaskPriorityLow, tasks[0].Priority)
	assert.Equal(t, TaskStatusNew, tasks[0].Status)
	assert.Equal(t, int64(1), tasks[0].AssigneeID)
	assert.Equal(t, int64(1), tasks[0].ProjectID)
	assert.WithinDuration(t, now, tasks[0].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[0].CompletionDate.Time, time.Second)
	assert.Equal(t, int64(2), tasks[1].ID)
	assert.Equal(t, "Test Task 2", tasks[1].Title)
	assert.Equal(t, "Description 2", tasks[1].Description)
	assert.Equal(t, TaskPriorityMedium, tasks[1].Priority)
	assert.Equal(t, TaskStatusNew, tasks[1].Status)
	assert.Equal(t, int64(2), tasks[1].AssigneeID)
	assert.Equal(t, int64(2), tasks[1].ProjectID)
	assert.WithinDuration(t, now, tasks[1].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[1].CompletionDate.Time, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSearchTasksByTitle(t *testing.T) {
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
		AddRow(2, "Test Task 2", "Description 2", TaskPriorityMedium, TaskStatusNew, 2, 2, now, completionDate)

	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE title ILIKE '%' \\|\\| \\$1 \\|\\| '%' ORDER BY creation_date ASC").
		WithArgs("Test").
		WillReturnRows(rows)

	// Create a sql.NullString for the title
	title := sql.NullString{String: "Test", Valid: true}

	tasks, err := queries.SearchTasksByTitle(context.Background(), title)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, int64(1), tasks[0].ID)
	assert.Equal(t, "Test Task 1", tasks[0].Title)
	assert.Equal(t, "Description 1", tasks[0].Description)
	assert.Equal(t, TaskPriorityLow, tasks[0].Priority)
	assert.Equal(t, TaskStatusNew, tasks[0].Status)
	assert.Equal(t, int64(1), tasks[0].AssigneeID)
	assert.Equal(t, int64(1), tasks[0].ProjectID)
	assert.WithinDuration(t, now, tasks[0].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[0].CompletionDate.Time, time.Second)
	assert.Equal(t, int64(2), tasks[1].ID)
	assert.Equal(t, "Test Task 2", tasks[1].Title)
	assert.Equal(t, "Description 2", tasks[1].Description)
	assert.Equal(t, TaskPriorityMedium, tasks[1].Priority)
	assert.Equal(t, TaskStatusNew, tasks[1].Status)
	assert.Equal(t, int64(2), tasks[1].AssigneeID)
	assert.Equal(t, int64(2), tasks[1].ProjectID)
	assert.WithinDuration(t, now, tasks[1].CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, tasks[1].CompletionDate.Time, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdateTask(t *testing.T) {
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
	completionDate := sql.NullTime{Time: now, Valid: true}

	// Mock expected rows to return
	rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee_id", "project_id", "creation_date", "completion_date"}).
		AddRow(1, "Updated Task", "Updated Description", TaskPriorityHigh, TaskStatusInProgress, int64(123), int64(456), now, now)

	// Expectation: QueryRowContext with expected arguments
	mock.ExpectQuery("UPDATE tasks SET title = \\$2, description = \\$3, priority = \\$4, status = \\$5, assignee_id = \\$6, project_id = \\$7, completion_date = \\$8 WHERE id = \\$1 RETURNING id, title, description, priority, status, assignee_id, project_id, creation_date, completion_date").
		WithArgs(int64(1), "Updated Task", "Updated Description", TaskPriorityHigh, TaskStatusInProgress, int64(123), int64(456), completionDate).
		WillReturnRows(rows)

	// Prepare input params
	params := UpdateTaskParams{
		ID:             1,
		Title:          "Updated Task",
		Description:    "Updated Description",
		Priority:       TaskPriorityHigh,
		Status:         TaskStatusInProgress,
		AssigneeID:     123,
		ProjectID:      456,
		CompletionDate: completionDate,
	}

	// Call the UpdateTask method
	task, err := queries.UpdateTask(context.Background(), params)

	// Verify results
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, int64(1), task.ID, "Expected task ID to match")
	assert.Equal(t, "Updated Task", task.Title, "Expected task title to match")
	assert.Equal(t, "Updated Description", task.Description, "Expected task description to match")
	assert.Equal(t, TaskPriorityHigh, task.Priority, "Expected task priority to match")
	assert.Equal(t, TaskStatusInProgress, task.Status, "Expected task status to match")
	assert.Equal(t, int64(123), task.AssigneeID, "Expected task assignee ID to match")
	assert.Equal(t, int64(456), task.ProjectID, "Expected task project ID to match")
	assert.WithinDuration(t, now, task.CreationDate, time.Second, "Expected task creation date to match")

	assert.True(t, task.CompletionDate.Valid, "Expected task completion date to be valid")
	assert.WithinDuration(t, now, task.CompletionDate.Time, time.Second, "Expected task completion date to match")

	// Assert all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}