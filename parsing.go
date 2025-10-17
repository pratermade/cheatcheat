package main

import (
	"os"

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
