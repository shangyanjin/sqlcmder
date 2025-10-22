package components

import (
	"github.com/rivo/tview"

	"sqlcmder/internal/app"
	"sqlcmder/internal/helpers/logger"
)

const (
	pageNameCommandModal = "command_modal"
)

// Global command line reference for displaying messages
var globalCommandLine *CommandLine

// ShowModal displays a modal dialog
func ShowModal(modal tview.Primitive, width, height int) {
	// Create a flex to center the modal
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modal, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)

	mainPages.AddPage(pageNameCommandModal, flex, true, true)
	App.SetFocus(modal)
}

// CloseModal closes the current modal dialog
func CloseModal() {
	mainPages.RemovePage(pageNameCommandModal)
}

// truncateMessage truncates message to fit in one line
func truncateMessage(message string, maxLen int) string {
	if len(message) <= maxLen {
		return message
	}
	return message[:maxLen-3] + "..."
}

// ShowError displays an error message in command line style
func ShowError(message string) {
	logger.Debug("ShowError called", map[string]any{
		"message":           message,
		"globalCommandLine": globalCommandLine != nil,
		"hasMessageView":    globalCommandLine != nil && globalCommandLine.MessageView != nil,
		"hasInputField":     globalCommandLine != nil && globalCommandLine.InputField != nil,
	})

	// Display in message row - simple direct update
	if globalCommandLine != nil && globalCommandLine.MessageView != nil {
		logger.Debug("Setting error message to MessageView", nil)
		// Truncate message to avoid breaking table layout (max ~100 chars)
		displayMsg := truncateMessage(message, 100)
		globalCommandLine.MessageView.SetText("[red]✗ " + displayMsg)
		globalCommandLine.InputField.SetText("")
		// Force UI redraw to prevent layout issues
		app.App.Draw()
		// Auto-restore focus to input field
		app.App.SetFocus(globalCommandLine.InputField)
		logger.Debug("Error message set and focus restored", nil)
	} else {
		logger.Debug("Cannot show error - globalCommandLine not initialized", nil)
	}
}

// ShowSuccess displays a success message in command line style
func ShowSuccess(message string) {
	// Display in message row - simple direct update
	if globalCommandLine != nil && globalCommandLine.MessageView != nil {
		// Truncate message to avoid breaking table layout
		displayMsg := truncateMessage(message, 100)
		globalCommandLine.MessageView.SetText("[green]✓ " + displayMsg)
		globalCommandLine.InputField.SetText("")
		// Force UI redraw to prevent layout issues
		app.App.Draw()
		// Auto-restore focus to input field
		app.App.SetFocus(globalCommandLine.InputField)
	}
}

// ShowInfo displays an information message in command line style
func ShowInfo(message string) {
	// Display in message row - simple direct update
	if globalCommandLine != nil && globalCommandLine.MessageView != nil {
		// Truncate message to avoid breaking table layout
		displayMsg := truncateMessage(message, 100)
		globalCommandLine.MessageView.SetText("[blue]ℹ " + displayMsg)
		globalCommandLine.InputField.SetText("")
		// Force UI redraw to prevent layout issues
		app.App.Draw()
		// Auto-restore focus to input field
		app.App.SetFocus(globalCommandLine.InputField)
	}
}

// RefreshTree refreshes the database tree view
// This is a placeholder - implement based on your tree component
func RefreshTree() {
	// TODO: Implement tree refresh logic
	// This should call your tree component's refresh method
}
