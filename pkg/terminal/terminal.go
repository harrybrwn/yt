package terminal

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
)

var (
	// Output is the default output for the terminal package
	Output io.Writer

	controlChar string
)

func init() {
	switch runtime.GOOS {
	case "windows":
		Output = ioutil.Discard
		controlChar = ""
	default:
		Output = os.Stdout
		controlChar = "\033"
	}
}

// CursorOn turns the cursor on
func CursorOn() { out(control("?25h")) }

// CursorOff turns the cursor off
func CursorOff() { out(control("?25l")) }

// Blue returns s but colored blue
func Blue(s string) string { return color(36, s) }

// Green returns s but colored green
func Green(s string) string { return color(32, s) }

// Red returns s but colored red
func Red(s string) string { return color(31, s) }

func color(color int, s string) string {
	if controlChar == "" {
		return s
	}
	return fmt.Sprintf("%[1]s[%dm%s%[1]s[0m", controlChar, color, s)
}

func out(s string) {
	fmt.Fprint(Output, s)
}

func control(code string) string {
	if controlChar == "" {
		return ""
	}
	return fmt.Sprintf("%s[%s", controlChar, code)
}
