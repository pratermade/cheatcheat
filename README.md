# cheatcheat

A lightweight Terminal User Interface (TUI) application for browsing interactive command cheatsheets. Built with Go and Bubble Tea, cheatcheat helps you quickly reference commands, options, and examples without leaving your terminal.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)

## Features

- **Cheatsheet Selector**: Browse and select from available cheatsheets at launch
- **Interactive TUI**: Navigate cheatsheets with intuitive keyboard controls
- **Tag-based Filtering**: Quickly filter commands by category tags
- **Live Search**: Real-time case-insensitive search through command names
- **Detailed Command View**: See syntax, examples, options, and notes for each command
- **Vim-style Navigation**: Use hjkl or arrow keys to navigate
- **YAML-based**: Easy to create and share cheatsheets
- **Syntax Highlighting**: Color-coded output for better readability
- **Switch on the Fly**: Open cheatsheet selector anytime with `o` key

## Installation

### From Source

```bash
git clone https://github.com/yourusername/cheatcheat.git
cd cheatcheat
go build -o cheatcheat
```

### Using Go Install

```bash
go install github.com/yourusername/cheatcheat@latest
```

## Quick Start

### Interactive Mode (Recommended)

Launch without arguments to browse and select from available cheatsheets:

```bash
./cheatcheat
```

This will display all `.yaml` files in the `cheatsheets/` directory (including subdirectories). Use arrow keys to navigate and press Enter to select.

### Direct Mode

Run cheatcheat with a specific YAML cheatsheet file:

```bash
./cheatcheat cheatsheets/kubectl.yaml
```

### Custom Cheatsheet Directory

Specify a different directory containing your cheatsheets:

```bash
./cheatcheat --dir /path/to/my/cheatsheets
```

## Usage

### Navigation

**Cheatsheet Selector:**
- `↑/k` or `↓/j` - Navigate through available cheatsheets
- `Enter` - Load selected cheatsheet
- `q` - Quit application

**List View:**
- `↑/k` or `↓/j` - Navigate through commands
- `←/h` or `→/l` - Switch between tag filters
- `/` - Activate search mode
- `Enter` - View detailed information for selected command
- `o` - Open cheatsheet selector
- `q` - Quit application

**Detail View:**
- `↑/k` or `↓/j` - Scroll through command details
- `Esc` - Return to command list
- `q` - Quit application

**Search Mode:**
- Type to search - Results update in real-time as you type
- `Backspace` - Delete characters from search query
- `Enter` - Apply search filter and exit search mode
- `Esc` - Cancel search and return to normal view

**Search Active:**
- `↑/k` or `↓/j` - Navigate through filtered results
- `Enter` - View detailed information for selected command
- `Esc` - Clear search filter and return to full list
- `o` - Open cheatsheet selector
- `q` - Quit application

### Search

Press `/` to activate search mode and start typing to filter commands by name. The search is:
- **Case-insensitive**: "git" matches "Git", "GIT", etc.
- **Live**: Results update in real-time as you type
- **Substring matching**: Searches anywhere in the command name

When search is active, tag navigation is disabled. Press `Esc` to clear the search and return to tag-based filtering.

### Tag Filtering

The tag menu at the top shows all available tags from your cheatsheet. Use `←/h` and `→/l` to switch between tags:
- Select "all" to see all commands
- Select specific tags to filter commands by category
- Indicators (« ») show when there are more tags to scroll through
- Tag navigation is disabled when search is active

## Creating Cheatsheets

Cheatsheets are defined in YAML format. Here's the structure:

```yaml
title: "Tool Name Cheat Sheet"
description: "Brief description of the tool"
category: "Category Name"

commands:
  - name: "command-name"
    shortDesc: "Brief one-line description"
    syntax: "command [options] <args>"
    tags: ["tag1", "tag2"]
    complexity: "beginner"  # or "intermediate", "advanced"
    examples:
      - code: "command --flag value"
        description: "What this example does"
      - code: "command --other-flag"
        description: "Another example"
    notes:
      - "Important note about the command"
      - "Tips or caveats to remember"
    options:
      - flag: "--flag-name"
        description: "What this flag does"
      - flag: "-f, --flag"
        description: "Short and long form"
    related: ["related-command", "another-command"]
```

### Required Fields

- `name`: Command name
- `shortDesc`: Brief description
- `syntax`: Command syntax

### Optional Fields

- `tags`: Array of category tags (enables filtering)
- `complexity`: Difficulty level indicator
- `examples`: Code examples with descriptions
- `notes`: Important information and tips
- `options`: Command flags and options
- `related`: Related commands for reference

### Example Cheatsheet

See the included cheatsheets for reference:
- `cheatsheets/kubectl.yaml` - Kubernetes CLI commands
- `cheatsheets/git.yaml` - Git version control commands
- `cheatsheets/databases/` - Database-related cheatsheets (MongoDB, MySQL)
- `cheatsheets/linux/` - Linux command cheatsheets

The selector will automatically discover all `.yaml` files in the `cheatsheets/` directory and its subdirectories.

## AI-Assisted Cheatsheet Creation

You can use AI assistants to generate comprehensive cheatsheets. Here's a prompt template:

```
Please create a comprehensive cheat sheet in YAML format for [TOOL/TECHNOLOGY] with the following structure:

title: "[TITLE]"
description: "[BRIEF DESCRIPTION]"
category: "[CATEGORY]"

For each command, include:
- name: The command name
- shortDesc: A brief one-line description
- syntax: The command syntax with placeholders
- tags: 2-4 descriptive tags for categorization
- complexity: One of "beginner", "intermediate", or "advanced"
- examples: 2-3 practical examples with both code and explanations
- notes: 2-3 important caveats, tips, or explanations
- options: The most important flags or options with descriptions
- related: Other related commands

Organize commands into logical groups and use consistent tags for proper filtering.
```

**Tips for better results:**
1. Be specific about which aspects of the tool to focus on
2. Mention preferred tag categories
3. Specify the desired comprehensiveness level
4. Request YAML validation

## Development

### Building

```bash
go build -o cheatcheat
```

### Running with Debug Logging

```bash
LogLevel=debug ./cheatcheat cheatsheets/kubectl.yaml
```

Debug logs are written to `logs/application_YYYY-MM-DD.log`.

### Running Tests

```bash
go test ./...
go test -v ./...        # verbose output
go test -cover ./...    # with coverage report
```

## Architecture

cheatcheat follows the Model-View-Update (MVU) pattern:

- **Model** (`model.go`): Application state
- **Update** (`main.go`): Event handling and state transitions
- **View** (`main.go`, `render.go`): UI rendering

Key components:
- `parsing.go`: YAML deserialization
- `render.go`: UI styling and layout
- `logging.go`: Debug logging utilities

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parsing
- [Logrus](https://github.com/sirupsen/logrus) - Structured logging

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

Built with the excellent [Charm](https://charm.sh) TUI libraries.
