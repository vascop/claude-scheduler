package cmd

import (
	"fmt"

	"github.com/vascop/claude-scheduler/internal/platform"
	"github.com/vascop/claude-scheduler/internal/task"
)

func Remove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: remove <id>")
	}
	id := args[0]

	if _, err := task.Load(id); err != nil {
		return fmt.Errorf("task '%s' not found", id)
	}

	if err := platform.NewScheduler().Uninstall(id); err != nil {
		return fmt.Errorf("uninstalling from scheduler: %w", err)
	}
	task.Delete(id)
	fmt.Printf("Removed task '%s'\n", id)
	return nil
}
