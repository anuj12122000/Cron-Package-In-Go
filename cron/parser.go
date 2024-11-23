package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CronSchedule struct {
	Second  []int
	Minute  []int
	Hour    []int
	Day     []int
	Month   []int
	Weekday []int
	NextRun time.Time
}

func (c *CronSchedule) Matches(t time.Time) bool {
	return contains(c.Second, t.Second()) &&
		contains(c.Minute, t.Minute()) &&
		contains(c.Hour, t.Hour()) &&
		contains(c.Day, t.Day()) &&
		contains(c.Month, int(t.Month()))
}

func ParseCron(expr string) (*CronSchedule, error) {
	fields := strings.Fields(expr)
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid cron expression: must have 6 fields")
	}

	schedule := &CronSchedule{
		Second:  parseField(fields[0], 0, 59),
		Minute:  parseField(fields[1], 0, 59),
		Hour:    parseField(fields[2], 0, 23),
		Day:     parseField(fields[3], 1, 31),
		Month:   parseField(fields[4], 1, 12),
		Weekday: parseField(fields[5], 0, 6),
	}

	return schedule, nil
}

func parseField(field string, min, max int) []int {
	if field == "*" {
		values := make([]int, max-min+1)
		for i := min; i <= max; i++ {
			values[i-min] = i
		}
		return values
	}

	// 	If the field starts with "*/", it represents a step value, e.g., "*/5" means every 5 units.
	// The function extracts the step value, iterates through the range using the step, and returns the corresponding values.
	if strings.HasPrefix(field, "*/") {
		step, _ := strconv.Atoi(field[2:]) // as 0th and 1st indec are */ , so we would get the value from second index ownwards
		var values []int
		for i := min; i <= max; i += step {
			values = append(values, i)
		}
		return values
	}

	parts := strings.Split(field, ",")
	var values []int
	for _, part := range parts {
		value, _ := strconv.Atoi(part)
		values = append(values, value)
	}
	return values
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
