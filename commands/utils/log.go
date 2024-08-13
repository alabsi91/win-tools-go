package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type logStyles struct {
	Success      lipgloss.Style
	Warning      lipgloss.Style
	Error        lipgloss.Style
	Info         lipgloss.Style
	Log          lipgloss.Style
	PaddingStyle lipgloss.Style
}

type log struct {
	Style logStyles
}

var Log = log{
	Style: logStyles{
		Success:      lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")),
		Warning:      lipgloss.NewStyle().Foreground(lipgloss.Color("#F1C40F")),
		Error:        lipgloss.NewStyle().Foreground(lipgloss.Color("#E74C3C")),
		Info:         lipgloss.NewStyle().Foreground(lipgloss.Color("#5eb9ff")),
		Log:          lipgloss.NewStyle(),
		PaddingStyle: lipgloss.NewStyle().Faint(true),
	},
}

const titleWidth = 13
const newlineSeparator = "- "
const newlinePaddingWidth = titleWidth + 3 // titleWidth + "|" + "|" + " "

// splitOnNewline splits a given string into leading newlines, content, and trailing newlines
func splitOnNewline(input string) (string, string, string) {
	// Check for leading newlines
	newlineStart := 0
	for newlineStart < len(input) && input[newlineStart] == '\n' {
		newlineStart++
	}

	// Check for trailing newlines
	newlineEnd := len(input)
	for newlineEnd > newlineStart && input[newlineEnd-1] == '\n' {
		newlineEnd--
	}

	// Return leading newlines, content, and trailing newlines
	return input[:newlineStart], input[newlineStart:newlineEnd], input[newlineEnd:]
}

// formatTitle formats a title with the given style
func formatTitle(title string, style lipgloss.Style) string {
	return style.Render("|") +
		lipgloss.NewStyle().Inherit(style).Reverse(true).Bold(true).Align(lipgloss.Center).Width(titleWidth).Render(title) +
		style.Render("|")
}

func (l log) printLog(title string, style lipgloss.Style, strs ...string) {
	title = formatTitle(title, style)

	fullString := strings.Join(strs, " ")
	leadingNewlines, msg, trailingNewlines := splitOnNewline(fullString)

	split := strings.Split(msg, "\n")

	// format new lines so that they are aligned
	var content string
	for i, s := range split {
		isFirst := i == 0
		if !isFirst {
			paddingString := strings.Repeat(newlineSeparator, newlinePaddingWidth/len(newlineSeparator))
			content += "\n" + l.Style.PaddingStyle.Render(paddingString)
		}
		content += style.Render(s)
	}

	fmt.Println(leadingNewlines+title, content+trailingNewlines)
}

func (l log) Success(strs ...string) {
	l.printLog("SUCCESS", l.Style.Success, strs...)
}

func (l log) Error(strs ...string) {
	l.printLog("ERROR", l.Style.Error, strs...)
}

// FATAL logs a fatal error and exits the program
func (l log) Fatal(strs ...string) {
	l.printLog("FATAL", l.Style.Error, strs...)
	os.Exit(1)
}

func (l log) Info(strs ...string) {
	l.printLog("INFO", l.Style.Info, strs...)
}

func (l log) Warning(strs ...string) {
	l.printLog("WARNING", l.Style.Warning, strs...)
}

func (l log) Log(strs ...string) {
	l.printLog("LOG", l.Style.Log, strs...)
}
