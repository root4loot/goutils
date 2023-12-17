package color

import "fmt"

const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	LightGrey = "\033[90m"
	Orange    = "\033[38;5;214m"
	// Add any other colors you need
)

// Colorize formats the text with the specified color and then resets it.
func Colorize(color string, a ...interface{}) string {
	var result string
	for i, item := range a {
		if i > 0 {
			result += " " // Add a space only between items
		}
		result += fmt.Sprint(item)
	}
	return color + result + Reset
}
