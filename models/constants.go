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

// Theme names
const (
	ThemeDark      = "dark"
	ThemeLight     = "light"
	ThemeSolarized = "solarized"
	ThemeGruvbox   = "gruvbox"
	ThemeNord      = "nord"
)

// ColorScheme defines a complete color scheme for the TUI
type ColorScheme struct {
	Name              string
	TextColor         tcell.Color // Default text color
	MutedText         tcell.Color // Secondary/muted text
	InputBg           tcell.Color // Input field background
	InputText         tcell.Color // Input field text
	Border            tcell.Color // Border color
	Highlight         tcell.Color // Highlight/accent color
	ButtonBg          tcell.Color // Button background color
	UnfocusedBorder   tcell.Color // Unfocused panel border (lighter/dimmer)
	UnfocusedText     tcell.Color // Unfocused panel text (slightly dimmed)
	AccentYellow      tcell.Color // Accent color for highlights (focused)
	AccentGreen       tcell.Color // Accent color for success/info (focused)
	UnfocusedAccent   tcell.Color // Dimmed accent color (unfocused)
	SelectedTextColor tcell.Color // Text color on selected row (high contrast with AccentYellow)
}

// Predefined color schemes using TrueColor to avoid terminal theme mapping issues
var ColorSchemes = map[string]*ColorScheme{
	ThemeDark: {
		Name:              "Dark",
		TextColor:         tcell.NewRGBColor(220, 220, 220), // Soft white (降低亮度)
		MutedText:         tcell.NewRGBColor(150, 150, 150), // Medium gray
		InputBg:           tcell.NewRGBColor(45, 45, 48),    // 柔和深灰 (略微提亮)
		InputText:         tcell.NewRGBColor(210, 210, 210), // 柔和白色
		Border:            tcell.NewRGBColor(100, 100, 100), // 适中灰色
		Highlight:         tcell.NewRGBColor(100, 180, 230), // 柔和青色 (降低饱和度)
		ButtonBg:          tcell.NewRGBColor(60, 60, 65),    // 深灰色按钮背景
		UnfocusedBorder:   tcell.NewRGBColor(80, 80, 80),    // 未聚焦面板边框 (略微变暗，差别不大)
		UnfocusedText:     tcell.NewRGBColor(160, 160, 160), // 未聚焦文字 (略微变暗)
		AccentYellow:      tcell.NewRGBColor(230, 220, 100), // 柔和黄色
		AccentGreen:       tcell.NewRGBColor(120, 200, 120), // 柔和绿色
		UnfocusedAccent:   tcell.NewRGBColor(140, 140, 140), // 失去焦点的彩色变灰
		SelectedTextColor: tcell.NewRGBColor(40, 40, 40),    // 选中行文字（深色，与黄色背景高对比）
	},
	ThemeLight: {
		Name:              "Light",
		TextColor:         tcell.NewRGBColor(50, 50, 50),    // 柔和深灰 (不是纯黑)
		MutedText:         tcell.NewRGBColor(120, 120, 120), // 中等灰色
		InputBg:           tcell.NewRGBColor(250, 250, 250), // 极浅灰 (不是纯白)
		InputText:         tcell.NewRGBColor(40, 40, 40),    // 深灰文字
		Border:            tcell.NewRGBColor(200, 200, 200), // 浅灰边框
		Highlight:         tcell.NewRGBColor(70, 130, 200),  // 柔和蓝色 (降低亮度)
		ButtonBg:          tcell.NewRGBColor(230, 230, 230), // 中灰色按钮背景
		UnfocusedBorder:   tcell.NewRGBColor(210, 210, 210), // 未聚焦面板边框 (略微变淡，差别不大)
		UnfocusedText:     tcell.NewRGBColor(120, 120, 120), // 未聚焦文字 (略微变淡)
		AccentYellow:      tcell.NewRGBColor(200, 160, 0),   // 深黄色
		AccentGreen:       tcell.NewRGBColor(50, 150, 50),   // 深绿色
		UnfocusedAccent:   tcell.NewRGBColor(140, 140, 140), // 失去焦点的彩色变灰
		SelectedTextColor: tcell.NewRGBColor(255, 255, 255), // 选中行文字（白色，与深黄色背景高对比）
	},
	ThemeSolarized: {
		Name:              "Solarized Dark",
		TextColor:         tcell.NewRGBColor(131, 148, 150), // Base0
		MutedText:         tcell.NewRGBColor(88, 110, 117),  // Base01
		InputBg:           tcell.NewRGBColor(0, 43, 54),     // Base03
		InputText:         tcell.NewRGBColor(147, 161, 161), // Base1
		Border:            tcell.NewRGBColor(88, 110, 117),  // Base01
		Highlight:         tcell.NewRGBColor(42, 161, 152),  // Cyan
		ButtonBg:          tcell.NewRGBColor(7, 54, 66),     // Base02
		UnfocusedBorder:   tcell.NewRGBColor(70, 90, 100),   // 略微变暗的 Base01
		UnfocusedText:     tcell.NewRGBColor(101, 123, 131), // 略微变暗的文字
		AccentYellow:      tcell.NewRGBColor(181, 137, 0),   // Yellow
		AccentGreen:       tcell.NewRGBColor(133, 153, 0),   // Green
		UnfocusedAccent:   tcell.NewRGBColor(88, 110, 117),  // Base01 dimmed
		SelectedTextColor: tcell.NewRGBColor(0, 43, 54),     // Base03 deep contrast
	},
	ThemeGruvbox: {
		Name:              "Gruvbox Dark",
		TextColor:         tcell.NewRGBColor(235, 219, 178), // fg
		MutedText:         tcell.NewRGBColor(168, 153, 132), // fg2
		InputBg:           tcell.NewRGBColor(40, 40, 40),    // bg0_h
		InputText:         tcell.NewRGBColor(235, 219, 178), // fg
		Border:            tcell.NewRGBColor(80, 73, 69),    // bg2
		Highlight:         tcell.NewRGBColor(131, 165, 152), // aqua
		ButtonBg:          tcell.NewRGBColor(60, 56, 54),    // bg1
		UnfocusedBorder:   tcell.NewRGBColor(70, 65, 62),    // 略微变暗的 bg2
		UnfocusedText:     tcell.NewRGBColor(200, 186, 155), // 略微变暗的 fg
		AccentYellow:      tcell.NewRGBColor(250, 189, 47),  // yellow
		AccentGreen:       tcell.NewRGBColor(184, 187, 38),  // green
		UnfocusedAccent:   tcell.NewRGBColor(168, 153, 132), // fg2 dimmed
		SelectedTextColor: tcell.NewRGBColor(40, 40, 40),    // bg0_h deep contrast
	},
	ThemeNord: {
		Name:              "Nord",
		TextColor:         tcell.NewRGBColor(236, 239, 244), // Snow Storm
		MutedText:         tcell.NewRGBColor(216, 222, 233), // Polar Night
		InputBg:           tcell.NewRGBColor(46, 52, 64),    // Polar Night 0
		InputText:         tcell.NewRGBColor(236, 239, 244), // Snow Storm
		Border:            tcell.NewRGBColor(76, 86, 106),   // Polar Night 2
		Highlight:         tcell.NewRGBColor(136, 192, 208), // Frost
		ButtonBg:          tcell.NewRGBColor(59, 66, 82),    // Polar Night 1
		UnfocusedBorder:   tcell.NewRGBColor(65, 75, 95),    // 略微变暗的 Polar Night 2
		UnfocusedText:     tcell.NewRGBColor(200, 205, 215), // 略微变暗的 Snow Storm
		AccentYellow:      tcell.NewRGBColor(235, 203, 139), // Aurora yellow
		AccentGreen:       tcell.NewRGBColor(163, 190, 140), // Aurora green
		UnfocusedAccent:   tcell.NewRGBColor(143, 157, 180), // Dimmed aurora
		SelectedTextColor: tcell.NewRGBColor(46, 52, 64),    // Polar Night 0 deep contrast
	},
}

// Current active color scheme (default: dark)
var ActiveColorScheme *ColorScheme

// GetColorScheme returns a color scheme by name, defaults to dark theme if not found
func GetColorScheme(name string) *ColorScheme {
	if scheme, ok := ColorSchemes[name]; ok {
		return scheme
	}
	return ColorSchemes[ThemeDark]
}

// SetActiveColorScheme sets the active color scheme
func SetActiveColorScheme(name string) {
	ActiveColorScheme = GetColorScheme(name)
}

// Legacy color constants for backward compatibility (using active color scheme)
var (
	ColorTextDefault = tcell.NewRGBColor(255, 255, 255) // Pure white text
	ColorTextMuted   = tcell.NewRGBColor(180, 180, 180) // Muted gray text
	ColorBackground  = tcell.NewRGBColor(25, 25, 25)    // Dark background
	ColorBorder      = tcell.NewRGBColor(90, 90, 90)    // Medium gray border
	ColorHighlight   = tcell.NewRGBColor(0, 200, 255)   // Cyan highlight
	ColorInputBg     = tcell.NewRGBColor(30, 30, 30)    // Input field background
	ColorInputText   = tcell.NewRGBColor(255, 255, 255) // Input field text
)
