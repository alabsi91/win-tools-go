package utils

import (
	"fmt"
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
		Info:         lipgloss.NewStyle().Foreground(lipgloss.Color("#3498DB")),
		Log:          lipgloss.NewStyle(),
		PaddingStyle: lipgloss.NewStyle().Faint(true),
	},
}

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

func (l log) printLog(title string, style lipgloss.Style, strs ...string) {
	formattedTitle := style.Render("|") +
		lipgloss.NewStyle().Inherit(style).Reverse(true).Bold(true).Align(lipgloss.Center).Width(13).Render(title) +
		style.Render("|")

	fullString := strings.Join(strs, " ")
	leadingNewlines, msg, trailingNewlines := splitOnNewline(fullString)

	split := strings.Split(msg, "\n")

	var res string
	for i, s := range split {
		isFirst := i == 0

		if !isFirst {
			res += "\n" + l.Style.PaddingStyle.Render(strings.Repeat("- ", 8))
		}

		res += style.Render(s)
	}

	fmt.Println(leadingNewlines+formattedTitle, res+trailingNewlines)
}

func (l log) Success(strs ...string) {
	l.printLog("SUCCESS", l.Style.Success, strs...)
}

func (l log) Error(strs ...string) {
	l.printLog("ERROR", l.Style.Error, strs...)
}

func (l log) Fatal(strs ...string) {
	l.printLog("FATAL", l.Style.Error, strs...)
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
