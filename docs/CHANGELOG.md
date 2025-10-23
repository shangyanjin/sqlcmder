# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### UI/UX Improvements

#### Command Palette Removal
- **Removed Command Palette functionality**
  - Deleted `ui/cmd_palette.go`, `ui/cmd_palette_database.go`, `ui/cmd_palette_table.go`
  - Removed `Ctrl+P/K` shortcuts for Command Palette
  - Simplified command line interface by removing complex command system
  - Updated status bar to remove Command Palette references

#### Command Line Removal
- **Removed Command Line input area**
  - Deleted `ui/cmd_line_input.go` file
  - Removed `Ctrl+\` shortcut for Command Line
  - Simplified UI by removing bottom command input area
  - Updated status bar to remove Command Line references
  - Simplified `ui/cmd_helpers.go` to only log messages instead of displaying in UI

#### Keyboard Shortcuts Updates
- **Updated global shortcuts**
  - Changed `Ctrl+F` to `Ctrl+\` for tree search functionality
  - Updated status bar hints to reflect new shortcut mappings
  - Updated README.md with current keyboard shortcuts

#### SQL Editor Improvements
- **Enhanced SQL editor placeholder text**
  - Changed from "Enter your SQL query here..." to "Input your SQL query here, press ctrl+R run, ESC return"
  - Added clear instructions for running queries (`Ctrl+R`) and returning (`ESC`)
  - Improved user guidance for SQL editor usage

#### Table Interface Enhancements
- **Added table shortcut hints in pagination area**
  - Displayed common table shortcuts: `c Edit`, `d Delete`, `o Add`, `<> Page`, `CTRL+s Commit`
  - Used color-coded shortcuts (yellow) with white action descriptions
  - Added separator between shortcuts and pagination info
  - Improved discoverability of table operations

#### Theme System Updates
- **Added PrimitiveBackgroundColor to dark theme**
  - Extended `ColorScheme` struct with `PrimitiveBackgroundColor` field
  - Defined specific `PrimitiveBackgroundColor` values for all themes
  - Updated theme initialization to use theme-defined colors instead of hardcoded values
  - Improved color consistency across UI components

#### Code Cleanup
- **Unified field edit colors**
  - Standardized database field editing colors to use `InverseTextColor` for background
  - Used `PrimaryTextColor` for foreground in field editing
  - Aligned field edit colors with dark theme settings
  - Improved visual consistency across table and sidebar editing

## [0.3.3] - 2025-10-22

### Theme System Overhaul

#### TrueColor Support
- **Fixed tview input box display anomaly in Linux dark terminals**
  - Replaced basic color constants (`tcell.ColorWhite`, `tcell.ColorBlack`) with TrueColor RGB values
  - Implemented `tcell.NewRGBColor(r, g, b)` for direct 24-bit color specification
  - Set environment variables to enable TrueColor: `TERM=xterm-256color`, `COLORTERM=truecolor`
  - Eliminated terminal theme mapping issues that caused white-on-white text

#### Multi-Theme System
- **Added comprehensive theme system with 5 built-in themes**
  - **Dark**: Soft whites and grays with reduced brightness for comfortable viewing
  - **Light**: Muted blacks and light grays with lower contrast
  - **Solarized Dark**: Ethan Schoonover's popular color scheme
  - **Gruvbox Dark**: Retro groove color scheme with warm tones
  - **Nord**: Arctic, north-bluish color palette
  
- **Theme configuration via `config.toml`**
  - Added `theme` field to `[application]` section
  - Supports runtime theme switching
  - Default theme: `dark`

#### Color Scheme Architecture
- **Created `ColorScheme` struct in `models/constants.go`**
  - 16 color fields covering all UI elements: text, borders, buttons, highlights, accents
  - Dedicated colors for focused/unfocused states
  - Specialized `SelectedTextColor` for table row selection with proper contrast
  
- **Unified color management**
  - `ColorSchemes` map stores all predefined themes
  - `ActiveColorScheme` tracks current theme
  - `SetActiveColorScheme()` function for theme switching
  - `GetColorScheme()` with fallback to dark theme

#### Enhanced Theme System in `cmd/app/app.go`
- **Extended `Theme` struct with custom color fields**
  - `ButtonBackgroundColor`: Darker button backgrounds (F1, F2, F3, Esc)
  - `UnfocusedBorderColor`: Lighter borders for inactive panels
  - `UnfocusedTextColor`: Slightly dimmed text for inactive panels
  - `UnfocusedAccentColor`: Unified gray for unfocused colored elements
  - `SelectedTextColor`: High-contrast text on selected table rows
  
- **Refactored color initialization**
  - New `initializeTheme()` function for centralized theme setup
  - `ApplyTheme()` method loads theme from configuration
  - Seamless integration with tview's theme system

#### UI Component Color Consistency
- **Unified unfocused state colors across all components**
  - Database tree panel: borders, graphics, title, node text, filter label
  - Table results: borders, title, data rows
  - SQL editor: border, text style
  - Filter inputs: labels, placeholder text, field text
  - Menu components: tabbed panes, table menus
  - Command palette and command line input
  
- **Improved focus state differentiation**
  - Minimal difference between focused/unfocused states
  - Consistent application of accent colors (yellow, green)
  - Unified gray tones for all unfocused colored text
  
- **Optimized table row selection contrast**
  - Yellow background with deep gray/dark text in dark themes
  - Yellow background with white text in light theme
  - Ensures readability while maintaining visual consistency

#### Color Refinements
- **Reduced excessive contrast in dark/light themes**
  - Dark: Softened whites to RGB(220,220,220), input backgrounds to RGB(45,45,48)
  - Light: Muted blacks to RGB(50,50,50), highlights to RGB(70,130,200)
  - Gentle borders and accents for comfortable long-term viewing
  
- **Button styling improvements**
  - Darker backgrounds for F1/F2/F3/Esc shortcut buttons
  - Better visual separation from surrounding elements
  
- **Placeholder text consistency**
  - Unfocused: `UnfocusedTextColor` for all input placeholders
  - Focused: `PrimaryTextColor` for clear visibility

### Technical Benefits
- **Eliminated terminal color mapping issues**: Direct RGB control prevents theme conflicts
- **Enhanced accessibility**: Multiple themes support different viewing preferences and environments
- **Improved visual comfort**: Reduced contrast prevents eye strain during extended use
- **Consistent color semantics**: Unified approach to focused/unfocused states across all UI components
- **Maintainable architecture**: Centralized color management simplifies future theme additions

## [0.3.2] - 2025-10-22

### Data Layer Reorganization

#### Storage Directory Restructuring
- **Renamed `storage/` to `data/` for better semantic clarity**
  - Improved directory naming to reflect data storage purpose
  - Enhanced project structure readability and maintainability
  - Aligned with common project organization conventions

#### Query Management Optimization
- **Renamed `saved/` subdirectory to `queries/`**
  - More descriptive naming for saved query functionality
  - Clearer separation between query storage and history tracking
  - Updated package declaration: `package saved` → `package queries`

#### Import Path Updates
- **Updated all import references across the codebase**
  - `sqlcmder/storage` → `sqlcmder/data`
  - `sqlcmder/data/saved` → `sqlcmder/data/queries`
  - Fixed package references in UI components (`saved.` → `queries.`)
  - Maintained backward compatibility and functionality

#### Directory Structure Enhancement
- **Final data layer structure:**
  ```
  data/
  ├── queries/        # Saved SQL queries management
  └── history/        # Query execution history tracking
  ```
- **Benefits:**
  - Clearer semantic meaning for data storage operations
  - Better separation of concerns between query persistence and history
  - Improved code maintainability and developer experience
  - Enhanced project structure professional appearance

## [0.3.1] - 2025-10-22

### Documentation Enhancement

#### MVC Architecture Planning
- **Added comprehensive MVC refactoring plan**
  - Created detailed `docs/plan.md` with English documentation
  - Outlined 4-phase refactoring approach (Service → Controller → View → Repository)
  - Defined clear separation of concerns for each layer
  - Provided implementation guidelines and risk mitigation strategies
  - Included success metrics and timeline estimates (10-13 weeks total)

#### Architecture Documentation Benefits
- **Clear roadmap for future development**
  - Identified current architectural issues and improvement opportunities
  - Proposed clean MVC structure with proper layer separation
  - Defined responsibilities for Model, View, Controller, Service, and Repository layers
  - Provided phased implementation approach to minimize disruption
  - Enhanced project maintainability and scalability planning

## [0.3.0] - 2025-10-22

### Major Package Structure Refactoring

#### Configuration Management Separation
- **Created dedicated `internal/config/` package**
  - Moved `config.go` from `internal/app/` to `internal/config/`
  - Separated configuration logic from application core
  - Updated package declaration: `package app` → `package config`
  - Resolved circular import issues between app and config packages

#### Keymap System Reorganization  
- **Moved keymap system to `internal/keymap/` package**
  - Renamed `keymap_system.go` to `config.go` for better clarity
  - Consolidated all keymap-related files under single package
  - Updated all UI components to import from `sqlcmder/internal/keymap`
  - Removed redundant `_keymap` suffixes from filenames

#### File Naming Standardization
- **Eliminated generic suffixes across all packages**
  - `cmd_palette_main.go` → `cmd_palette.go` (avoided `main` naming)
  - `manager_history.go` → `query_history.go` (business-specific naming)
  - `manager_saved.go` → `saved_queries.go` (business-specific naming)
  - `core.go` → `app.go` (removed generic `core` suffix)
  - `command_helper.go` → `command.go` (removed redundant `_helper`)
  - `utils_helper.go` → `utils.go` (removed redundant `_helper`)
  - `sqlite_driver.go` → merged into `sqlite.go` (eliminated duplicate)

#### Package Dependency Optimization
- **Resolved circular import cycles**
  - Modified `LoadConfig()` to accept config parameter instead of accessing global state
  - Added `GetConfig()` method to Application for accessing full configuration object
  - Updated all import paths and function calls throughout codebase
  - Maintained backward compatibility while improving architecture

#### Model Consolidation
- **Merged `saved_query.go` into `models.go`**
  - Consolidated all model definitions in single file
  - Reduced file fragmentation
  - Improved code organization

#### Build System Improvements
- **Fixed all compilation errors**
  - Resolved import path issues after package moves
  - Removed unused imports across all files
  - Ensured clean compilation with `go build`
  - Maintained full functionality after refactoring

### Technical Benefits
- **Better separation of concerns**: Configuration, keymaps, and application logic are now properly separated
- **Improved maintainability**: Clear package boundaries and consistent naming conventions
- **Reduced complexity**: Eliminated circular dependencies and generic naming
- **Enhanced readability**: Business-specific file names make code purpose immediately clear
- **Future-proof architecture**: Clean package structure supports easier feature additions

## [0.2.0] - 2025-10-22

### Major Refactoring - Project Structure Reorganization

#### Directory Structure
- **Renamed `components/` to `ui/`**
  - Better naming: explicitly indicates Terminal UI layer
  - Updated all package declarations: `package components` → `package ui`
  - Updated all import paths: `sqlcmder/components` → `sqlcmder/ui`
  
- **Unified command logic under `internal/commands/`**
  - Adopted `cmd_` prefix naming convention for consistency
  - `cmd_types.go` - Command enum types (moved from `commands/`)
  - `cmd_context.go` - Command execution context
  - `cmd_database.go` - Database command handlers
  - `cmd_table.go` - Table command handlers
  - `cmd_backup.go` - Backup/import functionality
  - `cmd_sql.go` - SQL execution handler
  - `cmd_utils.go` - Utility functions (e.g., Contains)
  
 - **Reorganized storage layer under `internal/storage/`**
   - `internal/storage/history/` - Query execution history (JSON)
   - `internal/storage/saved/` - Saved query templates (TOML)
   - Clear separation: storage vs business logic
  
- **Moved all core packages to `internal/`**
  - `internal/app/` - Application core (from `app/`)
  - `internal/keymap/` - Keyboard mappings (from `keymap/`)
  - `internal/lib/` - Utilities (from `lib/`)
  - `internal/helpers/` - Helper functions (from `helpers/`)
  - `internal/drivers/` - Database drivers (from `drivers/`)
  - Follows Go best practices: `internal/` packages not importable externally

#### Model Field Renaming (Standard Database Terminology)
- **Connection model fields:**
  - `URL` → `DSN` (Data Source Name - industry standard)
  - `URLParams` → `DSNParams` (DSN Parameters)
  - `Provider` → `Driver` (Database Driver - more accurate)
- **Config function renamed:**
  - `parseConfigURL()` → `parseConfigDSN()`
- **Benefits:**
  - Standard database terminology throughout codebase
  - Clearer, more professional naming
  - Better code documentation

#### Configuration Improvements
- **Relative paths:** `ConfigFile` now uses `./config.toml` instead of absolute path
  - Portable configuration across environments
  - Works regardless of installation location
  
#### UI/UX Improvements
- **Removed emojis and special characters** for better terminal compatibility
  - `✗` → `ERROR:`
  - `✓` → `OK:`
  - `ℹ` → `INFO:`
  - Sorting arrows → `ASC`/`DESC` text
  - Pure ASCII characters work in all terminals
  - No UTF-8 encoding issues
  
#### Code Quality
- **Fixed circular dependencies:**
  - Removed `helpers` import from `commands`
  - Created `cmd_utils.go` with `Contains()` function
- **Deleted duplicate code:**
  - Removed `components/commands/database_commands.go` (duplicate)
  - Consolidated command palette registrations
- **Consistent naming conventions:**
  - All command files use `cmd_` prefix
  - Clear separation between UI and business logic
  
#### Build & Compilation
- All changes compile successfully
- No breaking changes to functionality
- Improved code organization and maintainability

#### Package Consolidation
- **Merged internal/lib into internal/helpers**
  - Moved clipboard functionality from `internal/lib/clipboard.go` to `internal/helpers/clipboard.go`
  - Updated package declaration: `package lib` → `package helpers`
  - Updated all 6 UI component imports and usage: `lib.NewClipboard()` → `helpers.NewClipboard()`
  - Removed empty `internal/lib/` directory
  - Reduced package fragmentation: 8 packages → 7 packages
  - Consolidated utility functions for better organization

### Added
- **VI-Style Command Line (CMDER)** - Major new feature
  - Built-in command interpreter accessible via `Ctrl+\` or typing `:`
  - Two-row interface: system messages (row 1) + user input with `SQL#` prompt (row 2)
  - Real-time command execution for database operations
  - Command history navigation with Up/Down arrow keys
  - Auto-focus management: messages display and return focus to input automatically
  - Supports direct SQL execution and built-in commands
  
