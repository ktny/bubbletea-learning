// Package common provides common functionality used across multiple applications.
package common

import tea "github.com/charmbracelet/bubbletea"

// HandleQuitKeys handles common quit key combinations (Ctrl+C and 'q').
// Returns tea.Quit command if a quit key is pressed, nil otherwise.
func HandleQuitKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlC:
		return tea.Quit
	case tea.KeyRunes:
		if len(msg.Runes) == 1 && msg.Runes[0] == 'q' {
			return tea.Quit
		}
	}
	return nil
}

// IsQuitKey checks if the given key message is a quit key.
func IsQuitKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyCtrlC:
		return true
	case tea.KeyRunes:
		return len(msg.Runes) == 1 && msg.Runes[0] == 'q'
	}
	return false
}

// HandleCommonNavigation handles common navigation keys (arrows and vim-style).
// Returns direction as an integer: -1 (up/left), 1 (down/right), 0 (no movement).
func HandleCommonNavigation(msg tea.KeyMsg) int {
	switch msg.Type {
	case tea.KeyUp:
		return -1
	case tea.KeyDown:
		return 1
	case tea.KeyRunes:
		if len(msg.Runes) == 1 {
			switch msg.Runes[0] {
			case 'k': // vim up
				return -1
			case 'j': // vim down
				return 1
			}
		}
	}
	return 0
}

// IsTabKey checks if the given key message is a Tab key.
func IsTabKey(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyTab
}

// IsShiftTabKey checks if the given key message is a Shift+Tab key.
func IsShiftTabKey(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyShiftTab
}

// IsEnterKey checks if the given key message is an Enter key.
func IsEnterKey(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyEnter
}

// IsEscKey checks if the given key message is an Escape key.
func IsEscKey(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyEsc
}

// IsSpaceKey checks if the given key message is a Space key.
func IsSpaceKey(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeySpace || 
		(msg.Type == tea.KeyRunes && len(msg.Runes) == 1 && msg.Runes[0] == ' ')
}