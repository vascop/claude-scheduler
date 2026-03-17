package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/vascop/claude-scheduler/internal/task"
)

var scriptTmpl = template.Must(template.New("runner").Parse(`#!/usr/bin/env bash
set -euo pipefail

LOG_FILE="{{.LogPath}}"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

echo "" >> "$LOG_FILE"
echo "=== Run: $TIMESTAMP ===" >> "$LOG_FILE"
{{- if .Worktree}}

mkdir -p "$HOME/.claude-scheduler/worktrees"
WORK_DIR=$(mktemp -d "$HOME/.claude-scheduler/worktrees/XXXXXX")
BRANCH="claude-task/{{.ID}}-$(date +%s)"
if ! git worktree add -b "$BRANCH" "$WORK_DIR" HEAD 2>>"$LOG_FILE"; then
  echo "=== Exit: failed to create worktree ===" >> "$LOG_FILE"
  exit 1
fi
cd "$WORK_DIR"

cleanup_worktree() {
  git worktree remove --force "$WORK_DIR" 2>/dev/null || true
}
trap cleanup_worktree EXIT
{{- end}}

"{{.ClaudeBin}}"{{.ClaudeFlags}} -p {{.QuotedPrompt}} >> "$LOG_FILE" 2>&1
EXIT_CODE=$?
{{- if .Worktree}}

if [[ $EXIT_CODE -eq 0 ]] && ! git diff --quiet HEAD 2>/dev/null; then
  git add -A && git commit -m "claude-scheduler: automated changes" && git push origin HEAD 2>>"$LOG_FILE"
fi
{{- end}}

echo "=== Exit: $EXIT_CODE ===" >> "$LOG_FILE"
`))

type scriptData struct {
	ID           string
	LogPath      string
	ClaudeBin    string
	ClaudeFlags  string
	QuotedPrompt string
	Worktree     bool
}

// WriteScript generates and writes the runner script for a task.
func WriteScript(t *task.Task) (string, error) {
	claudeBin := findClaude()

	var flags string
	if t.Autonomous {
		flags = " --dangerously-skip-permissions"
	}

	var buf bytes.Buffer
	err := scriptTmpl.Execute(&buf, scriptData{
		ID:           t.ID,
		LogPath:      task.LogPath(t.ID),
		ClaudeBin:    claudeBin,
		ClaudeFlags:  flags,
		QuotedPrompt: fmt.Sprintf("%q", t.Prompt),
		Worktree:     t.Worktree,
	})
	if err != nil {
		return "", err
	}

	path := task.ScriptPath(t.ID)
	if err := os.WriteFile(path, buf.Bytes(), 0755); err != nil {
		return "", err
	}
	return path, nil
}

func findClaude() string {
	if p, err := exec.LookPath("claude"); err == nil {
		return p
	}
	return "/usr/local/bin/claude"
}
