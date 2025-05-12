package color

import "fmt"

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

// GetColorForUsername returns a consistent color for a given username
func GetColorForUsername(username string) string {
	// Simple hash function to get a consistent color for a username
	var hash int
	for _, char := range username {
		hash = 31*hash + int(char)
	}

	// Use the hash to select a color
	colors := []string{Red, Green, Yellow, Blue, Purple, Cyan}
	return colors[hash%len(colors)]
}

// ColorizeUsername returns the username in its assigned color
func ColorizeUsername(username string) string {
	return GetColorForUsername(username) + username + Reset
}

// ColorizeMessage returns the message with the username in color
func ColorizeMessage(username, content string) string {
	return fmt.Sprintf("%s: %s", ColorizeUsername(username), content)
}
