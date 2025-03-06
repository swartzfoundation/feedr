package model

import (
	"strings"
	"testing"
)

func TestNewID(t *testing.T) {
	expectedLen := 9
	id := NewID()
	if len(id) != expectedLen {
		t.Errorf("NewID length error expected %d, got %d", expectedLen, len(id))
	}
}

func TestSentenceCase(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		// the table itself
		{"table-driven test", "Table-Driven Test"},
		{"capitalization rules are explained in more detail in the next section", "Capitalization Rules Are Explained In More Detail In The Next Section"},
		{"The following section provides an overview of the title case rules.", "The Following Section Provides An Overview Of The Title Case Rules."},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToTitleCase(tt.input)
			if isEqual := strings.Compare(got, tt.want); isEqual != 0 {
				t.Errorf("got %s expected %s", got, tt.want)
			}
		})
	}
}

func TestPasswordCompare(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"increase your rounds", "$2a$12$AwZfRBJqsSM9MljlOYYFWO1d2O6MNGxD8iRWm/wbvVxdH2YZxgvNy"},
		{"hashing function", "$2a$12$/cm.1QQeEH/agAqaOo6uTuEW5nFVYEKAKe8JEOQ8ZKiLCCUteSrkK"},
		{"secure?", "$2a$12$dxQJgH2HiIi2b2zkKHOwxuM7ksAUQscx/WmcFMf8jB5ux/7G7EnZm"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			hash, err := HashPassword(tt.input)
			if err != nil {
				t.Errorf("password hashing fail for %s", tt.input)
			}
			if pass := CheckPassword(hash, tt.input); !pass {
				t.Errorf("password compare for %s failed", tt.input)
			}
		})
	}
}