- **Database Quick Commands**
  - `db create <name>` - Create database
  - `db drop <name>` - Drop database  
  - `db use <name>` - Switch database
  - `db list` - List all databases
  - `db backup <file>` - Backup current database to ./backup/ directory
  - `db import <file>` - Import SQL from backup file

- **Table Quick Commands**
  - `table create <name>` - Create table (interactive)
  - `table drop <name>` - Drop table
  - `table truncate <name>` - Clear table data
  - `table rename <old> <new>` - Rename table

- **Database Backup & Import**
  - Cross-platform backup and restore functionality for all supported databases
  - MySQL: Uses `mysqldump` and `mysql` client tools
  - PostgreSQL: Uses `pg_dump` and `psql` client tools
  - SQLite: Direct file copy (no external dependencies)
  - MSSQL: Uses `sqlcmd` command-line tool
  - Automatic backup directory creation (./backup/)
  - Smart file lookup: checks current directory and ./backup/ automatically
  - Full context help: `help backup` and `help import` for detailed documentation

- **Comprehensive Help System**
  - Multiple help triggers: `help`, `?`, `/?`, `/help`
  - Detailed syntax help: `help insert`, `help update`, `help delete`, `help select`, `help db`, `help table`, `help backup`, `help import`
  - Modal dialog with scrollable content showing SQL syntax, examples, and best practices
  - Context-aware command suggestions

