package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
	"sqlcmder/drivers"
	"sqlcmder/models"
)

// CommandContext holds the current context for command execution
type CommandContext struct {
	DB              drivers.Driver
	CurrentDatabase string
	CurrentTable    string
	Connection      string
	ConnectionModel *models.Connection // Full connection details for backup/import
}

// Command represents a command in the palette
type Command struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Category    string
	Handler     func(ctx CommandContext) error
	Enabled     func(ctx CommandContext) bool
}

// CommandPalette is a modal dialog for executing commands
type CommandPalette struct {
	*tview.Flex
	InputField *tview.InputField
	List       *tview.List
	Commands   []Command
	Context    CommandContext
	OnClose    func()
}

// NewCommandPalette creates a new command palette
func NewCommandPalette() *CommandPalette {
	cp := &CommandPalette{
		Flex:     tview.NewFlex(),
		Commands: []Command{},
	}

	cp.SetDirection(tview.FlexRow)
	cp.SetBorder(true)
	cp.SetTitle(" Command Palette (Ctrl+P) ")
	cp.SetBorderColor(tcell.ColorYellow)

	// Input field for search
	cp.InputField = tview.NewInputField()
	cp.InputField.SetLabel("> ")
	cp.InputField.SetFieldWidth(0)
	cp.InputField.SetFieldBackgroundColor(app.Styles.InverseTextColor)
	cp.InputField.SetLabelColor(tcell.ColorYellow)

	// Command list
	cp.List = tview.NewList()
	cp.List.ShowSecondaryText(false)
	cp.List.SetHighlightFullLine(true)
	cp.List.SetSelectedBackgroundColor(tcell.ColorBlue)

	cp.AddItem(cp.InputField, 1, 0, true)
	cp.AddItem(cp.List, 0, 1, false)

	// Setup input handler
	cp.InputField.SetChangedFunc(func(text string) {
		cp.FilterCommands(text)
	})

	cp.InputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if cp.OnClose != nil {
				cp.OnClose()
			}
			return nil
		case tcell.KeyDown:
			App.SetFocus(cp.List)
			return nil
		case tcell.KeyEnter:
			if cp.List.GetItemCount() > 0 {
				cp.ExecuteSelected()
			}
			return nil
		}
		return event
	})

	cp.List.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if cp.OnClose != nil {
				cp.OnClose()
			}
			return nil
		case tcell.KeyUp:
			if cp.List.GetCurrentItem() == 0 {
				App.SetFocus(cp.InputField)
				return nil
			}
		case tcell.KeyEnter:
			cp.ExecuteSelected()
			return nil
		}
		return event
	})

	return cp
}

// RegisterCommand adds a command to the palette
func (cp *CommandPalette) RegisterCommand(cmd Command) {
	cp.Commands = append(cp.Commands, cmd)
}

// SetContext updates the command context
func (cp *CommandPalette) SetContext(ctx CommandContext) {
	cp.Context = ctx
	cp.RefreshCommands()
}

// RefreshCommands rebuilds the command list
func (cp *CommandPalette) RefreshCommands() {
	cp.List.Clear()

	for _, cmd := range cp.Commands {
		// Check if command is enabled in current context
		if cmd.Enabled != nil && !cmd.Enabled(cp.Context) {
			continue
		}

		displayText := cmd.Icon + " " + cmd.Name
		if cmd.Description != "" {
			displayText += " - " + cmd.Description
		}

		cp.List.AddItem(displayText, "", 0, nil)
	}

	if cp.List.GetItemCount() > 0 {
		cp.List.SetCurrentItem(0)
	}
}

// FilterCommands filters commands based on search text
func (cp *CommandPalette) FilterCommands(searchText string) {
	cp.List.Clear()

	searchLower := strings.ToLower(searchText)

	for _, cmd := range cp.Commands {
		// Check if command is enabled
		if cmd.Enabled != nil && !cmd.Enabled(cp.Context) {
			continue
		}

		// Fuzzy match
		nameLower := strings.ToLower(cmd.Name)
		descLower := strings.ToLower(cmd.Description)

		if searchText == "" || strings.Contains(nameLower, searchLower) || strings.Contains(descLower, searchLower) {
			displayText := cmd.Icon + " " + cmd.Name
			if cmd.Description != "" {
				displayText += " - " + cmd.Description
			}

			cp.List.AddItem(displayText, "", 0, nil)
		}
	}

	if cp.List.GetItemCount() > 0 {
		cp.List.SetCurrentItem(0)
	}
}

// ExecuteSelected executes the currently selected command
func (cp *CommandPalette) ExecuteSelected() {
	currentIndex := cp.List.GetCurrentItem()
	if currentIndex < 0 {
		return
	}

	// Find the actual command (accounting for filtering)
	searchText := strings.ToLower(cp.InputField.GetText())
	visibleIndex := 0

	for _, cmd := range cp.Commands {
		if cmd.Enabled != nil && !cmd.Enabled(cp.Context) {
			continue
		}

		nameLower := strings.ToLower(cmd.Name)
		descLower := strings.ToLower(cmd.Description)

		if searchText == "" || strings.Contains(nameLower, searchText) || strings.Contains(descLower, searchText) {
			if visibleIndex == currentIndex {
				// Execute the command
				if cmd.Handler != nil {
					if err := cmd.Handler(cp.Context); err != nil {
						// Show error
						ShowError(err.Error())
					}
				}

				// Close the palette
				if cp.OnClose != nil {
					cp.OnClose()
				}
				return
			}
			visibleIndex++
		}
	}
}

// Show displays the command palette
func (cp *CommandPalette) Show() {
	cp.InputField.SetText("")
	cp.RefreshCommands()
	App.SetFocus(cp.InputField)
}
