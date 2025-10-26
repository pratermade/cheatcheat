# cheatcheat
Cheatsheet manager and tui
Please create a comprehensive cheat sheet in YAML format for [TOOL/TECHNOLOGY] with the following structure:

title: "[TITLE]"
description: "[BRIEF DESCRIPTION OF THE TOOL]"
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

Please organize the commands into logical groups based on functionality, and ensure the tags are consistent to enable proper grouping in my TUI application.

For reference, here's an example of a properly formatted command:

```yaml
commands:
  - name: "git clone"
    shortDesc: "Clone a repository into a new directory"
    syntax: "git clone <repository> [directory]"
    tags: ["setup", "repository"]
    complexity: "beginner"
    examples:
      - code: "git clone https://github.com/user/repo.git"
        description: "Clone a repository into a directory named after the repository"
      - code: "git clone --depth=1 https://github.com/user/repo.git"
        description: "Create a shallow clone with only the latest revision"
    notes:
      - "Clones the repository and creates remote-tracking branches"
      - "Automatically sets up the 'origin' remote"
    options:
      - flag: "--depth=<depth>"
        description: "Create a shallow clone with limited revision history"
      - flag: "--branch, -b <branch>"
        description: "Clone the specified branch instead of the default"
    related: ["git init", "git remote"]



## How to Use This Prompt

1. Copy the entire prompt above.
2. Replace `[TOOL/TECHNOLOGY]` with the specific tool you want a cheat sheet for (e.g., Docker, SQL, Vim, etc.)
3. Replace `[TITLE]` with a descriptive title for the cheat sheet.
4. Replace `[BRIEF DESCRIPTION OF THE TOOL]` with a short description.
5. Replace `[CATEGORY]` with an appropriate category (e.g., "Developer Tools", "Database", "Editor").
6. Share this with ChatGPT or another AI assistant.

## Tips for Getting Better Results

1. **Be Specific**: Mention particular aspects of the tool you want to focus on.
2. **Specify Tag Categories**: If you have particular tag groups in mind, mention them.
3. **Request Size**: Specify if you want a small, medium, or comprehensive cheat sheet.
4. **Ask for Validation**: Request that the AI validate the YAML format to ensure it will parse correctly.

Example usage:
"Please create a comprehensive cheat sheet in YAML format for Docker with a focus on container management and networking..."

The resulting YAML can be saved directly to a file and used with your TUI cheat sheet application.
