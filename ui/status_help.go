package ui

import (
	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
	"sqlcmder/keymap"
)

type HelpStatus struct {
	*tview.TextView
}

func NewHelpStatus() HelpStatus {
	status := HelpStatus{tview.NewTextView().SetTextColor(app.Styles.TertiaryTextColor)}

	status.SetStatusOnTree()

	return status
}

func (status *HelpStatus) UpdateText(binds []keymap.Bind) {
	newtext := ""

	for i, key := range binds {

		newtext += key.Cmd.String()

		newtext += ": "

		newtext += key.Key.String()

		islast := i == len(binds)-1

		if !islast {
			newtext += " | "
		}

	}

	status.SetText(newtext)
}

func (status *HelpStatus) SetStatusOnTree() {
	status.UpdateText(keymap.Keymaps.Global)
}

func (status *HelpStatus) SetStatusOnEditorView() {
	status.UpdateText(keymap.Keymaps.Group(keymap.EditorGroup))
}

func (status *HelpStatus) SetStatusOnTableView() {
	status.UpdateText(keymap.Keymaps.Group(keymap.TableGroup))
}
