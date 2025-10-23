package ui

import (
	"github.com/rivo/tview"

	"sqlcmder/logger"
)

const (
	pageNameCommandModal = "command_modal"
)

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

// ShowError displays an error message (simplified without command line)
func ShowError(message string) {
	logger.Error("Error", map[string]any{"message": message})
	// For now, just log the error since command line is removed
}

// ShowSuccess displays a success message (simplified without command line)
func ShowSuccess(message string) {
	logger.Info("Success", map[string]any{"message": message})
	// For now, just log the success since command line is removed
}

// ShowInfo displays an information message (simplified without command line)
func ShowInfo(message string) {
	logger.Info("Info", map[string]any{"message": message})
	// For now, just log the info since command line is removed
}

// RefreshTree refreshes the database tree view
// This is a placeholder - implement based on your tree component
func RefreshTree() {
	// TODO: Implement tree refresh logic
	// This should call your tree component's refresh method
}
