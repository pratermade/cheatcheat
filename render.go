package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Define styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(1).
			PaddingRight(1)

	headingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			PaddingBottom(1)

	syntaxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5AF78E"))

	optionFlagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F3A922"))

	optionDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD"))

	noteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A8A8A8")).
			Italic(true)

	selectedCommandStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#3C3836")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)

	normalCommandStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingRight(1)

	commandNumberStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#777777"))

	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89DDFF")).
			Italic(true)

	complexityStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B5E8B5"))

	codeBlockStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#282828")).
			Foreground(lipgloss.Color("#B8BB26")).
			Padding(0, 2)
)

// RenderCommandList renders a styled list of commands with the current selection highlighted
func RenderCommandList(sheet CheatSheet, selectedIdx int) string {
	var b strings.Builder

	// Description
	b.WriteString(sheet.Description)
	b.WriteString("\n\n")

	// Commands
	for i, cmd := range sheet.Commands {
		// Format the command number
		cmdNum := commandNumberStyle.Render(fmt.Sprintf("%d.", i+1))

		// Format the command name and description
		cmdText := fmt.Sprintf("%s %s - %s", cmdNum, cmd.Name, cmd.ShortDesc)

		// Apply the appropriate style based on whether this is the selected command
		var styledCmd string
		if i == selectedIdx {
			styledCmd = selectedCommandStyle.Render(cmdText)
		} else {
			styledCmd = normalCommandStyle.Render(cmdText)
		}

		b.WriteString(styledCmd)

		// Add tags if present
		if len(cmd.Tags) > 0 {
			tags := fmt.Sprintf(" [%s]", strings.Join(cmd.Tags, ", "))
			b.WriteString(tagStyle.Render(tags))
		}

		b.WriteString("\n\n")
	}

	return b.String()
}

// RenderCommandDetail renders styled details for a command with enhanced formatting
func RenderCommandDetail(cmd Command) string {
	var b strings.Builder

	// Description with nice styling
	b.WriteString(cmd.ShortDesc)
	b.WriteString("\n\n")

	// Syntax with nice code block styling
	b.WriteString("Syntax: ")
	b.WriteString("\n")
	b.WriteString(codeBlockStyle.Render(cmd.Syntax))
	b.WriteString("\n\n")

	// Complexity with color coding
	if cmd.Complexity != "" {
		complexity := cmd.Complexity
		complexityText := fmt.Sprintf("Complexity: %s", complexity)
		b.WriteString(complexityStyle.Render(complexityText))
		b.WriteString("\n\n")
	}

	// Tags
	if len(cmd.Tags) > 0 {
		tagText := fmt.Sprintf("Tags: %s", strings.Join(cmd.Tags, ", "))
		b.WriteString(tagStyle.Render(tagText))
		b.WriteString("\n\n")
	}

	// Options
	if len(cmd.Options) > 0 {
		b.WriteString(headingStyle.Render("Options:"))
		b.WriteString("\n")
		for _, opt := range cmd.Options {
			flagLine := fmt.Sprintf("  %s\n    %s",
				optionFlagStyle.Render(opt.Flag),
				optionDescStyle.Render(opt.Description))
			b.WriteString(flagLine)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Examples with code blocks
	if len(cmd.Examples) > 0 {
		b.WriteString(headingStyle.Render("Examples:"))
		b.WriteString("\n")
		for i, ex := range cmd.Examples {
			b.WriteString(fmt.Sprintf("  Example %d:\n", i+1))
			b.WriteString("  ")
			b.WriteString(codeBlockStyle.Render(fmt.Sprintf("$ %s", ex.Code)))
			b.WriteString("\n")
			b.WriteString(fmt.Sprintf("    %s", ex.Description))
			b.WriteString("\n\n")
		}
	}

	// Notes with bullet points
	if len(cmd.Notes) > 0 {
		b.WriteString(headingStyle.Render("Notes:"))
		b.WriteString("\n")
		for _, note := range cmd.Notes {
			noteLine := fmt.Sprintf("  • %s",
				noteStyle.Render(note))
			b.WriteString(noteLine)
			b.WriteString("\n")
		}
	}

	// Related commands
	if len(cmd.Related) > 0 {
		b.WriteString("\n")
		b.WriteString(headingStyle.Render("Related Commands:"))
		b.WriteString("\n")

		var relatedCommands []string
		for _, related := range cmd.Related {
			relatedCommands = append(relatedCommands,
				syntaxStyle.Render(related))
		}

		b.WriteString("  " + strings.Join(relatedCommands, ", "))
		b.WriteString("\n")
	}

	return b.String()
}
