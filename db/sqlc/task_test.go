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
	completionDate := sql.NullTime{Time: now.AddDate(0, 0, 1), Valid: true}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee_id", "project_id", "creation_date", "completion_date"}).
		AddRow(1, "Test Task", "Description", TaskPriorityLow, TaskStatusNew, 1, 1, now, completionDate)

	mock.ExpectQuery("INSERT INTO tasks").
		WithArgs("Test Task", "Description", TaskPriorityLow, TaskStatusNew, 1, 1, now, completionDate).
		WillReturnRows(rows)

	params := CreateTaskParams{
		Title:          "Test Task",
		Description:    "Description",
		Priority:       TaskPriorityLow,
		Status:         TaskStatusNew,
		AssigneeID:     1,
		ProjectID:      1,
		CreationDate:   now,
		CompletionDate: completionDate,
	}

	task, err := queries.CreateTask(context.Background(), params)

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

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	mock.ExpectExec("DELETE FROM tasks").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = queries.DeleteTask(context.Background(), 1)
	assert.NoError(t, err)

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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	queries := New(db)

	now := time.Now()
	completionDate := sql.NullTime{Time: now.AddDate(0, 0, 1), Valid: true}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee_id", "project_id", "creation_date", "completion_date"}).
		AddRow(1, "Updated Task", "Updated Description", TaskPriorityMedium, TaskStatusInProgress, 1, 1, now, completionDate)

	mock.ExpectQuery("UPDATE tasks SET").
		WithArgs(1, "Updated Task", "Updated Description", TaskPriorityMedium, TaskStatusInProgress, 1, 1, now, completionDate).
		WillReturnRows(rows)

	params := UpdateTaskParams{
		ID:             1,
		Title:          "Updated Task",
		Description:    "Updated Description",
		Priority:       TaskPriorityMedium,
		Status:         TaskStatusInProgress,
		AssigneeID:     1,
		ProjectID:      1,
		CreationDate:   now,
		CompletionDate: completionDate,
	}

	task, err := queries.UpdateTask(context.Background(), params)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), task.ID)
	assert.Equal(t, "Updated Task", task.Title)
	assert.Equal(t, "Updated Description", task.Description)
	assert.Equal(t, TaskPriorityMedium, task.Priority)
	assert.Equal(t, TaskStatusInProgress, task.Status)
	assert.Equal(t, int64(1), task.AssigneeID)
	assert.Equal(t, int64(1), task.ProjectID)
	assert.WithinDuration(t, now, task.CreationDate, time.Second)
	assert.WithinDuration(t, completionDate.Time, task.CompletionDate.Time, time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}