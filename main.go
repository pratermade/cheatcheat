package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// Parse command-line arguments
	flag.Parse()
	args := flag.Args()

	// Check if a filepath was provided
	if len(args) < 1 {
		fmt.Println("Please provide a path to a cheat sheet YAML file")
		os.Exit(1)
	}

	// Create and run the Bubble Tea program
	p := tea.NewProgram(initialModel(args[0]), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func initialModel(filePath string) model {
	// Initial viewport
	vp := viewport.New(80, 24)
	vp.SetContent("Loading cheat sheet...")

	// Create initial model
	m := model{
		viewport:       vp,
		currentCommand: 0,
		showDetail:     false,
	}

	// Load the cheat sheet in the Init function
	return m
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return loadCheatSheetMsg(flag.Arg(0))
	}
}

// Custom message types
type errorMsg struct{ err error }
type cheatSheetLoadedMsg CheatSheet

func (e errorMsg) Error() string { return e.err.Error() }

// Command to load a cheat sheet
func loadCheatSheetMsg(filePath string) tea.Msg {
	sheet, err := LoadCheatSheet(filePath)
	if err != nil {
		return errorMsg{err}
	}
	return cheatSheetLoadedMsg(sheet)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Enter):
			if !m.showDetail && len(m.commands) > 0 {
				// Show detail view of the selected command
				m.showDetail = true
				content := RenderCommandDetail(m.commands[m.currentCommand])
				m.viewport.SetContent(content)
				m.viewport.GotoTop()
			}

		case key.Matches(msg, keys.Back):
			if m.showDetail {
				// Go back to command list view
				m.showDetail = false
				content := RenderCommandList(m.cheatSheet, m.currentCommand)
				m.viewport.SetContent(content)
				m.viewport.GotoTop()
			}

		case key.Matches(msg, keys.Up):
			if !m.showDetail {
				// Navigate up in the command list
				if m.currentCommand > 0 {
					m.currentCommand--
					// Update the content to reflect the new selection
					content := RenderCommandList(m.cheatSheet, m.currentCommand)
					m.viewport.SetContent(content)
				}
			} else {
				// Scroll up in the viewport
				m.viewport.LineUp(1)
			}

		case key.Matches(msg, keys.Down):
			if !m.showDetail {
				// Navigate down in the command list
				if m.currentCommand < len(m.commands)-1 {
					m.currentCommand++
					// Update the content to reflect the new selection
					content := RenderCommandList(m.cheatSheet, m.currentCommand)
					m.viewport.SetContent(content)
				}
			} else {
				// Scroll down in the viewport
				m.viewport.LineDown(1)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 4 // Reserve space for header/footer
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight

	case cheatSheetLoadedMsg:
		// Handle the loaded cheat sheet
		m.cheatSheet = CheatSheet(msg)
		m.commands = m.cheatSheet.Commands

		// Update the view with the command list
		if len(m.commands) > 0 {
			content := RenderCommandList(m.cheatSheet, m.currentCommand)
			m.viewport.SetContent(content)
		} else {
			m.viewport.SetContent("No commands found in the cheat sheet.")
		}

	case errorMsg:
		// Handle errors
		m.err = msg.err
		m.viewport.SetContent(fmt.Sprintf("Error: %v", msg.err))
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}
func (m model) View() string {
	// If there's an error, display it
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)
	}

	// Create a footer with help text
	helpView := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("↑/↓: Navigate • Enter: View details • Esc: Back • q: Quit")

	// Create a header
	var header string
	if m.showDetail && len(m.commands) > 0 {
		header = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1).
			Render(m.commands[m.currentCommand].Name)
	} else {
		header = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#2D9CDB")).
			Padding(0, 1).
			Render(m.cheatSheet.Title)
	}

	// Combine the parts
	return fmt.Sprintf("%s\n\n%s\n\n%s", header, m.viewport.View(), helpView)
}
