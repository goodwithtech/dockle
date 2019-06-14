package color

import "fmt"

// Foreground colors.
const (
	Red Color = iota + 31
	Green
	Yellow
	Blue
	Magenta
	Cyan
)

// Color represents a text color.
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
