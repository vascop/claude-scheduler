package cmd

import (
	"fmt"

	"github.com/vascop/claude-scheduler/internal/platform"
	"github.com/vascop/claude-scheduler/internal/task"
)

func List(args []string) error {
	tasks, err := task.ListAll()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No scheduled tasks.")
		return nil
	}

	sched := platform.NewScheduler()
	for _, t := range tasks {
		status := "unloaded"
		if sched.IsInstalled(t.ID) {
			status = "loaded"
		}
		fmt.Printf("%-20s %-18s %-10s %-10v %s\n", t.ID, t.Cron, status, t.Autonomous, t.Prompt)
	}
	return nil
}
