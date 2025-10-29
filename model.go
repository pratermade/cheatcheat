package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

// Model for our application
type model struct {
	cheatSheet            CheatSheet
	commands              []Command
	tagMenu               []string
	currentCommand        int
	currentTag            int
	tagViewPort           viewport.Model
	viewport              viewport.Model
	width                 int
	height                int
	showDetail            bool
	searchMode            bool   // true when user is typing search query
	searchQuery           string // current search input text
	searchActive          bool   // true when search filter is applied
	err                   error  // Store any error that occurs
	cheatsheets           []string // list of discovered cheatsheet files
	currentCheatsheet     int      // selected index in cheatsheet selector
	showCheatsheetSelector bool    // true when showing cheatsheet selector
	cheatsheetDir         string   // base directory for cheatsheets
}

// Define key mappings
type keyMap struct {
	Up           key.Binding
	Down         key.Binding
	Enter        key.Binding
	Back         key.Binding
	Quit         key.Binding
	Search       key.Binding
	Left         key.Binding
	Right        key.Binding
	OpenSelector key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view details"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous tag"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next tag"),
	),
	OpenSelector: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open cheatsheet"),
	),
}
