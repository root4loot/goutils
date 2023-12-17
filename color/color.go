package color

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

// Colorize applies the specified color to a string and then resets it.
func Colorize(color, text string) string {
	return color + text + Reset
}
