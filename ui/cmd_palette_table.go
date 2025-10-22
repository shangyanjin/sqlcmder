package ui

import (
	"fmt"

	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
)

// RegisterTableCommands registers all table-related commands
func RegisterTableCommands(cp *CommandPalette) {
	// Create Table
	cp.RegisterCommand(Command{
		ID:          "table.create",
		Name:        "Create Table",
		Description: "Create a new table",
		Icon:        "[TBL]",
		Category:    "Table",
		Handler:     handleCreateTable,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil && ctx.CurrentDatabase != ""
		},
	})

	// Drop Table
	cp.RegisterCommand(Command{
		ID:          "table.drop",
		Name:        "Drop Table",
		Description: "Delete current table",
		Icon:        "[DEL]",
		Category:    "Table",
		Handler:     handleDropTable,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil && ctx.CurrentTable != ""
		},
	})

	// Rename Table
	cp.RegisterCommand(Command{
		ID:          "table.rename",
		Name:        "Rename Table",
		Description: "Rename current table",
		Icon:        "[REN]",
		Category:    "Table",
		Handler:     handleRenameTable,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil && ctx.CurrentTable != ""
		},
	})

	// Truncate Table
	cp.RegisterCommand(Command{
		ID:          "table.truncate",
		Name:        "Truncate Table",
		Description: "Delete all rows from table",
		Icon:        "[TRU]",
		Category:    "Table",
		Handler:     handleTruncateTable,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil && ctx.CurrentTable != ""
		},
	})

	// Copy Table Structure
	cp.RegisterCommand(Command{
		ID:          "table.copy",
		Name:        "Copy Table Structure",
		Description: "Create a copy of table structure",
		Icon:        "[CPY]",
		Category:    "Table",
		Handler:     handleCopyTable,
		Enabled: func(ctx CommandContext) bool {
			return ctx.DB != nil && ctx.CurrentTable != ""
		},
	})
}

func handleCreateTable(ctx CommandContext) error {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Create Table ")
	form.SetTitleAlign(tview.AlignLeft)

	var tableName string
	var withID bool = true

	form.AddInputField("Table Name", "", 30, nil, func(text string) {
		tableName = text
	})

	form.AddCheckbox("Add ID column (auto-increment)", true, func(checked bool) {
		withID = checked
	})

	form.AddButton("Create", func() {
		if tableName == "" {
			ShowError("Table name cannot be empty")
			return
		}

		// Build basic CREATE TABLE SQL
		var sql string
		if withID {
			sql = fmt.Sprintf(`CREATE TABLE %s (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`, tableName)
		} else {
			sql = fmt.Sprintf(`CREATE TABLE %s (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`, tableName)
		}

		// Execute
		_, err := ctx.DB.ExecuteDMLStatement(sql)
		if err != nil {
			ShowError(fmt.Sprintf("Failed to create table: %v", err))
			return
		}

		ShowSuccess(fmt.Sprintf("Table '%s' created successfully", tableName))
		CloseModal()
		RefreshTree()
	})

	form.AddButton("Cancel", func() {
		CloseModal()
	})

	form.SetButtonsAlign(tview.AlignCenter)
	form.SetFieldBackgroundColor(app.Styles.InverseTextColor)

	ShowModal(form, 60, 12)
	return nil
}

func handleDropTable(ctx CommandContext) error {
	tableName := ctx.CurrentTable

	modal := tview.NewModal()
	modal.SetText(fmt.Sprintf("Are you sure you want to drop table '%s'?\n\nAll data will be lost!", tableName))
	modal.AddButtons([]string{"Drop", "Cancel"})
	modal.SetButtonBackgroundColor(app.Styles.InverseTextColor)
	modal.SetButtonTextColor(app.Styles.PrimaryTextColor)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Drop" {
			sql := fmt.Sprintf("DROP TABLE `%s`", tableName)
			_, err := ctx.DB.ExecuteDMLStatement(sql)
			if err != nil {
				ShowError(fmt.Sprintf("Failed to drop table: %v", err))
				return
			}

			ShowSuccess(fmt.Sprintf("Table '%s' dropped successfully", tableName))
			RefreshTree()
		}
		CloseModal()
	})

	ShowModal(modal, 60, 8)
	return nil
}

