//go:build darwin

package platform

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/vascop/claude-scheduler/internal/cron"
)

const plistPrefix = "com.claude.scheduler"

// Launchd implements Scheduler using macOS launchd.
type Launchd struct{}

func NewScheduler() Scheduler {
	return &Launchd{}
}

func (l *Launchd) Name() string { return "launchd" }

func (l *Launchd) label(id string) string {
	return plistPrefix + "." + id
}

func (l *Launchd) plistPath(id string) string {
	return filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", l.label(id)+".plist")
}

func (l *Launchd) Describe(id string) string {
	return l.plistPath(id)
}

func (l *Launchd) Install(id string, cronExpr string, schedule *cron.Schedule, scriptPath, logPath string) error {
	content, err := l.generate(id, schedule, scriptPath, logPath)
	if err != nil {
		return err
	}

	path := l.plistPath(id)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return err
	}
	exec.Command("launchctl", "unload", path).Run()
	return exec.Command("launchctl", "load", path).Run()
}

func (l *Launchd) Uninstall(id string) error {
	path := l.plistPath(id)
	// Unload from launchd (ignore error — may not be loaded)
	exec.Command("launchctl", "unload", path).Run()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (l *Launchd) IsInstalled(id string) bool {
	return exec.Command("launchctl", "list", l.label(id)).Run() == nil
}

var plistTmpl = template.Must(template.New("plist").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>{{.Label}}</string>
  <key>ProgramArguments</key>
  <array>
    <string>/bin/bash</string>
    <string>{{.ScriptPath}}</string>
  </array>
  <key>StartCalendarInterval</key>
  <array>
{{.CalendarIntervals}}  </array>
  <key>StandardOutPath</key>
  <string>{{.LogPath}}</string>
  <key>StandardErrorPath</key>
  <string>{{.LogPath}}</string>
  <key>EnvironmentVariables</key>
  <dict>
    <key>PATH</key>
    <string>{{.EnvPath}}</string>
  </dict>
</dict>
</plist>
`))

type plistData struct {
	Label             string
	ScriptPath        string
	LogPath           string
	CalendarIntervals string
	EnvPath           string
}

func (l *Launchd) generate(id string, schedule *cron.Schedule, scriptPath, logPath string) ([]byte, error) {
	intervals := buildCalendarIntervals(schedule)
	envPath := fmt.Sprintf("/usr/local/bin:/usr/bin:/bin:/opt/homebrew/bin:%s/.local/bin", os.Getenv("HOME"))

	var buf bytes.Buffer
	err := plistTmpl.Execute(&buf, plistData{
		Label:             l.label(id),
		ScriptPath:        scriptPath,
		LogPath:           logPath,
		CalendarIntervals: intervals,
		EnvPath:           envPath,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func buildCalendarIntervals(s *cron.Schedule) string {
	type entry struct {
		key string
		val int
	}

	fields := []struct {
		key    string
		values []int
	}{
		{"Minute", fieldValues(s.Minute)},
		{"Hour", fieldValues(s.Hour)},
		{"Day", fieldValues(s.DayOfMonth)},
		{"Month", fieldValues(s.Month)},
		{"Weekday", fieldValues(s.DayOfWeek)},
	}

	type dict []entry
	combos := []dict{{}}

	for _, f := range fields {
		if f.values == nil {
			continue
		}
		var next []dict
		for _, combo := range combos {
			for _, v := range f.values {
				newCombo := make(dict, len(combo), len(combo)+1)
				copy(newCombo, combo)
				newCombo = append(newCombo, entry{f.key, v})
				next = append(next, newCombo)
			}
		}
		combos = next
	}

	var buf bytes.Buffer
	for _, combo := range combos {
		buf.WriteString("    <dict>\n")
		for _, e := range combo {
			fmt.Fprintf(&buf, "      <key>%s</key>\n      <integer>%d</integer>\n", e.key, e.val)
		}
		buf.WriteString("    </dict>\n")
	}
	return buf.String()
}

func fieldValues(f cron.Field) []int {
	if f.IsWildcard() {
		return nil
	}
	return f.Values
}
