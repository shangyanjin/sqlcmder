package ui

import (
	"sqlcmder/internal/app"
	"sqlcmder/models"
)

var App = app.App

// Page name aliases from models package for backward compatibility
const (
	pageNameHelp                   = models.PageNameHelp
	pageNameConfirmation           = models.PageNameConfirmation
	pageNameConnections            = models.PageNameConnections
	pageNameDMLPreview             = models.PageNameDMLPreview
	pageNameErrorModal             = models.PageNameErrorModal
	pageNameTable                  = models.PageNameTable
	pageNameTableError             = models.PageNameTableError
	pageNameTableLoading           = models.PageNameTableLoading
	pageNameTableEditorTable       = models.PageNameTableEditorTable
	pageNameTableEditorResultsInfo = models.PageNameTableEditorResultsInfo
	pageNameTableEditCell          = models.PageNameTableEditCell
	pageNameQueryPreviewError      = models.PageNameQueryPreviewError
	pageNameJSONViewer             = models.PageNameJSONViewer
	pageNameSidebar                = models.PageNameSidebar
	pageNameConnectionSelection    = models.PageNameConnectionSelection
	pageNameConnectionForm         = models.PageNameConnectionForm
	pageNameSetValue               = models.PageNameSetValue
	pageNameQueryHistory           = models.PageNameQueryHistory
	pageNameSaveQuery              = models.PageNameSaveQuery
	pageNameSavedQueryDelete       = models.PageNameSavedQueryDelete
	pageNameCommandPalette         = models.PageNameCommandPalette
)

// Tab name aliases from models package
const (
	tabNameEditor            = models.TabNameEditor
	savedQueryTabReference   = models.SavedQueryTabReference
	queryHistoryTabReference = models.QueryHistoryTabReference
)

// Event name aliases from models package
const (
	eventSidebarEditing        = models.EventSidebarEditing
	eventSidebarUnfocusing     = models.EventSidebarUnfocusing
	eventSidebarToggling       = models.EventSidebarToggling
	eventSidebarCommitEditing  = models.EventSidebarCommitEditing
	eventSidebarError          = models.EventSidebarError
	eventSQLEditorQuery        = models.EventSQLEditorQuery
	eventSQLEditorEscape       = models.EventSQLEditorEscape
	eventResultsTableFiltering = models.EventResultsTableFiltering
	eventTreeSelectedDatabase  = models.EventTreeSelectedDatabase
	eventTreeSelectedTable     = models.EventTreeSelectedTable
	eventTreeIsFiltering       = models.EventTreeIsFiltering
)

// Menu item aliases from models package
const (
	menuRecords     = models.MenuRecords
	menuColumns     = models.MenuColumns
	menuConstraints = models.MenuConstraints
	menuForeignKeys = models.MenuForeignKeys
	menuIndexes     = models.MenuIndexes
)

// Action aliases from models package
const (
	actionNewConnection  = models.ActionNewConnection
	actionEditConnection = models.ActionEditConnection
)

// Focus and color aliases from models package
const (
	focusedWrapperLeft  = models.FocusedWrapperLeft
	focusedWrapperRight = models.FocusedWrapperRight
	colorTableChange    = models.ColorTableChange
	colorTableInsert    = models.ColorTableInsert
	colorTableDelete    = models.ColorTableDelete
)