func handleRenameTable(ctx CommandContext) error {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Rename Table ")

	oldName := ctx.CurrentTable
	var newName string

	form.AddInputField("Current Name", oldName, 30, nil, nil).GetFormItem(0).(*tview.InputField).SetDisabled(true)
	form.AddInputField("New Name", "", 30, nil, func(text string) {
		newName = text
	})

	form.AddButton("Rename", func() {
		if newName == "" {
			ShowError("New table name cannot be empty")
			return
		}

		sql := fmt.Sprintf("ALTER TABLE `%s` RENAME TO `%s`", oldName, newName)
		_, err := ctx.DB.ExecuteDMLStatement(sql)
		if err != nil {
			ShowError(fmt.Sprintf("Failed to rename table: %v", err))
			return
		}

		ShowSuccess(fmt.Sprintf("Table renamed to '%s'", newName))
		CloseModal()
		RefreshTree()
	})

	form.AddButton("Cancel", func() {
		CloseModal()
	})

	form.SetButtonsAlign(tview.AlignCenter)
	form.SetFieldBackgroundColor(app.Styles.InverseTextColor)

	ShowModal(form, 60, 10)
	return nil
}

func handleTruncateTable(ctx CommandContext) error {
	tableName := ctx.CurrentTable

	modal := tview.NewModal()
	modal.SetText(fmt.Sprintf("Truncate table '%s'?\n\nThis will delete all rows but keep the table structure.", tableName))
	modal.AddButtons([]string{"Truncate", "Cancel"})
	modal.SetButtonBackgroundColor(app.Styles.InverseTextColor)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Truncate" {
			sql := fmt.Sprintf("TRUNCATE TABLE `%s`", tableName)
			_, err := ctx.DB.ExecuteDMLStatement(sql)
			if err != nil {
				ShowError(fmt.Sprintf("Failed to truncate table: %v", err))
				return
			}

			ShowSuccess(fmt.Sprintf("Table '%s' truncated successfully", tableName))
			RefreshTree()
		}
		CloseModal()
	})

	ShowModal(modal, 60, 8)
	return nil
}

func handleCopyTable(ctx CommandContext) error {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Copy Table ")

	sourceName := ctx.CurrentTable
	var targetName string
	var copyData bool = false

	form.AddInputField("Source Table", sourceName, 30, nil, nil).GetFormItem(0).(*tview.InputField).SetDisabled(true)
	form.AddInputField("New Table Name", sourceName+"_copy", 30, nil, func(text string) {
		targetName = text
	})
	form.AddCheckbox("Copy data", false, func(checked bool) {
		copyData = checked
	})

	form.AddButton("Copy", func() {
		if targetName == "" {
			ShowError("Target table name cannot be empty")
			return
		}

		// Create table with same structure
		sql := fmt.Sprintf("CREATE TABLE `%s` LIKE `%s`", targetName, sourceName)
		_, err := ctx.DB.ExecuteDMLStatement(sql)
		if err != nil {
			ShowError(fmt.Sprintf("Failed to copy table structure: %v", err))
			return
		}

		// Copy data if requested
		if copyData {
			sql = fmt.Sprintf("INSERT INTO `%s` SELECT * FROM `%s`", targetName, sourceName)
			_, err = ctx.DB.ExecuteDMLStatement(sql)
			if err != nil {
				ShowError(fmt.Sprintf("Failed to copy table data: %v", err))
				return
			}
		}

		ShowSuccess(fmt.Sprintf("Table copied to '%s'", targetName))
		CloseModal()
		RefreshTree()
	})

	form.AddButton("Cancel", func() {
		CloseModal()
	})

	form.SetButtonsAlign(tview.AlignCenter)
	form.SetFieldBackgroundColor(app.Styles.InverseTextColor)

	ShowModal(form, 60, 12)
	return nil
}
