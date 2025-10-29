package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/sirupsen/logrus"
)

func main() {
	startLogging()
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
	width, _, _ := term.GetSize(os.Stdout.Fd())
	vp := viewport.New(width, 24)
	vp.SetContent("Loading cheat sheet...")
	tagVp := viewport.New(width, 3)
	tagVp.SetContent("Loading tags...")
	// Create initial model
	m := model{
		tagViewPort:    tagVp,
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

// Return a list of unique tags from commands
func UniqueTags(commands []Command) []string {
	tagSet := make(map[string]struct{})
	for _, cmd := range commands {
		for _, tag := range cmd.Tags {
			tagSet[tag] = struct{}{}
		}
	}
	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i] < tags[j]
	})
	all := []string{"all"}
	tags = append(all, tags...)
	return tags
}

func filterCommandsByTag(commands []Command, tag string) []Command {
    if tag == "all" {
        return commands
    }

    var filtered []Command
    for _, cmd := range commands {
        for _, t := range cmd.Tags {
            if t == tag {
                filtered = append(filtered, cmd)
                break
            }
        }
    }
    return filtered
}

func filterCommandsBySearch(commands []Command, query string) []Command {
	if query == "" {
		return commands
	}

	var filtered []Command
	lowerQuery := strings.ToLower(query)
	for _, cmd := range commands {
		if strings.Contains(strings.ToLower(cmd.Name), lowerQuery) {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode input
		if m.searchMode {
			switch {
			case key.Matches(msg, keys.Enter):
				// Exit search mode (filtering already done live)
				m.searchMode = false
				m.searchActive = true
				return m, nil
			case key.Matches(msg, keys.Back):
				// Cancel search
				m.searchMode = false
				m.searchQuery = ""
				return m, nil
			case msg.Type == tea.KeyBackspace:
				// Remove last character from search query and filter live
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
				m.commands = filterCommandsBySearch(m.cheatSheet.Commands, m.searchQuery)
				m.currentCommand = 0
				content := RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
				m.viewport.SetContent(content)
				return m, nil
			case msg.Type == tea.KeyRunes:
				// Add character to search query and filter live
				m.searchQuery += string(msg.Runes)
				m.commands = filterCommandsBySearch(m.cheatSheet.Commands, m.searchQuery)
				m.currentCommand = 0
				content := RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
				m.viewport.SetContent(content)
				return m, nil
			}
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Search):
			if !m.showDetail {
				// Enter search mode
				m.searchMode = true
				m.searchQuery = ""
				return m, nil
			}

		case key.Matches(msg, keys.Enter):
			if !m.showDetail && len(m.commands) > 0 {
				// Show detail view of the selected command
				m.showDetail = true
				content := RenderCommandDetail(m.commands[m.currentCommand])
				m.viewport.SetContent(content)
				m.viewport.GotoTop()
			}

		case key.Matches(msg, keys.Back):
			if m.searchActive {
				// Clear search filter and restore full list
				m.searchActive = false
				m.searchQuery = ""
				m.commands = filterCommandsByTag(m.cheatSheet.Commands, m.tagMenu[m.currentTag])
				m.currentCommand = 0
				content := RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
				m.viewport.SetContent(content)
			} else if m.showDetail {
				// Go back to command list view
				m.showDetail = false
				content := RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
				m.viewport.SetContent(content)
				m.viewport.GotoTop()
			}

		case key.Matches(msg, keys.Up):
			if !m.showDetail {
				// Navigate up in the command list
				if m.currentCommand > 0 {
					m.currentCommand--
					// Update the content to reflect the new selection
					content := RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
					m.viewport.SetContent(content)
				}
			} else {
				// Scroll up in the viewport
				m.viewport.ScrollUp(1)
			}

		case key.Matches(msg, keys.Down):
			if !m.showDetail {
				// Navigate down in the command list
				if m.currentCommand < len(m.commands)-1 {
					m.currentCommand++
					// Update the content to reflect the new selection
					content := RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
					m.viewport.SetContent(content)

				}
			} else {
				// Scroll down in the viewport
				m.viewport.ScrollDown(1)
			}

		case key.Matches(msg, keys.Right):
			if !m.showDetail && !m.searchActive {
				// Navigate right on the tags (disabled when search is active)
				if m.currentTag < len(m.tagMenu)-1 {
					m.currentTag++
					m.commands = filterCommandsByTag(m.cheatSheet.Commands, m.tagMenu[m.currentTag])
					m.currentCommand = 0 // Reset the cursor
					content := lipgloss.Style(boxedViewportStyle).Render(RenderTagMenu(m.tagMenu, m.currentTag, m.width))
					m.tagViewPort.SetContent(content)
					content = RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
					m.viewport.SetContent(content)
				}
			} else if m.showDetail {
				m.tagViewPort.ScrollRight(10)
			}
		case key.Matches(msg, keys.Left):
			if !m.showDetail && !m.searchActive {
				// Navigate left on the tags (disabled when search is active)
				if m.currentTag > 0 {
					m.currentTag--
					m.commands = filterCommandsByTag(m.cheatSheet.Commands, m.tagMenu[m.currentTag])
					m.currentCommand = 0 // Reset the cursor
					logrus.Debugf("box size set: width=%d, height=%d", m.width, m.height)
					content := lipgloss.Style(boxedViewportStyle).Render(RenderTagMenu(m.tagMenu, m.currentTag, m.width))
					m.tagViewPort.SetContent(content)
					content = RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
					m.viewport.SetContent(content)
				}
			} else if m.showDetail {
				m.tagViewPort.ScrollLeft(10)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		logrus.Debugf("Window size changed: width=%d, height=%d", msg.Width, msg.Height)
		headerHeight := 9 // Reserve space for header/footer
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight
		m.tagViewPort.Width = msg.Width

		if m.tagMenu != nil && len(m.commands) > 0 {
			content := RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
			m.viewport.SetContent(content)
			tagsContent := lipgloss.Style(boxedViewportStyle).Render(RenderTagMenu(m.tagMenu, m.currentTag, m.width))
			m.tagViewPort.SetContent(tagsContent)
		}

	case cheatSheetLoadedMsg:
		// Handle the loaded cheat sheet
		m.cheatSheet = CheatSheet(msg)
		m.commands = m.cheatSheet.Commands
		m.tagMenu = UniqueTags(m.commands)
		m.currentTag = 0
		// Update the view with the command list
		logrus.Debugf("Window size set: width=%d, height=%d", m.width, m.height)
		if len(m.commands) > 0 {
			content := RenderCommandList(m.cheatSheet.Description, m.commands, m.currentCommand)
			tagsContent := lipgloss.Style(boxedViewportStyle).Render(RenderTagMenu(m.tagMenu, m.currentTag, m.width))
			m.tagViewPort.SetContent(tagsContent)
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

	// Create help text based on current mode
	var helpText string
	if m.searchMode {
		helpText = "Type to search • Enter: Apply • Esc: Cancel • q: Quit"
	} else if m.searchActive {
		helpText = "↑/↓: Navigate • Enter: View details • Esc: Clear search • q: Quit"
	} else {
		helpText = "↑/↓: Navigate • ←/→: Tag Filter • /: Search • Enter: View details • Esc: Back • q: Quit"
	}

	helpView := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render(helpText)

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

	logrus.Debug(m.tagMenu)

	// Build the view based on current mode
	var parts []string
	parts = append(parts, header)

	// Show search bar when in search mode, or search indicator when search is active
	if m.searchMode {
		searchBar := RenderSearchBar(m.searchQuery, m.width)
		parts = append(parts, searchBar)
	} else if m.searchActive {
		searchIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5AF78E")).
			Render(fmt.Sprintf("🔍 Search: %s (%d results)", m.searchQuery, len(m.commands)))
		parts = append(parts, searchIndicator)
	} else {
		// Show tag viewport only when not in search mode or active
		parts = append(parts, m.tagViewPort.View())
	}

	// Add main viewport and help
	parts = append(parts, m.viewport.View())
	parts = append(parts, helpView)

	return strings.Join(parts, "\n\n")
}
