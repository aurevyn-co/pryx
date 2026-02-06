package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"pryx-core/internal/store"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

// TaskType defines the type of scheduled task
type TaskType string

const (
	TaskTypeMessage  TaskType = "message"
	TaskTypeWorkflow TaskType = "workflow"
	TaskTypeReminder TaskType = "reminder"
	TaskTypeWebhook  TaskType = "webhook"
)

// TaskStatus defines the status of a scheduled task
type TaskStatus string

const (
	TaskStatusActive TaskStatus = "active"
	TaskStatusPaused TaskStatus = "paused"
	TaskStatusError  TaskStatus = "error"
)

// RunStatus defines the status of a task execution
type RunStatus string

const (
	RunStatusPending RunStatus = "pending"
	RunStatusRunning RunStatus = "running"
	RunStatusSuccess RunStatus = "success"
	RunStatusFailed  RunStatus = "failed"
)

// ScheduledTask represents a scheduled task in the database
type ScheduledTask struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	CronExpression string     `json:"cron_expression"`
	TaskType       TaskType   `json:"task_type"`
	Payload        string     `json:"payload"`
	Timezone       string     `json:"timezone"`
	Enabled        bool       `json:"enabled"`
	LastRunAt      *time.Time `json:"last_run_at,omitempty"`
	LastRunStatus  string     `json:"last_run_status,omitempty"`
	LastRunError   string     `json:"last_run_error,omitempty"`
	NextRunAt      *time.Time `json:"next_run_at,omitempty"`
	RunCount       int        `json:"run_count"`
	UserID         string     `json:"user_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TaskRun represents a single execution of a scheduled task
type TaskRun struct {
	ID          string     `json:"id"`
	TaskID      string     `json:"task_id"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Status      RunStatus  `json:"status"`
	Error       string     `json:"error,omitempty"`
	Output      string     `json:"output,omitempty"`
}

// TaskExecutor defines the interface for executing scheduled tasks
type TaskExecutor interface {
	Execute(ctx context.Context, task *ScheduledTask) (string, error)
}

// Scheduler manages scheduled tasks and their execution
type Scheduler struct {
	db        *sql.DB
	store     *store.Store
	cron      *cron.Cron
	executors map[TaskType]TaskExecutor
	tasks     map[string]cron.EntryID
	mu        sync.RWMutex
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

// New creates a new Scheduler instance
func New(db *sql.DB) *Scheduler {
	return &Scheduler{
		db:        db,
		store:     store.NewFromDB(db),
		cron:      cron.New(cron.WithSeconds()),
		executors: make(map[TaskType]TaskExecutor),
		tasks:     make(map[string]cron.EntryID),
		stopChan:  make(chan struct{}),
	}
}

// RegisterExecutor registers a task executor for a specific task type
func (s *Scheduler) RegisterExecutor(taskType TaskType, executor TaskExecutor) {
	s.executors[taskType] = executor
}

// Start begins the scheduler and loads all enabled tasks
func (s *Scheduler) Start(ctx context.Context) error {
	s.wg.Add(1)
	go s.run(ctx)

	// Load and schedule all enabled tasks
	tasks, err := s.loadEnabledTasks()
	if err != nil {
		return fmt.Errorf("failed to load enabled tasks: %w", err)
	}

	for _, task := range tasks {
		if err := s.scheduleTask(task); err != nil {
			log.Printf("Failed to schedule task %s: %v", task.ID, err)
		}
	}

	log.Printf("Scheduler started with %d tasks", len(tasks))
	return nil
}

// Stop gracefully shuts down the scheduler
func (s *Scheduler) Stop() {
	close(s.stopChan)
	s.cron.Stop()
	s.wg.Wait()
	log.Println("Scheduler stopped")
}

// run is the main scheduler loop
func (s *Scheduler) run(ctx context.Context) {
	defer s.wg.Done()

	// Periodic task refresh (every 5 minutes)
	refreshTicker := time.NewTicker(5 * time.Minute)
	defer refreshTicker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ctx.Done():
			return
		case <-refreshTicker.C:
			s.refreshTasks()
		}
	}
}

// refreshTasks reloads tasks from the database
func (s *Scheduler) refreshTasks() {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.loadEnabledTasks()
	if err != nil {
		log.Printf("Failed to refresh tasks: %v", err)
		return
	}

	// Add new tasks
	for _, task := range tasks {
		if _, exists := s.tasks[task.ID]; !exists {
			if err := s.scheduleTask(task); err != nil {
				log.Printf("Failed to schedule task %s: %v", task.ID, err)
			}
		}
	}
}

