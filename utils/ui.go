package utils

import (
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

// noBellStdout wraps readline's stdout and discards bell characters (\a)
// so that promptui.Select does not trigger the terminal beep on every
// keystroke (see manifoldco/promptui#49).
type noBellStdout struct {
	w io.Writer
}

func (n *noBellStdout) Write(p []byte) (int, error) {
	if len(p) == 1 && p[0] == readline.CharBell {
		return len(p), nil
	}
	return n.w.Write(p)
}

func (n *noBellStdout) Close() error {
	return nil
}

// NoBellStdout is passed to promptui.Select.Stdout to suppress terminal bells.
var NoBellStdout io.WriteCloser = &noBellStdout{w: readline.Stdout}

func PrintSuccess(message string) {
	fmt.Printf("✓ %s\n", color.GreenString(message))
}

func PrintError(message string) {
	fmt.Printf("✗ %s\n", color.RedString(message))
}

func PrintInfo(message string) {
	fmt.Printf("ℹ %s\n", color.BlueString(message))
}

func PrintWarning(message string) {
	fmt.Printf("⚠ %s\n", color.YellowString(message))
}

func PrintHeader(title string) {
	fmt.Printf("\n%s\n", color.CyanString(title))
	fmt.Printf("%s\n", color.CyanString(generateDivider(len(title))))
}

func generateDivider(length int) string {
	divider := ""
	for i := 0; i < length; i++ {
		divider += "="
	}
	return divider
}

func HighlightText(text string) string {
	return color.YellowString(text)
}

func SuccessText(text string) string {
	return color.GreenString(text)
}

func ErrorText(text string) string {
	return color.RedString(text)
}

func InfoText(text string) string {
	return color.BlueString(text)
}