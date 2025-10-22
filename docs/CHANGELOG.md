# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

- **Table Quick Commands**
  - `table create <name>` - Create table (interactive)
  - `table drop <name>` - Drop table
  - `table truncate <name>` - Clear table data
  - `table rename <old> <new>` - Rename table

- **Comprehensive Help System**
  - Multiple help triggers: `help`, `?`, `/?`, `/help`
  - Detailed syntax help: `help insert`, `help update`, `help delete`, `help select`, `help db`, `help table`
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

