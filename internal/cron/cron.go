package cron

import (
	"fmt"
	"strconv"
	"strings"
)

// Field represents a single cron field. A nil Values slice means wildcard.
type Field struct {
	Values []int
}

// IsWildcard returns true if this field matches all values.
func (f Field) IsWildcard() bool {
	return f.Values == nil
}

// Schedule represents a parsed 5-field cron expression.
type Schedule struct {
	Minute     Field
	Hour       Field
	DayOfMonth Field
	Month      Field
	DayOfWeek  Field
}

// Parse parses a 5-field cron expression (minute hour dom month dow).
func Parse(expr string) (*Schedule, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	minute, err := parseField(fields[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("minute: %w", err)
	}
	hour, err := parseField(fields[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("hour: %w", err)
	}
	dom, err := parseField(fields[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("day-of-month: %w", err)
	}
	month, err := parseField(fields[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("month: %w", err)
	}
	dow, err := parseField(fields[4], 0, 6)
	if err != nil {
		return nil, fmt.Errorf("day-of-week: %w", err)
	}

	return &Schedule{
		Minute:     minute,
		Hour:       hour,
		DayOfMonth: dom,
		Month:      month,
		DayOfWeek:  dow,
	}, nil
}

func parseField(s string, min, max int) (Field, error) {
	if s == "*" {
		return Field{}, nil
	}

	// Step: */N
	if strings.HasPrefix(s, "*/") {
		step, err := strconv.Atoi(s[2:])
		if err != nil || step <= 0 {
			return Field{}, fmt.Errorf("invalid step %q", s)
		}
		var vals []int
		for i := min; i <= max; i += step {
			vals = append(vals, i)
		}
		return Field{Values: vals}, nil
	}

	// Comma-separated list (each element can be a range)
	var vals []int
	for _, part := range strings.Split(s, ",") {
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			lo, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return Field{}, fmt.Errorf("invalid range %q", part)
			}
			hi, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return Field{}, fmt.Errorf("invalid range %q", part)
			}
			if lo < min || hi > max || lo > hi {
				return Field{}, fmt.Errorf("range %d-%d out of bounds [%d,%d]", lo, hi, min, max)
			}
			for i := lo; i <= hi; i++ {
				vals = append(vals, i)
			}
		} else {
			v, err := strconv.Atoi(part)
			if err != nil {
				return Field{}, fmt.Errorf("invalid value %q", part)
			}
			if v < min || v > max {
				return Field{}, fmt.Errorf("value %d out of bounds [%d,%d]", v, min, max)
			}
			vals = append(vals, v)
		}
	}
	return Field{Values: vals}, nil
}
