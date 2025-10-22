package ui

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
	"sqlcmder/cli"
	"sqlcmder/drivers"
	"sqlcmder/keymap"
	"sqlcmder/logger"
	"sqlcmder/models"
	"sqlcmder/data/history"
)

type Home struct {
	*tview.Flex
	Tree                 *Tree
	TabbedPane           *TabbedPane
	LeftWrapper          *tview.Flex
	RightWrapper         *tview.Flex
	HelpStatus           HelpStatus
	HelpModal            *HelpModal
	QueryHistoryModal    *QueryHistoryModal
	CommandPalette       *CommandPalette
	CommandLine          *CommandLine
	CommandStatusBar     *tview.TextView
	DBDriver             drivers.Driver
	FocusedWrapper       string
	ListOfDBChanges      []models.DBDMLChange
	ConnectionIdentifier string
	ConnectionURL        string
	CurrentDatabase      string
	CurrentTable         string
	Connection           models.Connection // Full connection details
}

func NewHomePage(connection models.Connection, dbdriver drivers.Driver) *Home {
	tree := NewTree(connection.DBName, dbdriver)
	leftWrapper := tview.NewFlex()
	rightWrapper := tview.NewFlex()

	maincontent := tview.NewFlex()

	connectionIdentifier := connection.Name
	if connectionIdentifier == "" {
		parsedURL, err := url.Parse(connection.DSN)
		if err == nil {
			connectionIdentifier = history.SanitizeFilename(parsedURL.Host + strings.ReplaceAll(parsedURL.Path, "/", "_"))
		} else {
			connectionIdentifier = "unnamed_or_invalid_url_connection"
		}
	}

	home := &Home{
		Flex:         tview.NewFlex().SetDirection(tview.FlexRow),
		Tree:         tree,
		LeftWrapper:  leftWrapper,
		RightWrapper: rightWrapper,
		HelpStatus:   NewHelpStatus(),
		HelpModal:    NewHelpModal(),

		DBDriver:             dbdriver,
		ListOfDBChanges:      []models.DBDMLChange{},
		ConnectionIdentifier: connectionIdentifier,
		ConnectionURL:        connection.DSN,
		Connection:           connection, // Store full connection
	}

	tabbedPane := NewTabbedPane()

	home.TabbedPane = tabbedPane

	qhm := NewQueryHistoryModal(connectionIdentifier, func(selectedQuery string) {
		home.createOrFocusEditorTab()

		currentTab := home.TabbedPane.GetCurrentTab()
		if currentTab != nil {
			table := currentTab.Content.(*ResultsTable)
			table.Editor.SetText(selectedQuery, true)
		}
	})

	home.QueryHistoryModal = qhm

	// Initialize command palette
	commandPalette := NewCommandPalette()
	commandPalette.OnClose = func() {
		mainPages.RemovePage(pageNameCommandPalette)
		app.App.SetFocus(home.Tree)
	}

	ctx := CommandContext{
		DB:              dbdriver,
		CurrentDatabase: connection.DBName,
		Connection:      connectionIdentifier,
		ConnectionModel: &connection,
	}
	commandPalette.SetContext(ctx)

	// Register commands
	RegisterDatabaseCommands(commandPalette)
	RegisterTableCommands(commandPalette)

	home.CommandPalette = commandPalette
	home.CurrentDatabase = connection.DBName

	// Create command line
	commandLine := NewCommandLine()
	commandLine.OnCommand = func(cmd string) {
		ctx := CommandContext{
			DB:              home.DBDriver,
			CurrentDatabase: home.CurrentDatabase,
			CurrentTable:    home.CurrentTable,
			Connection:      home.ConnectionIdentifier,
			ConnectionModel: &home.Connection,
		}
		commandLine.ExecuteCommand(cmd, ctx)
	}
	commandLine.OnCancel = func() {
		logger.Debug("CommandLine OnCancel - refocus to table", nil)
		// Just refocus back to table
		tab := home.TabbedPane.GetCurrentTab()
		if tab != nil {
			table := tab.Content.(*ResultsTable)
			app.App.SetFocus(table)
		}
	}
	home.CommandLine = commandLine

	// Create command status bar
	commandStatusBar := tview.NewTextView()
	commandStatusBar.SetDynamicColors(true)
	commandStatusBar.SetText(" [yellow]Ctrl+Left/Right[white]: Switch Panel | [yellow]Ctrl+P/K[white]: Command Palette | [yellow]Ctrl+\\[white]: Command Line | [yellow]Ctrl+F[white]: Search | [yellow]?[white]: Help")
	commandStatusBar.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)
	commandStatusBar.SetTextColor(app.Styles.PrimaryTextColor)
	home.CommandStatusBar = commandStatusBar

	go home.subscribeToTreeChanges()

	leftWrapper.SetBorderColor(app.Styles.InverseTextColor)
	leftWrapper.AddItem(tree.Wrapper, 0, 1, true)

	rightWrapper.SetBorderColor(app.Styles.InverseTextColor)
	rightWrapper.SetBorder(true) // Overall border like left panel
	rightWrapper.SetDirection(tview.FlexColumnCSS)
	rightWrapper.SetInputCapture(home.rightWrapperInputCapture)
	rightWrapper.AddItem(tabbedPane.HeaderContainer, 1, 0, false)
	rightWrapper.AddItem(tabbedPane.Pages, 0, 1, false)
	rightWrapper.AddItem(commandLine, 2, 0, false)      // Command line 2 rows, always visible
	rightWrapper.AddItem(commandStatusBar, 1, 0, false) // Status bar always visible

	maincontent.AddItem(leftWrapper, 30, 1, false)
	maincontent.AddItem(rightWrapper, 0, 5, false)

	home.AddItem(maincontent, 0, 1, false)
	// home.AddItem(home.HelpStatus, 1, 1, false)

	home.SetInputCapture(home.homeInputCapture)

	home.SetFocusFunc(func() {
		if home.FocusedWrapper == focusedWrapperLeft || home.FocusedWrapper == "" {
			home.focusLeftWrapper()
		} else {
			home.focusRightWrapper()
		}
	})

	mainPages.AddPage(connection.DSN, home, true, false)
	return home
}

