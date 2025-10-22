package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/internal/app"
	"sqlcmder/components"
)

// RegisterDatabaseCommands registers all database-related commands
func RegisterDatabaseCommands(cp *CommandPalette) {
	// Create Database
	cp.RegisterCommand(Command{
		ID:          "db.create",
		Name:        "Create Database",
		Description: "Create a new database",
		Icon:        "📊",
		Category:    "Database",
		Handler:     handleCreateDatabase,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil
		},
	})

	// Drop Database
	cp.RegisterCommand(Command{
		ID:          "db.drop",
		Name:        "Drop Database",
		Description: "Delete current database",
		Icon:        "🗑�?,
		Category:    "Database",
		Handler:     handleDropDatabase,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil && ctx.CurrentDatabase != ""
		},
	})

	// Switch Database
	cp.RegisterCommand(Command{
		ID:          "db.switch",
		Name:        "Switch Database",
		Description: "Change to another database",
		Icon:        "🔄",
		Category:    "Database",
		Handler:     handleSwitchDatabase,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil
		},
	})

	// Refresh Database List
	cp.RegisterCommand(Command{
		ID:          "db.refresh",
		Name:        "Refresh Database List",
		Description: "Reload database structure",
		Icon:        "♻️",
		Category:    "Database",
		Handler:     handleRefreshDatabase,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil
		},
	})
}

func handleCreateDatabase(ctx CommandContext) error {
	// Create input form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Create Database ")
	form.SetTitleAlign(tview.AlignLeft)

	var dbName string
	var charset string

	form.AddInputField("Database Name", "", 30, nil, func(text string) {
		dbName = text
	})

	form.AddDropDown("Character Set", []string{"utf8mb4", "utf8", "latin1"}, 0, func(option string, index int) {
		charset = option
	})

	form.AddButton("Create", func() {
		if dbName == "" {
			components.ShowError("Database name cannot be empty")
			return
		}

		// Build SQL
		sql := fmt.Sprintf("CREATE DATABASE `%s`", dbName)
		if charset != "" {
			sql += fmt.Sprintf(" CHARACTER SET %s", charset)
		}

		// Execute
		_, err := ctx.DB.ExecuteDMLStatement(sql)
		if err != nil {
			ShowError(fmt.Sprintf("Failed to create database: %v", err))
			return
		}

		ShowSuccess(fmt.Sprintf("Database '%s' created successfully", dbName))

		// Close modal and refresh
		CloseModal()
		RefreshTree()
	})

	form.AddButton("Cancel", func() {
		CloseModal()
	})

	form.SetButtonsAlign(tview.AlignCenter)
	form.SetFieldBackgroundColor(app.Styles.InverseTextColor)

	ShowModal(form, 50, 10)
	return nil
}

func handleDropDatabase(ctx CommandContext) error {
	dbName := ctx.CurrentDatabase

	// Create confirmation dialog
	modal := tview.NewModal()
	modal.SetText(fmt.Sprintf("Are you sure you want to drop database '%s'?\n\nThis action cannot be undone!", dbName))
	modal.AddButtons([]string{"Drop", "Cancel"})
	modal.SetButtonBackgroundColor(app.Styles.InverseTextColor)
	modal.SetButtonTextColor(app.Styles.PrimaryTextColor)
	modal.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Drop" {
			// Execute DROP DATABASE
			sql := fmt.Sprintf("DROP DATABASE `%s`", dbName)
			_, err := ctx.DB.ExecuteDMLStatement(sql)
			if err != nil {
				ShowError(fmt.Sprintf("Failed to drop database: %v", err))
				return
			}

			ShowSuccess(fmt.Sprintf("Database '%s' dropped successfully", dbName))
			RefreshTree()
		}
		CloseModal()
	})

	ShowModal(modal, 60, 8)
	return nil
}

func handleSwitchDatabase(ctx CommandContext) error {
	// Get list of databases
	databases, err := ctx.DB.GetDatabases()
	if err != nil {
		return fmt.Errorf("failed to get databases: %w", err)
	}

	// Create selection list
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Select Database ")
	list.ShowSecondaryText(false)

	for _, db := range databases {
		dbName := db
		list.AddItem(fmt.Sprintf("📁 %s", dbName), "", 0, func() {
			// Switch database by executing USE statement
			sql := fmt.Sprintf("USE `%s`", dbName)
			_, err := ctx.DB.ExecuteDMLStatement(sql)
			if err != nil {
				ShowError(fmt.Sprintf("Failed to switch database: %v", err))
				return
			}

			ShowSuccess(fmt.Sprintf("Switched to database '%s'", dbName))
			CloseModal()
			RefreshTree()
		})
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			CloseModal()
			return nil
		}
		return event
	})

	ShowModal(list, 50, 15)
	return nil
}

func handleRefreshDatabase(ctx CommandContext) error {
	RefreshTree()
	ShowSuccess("Database list refreshed")
	return nil
}
