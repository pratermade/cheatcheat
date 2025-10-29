package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

// Define styles
var (
	boxedViewportStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).         // Choose from NormalBorder(), RoundedBorder(), DoubleBorder(), etc.
				BorderForeground(lipgloss.Color("#69C")). // Border color
				Padding(0, 1).                            // Inner padding (vertical, horizontal)
				BorderTop(true).
				BorderLeft(true).
				BorderRight(true).
				BorderBottom(true)

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

func RenderTagMenu(tags []string, selectedIndex int, termWidth int) string {

	// Create a menu bar style
	menuBarStyle := lipgloss.NewStyle().
		Width(termWidth-4).
		Padding(0, 1)

	// Get the visible tags with proper styling and indicators
	visibleTags := GetVisibleTags(tags, selectedIndex, termWidth)

	// Join them horizontally
	menu := lipgloss.JoinHorizontal(lipgloss.Top, visibleTags...)

	logrus.Debugf("Rendered tag menu: %s (selected: %d)", menu, selectedIndex)
	return menuBarStyle.Render(menu)
}

// GetVisibleTags returns the visible tags for a menu based on the selected index and available width
// It returns the slice of styled tag strings (including any scroll indicators)
func GetVisibleTags(tags []string, selectedIndex, termWidth int) []string {
	// Ensure the selected index is within bounds
	if selectedIndex < 0 {
		selectedIndex = 0
	} else if selectedIndex >= len(tags) {
		selectedIndex = len(tags) - 1
	}

	// Calculate the total available width for tags
	availableWidth := termWidth - 8 // Subtract some padding for the outer container

	// We'll build our result as we go
	var result []string
	var currentWidth int

	// We'll track which tags we've added
	addedTags := make(map[int]bool)

	// Always include the selected tag first to ensure it's visible
	selectedTagStyle := selectedCommandStyle.MarginRight(1)
	selectedTag := selectedTagStyle.Render(tags[selectedIndex])
	selectedTagWidth := lipgloss.Width(selectedTag)

	// Make sure we have room for the selected tag
	if currentWidth+selectedTagWidth <= availableWidth {
		result = append(result, selectedTag)
		currentWidth += selectedTagWidth
		addedTags[selectedIndex] = true
	} else {
		// If we can't even fit the selected tag, just return it alone
		return []string{selectedTag}
	}

	// Now add tags to the left of the selected tag, starting from the closest one
	for i := selectedIndex - 1; i >= 0; i-- {
		normalTagStyle := normalCommandStyle.MarginRight(1)
		styledTag := normalTagStyle.Render(tags[i])
		tagWidth := lipgloss.Width(styledTag)

		if currentWidth+tagWidth > availableWidth {
			break
		}

		// Insert at beginning to maintain left-to-right order
		result = append([]string{styledTag}, result...)
		currentWidth += tagWidth
		addedTags[i] = true
	}

	// Then add tags to the right of the selected tag
	for i := selectedIndex + 1; i < len(tags); i++ {
		normalTagStyle := normalCommandStyle.MarginRight(1)
		styledTag := normalTagStyle.Render(tags[i])
		tagWidth := lipgloss.Width(styledTag)

		if currentWidth+tagWidth > availableWidth {
			break
		}

		result = append(result, styledTag)
		currentWidth += tagWidth
		addedTags[i] = true
	}

	// Add left scroll indicator if needed - only if the first tag isn't shown
	showLeftIndicator := !addedTags[0]
	if showLeftIndicator && len(result) > 0 {
		leftIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")).
			Render("« ")

		// Insert at the beginning
		result = append([]string{leftIndicator}, result...)
	}

	// Add right scroll indicator if needed - only if the last tag isn't shown
	showRightIndicator := !addedTags[len(tags)-1]
	if showRightIndicator && len(result) > 0 {
		rightIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")).
			Render(" »")

		result = append(result, rightIndicator)
	}

	return result
}

// RenderCommandList renders a styled list of commands with the current selection highlighted
func RenderCommandList(description string, commands []Command, selectedIdx int) string {
	var b strings.Builder

	// Description
	b.WriteString(description)
	b.WriteString("\n\n")

	// Commands
	for i, cmd := range commands {
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

// RenderSearchBar renders the search input bar
func RenderSearchBar(query string, width int) string {
	searchStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5AF78E")).
		Bold(true)

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5AF78E"))

	// Create the search prompt with query and cursor
	searchPrompt := searchStyle.Render("Search: ") + query + cursorStyle.Render(" ")

	return searchPrompt
}

// function to append to a debug log file
func Debug(entry string) {
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	if _, err := f.WriteString(entry + "\n"); err != nil {
		return
	}
}
