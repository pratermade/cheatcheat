package main

import (
	"testing"
)

func TestLoadCheatSheet(t *testing.T) {
	sheet, err := LoadCheatSheet("cheatsheets/git.yaml")
	if err != nil {
		t.Fatalf("Failed to load cheatsheet: %v", err)
	}

	if sheet.Title != "Git Commands" {
		t.Errorf("Expected title 'Git Commands', got '%s'", sheet.Title)
	}
	if len(sheet.Commands) == 0 {
		t.Errorf("Expected at least one command, got 0")
	}
}
