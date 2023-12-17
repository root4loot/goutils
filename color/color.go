package color

import "fmt"

const (
	Reset     = "\033[0m"
	LightGrey = "\033[90m"
	Red       = "\033[91m"
	Yellow    = "\033[93m"
	Blue      = "\033[94m"
	White     = "\033[37m"
	Green     = "\033[92m"
	Purple    = "\033[95m"
	Cyan      = "\033[96m"
	Orange    = "\033[38;5;214m"
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
