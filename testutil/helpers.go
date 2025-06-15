// Package testutil provides common testing utilities for the Bubble Tea applications.
package testutil

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Contains checks if a string contains a substring.
// This is a simple wrapper around strings.Contains for consistency.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// AssertContains asserts that a string contains a substring.
func AssertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !Contains(s, substr) {
		t.Errorf("Expected string to contain '%s', but it didn't. String: %s", substr, s)
	}
}

// AssertNotContains asserts that a string does not contain a substring.
func AssertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if Contains(s, substr) {
		t.Errorf("Expected string not to contain '%s', but it did. String: %s", substr, s)
	}
}

// AssertViewContains asserts that a model's View() contains specific text.
func AssertViewContains(t *testing.T, model tea.Model, expectedText string) {
	t.Helper()
	view := model.View()
	AssertContains(t, view, expectedText)
}

// AssertViewNotContains asserts that a model's View() does not contain specific text.
func AssertViewNotContains(t *testing.T, model tea.Model, unexpectedText string) {
	t.Helper()
	view := model.View()
	AssertNotContains(t, view, unexpectedText)
}

// SendKeyAndUpdate sends a key message to a model and returns the updated model.
func SendKeyAndUpdate(model tea.Model, key tea.KeyMsg) (tea.Model, tea.Cmd) {
	return model.Update(key)
}

// SendKeyTypeAndUpdate sends a key message with a specific type and returns the updated model.
func SendKeyTypeAndUpdate(model tea.Model, keyType tea.KeyType) (tea.Model, tea.Cmd) {
	return model.Update(tea.KeyMsg{Type: keyType})
}

// SendKeyRuneAndUpdate sends a key message with specific runes and returns the updated model.
func SendKeyRuneAndUpdate(model tea.Model, runes string) (tea.Model, tea.Cmd) {
	return model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(runes)})
}

// AssertKeyResult sends a key to a model and asserts the result using a custom check function.
func AssertKeyResult(t *testing.T, model tea.Model, key tea.KeyMsg, check func(tea.Model)) {
	t.Helper()
	newModel, _ := model.Update(key)
	check(newModel)
}

// AssertKeyTypeResult sends a key type to a model and asserts the result.
func AssertKeyTypeResult(t *testing.T, model tea.Model, keyType tea.KeyType, check func(tea.Model)) {
	t.Helper()
	AssertKeyResult(t, model, tea.KeyMsg{Type: keyType}, check)
}

// AssertKeyRuneResult sends key runes to a model and asserts the result.
func AssertKeyRuneResult(t *testing.T, model tea.Model, runes string, check func(tea.Model)) {
	t.Helper()
	AssertKeyResult(t, model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(runes)}, check)
}

// AssertInitCommand asserts that a model's Init() returns a non-nil command.
func AssertInitCommand(t *testing.T, model tea.Model) {
	t.Helper()
	cmd := model.Init()
	if cmd == nil {
		t.Error("Expected Init() to return a non-nil command")
	}
}

// AssertNoInitCommand asserts that a model's Init() returns nil.
func AssertNoInitCommand(t *testing.T, model tea.Model) {
	t.Helper()
	cmd := model.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil")
	}
}

// AssertNonEmptyView asserts that a model's View() returns non-empty content.
func AssertNonEmptyView(t *testing.T, model tea.Model) {
	t.Helper()
	view := model.View()
	if view == "" {
		t.Error("Expected View() to return non-empty content")
	}
}

// CommonKeyTestCases provides common key test cases for models that support quit functionality.
type CommonKeyTestCase struct {
	Name     string
	Key      tea.KeyMsg
	ShouldQuit bool
}

// GetQuitKeyTestCases returns common test cases for quit keys.
func GetQuitKeyTestCases() []CommonKeyTestCase {
	return []CommonKeyTestCase{
		{
			Name:       "Ctrl+Cで終了",
			Key:        tea.KeyMsg{Type: tea.KeyCtrlC},
			ShouldQuit: true,
		},
		{
			Name:       "qキーで終了",
			Key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			ShouldQuit: true,
		},
		{
			Name:       "その他のキーは終了しない",
			Key:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			ShouldQuit: false,
		},
	}
}

// AssertQuitBehavior tests common quit key behavior.
func AssertQuitBehavior(t *testing.T, model tea.Model, testCases []CommonKeyTestCase) {
	t.Helper()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, cmd := model.Update(tc.Key)
			hasQuitCmd := cmd != nil
			if hasQuitCmd != tc.ShouldQuit {
				if tc.ShouldQuit {
					t.Error("Expected quit command but got none")
				} else {
					t.Error("Expected no quit command but got one")
				}
			}
		})
	}
}