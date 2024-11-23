package cron

// Task represents a single task to be scheduled
type Task struct {
	ID       string
	Schedule *CronSchedule
	Execute  func() error
}