func (home *Home) subscribeToTreeChanges() {
	ch := home.Tree.Subscribe()

	for stateChange := range ch {
		switch stateChange.Key {
		case eventTreeSelectedTable:
			databaseName := home.Tree.GetSelectedDatabase()
			tableName := stateChange.Value.(string)

			// Update context for command palette
			home.CurrentDatabase = databaseName
			home.CurrentTable = tableName
			home.UpdateCommandContext()

			tabReference := fmt.Sprintf("%s.%s", databaseName, tableName)

			tab := home.TabbedPane.GetTabByReference(tabReference)

			var table *ResultsTable

			if tab != nil {
				table = tab.Content.(*ResultsTable)
				home.TabbedPane.SwitchToTabByReference(tab.Reference)
			} else {
				table = NewResultsTable(&home.ListOfDBChanges, home.Tree, home.DBDriver, home.ConnectionIdentifier, home.ConnectionURL).WithFilter()
				table.SetDatabaseName(databaseName)
				table.SetTableName(tableName)

				home.TabbedPane.AppendTab(tableName, table, tabReference)
			}

			results := table.FetchRecords(func() {
				home.focusLeftWrapper()
			})

			// Show sidebar if there is more then 1 row (row 0 are
			// the column names) and the sidebar is not disabled.
			if !app.App.Config().DisableSidebar && len(results) > 1 && !table.GetShowSidebar() {
				table.ShowSidebar(true)
			}

			if table.state.error == "" {
				home.focusRightWrapper()
			}

			app.App.ForceDraw()
		case eventTreeIsFiltering:
			isFiltering := stateChange.Value.(bool)
			if isFiltering {
				home.SetInputCapture(nil)
			} else {
				home.SetInputCapture(home.homeInputCapture)
			}
		}
	}
}

func (home *Home) focusRightWrapper() {
	logger.Debug("Focus right wrapper", nil)
	home.Tree.RemoveHighlight()

	home.RightWrapper.SetBorderColor(app.Styles.PrimaryTextColor)
	home.LeftWrapper.SetBorderColor(app.Styles.InverseTextColor)
	home.TabbedPane.Highlight()
	tab := home.TabbedPane.GetCurrentTab()

	if tab != nil {
		home.focusTab(tab)
	}

	home.FocusedWrapper = focusedWrapperRight
}

