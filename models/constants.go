package models

import "github.com/gdamore/tcell/v2"

// Page names for navigation
const (
	// General pages
	PageNameHelp         = "Help"
	PageNameConfirmation = "Confirmation"
	PageNameConnections  = "Connections"
	PageNameDMLPreview   = "DMLPreview"
	PageNameErrorModal   = "ErrorModal"

	// Results table pages
	PageNameTable                  = "Table"
	PageNameTableError             = "TableError"
	PageNameTableLoading           = "TableLoading"
	PageNameTableEditorTable       = "TableEditorTable"
	PageNameTableEditorResultsInfo = "TableEditorResultsInfo"
	PageNameTableEditCell          = "TableEditCell"
	PageNameQueryPreviewError      = "QueryPreviewError"
	PageNameJSONViewer             = "json_viewer"

	// Sidebar page
	PageNameSidebar = "Sidebar"

	// Connection pages
	PageNameConnectionSelection = "ConnectionSelection"
	PageNameConnectionForm      = "ConnectionForm"

	// SetValueList page
	PageNameSetValue = "SetValue"

	// Query History pages
	PageNameQueryHistory     = "QueryHistoryModal"
	PageNameSaveQuery        = "SaveQueryModal"
	PageNameSavedQueryDelete = "SavedQueryDeleteModal"

	// Command Palette page
	PageNameCommandPalette = "CommandPalette"
)

// Tab names
const (
	TabNameEditor = "Editor"

	SavedQueryTabReference   = "saved_queries"
	QueryHistoryTabReference = "query_history"
)

// Event names
const (
	EventSidebarEditing       = "EditingSidebar"
	EventSidebarUnfocusing    = "UnfocusingSidebar"
	EventSidebarToggling      = "TogglingSidebar"
	EventSidebarCommitEditing = "CommitEditingSidebar"
	EventSidebarError         = "ErrorSidebar"

	EventSQLEditorQuery  = "Query"
	EventSQLEditorEscape = "Escape"

	EventResultsTableFiltering = "FilteringResultsTable"

	EventTreeSelectedDatabase = "SelectedDatabase"
	EventTreeSelectedTable    = "SelectedTable"
	EventTreeIsFiltering      = "IsFiltering"
)

// Results table menu items
const (
	MenuRecords     = "Records"
	MenuColumns     = "Columns"
	MenuConstraints = "Constraints"
	MenuForeignKeys = "Foreign Keys"
	MenuIndexes     = "Indexes"
)

// Connection actions
const (
	ActionNewConnection  = "NewConnection"
	ActionEditConnection = "EditConnection"
)

// Focus and UI state constants
const (
	FocusedWrapperLeft  = "left"
	FocusedWrapperRight = "right"

	ColorTableChange = tcell.ColorOrange
	ColorTableInsert = tcell.ColorDarkGreen
	ColorTableDelete = tcell.ColorRed
)
