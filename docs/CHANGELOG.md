# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

