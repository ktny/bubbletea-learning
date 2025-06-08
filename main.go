package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// model represents the application state
type model struct {
	message string
}

// Init is called when the program starts and returns an optional initial command
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and returns an updated model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the UI based on the model
func (m model) View() string {
	return fmt.Sprintf("%s\n\nPress 'q' or Ctrl+C to quit.\n", m.message)
}

func main() {
	// Initialize the counter model
	initialModel := NewCounterModel()

	// Create a new program
	p := tea.NewProgram(initialModel)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}