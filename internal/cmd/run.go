package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/vascop/claude-scheduler/internal/task"
)

func Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: run <id>")
	}
	id := args[0]

	scriptPath := task.ScriptPath(id)
	if _, err := os.Stat(scriptPath); err != nil {
		return fmt.Errorf("task '%s' not found", id)
	}

	fmt.Printf("Running task '%s' now...\n", id)
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
