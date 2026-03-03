package service

import (
	"fmt"
	"strings"
	"time"
)

// normalizeTime ensures time is stored as HH:MM
func normalizeTime(input string) (string, error) {
	// Trim spaces
	input = strings.TrimSpace(input)

	// Try multiple layouts to be flexible with input formats
	layouts := []string{
		"15:04",
		"15:04:05",
		"3 PM",
		"3:04 PM",
		"3:04PM",
		"15",
	}

	var parsed time.Time
	var err error

	// Try parsing with each layout
	for _, layout := range layouts {
		parsed, err = time.Parse(layout, input)
		if err == nil {
			return parsed.Format("15:04"), nil
		}
	}

	//
	return "", fmt.Errorf("invalid time format: %s", input)
}