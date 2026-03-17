package main

import (
	"fmt"
	"os"

	"github.com/vascop/claude-scheduler/internal/cmd"
	"github.com/vascop/claude-scheduler/internal/task"
)

func main() {
	if err := os.MkdirAll(task.TasksDir(), 0755); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(task.LogsDir(), 0755); err != nil {
		fatal(err)
	}

	commands := map[string]func([]string) error{
		"add":    cmd.Add,
		"remove": cmd.Remove,
		"list":   cmd.List,
		"logs":   cmd.Logs,
		"run":    cmd.Run,
		"status": cmd.Status,
	}

	if len(os.Args) < 2 {
		usage()
	}

	fn, ok := commands[os.Args[1]]
	if !ok {
		usage()
	}

	if err := fn(os.Args[2:]); err != nil {
		fatal(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `claude-scheduler — schedule Claude Code tasks via native OS scheduler

Uses launchd on macOS, crontab on Linux.

Commands:
  add <id> <cron> <prompt> [--autonomous] [--worktree]
      Schedule a new task. Cron is standard 5-field (min hour dom mon dow).
      --autonomous: run with --dangerously-skip-permissions
      --worktree:   run in isolated git worktree

  remove <id>       Remove a scheduled task
  list              List all tasks
  logs [id] [-n N]  Show execution logs (default: last 30 lines)
  run <id>          Run a task immediately
  status            Check which tasks are installed`)
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
