package platform

import (
	"strings"
	"testing"
)

func TestCrontabMarker(t *testing.T) {
	got := CrontabMarker("sync-repos")
	want := "# claude-scheduler:sync-repos"
	if got != want {
		t.Errorf("CrontabMarker() = %q, want %q", got, want)
	}
}

func TestCrontabEntry(t *testing.T) {
	got := CrontabEntry("sync-repos", "3 9 * * *", "/home/user/.claude-scheduler/tasks/sync-repos.sh", "/home/user/.claude-scheduler/logs/sync-repos.log")
	if !strings.Contains(got, "3 9 * * *") {
		t.Error("entry should contain cron expression")
	}
	if !strings.Contains(got, "/bin/bash") {
		t.Error("entry should use /bin/bash")
	}
	if !strings.Contains(got, "# claude-scheduler:sync-repos") {
		t.Error("entry should contain marker")
	}
	if !strings.Contains(got, ">> /home/user/.claude-scheduler/logs/sync-repos.log 2>&1") {
		t.Error("entry should redirect output to log file")
	}
}

func TestCrontabAddToEmpty(t *testing.T) {
	result := CrontabAdd("", "sync-repos", "3 9 * * *", "/path/to/script.sh", "/path/to/log.log")
	if !strings.Contains(result, "3 9 * * *") {
		t.Error("should contain cron expression")
	}
	if !strings.Contains(result, "# claude-scheduler:sync-repos") {
		t.Error("should contain marker")
	}
	if !strings.HasSuffix(result, "\n") {
		t.Error("should end with newline")
	}
}

func TestCrontabAddToExisting(t *testing.T) {
	existing := "0 * * * * /usr/bin/some-other-job\n"
	result := CrontabAdd(existing, "daily-review", "3 9 * * 1-5", "/path/to/script.sh", "/path/to/log.log")

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[0], "some-other-job") {
		t.Error("should preserve existing entries")
	}
	if !strings.Contains(lines[1], "claude-scheduler:daily-review") {
		t.Error("should add new entry")
	}
}

func TestCrontabAddReplacesExisting(t *testing.T) {
	existing := "3 9 * * * /bin/bash /old/script.sh >> /old/log.log 2>&1 # claude-scheduler:sync-repos\n"
	result := CrontabAdd(existing, "sync-repos", "7 10 * * *", "/new/script.sh", "/new/log.log")

	if strings.Contains(result, "/old/") {
		t.Error("should remove old entry")
	}
	if !strings.Contains(result, "7 10 * * *") {
		t.Error("should contain new cron expression")
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line after replace, got %d: %v", len(lines), lines)
	}
}

func TestCrontabRemove(t *testing.T) {
	content := "0 * * * * /usr/bin/job-a\n3 9 * * * /bin/bash /path/script.sh >> /path/log.log 2>&1 # claude-scheduler:sync-repos\n30 18 * * * /usr/bin/job-b\n"
	result := CrontabRemove(content, "sync-repos")

	if strings.Contains(result, "claude-scheduler:sync-repos") {
		t.Error("should remove the marked line")
	}
	if !strings.Contains(result, "job-a") {
		t.Error("should preserve other entries")
	}
	if !strings.Contains(result, "job-b") {
		t.Error("should preserve other entries")
	}
}

func TestCrontabRemoveNonExistent(t *testing.T) {
	content := "0 * * * * /usr/bin/job-a\n"
	result := CrontabRemove(content, "does-not-exist")

	if result != content {
		t.Errorf("should not modify content when task doesn't exist, got %q", result)
	}
}

func TestCrontabHasTask(t *testing.T) {
	content := "0 * * * * /usr/bin/job-a\n3 9 * * * /bin/bash /path/script.sh >> /path/log.log 2>&1 # claude-scheduler:sync-repos\n"

	if !CrontabHasTask(content, "sync-repos") {
		t.Error("should find sync-repos")
	}
	if CrontabHasTask(content, "nonexistent") {
		t.Error("should not find nonexistent task")
	}
}

func TestCrontabAddMultipleTasks(t *testing.T) {
	result := CrontabAdd("", "task-a", "0 9 * * *", "/path/a.sh", "/path/a.log")
	result = CrontabAdd(result, "task-b", "0 18 * * *", "/path/b.sh", "/path/b.log")
	result = CrontabAdd(result, "task-c", "*/15 * * * *", "/path/c.sh", "/path/c.log")

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	if !CrontabHasTask(result, "task-a") {
		t.Error("should contain task-a")
	}
	if !CrontabHasTask(result, "task-b") {
		t.Error("should contain task-b")
	}
	if !CrontabHasTask(result, "task-c") {
		t.Error("should contain task-c")
	}
}