func (home *Home) focusTab(tab *Tab) {
	if tab != nil {
		table := tab.Content.(*ResultsTable)
		table.HighlightAll()

		if table.GetIsFiltering() {
			go func() {
				if table.Filter != nil {
					app.App.SetFocus(table.Filter.Input)
					table.Filter.HighlightLocal()
				} else if table.Editor != nil {
					app.App.SetFocus(table.Editor)
					table.Editor.Highlight()
				}

				table.RemoveHighlightTable()
				app.App.Draw()
			}()
		} else {
			table.SetInputCapture(table.tableInputCapture)
			app.App.SetFocus(table)
		}

		if tab.Name == tabNameEditor {
			home.HelpStatus.SetStatusOnEditorView()
		} else {
			home.HelpStatus.SetStatusOnTableView()
		}
	}
}

func (home *Home) focusLeftWrapper() {
	logger.Debug("Focus left wrapper", nil)
	home.Tree.Highlight()

	home.RightWrapper.SetBorderColor(app.Styles.InverseTextColor)
	home.LeftWrapper.SetBorderColor(app.Styles.PrimaryTextColor)

	tab := home.TabbedPane.GetCurrentTab()

	if tab != nil {
		table := tab.Content.(*ResultsTable)

		table.RemoveHighlightAll()

	}

	home.TabbedPane.SetBlur()

	app.App.SetFocus(home.Tree)

	home.FocusedWrapper = focusedWrapperLeft
}

func (home *Home) rightWrapperInputCapture(event *tcell.EventKey) *tcell.EventKey {
	var tab *Tab

	command := keymap.Keymaps.Group(keymap.TableGroup).Resolve(event)

	switch command {
	case commands.TabPrev:

		tab := home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)
			if !table.GetIsEditing() && !table.GetIsFiltering() {
				home.TabbedPane.SwitchToPreviousTab()
				// home.focusTab(home.TabbedPane.SwitchToPreviousTab())
				return nil
			}

		}

		return event
	case commands.TabNext:
		tab := home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)
			if !table.GetIsEditing() && !table.GetIsFiltering() {
				home.TabbedPane.SwitchToNextTab()
				// home.focusTab(home.TabbedPane.SwitchToNextTab())
				return nil
			}
		}

		return event
	case commands.TabFirst:
		home.TabbedPane.SwitchToFirstTab()
		// home.focusTab(home.TabbedPane.SwitchToFirstTab())
		return nil
	case commands.TabLast:
		home.TabbedPane.SwitchToLastTab()
		// home.focusTab(home.TabbedPane.SwitchToLastTab())
		return nil
	case commands.TabClose:
		tab = home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)

			if !table.GetIsFiltering() && !table.GetIsEditing() && !table.GetIsLoading() {
				home.TabbedPane.RemoveCurrentTab()

				if home.TabbedPane.GetLength() == 0 {
					home.focusLeftWrapper()
					return nil
				}
			}
		}
	case commands.PagePrev:
		tab = home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)

			if ((table.Menu != nil && table.Menu.GetSelectedOption() == 1) ||
				table.Menu == nil) && !table.Pagination.GetIsFirstPage() && !table.GetIsLoading() {
				table.Pagination.SetOffset(table.Pagination.GetOffset() - table.Pagination.GetLimit())
				table.FetchRecords(nil)
			}
		}

	case commands.PageNext:
		tab = home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)

			if ((table.Menu != nil && table.Menu.GetSelectedOption() == 1) ||
				table.Menu == nil) && !table.Pagination.GetIsLastPage() && !table.GetIsLoading() {
				table.Pagination.SetOffset(table.Pagination.GetOffset() + table.Pagination.GetLimit())
				table.FetchRecords(nil)
			}
		}
	}

	return event
}

