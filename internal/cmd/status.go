package cmd

import (
	"fmt"

	"github.com/vascop/claude-scheduler/internal/platform"
	"github.com/vascop/claude-scheduler/internal/task"
)

func Status(args []string) error {
	tasks, err := task.ListAll()
	if err != nil {
		return err
	}

	sched := platform.NewScheduler()
	fmt.Printf("claude-scheduler tasks (%s):\n", sched.Name())
	fmt.Println()

	if len(tasks) == 0 {
		fmt.Println("  (none)")
		return nil
	}

	for _, t := range tasks {
		mark := "✗ not installed"
		if sched.IsInstalled(t.ID) {
			mark = "✓ installed"
		}
		fmt.Printf("  %s: %s\n", t.ID, mark)
	}
	return nil
}
