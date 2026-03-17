//go:build linux

package platform

import (
	"os/exec"
	"strings"

	"github.com/vascop/claude-scheduler/internal/cron"
)

// Crontab implements Scheduler using Linux crontab.
type Crontab struct{}

func NewScheduler() Scheduler {
	return &Crontab{}
}

func (c *Crontab) Name() string { return "crontab" }

func (c *Crontab) Describe(id string) string {
	return "crontab entry (marker: claude-scheduler:" + id + ")"
}

func (c *Crontab) Install(id string, cronExpr string, _ *cron.Schedule, scriptPath, logPath string) error {
	existing, _ := c.readCrontab()
	newCrontab := CrontabAdd(existing, id, cronExpr, scriptPath, logPath)
	return c.writeCrontab(newCrontab)
}

func (c *Crontab) Uninstall(id string) error {
	existing, err := c.readCrontab()
	if err != nil || existing == "" {
		return nil
	}

	newCrontab := CrontabRemove(existing, id)
	if strings.TrimSpace(newCrontab) == "" {
		return exec.Command("crontab", "-r").Run()
	}
	return c.writeCrontab(newCrontab)
}

func (c *Crontab) IsInstalled(id string) bool {
	existing, err := c.readCrontab()
	if err != nil {
		return false
	}
	return CrontabHasTask(existing, id)
}

func (c *Crontab) readCrontab() (string, error) {
	out, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		return "", nil
	}
	return string(out), nil
}

func (c *Crontab) writeCrontab(content string) error {
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(content)
	return cmd.Run()
}