func (home *Home) homeInputCapture(event *tcell.EventKey) *tcell.EventKey {
	tab := home.TabbedPane.GetCurrentTab()

	var table *ResultsTable

	if tab != nil {
		table = tab.Content.(*ResultsTable)
	}

	// Log key events at debug level for Ctrl/Alt keys
	if event.Modifiers()&tcell.ModCtrl != 0 || event.Modifiers()&tcell.ModAlt != 0 {
		logger.Debug("Key event", map[string]any{
			"key":          event.Key(),
			"rune":         string(event.Rune()),
			"runeCode":     int(event.Rune()),
			"modifiers":    event.Modifiers(),
			"focusWrapper": home.FocusedWrapper,
			"isEditing":    table != nil && table.GetIsEditing(),
			"isFiltering":  table != nil && table.GetIsFiltering(),
		})
	}

	// Special logging for backslash-related keys
	if event.Rune() == '\\' || event.Rune() == 28 || event.Key() == tcell.KeyCtrlBackslash {
		logger.Debug("Backslash key detected", map[string]any{
			"key":       event.Key(),
			"rune":      string(event.Rune()),
			"runeCode":  int(event.Rune()),
			"modifiers": event.Modifiers(),
		})
	}

	// Only handle Ctrl+Arrow when not editing/filtering
	if table != nil && (table.GetIsEditing() || table.GetIsFiltering()) {
		logger.Debug("Skip key handling - editing or filtering", nil)
		return event
	}

	// Ctrl+Right: cycle focus right (tree -> table -> sidebar -> tree)
	if event.Key() == tcell.KeyRight && event.Modifiers()&tcell.ModCtrl != 0 {
		logger.Debug("Ctrl+Right pressed", map[string]any{"focusWrapper": home.FocusedWrapper})
		if home.FocusedWrapper == focusedWrapperLeft {
			// From tree to table
			home.focusRightWrapper()
		} else if home.FocusedWrapper == focusedWrapperRight {
			// From table to sidebar (if visible), otherwise to tree
			if table != nil && table.GetShowSidebar() {
				app.App.SetFocus(table.Sidebar)
			} else {
				home.focusLeftWrapper()
			}
		}
		return nil
	}

	// Ctrl+Left: cycle focus left (tree <- table <- sidebar <- tree)
	if event.Key() == tcell.KeyLeft && event.Modifiers()&tcell.ModCtrl != 0 {
		// Check if sidebar has focus
		hasSidebarFocus := false
		if table != nil && table.GetShowSidebar() {
			// Try to detect if sidebar has focus by checking current primitive
			currentFocus := app.App.GetFocus()
			if currentFocus == table.Sidebar {
				hasSidebarFocus = true
			}
		}

		if hasSidebarFocus {
			// From sidebar to table
			app.App.SetFocus(table)
		} else if home.FocusedWrapper == focusedWrapperRight {
			// From table to tree
			home.focusLeftWrapper()
		} else if home.FocusedWrapper == focusedWrapperLeft {
			// From tree to sidebar (if visible), otherwise to table
			if table != nil && table.GetShowSidebar() {
				app.App.SetFocus(table.Sidebar)
			} else {
				home.focusRightWrapper()
			}
		}
		return nil
	}

	// Ctrl+P or Ctrl+K to open command palette
	if event.Key() == tcell.KeyCtrlP || event.Key() == tcell.KeyCtrlK {
		if table == nil || (!table.GetIsEditing() && !table.GetIsFiltering()) {
			logger.Debug("Opening command palette", nil)
			home.ShowCommandPalette()
			return nil
		}
	}

	// Ctrl+\ to focus command line - SIMPLE FOCUS SWITCH
	if event.Key() == tcell.KeyCtrlBackslash {
		logger.Debug("Ctrl+\\ - Focus command line", nil)
		if table == nil || (!table.GetIsEditing() && !table.GetIsFiltering()) {
			// Clear and focus command line input field
			home.CommandLine.SetText("")
			app.App.SetFocus(home.CommandLine.InputField)
			logger.Debug("Command line focused", nil)
			return nil
		}
	}

	command := keymap.Keymaps.Group(keymap.HomeGroup).Resolve(event)

	if command != commands.Noop {
		logger.Debug("Command resolved", map[string]any{
			"command": command.String(),
			"key":     event.Key(),
		})
	}

	// Handle commands
	switch command {
	case commands.MoveLeft:
		if table != nil && !table.GetIsEditing() && !table.GetIsFiltering() && home.FocusedWrapper == focusedWrapperRight {
			home.focusLeftWrapper()
			return nil
		}
	case commands.MoveRight:
		if table != nil && !table.GetIsEditing() && !table.GetIsFiltering() && home.FocusedWrapper == focusedWrapperLeft {
			home.focusRightWrapper()
			return nil
		}
	case commands.SwitchToEditorView:
		home.createOrFocusEditorTab()
		return nil
	case commands.SwitchToConnectionsView:
		if (table != nil && !table.GetIsEditing() && !table.GetIsFiltering() && !table.GetIsLoading()) || table == nil {
			mainPages.SwitchToPage(pageNameConnections)
			return nil
		}
	case commands.Quit:
		if tab == nil || (!table.GetIsEditing() && !table.GetIsFiltering()) {
			app.App.Stop()
			return nil
		}
	case commands.Save:
		if (len(home.ListOfDBChanges) > 0) && !table.GetIsEditing() {
			queryPreviewModal := NewQueryPreviewModal(&home.ListOfDBChanges, home.DBDriver, func() {
				for _, change := range home.ListOfDBChanges {
					queryString, err := home.DBDriver.DMLChangeToQueryString(change)
					if err != nil {
						logger.Error("Failed to convert DML change to query string", map[string]any{"error": err})
						continue
					}
					err = history.AddQueryToHistory(home.ConnectionIdentifier, queryString)
					if err != nil {
						logger.Error("Failed to add query to history", map[string]any{"error": err})
					}
				}
				home.ListOfDBChanges = []models.DBDMLChange{}
				table.FetchRecords(nil)
				home.Tree.ForceRemoveHighlight()
			})

			mainPages.AddPage(pageNameDMLPreview, queryPreviewModal, true, true)
			return nil
		}
	case commands.HelpPopup:
		if table == nil || !table.GetIsEditing() {
			mainPages.AddPage(pageNameHelp, home.HelpModal, true, true)
			return nil
		}
	case commands.SearchGlobal:
		if table != nil && !table.GetIsEditing() && !table.GetIsFiltering() && !table.GetIsLoading() && home.FocusedWrapper == focusedWrapperRight {
			home.focusLeftWrapper()
		}

		home.Tree.ForceRemoveHighlight()
		home.Tree.ClearSearch()
		app.App.SetFocus(home.Tree.Filter)
		home.Tree.SetIsFiltering(true)
		return nil
	case commands.ToggleQueryHistory:
		if mainPages.HasPage(pageNameQueryHistory) {
			mainPages.SwitchToPage(pageNameQueryHistory)
		} else {
			mainPages.AddPage(pageNameQueryHistory, home.QueryHistoryModal, true, true)
		}

		home.QueryHistoryModal.queryHistoryComponent.LoadHistory(home.ConnectionIdentifier)
		return nil
	}

	return event
}

