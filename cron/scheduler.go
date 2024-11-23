package cron

import (
	"fmt"
	"sync"
	"time"
)

// Scheduler orchestrates task management and execution
type Scheduler struct {
	mu    sync.Mutex
	tasks map[string]*Task
	stop  chan struct{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*Task),
		stop:  make(chan struct{}),
	}
}

// AddTask registers a task with a cron schedule
func (s *Scheduler) AddTask(id, cronExpr string, execute func() error) error {
	schedule, err := ParseCron(cronExpr)
	fmt.Println("TASK ADDED ", schedule.Second, schedule.Minute)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[id]; exists {
		return fmt.Errorf("task with ID '%s' already exists", id)
	}

	s.tasks[id] = &Task{
		ID:       id,
		Schedule: schedule,
		Execute:  execute,
	}
	return nil
}

// RemoveTask removes a task by ID
func (s *Scheduler) RemoveTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("task with ID '%s' not found", id)
	}
	delete(s.tasks, id)
	return nil
}

// Start initiates the scheduler's execution loop
func (s *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // every 10 seconds so that load on cpu is less
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.runDueTasks(time.Now())
			case <-s.stop:
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

// runDueTasks executes tasks due at the current time
func (s *Scheduler) runDueTasks(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		if task.Schedule.Matches(now) {
			go func(t *Task) {
				if err := t.Execute(); err != nil {
					fmt.Printf("Task '%s' execution failed: %v\n", t.ID, err)
				}
			}(task)
		}
	}
}


