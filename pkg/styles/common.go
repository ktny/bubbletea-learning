// Package styles provides common style definitions for the Bubble Tea applications.
package styles

import "github.com/charmbracelet/lipgloss"

// Common colors
var (
	// Primary colors
	PrimaryColor   = lipgloss.Color("12") // Blue
	SuccessColor   = lipgloss.Color("10") // Green
	ErrorColor     = lipgloss.Color("9")  // Red
	WarningColor   = lipgloss.Color("11") // Yellow
	SecondaryColor = lipgloss.Color("14") // Cyan
	GrayColor      = lipgloss.Color("8")  // Gray
)

// Common styles used across multiple applications
var (
	// TitleStyle is used for application titles
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		MarginBottom(1)

	// HelpStyle is used for help text and instructions
	HelpStyle = lipgloss.NewStyle().
		Foreground(GrayColor).
		Italic(true)

	// ErrorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true)

	// SuccessStyle is used for success messages
	SuccessStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true)

	// BorderStyle is the default border style
	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2)

	// ActiveBorderStyle is used for active/focused elements
	ActiveBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2)

	// InactiveBorderStyle is used for inactive elements
	InactiveBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GrayColor).
		Padding(1, 2)

	// LabelStyle is used for form labels
	LabelStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true).
		MarginRight(1)

	// ValueStyle is used for displaying values
	ValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")) // White

	// DimmedStyle is used for less important information
	DimmedStyle = lipgloss.NewStyle().
		Foreground(GrayColor)
)

// Counter specific styles
var (
	// CounterPositiveStyle for positive counter values
	CounterPositiveStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true)

	// CounterNegativeStyle for negative counter values
	CounterNegativeStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true)

	// CounterZeroStyle for zero counter value
	CounterZeroStyle = lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true)
)

// Dashboard specific styles
var (
	// DashboardTitleStyle for the main dashboard title
	DashboardTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	// DashboardHelpStyle for dashboard help text
	DashboardHelpStyle = lipgloss.NewStyle().
		Foreground(GrayColor).
		Background(lipgloss.Color("235")).
		Padding(0, 1)
)