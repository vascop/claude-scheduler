package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/vascop/claude-scheduler/internal/cron"
	"github.com/vascop/claude-scheduler/internal/platform"
	"github.com/vascop/claude-scheduler/internal/runner"
	"github.com/vascop/claude-scheduler/internal/task"
)

func Add(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: add <id> <cron> <prompt> [--autonomous] [--worktree]")
	}

	id := args[0]
	cronExpr := args[1]

	var promptParts []string
	autonomous := false
	worktree := false

	for _, arg := range args[2:] {
		switch arg {
		case "--autonomous":
			autonomous = true
		case "--worktree":
			worktree = true
		default:
			promptParts = append(promptParts, arg)
		}
	}

	prompt := strings.Join(promptParts, " ")
	if prompt == "" {
		return fmt.Errorf("prompt cannot be empty")
	}

	schedule, err := cron.Parse(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression %q: %w", cronExpr, err)
	}

	t := &task.Task{
		ID:         id,
		Cron:       cronExpr,
		Prompt:     prompt,
		Autonomous: autonomous,
		Worktree:   worktree,
		Created:    time.Now().UTC().Format(time.RFC3339),
	}

	scriptPath, err := runner.WriteScript(t)
	if err != nil {
		return fmt.Errorf("writing runner script: %w", err)
	}

	sched := platform.NewScheduler()
	if err := sched.Install(id, cronExpr, schedule, scriptPath, task.LogPath(id)); err != nil {
		return fmt.Errorf("installing schedule: %w", err)
	}

	if err := task.Save(t); err != nil {
		return fmt.Errorf("saving task: %w", err)
	}

	fmt.Printf("Scheduled task '%s'\n", id)
	fmt.Printf("  Cron:       %s\n", cronExpr)
	fmt.Printf("  Prompt:     %s\n", prompt)
	fmt.Printf("  Autonomous: %v\n", autonomous)
	fmt.Printf("  Worktree:   %v\n", worktree)
	fmt.Printf("  Backend:    %s\n", sched.Name())
	fmt.Printf("  Location:   %s\n", sched.Describe(id))
	fmt.Printf("  Logs:       %s\n", task.LogPath(id))
	return nil
}
