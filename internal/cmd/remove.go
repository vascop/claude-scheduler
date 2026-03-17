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
	platform.NewScheduler().Uninstall(id)
	task.Delete(id)
	fmt.Printf("Removed task '%s'\n", id)
	return nil
}
