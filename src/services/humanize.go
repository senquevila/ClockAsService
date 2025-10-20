package services

import (
	"fmt"
	"math"
)

// HumanizeDuration converts seconds to a human-readable string
// Examples:
//   - 65 seconds -> "1 minute, 5 seconds"
//   - 3661 seconds -> "1 hour, 1 minute, 1 second"
//   - 90061 seconds -> "1 day, 1 hour, 1 minute, 1 second"
func HumanizeDuration(seconds float64) string {
	if seconds < 0 {
		seconds = 0
	}

	totalSeconds := int(math.Round(seconds))

	if totalSeconds == 0 {
		return "0 seconds"
	}

	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	secs := totalSeconds % 60

	parts := []string{}

	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}

	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}

	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	if secs > 0 {
		if secs == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", secs))
		}
	}

	// Join parts with ", "
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += ", "
		}
		result += part
	}

	return result
}
