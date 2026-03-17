# claude-scheduler

Schedule Claude Code tasks that run automatically via your OS's native scheduler. Tasks persist across restarts and run even when Claude isn't open.

- **macOS**: launchd (plist in `~/Library/LaunchAgents/`)
- **Linux**: crontab

## Install

```bash
cd ~/claude-scheduler
go build -o bin/claude-scheduler .
```

Optionally add to PATH:

```bash
export PATH="$HOME/claude-scheduler/bin:$PATH"
```

## Usage

```
claude-scheduler add <id> <cron> <prompt> [--autonomous] [--worktree]
claude-scheduler remove <id>
claude-scheduler list
claude-scheduler logs [id] [-n N]
claude-scheduler run <id>
claude-scheduler status
```

### Add a task

```bash
# Sync repos daily at 9:03am
claude-scheduler add sync-repos "3 9 * * *" /sync-repos

# Code review every weekday at 9am (autonomous — can edit files)
claude-scheduler add daily-review "3 9 * * 1-5" "review yesterday's commits" --autonomous

# Security scan weekly in isolated worktree
claude-scheduler add security-scan "7 10 * * 1" "scan for vulnerabilities" --worktree
```

### Flags

- `--autonomous` — runs Claude with `--dangerously-skip-permissions`, allowing file edits, commands, and commits without prompts.
- `--worktree` — runs in an isolated git worktree. Changes are committed to a new branch and pushed for review.

### Cron format

Standard 5-field: `minute hour day-of-month month day-of-week`

| Expression | Meaning |
|---|---|
| `3 9 * * *` | Daily at 9:03 AM |
| `3 9 * * 1-5` | Weekdays at 9:03 AM |
| `17 */4 * * *` | Every 4 hours at :17 |
| `0 0 */2 * *` | Every 2 days at midnight |

## How it works

1. Task metadata is saved to `~/.claude-scheduler/tasks/<id>.json`
2. A runner shell script is generated at `~/.claude-scheduler/tasks/<id>.sh`
3. The task is registered with the OS scheduler:
   - **macOS**: a launchd plist is written to `~/Library/LaunchAgents/` and loaded
   - **Linux**: a crontab entry is added with a marker comment for tracking
4. At the scheduled time, the OS runs the script, which invokes `claude -p "<prompt>"`
5. Output is logged to `~/.claude-scheduler/logs/<id>.log`

## Claude Code skill

A companion skill at `~/.claude/skills/scheduler/` lets you manage tasks with natural language from within Claude Code:

```
/scheduler every weekday at 9am /sync-repos
/scheduler list
/scheduler remove sync-repos
```