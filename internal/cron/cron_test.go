package cron

import (
	"reflect"
	"testing"
)

func TestParseField(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		min     int
		max     int
		want    []int
		wantErr bool
	}{
		{"wildcard", "*", 0, 59, nil, false},
		{"step", "*/15", 0, 59, []int{0, 15, 30, 45}, false},
		{"single", "9", 0, 23, []int{9}, false},
		{"comma", "1,3,5", 0, 6, []int{1, 3, 5}, false},
		{"range", "1-5", 0, 6, []int{1, 2, 3, 4, 5}, false},
		{"out of range", "25", 0, 23, nil, true},
		{"bad step", "*/0", 0, 59, nil, true},
		{"bad syntax", "abc", 0, 59, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseField(tt.input, tt.min, tt.max)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseField(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got.Values, tt.want) {
				t.Errorf("parseField(%q) = %v, want %v", tt.input, got.Values, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{"every day at 9:03", "3 9 * * *", false},
		{"weekdays at 9:03", "3 9 * * 1-5", false},
		{"every 4 hours", "17 */4 * * *", false},
		{"too few fields", "3 9 *", true},
		{"too many fields", "3 9 * * * *", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
			}
		})
	}
}
