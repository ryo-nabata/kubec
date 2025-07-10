package utils

import (
	"fmt"
	"github.com/fatih/color"
)

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