func (home *Home) createOrFocusEditorTab() {
	tab := home.TabbedPane.GetTabByName(tabNameEditor)

	if tab != nil {
		home.TabbedPane.SwitchToTabByName(tabNameEditor)
		table := tab.Content.(*ResultsTable)
		table.SetIsFiltering(true)
	} else {
		tableWithEditor := NewResultsTable(&home.ListOfDBChanges, home.Tree, home.DBDriver, home.ConnectionIdentifier, home.ConnectionURL).WithEditor()
		home.TabbedPane.AppendTab(tabNameEditor, tableWithEditor, tabNameEditor)
		tableWithEditor.SetIsFiltering(true)
		home.TabbedPane.GetCurrentTab()
	}

	home.HelpStatus.SetStatusOnEditorView()
	home.focusRightWrapper()
	App.ForceDraw()
}

func (home *Home) UpdateCommandContext() {
	ctx := CommandContext{
		DB:              home.DBDriver,
		CurrentDatabase: home.CurrentDatabase,
		CurrentTable:    home.CurrentTable,
		Connection:      home.ConnectionIdentifier,
		ConnectionModel: &home.Connection,
	}
	home.CommandPalette.SetContext(ctx)
}

func (home *Home) ShowCommandPalette() {
	home.UpdateCommandContext()
	home.CommandPalette.Show()
	mainPages.AddPage(pageNameCommandPalette, home.CommandPalette, true, true)
}
