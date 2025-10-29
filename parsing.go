package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Option struct {
	Flag        string `yaml:"flag"`
	Description string `yaml:"description"`
}

type Example struct {
	Code        string `yaml:"code"`
	Description string `yaml:"description"`
}

type Command struct {
	Name       string    `yaml:"name"`
	ShortDesc  string    `yaml:"shortDesc"`
	Syntax     string    `yaml:"syntax"`
	Tags       []string  `yaml:"tags"`
	Complexity string    `yaml:"complexity"`
	Examples   []Example `yaml:"examples"`
	Notes      []string  `yaml:"notes"`
	Options    []Option  `yaml:"options"`
	Related    []string  `yaml:"related"`
}

type CheatSheet struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Category    string    `yaml:"category"`
	Commands    []Command `yaml:"commands"`
}

// Load a cheatsheet from file
func LoadCheatSheet(filename string) (CheatSheet, error) {
	var sheet CheatSheet
	data, err := os.ReadFile(filename)
	if err != nil {
		return sheet, err
	}
	err = yaml.Unmarshal(data, &sheet)
	return sheet, err
}

// DiscoverCheatsheets recursively scans a directory for .yaml files
// and returns their relative paths sorted alphabetically
func DiscoverCheatsheets(dir string) ([]string, error) {
	var cheatsheets []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include .yaml files
		if strings.HasSuffix(strings.ToLower(path), ".yaml") {
			// Get relative path from base directory
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			cheatsheets = append(cheatsheets, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort alphabetically for consistent display
	sort.Strings(cheatsheets)

	return cheatsheets, nil
}
