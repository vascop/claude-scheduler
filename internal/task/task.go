package task

// Task represents a scheduled Claude Code task.
type Task struct {
	ID         string `json:"id"`
	Cron       string `json:"cron"`
	Prompt     string `json:"prompt"`
	Autonomous bool   `json:"autonomous"`
	Worktree   bool   `json:"worktree"`
	Created    string `json:"created"`
}