- Connection selection screen enhancements
  - Hint bar showing available shortcuts: "Up/Down Select, Enter Connect, New, Edit, Delete, Quit"
  - Selected connection marked with `*` prefix in yellow
- Navicat-style two-column connection form layout for better space utilization
- Auto-generated DSN field that updates in real-time as form fields change
- Database preset shortcuts: Alt+P (PostgreSQL), Alt+M (MySQL), Alt+S (SQLite), Alt+Q (SQL Server)
- Smart tab navigation: Tab key cycles through fields row by row (left to right)
- Default credentials for common databases:
  - PostgreSQL: `postgres/postgres`
  - MySQL: `root/root`
  - SQLite: default path `./sqlite.db`
- Save confirmation with config file path display
- Support for `.exe` and `.exe~` files in `.gitignore`

### Changed
- README installation instructions: changed from binary downloads to git clone + build from source
- Added prominent disclaimer: marked as BETA/TEST version with warnings about production use
- Renamed project from LazySQL to SQLCmder throughout codebase
- UI layout improvements for cleaner interface
  - Middle window has single overall border (similar to left panel)
  - Internal components (table, menu, filter, pagination) now borderless
  - Right sidebar border properly displayed
- Database Type changed from dropdown to input field for more flexibility
- Password field now displays text instead of masking (can be edited freely)
- Connection form fields reorganized:
  - Left column: Connection Name, Username, Password, DB Name
  - Right column: Database Type, Hostname, Port, DSN (Auto)
- F1 Save now shows warnings but saves anyway (non-blocking validation)
- Status hints improved: "Preset: [type] | Use Tab to navigate between fields"
- Command line messages now truncated to 100 characters to prevent layout breaking

### Fixed
- Command line interaction no longer causes right panel display issues (forced UI redraw)
- SQL error messages properly truncated to avoid breaking table layout
- Config file save functionality - now properly saves to executable directory
- Tab navigation between left and right column forms
- Field focus defaults to Connection Name on form open
- Connection selection hint text compatibility (changed from arrows to "Up/Down" for better terminal support)

## [1.0.0] - 2025-10-22

### Changed
- Config file default location changed from system config directory to executable directory
  - Windows: `%APPDATA%\config.toml` → `.\config.toml`
  - Linux/macOS: `~/.config/config.toml` → `./config.toml`
  - Enables portable deployment without system dependencies

