package helper

import (
	"fmt"
	"strings"

	"github.com/olekukonko/ts"
)

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func FormatMessageBody(message string) string {
	var result strings.Builder
	words := strings.Fields(message)
	currentLine := ""

	size, err := ts.GetSize()
	if err != nil {
		return message
	}

	maxWidth := size.Col() - 20

	for _, word := range words {
		// If adding the word would exceed max width, start a new line
		if len(currentLine)+len(word)+1 > maxWidth {
			// Trim trailing space and add current line to result
			result.WriteString(strings.TrimSpace(currentLine) + "\n ")
			currentLine = ""
		}

		// Add word to current line
		if currentLine != "" {
			currentLine += " "
		}
		currentLine += word
	}

	// Add final line
	if currentLine != "" {
		result.WriteString(currentLine)
	}

	return result.String()
}
