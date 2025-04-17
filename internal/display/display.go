package display

import (
	"strings"
)

// Constants for display formatting
const (
	maxNameLength    = 80
	maxTagsLength    = 60
	maxMembersLength = 60
	formatJSON       = "json"
	formatText       = "text"
	iconFixedWidth   = 4 // Visual width allocated for icon characters
)

// Color constants for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// GetPaddedIcon formats icon characters to a fixed width
// This ensures consistent spacing regardless of emoji or Unicode character width
func GetPaddedIcon(icon string, width int) string {
	if icon == "" {
		return strings.Repeat(" ", width) // Return empty space of desired width
	}

	// For safety, return the icon without padding if width is invalid
	if width <= 0 {
		return icon
	}

	// For complex emoji handling, use a more robust approach
	iconRunes := []rune(icon)
	iconLength := len(iconRunes)

	// Make sure padding is never negative
	padSize := width - iconLength
	if padSize < 0 {
		padSize = 0
	}

	// Create padding string safely
	padding := ""
	if padSize > 0 {
		padding = strings.Repeat(" ", padSize)
	}

	return icon + padding
}
