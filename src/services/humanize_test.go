package services

import "testing"

func TestHumanizeDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{
			name:     "zero seconds",
			seconds:  0,
			expected: "0 seconds",
		},
		{
			name:     "negative seconds (clamped to zero)",
			seconds:  -10,
			expected: "0 seconds",
		},
		{
			name:     "one second",
			seconds:  1,
			expected: "1 second",
		},
		{
			name:     "multiple seconds",
			seconds:  45,
			expected: "45 seconds",
		},
		{
			name:     "one minute",
			seconds:  60,
			expected: "1 minute",
		},
		{
			name:     "one minute and seconds",
			seconds:  65,
			expected: "1 minute, 5 seconds",
		},
		{
			name:     "multiple minutes",
			seconds:  125,
			expected: "2 minutes, 5 seconds",
		},
		{
			name:     "one hour",
			seconds:  3600,
			expected: "1 hour",
		},
		{
			name:     "one hour, one minute, one second",
			seconds:  3661,
			expected: "1 hour, 1 minute, 1 second",
		},
		{
			name:     "multiple hours",
			seconds:  7265,
			expected: "2 hours, 1 minute, 5 seconds",
		},
		{
			name:     "one day",
			seconds:  86400,
			expected: "1 day",
		},
		{
			name:     "one day, one hour, one minute, one second",
			seconds:  90061,
			expected: "1 day, 1 hour, 1 minute, 1 second",
		},
		{
			name:     "multiple days",
			seconds:  172805,
			expected: "2 days, 5 seconds",
		},
		{
			name:     "fractional seconds (rounds)",
			seconds:  65.7,
			expected: "1 minute, 6 seconds",
		},
		{
			name:     "exact minutes (no seconds component)",
			seconds:  120,
			expected: "2 minutes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HumanizeDuration(tt.seconds)
			if result != tt.expected {
				t.Errorf("HumanizeDuration(%v) = %v, want %v", tt.seconds, result, tt.expected)
			}
		})
	}
}
