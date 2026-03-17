package platform

import (
	"fmt"
	"strings"
)

// CrontabEntry builds the crontab line for a task.
func CrontabEntry(id, cronExpr, scriptPath, logPath string) string {
	return fmt.Sprintf("%s /bin/bash %s >> %s 2>&1 %s", cronExpr, scriptPath, logPath, CrontabMarker(id))
}

// CrontabMarker returns the comment marker used to identify a task in crontab.
func CrontabMarker(id string) string {
	return fmt.Sprintf("# claude-scheduler:%s", id)
}

// CrontabAdd adds a task entry to the given crontab content, removing any existing entry for the same id.
func CrontabAdd(existing, id, cronExpr, scriptPath, logPath string) string {
	cleaned := CrontabRemove(existing, id)
	line := CrontabEntry(id, cronExpr, scriptPath, logPath)

	if strings.TrimSpace(cleaned) == "" {
		return line + "\n"
	}
	return strings.TrimRight(cleaned, "\n") + "\n" + line + "\n"
}

// CrontabRemove removes all lines matching the task marker from crontab content.
func CrontabRemove(content, id string) string {
	marker := CrontabMarker(id)
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		if !strings.Contains(line, marker) {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

// CrontabHasTask checks whether the crontab content contains the given task.
func CrontabHasTask(content, id string) bool {
	return strings.Contains(content, CrontabMarker(id))
}
