package task

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const dataDir = ".claude-scheduler"

func TasksDir() string {
	return filepath.Join(os.Getenv("HOME"), dataDir, "tasks")
}

func LogsDir() string {
	return filepath.Join(os.Getenv("HOME"), dataDir, "logs")
}

func LogPath(id string) string {
	return filepath.Join(LogsDir(), id+".log")
}

func ScriptPath(id string) string {
	return filepath.Join(TasksDir(), id+".sh")
}

func Save(t *Task) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(TasksDir(), t.ID+".json"), data, 0644)
}

func Load(id string) (*Task, error) {
	data, err := os.ReadFile(filepath.Join(TasksDir(), id+".json"))
	if err != nil {
		return nil, err
	}
	var t Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func Delete(id string) error {
	os.Remove(filepath.Join(TasksDir(), id+".json"))
	os.Remove(ScriptPath(id))
	return nil
}

func ListAll() ([]*Task, error) {
	entries, err := filepath.Glob(filepath.Join(TasksDir(), "*.json"))
	if err != nil {
		return nil, err
	}
	var tasks []*Task
	for _, path := range entries {
		id := strings.TrimSuffix(filepath.Base(path), ".json")
		t, err := Load(id)
		if err != nil {
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
