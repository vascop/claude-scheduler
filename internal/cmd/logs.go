package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vascop/claude-scheduler/internal/task"
)

func Logs(args []string) error {
	id := ""
	lines := 30

	for i := 0; i < len(args); i++ {
		if args[i] == "-n" && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				return fmt.Errorf("invalid line count %q", args[i+1])
			}
			lines = n
			i++
		} else if id == "" {
			id = args[i]
		}
	}

	if id != "" {
		return showLog(task.LogPath(id), lines)
	}

	// Show all logs
	entries, err := filepath.Glob(filepath.Join(task.LogsDir(), "*.log"))
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("No logs found.")
		return nil
	}
	for _, path := range entries {
		taskID := strings.TrimSuffix(filepath.Base(path), ".log")
		fmt.Printf("--- %s ---\n", taskID)
		showLog(path, 5)
		fmt.Println()
	}
	return nil
}

func showLog(path string, n int) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No logs at %s\n", path)
			return nil
		}
		return err
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	start := len(lines) - n
	if start < 0 {
		start = 0
	}
	for _, line := range lines[start:] {
		fmt.Println(line)
	}
	return nil
}
