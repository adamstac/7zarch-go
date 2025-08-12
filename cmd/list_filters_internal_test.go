package cmd

import (
	"testing"
)

func Test_parseHumanDuration(t *testing.T) {
	cases := []struct {
		in        string
		wantHours int64
	}{
		{"1h", 1},
		{"30m", 0},
		{"24h", 24},
		{"1d", 24},
		{"2d", 48},
		{"1w", 168},
	}
	for _, c := range cases {
		d, err := parseHumanDuration(c.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", c.in, err)
		}
		if int64(d.Hours()) != c.wantHours {
			t.Fatalf("for %q got %d hours, want %d", c.in, int64(d.Hours()), c.wantHours)
		}
	}
}
