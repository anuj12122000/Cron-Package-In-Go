package main

import (
	"fmt"
	"log"
	"task_scheduler/cron"
	"time"
)

func main() {
	var s = cron.NewScheduler()

	// Define a task
	taskID := "example_task"
	taskCron := "*/10 * *  * * *" // Every 10 seconds

	// first is seconds,minute,hour,day,month,weekday
	taskFunc := func() error {
		fmt.Printf("Task '%s' executed at %s\n", taskID, time.Now())
		return nil
	}

	// Add the task to the scheduler
	if err := s.AddTask(taskID, taskCron, taskFunc); err != nil {
		log.Fatalf("Failed to add task: %v", err)
	}

	s.Start()

	select {}
}
