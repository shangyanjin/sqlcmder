package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
	"sqlcmder/logger"
)

type ConfirmationModal struct {
	*tview.Modal
}

func NewConfirmationModal(confirmationText string) *ConfirmationModal {
	modal := tview.NewModal()
	if confirmationText != "" {
		modal.SetText(confirmationText)
	} else {
		modal.SetText("Are you sure?")
	}
	modal.AddButtons([]string{"Yes", "No"})
	modal.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)
	modal.SetBorderStyle(tcell.StyleDefault.Background(app.Styles.PrimitiveBackgroundColor))

	// Unselected button style - use dark text for low-key appearance
	modal.SetButtonStyle(tcell.StyleDefault.
		Background(app.Styles.ButtonUnselectedBgColor).
		Foreground(app.Styles.SelectedTextColor), // Dark text for unselected buttons
	)

	// Selected/Activated button style - use bright yellow text for high visibility
	modal.SetButtonActivatedStyle(tcell.StyleDefault.
		Background(app.Styles.ButtonBackgroundColor).
		Foreground(app.Styles.ButtonTextColor), // Bright yellow text for selected buttons
	)

	modal.SetTextColor(app.Styles.PrimaryTextColor)

	// Log when modal is created
	logger.Info("Confirmation modal created", map[string]any{
		"text": confirmationText,
	})

	return &ConfirmationModal{
		Modal: modal,
	}
}
