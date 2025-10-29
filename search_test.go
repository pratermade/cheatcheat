package main

import (
	"testing"
)

func TestFilterCommandsBySearch(t *testing.T) {
	// Create test commands
	commands := []Command{
		{Name: "git status", ShortDesc: "Show working tree status"},
		{Name: "git commit", ShortDesc: "Record changes to repository"},
		{Name: "git push", ShortDesc: "Update remote refs"},
		{Name: "kubectl get pods", ShortDesc: "List pods"},
		{Name: "kubectl apply", ShortDesc: "Apply configuration"},
	}

	// Test case 1: Empty query should return all commands
	result := filterCommandsBySearch(commands, "")
	if len(result) != len(commands) {
		t.Errorf("Expected %d commands with empty query, got %d", len(commands), len(result))
	}

	// Test case 2: Search for "git" should return 3 commands
	result = filterCommandsBySearch(commands, "git")
	if len(result) != 3 {
		t.Errorf("Expected 3 commands with 'git' query, got %d", len(result))
	}

	// Test case 3: Case-insensitive search for "GIT" should return 3 commands
	result = filterCommandsBySearch(commands, "GIT")
	if len(result) != 3 {
		t.Errorf("Expected 3 commands with 'GIT' query (case-insensitive), got %d", len(result))
	}

	// Test case 4: Search for "kubectl" should return 2 commands
	result = filterCommandsBySearch(commands, "kubectl")
	if len(result) != 2 {
		t.Errorf("Expected 2 commands with 'kubectl' query, got %d", len(result))
	}

	// Test case 5: Search for "push" should return 1 command
	result = filterCommandsBySearch(commands, "push")
	if len(result) != 1 {
		t.Errorf("Expected 1 command with 'push' query, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "git push" {
		t.Errorf("Expected 'git push', got '%s'", result[0].Name)
	}

	// Test case 6: Search for non-existent command should return 0 commands
	result = filterCommandsBySearch(commands, "nonexistent")
	if len(result) != 0 {
		t.Errorf("Expected 0 commands with 'nonexistent' query, got %d", len(result))
	}

	// Test case 7: Partial match should work
	result = filterCommandsBySearch(commands, "get")
	if len(result) != 1 {
		t.Errorf("Expected 1 command with 'get' query, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "kubectl get pods" {
		t.Errorf("Expected 'kubectl get pods', got '%s'", result[0].Name)
	}
}
