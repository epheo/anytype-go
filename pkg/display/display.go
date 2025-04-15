package display

// Output formats
const (
	formatJSON = "json"
	formatText = "text"
)

// Constants for text formatting
const (
	// Color escape sequences
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBlue  = "\033[34m"

	// Maximum lengths for truncation
	maxMembersLength = 50
	maxNameLength    = 37
	maxTagsLength    = 30
)
