package platform

import "github.com/vascop/claude-scheduler/internal/cron"

// Scheduler is the interface for OS-native task scheduling backends.
type Scheduler interface {
	// Install registers a scheduled task with the OS scheduler.
	Install(id string, cronExpr string, schedule *cron.Schedule, scriptPath, logPath string) error

	// Uninstall removes a scheduled task from the OS scheduler.
	Uninstall(id string) error

	// IsInstalled checks whether a task is registered with the OS scheduler.
	IsInstalled(id string) bool

	// Describe returns a human-readable description of where the task is registered.
	// For example, the plist path on macOS or "crontab" on Linux.
	Describe(id string) string

	// Name returns the scheduler backend name (e.g. "launchd", "crontab").
	Name() string
}