// loadEnabledTasks loads all enabled tasks from the database
func (s *Scheduler) loadEnabledTasks() ([]*ScheduledTask, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, cron_expression, task_type, payload,
		       timezone, enabled, last_run_at, last_run_status, last_run_error,
		       next_run_at, run_count, user_id, created_at, updated_at
		FROM scheduled_tasks
		WHERE enabled = 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*ScheduledTask
	for rows.Next() {
		task := &ScheduledTask{}
		err := rows.Scan(
			&task.ID, &task.Name, &task.Description, &task.CronExpression,
			&task.TaskType, &task.Payload, &task.Timezone, &task.Enabled,
			&task.LastRunAt, &task.LastRunStatus, &task.LastRunError,
			&task.NextRunAt, &task.RunCount, &task.UserID,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// scheduleTask adds a task to the cron scheduler
func (s *Scheduler) scheduleTask(task *ScheduledTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Parse cron expression
	if _, err := cron.ParseStandard(task.CronExpression); err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	// Create runner function
	runner := func() {
		s.executeTask(task)
	}

	// Add to cron scheduler
	entryID, err := s.cron.AddFunc(task.CronExpression, runner)
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	s.tasks[task.ID] = entryID
	log.Printf("Scheduled task %s (%s) with cron: %s", task.ID, task.Name, task.CronExpression)

	return nil
}

// removeTask removes a task from the cron scheduler
func (s *Scheduler) removeTask(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.tasks[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.tasks, taskID)
		log.Printf("Removed task %s from scheduler", taskID)
	}
}

// executeTask runs a single scheduled task
func (s *Scheduler) executeTask(task *ScheduledTask) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Create run record
	runID := uuid.New().String()
	run := &TaskRun{
		ID:        runID,
		TaskID:    task.ID,
		StartedAt: time.Now(),
		Status:    RunStatusRunning,
	}

	// Save run start
	if err := s.saveRun(run); err != nil {
		log.Printf("Failed to save run start: %v", err)
	}

	// Get executor
	executor, exists := s.executors[task.TaskType]
	if !exists {
		run.Status = RunStatusFailed
		run.Error = fmt.Sprintf("no executor for task type: %s", task.TaskType)
		s.completeRun(run, task)
		return
	}

	// Execute task
	output, err := executor.Execute(ctx, task)
	now := time.Now()
	run.CompletedAt = &now

	if err != nil {
		run.Status = RunStatusFailed
		run.Error = err.Error()
	} else {
		run.Status = RunStatusSuccess
		run.Output = output
	}

	s.completeRun(run, task)
}

// completeRun updates the task and run records after execution
func (s *Scheduler) completeRun(run *TaskRun, task *ScheduledTask) {
	// Save run completion
	if err := s.saveRun(run); err != nil {
		log.Printf("Failed to save run completion: %v", err)
	}

	// Update task status
	now := time.Now()
	nextRun := s.getNextRunTime(task.CronExpression)

	_, err := s.db.Exec(`
		UPDATE scheduled_tasks
		SET last_run_at = ?, last_run_status = ?, last_run_error = ?,
		    next_run_at = ?, run_count = run_count + 1, updated_at = ?
		WHERE id = ?
	`,
		run.StartedAt, run.Status, run.Error,
		nextRun, now, task.ID,
	)

	if err != nil {
		log.Printf("Failed to update task: %v", err)
	}
}

// saveRun saves a task run record
func (s *Scheduler) saveRun(run *TaskRun) error {
	_, err := s.db.Exec(`
		INSERT INTO scheduled_task_runs (id, task_id, started_at, completed_at, status, error, output)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		run.ID, run.TaskID, run.StartedAt, run.CompletedAt, run.Status, run.Error, run.Output,
	)
	return err
}

// getNextRunTime calculates the next run time from a cron expression
func (s *Scheduler) getNextRunTime(cronExpr string) *time.Time {
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return nil
	}
	next := schedule.Next(time.Now())
	return &next
}

// CreateTask creates a new scheduled task
func (s *Scheduler) CreateTask(task *ScheduledTask) error {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	// Calculate next run time
	nextRun := s.getNextRunTime(task.CronExpression)
	task.NextRunAt = nextRun

	// Insert into database
	_, err := s.db.Exec(`
		INSERT INTO scheduled_tasks (
			id, name, description, cron_expression, task_type, payload,
			timezone, enabled, next_run_at, run_count, user_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.ID, task.Name, task.Description, task.CronExpression,
		task.TaskType, task.Payload, task.Timezone, task.Enabled,
		task.NextRunAt, task.RunCount, task.UserID, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Schedule if enabled
	if task.Enabled {
		if err := s.scheduleTask(task); err != nil {
			return err
		}
	}

	return nil
}

// GetTask retrieves a task by ID
func (s *Scheduler) GetTask(id string) (*ScheduledTask, error) {
	task := &ScheduledTask{}
	err := s.db.QueryRow(`
		SELECT id, name, description, cron_expression, task_type, payload,
		       timezone, enabled, last_run_at, last_run_status, last_run_error,
		       next_run_at, run_count, user_id, created_at, updated_at
		FROM scheduled_tasks WHERE id = ?
	`, id).Scan(
		&task.ID, &task.Name, &task.Description, &task.CronExpression,
		&task.TaskType, &task.Payload, &task.Timezone, &task.Enabled,
		&task.LastRunAt, &task.LastRunStatus, &task.LastRunError,
		&task.NextRunAt, &task.RunCount, &task.UserID,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

// ListTasks lists all tasks for a user
func (s *Scheduler) ListTasks(userID string) ([]*ScheduledTask, error) {
	var query string
	var args []interface{}

	if userID == "" {
		query = `SELECT id, name, description, cron_expression, task_type, payload,
		        timezone, enabled, last_run_at, last_run_status, last_run_error,
		        next_run_at, run_count, user_id, created_at, updated_at
		    FROM scheduled_tasks ORDER BY created_at DESC`
	} else {
		query = `SELECT id, name, description, cron_expression, task_type, payload,
		        timezone, enabled, last_run_at, last_run_status, last_run_error,
		        next_run_at, run_count, user_id, created_at, updated_at
		    FROM scheduled_tasks WHERE user_id = ? ORDER BY created_at DESC`
		args = []interface{}{userID}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*ScheduledTask
	for rows.Next() {
		task := &ScheduledTask{}
		err := rows.Scan(
			&task.ID, &task.Name, &task.Description, &task.CronExpression,
			&task.TaskType, &task.Payload, &task.Timezone, &task.Enabled,
			&task.LastRunAt, &task.LastRunStatus, &task.LastRunError,
			&task.NextRunAt, &task.RunCount, &task.UserID,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (s *Scheduler) UpdateTask(task *ScheduledTask) error {
	task.UpdatedAt = time.Now()

	// Recalculate next run
	nextRun := s.getNextRunTime(task.CronExpression)
	task.NextRunAt = nextRun

	_, err := s.db.Exec(`
		UPDATE scheduled_tasks
		SET name = ?, description = ?, cron_expression = ?, task_type = ?,
		    payload = ?, timezone = ?, enabled = ?, next_run_at = ?, updated_at = ?
		WHERE id = ?
	`,
		task.Name, task.Description, task.CronExpression, task.TaskType,
		task.Payload, task.Timezone, task.Enabled, task.NextRunAt,
		task.UpdatedAt, task.ID,
	)
	if err != nil {
		return err
	}

	// Update scheduler
	s.removeTask(task.ID)
	if task.Enabled {
		if err := s.scheduleTask(task); err != nil {
			return err
		}
	}

	return nil
}

// DeleteTask deletes a task and its runs
func (s *Scheduler) DeleteTask(id string) error {
	s.removeTask(id)

	// Delete runs first (foreign key constraint)
	_, err := s.db.Exec("DELETE FROM scheduled_task_runs WHERE task_id = ?", id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("DELETE FROM scheduled_tasks WHERE id = ?", id)
	return err
}

// EnableTask enables a disabled task
func (s *Scheduler) EnableTask(id string) error {
	_, err := s.db.Exec("UPDATE scheduled_tasks SET enabled = 1, updated_at = ? WHERE id = ?",
		time.Now(), id)
	if err != nil {
		return err
	}

	task, err := s.GetTask(id)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task not found: %s", id)
	}

	return s.scheduleTask(task)
}

// DisableTask disables an active task
func (s *Scheduler) DisableTask(id string) error {
	s.removeTask(id)

	_, err := s.db.Exec("UPDATE scheduled_tasks SET enabled = 0, updated_at = ? WHERE id = ?",
		time.Now(), id)
	return err
}

// GetTaskRuns returns the execution history for a task
func (s *Scheduler) GetTaskRuns(taskID string, limit int) ([]*TaskRun, error) {
	rows, err := s.db.Query(`
		SELECT id, task_id, started_at, completed_at, status, error, output
		FROM scheduled_task_runs
		WHERE task_id = ?
		ORDER BY started_at DESC
		LIMIT ?
	`, taskID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []*TaskRun
	for rows.Next() {
		run := &TaskRun{}
		err := rows.Scan(
			&run.ID, &run.TaskID, &run.StartedAt, &run.CompletedAt,
			&run.Status, &run.Error, &run.Output,
		)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}

	return runs, nil
}

// ValidateCronExpression validates a cron expression
func ValidateCronExpression(expr string) error {
	_, err := cron.ParseStandard(expr)
	return err
}

// ParseTaskPayload parses the payload JSON for a task
func ParseTaskPayload(payload string, target interface{}) error {
	return json.Unmarshal([]byte(payload), target)
}

// MarshalTaskPayload marshals a struct to JSON for task payload
func MarshalTaskPayload(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
