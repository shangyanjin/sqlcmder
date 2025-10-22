# SQLCmder

A cross-platform Terminal UI database management tool with built-in command line interface (CMDER), written in Go. Personally customized and enhanced based on LazySQL.

## Features

### VI-Style Command Line (CMDER)
- Built-in command interpreter with `SQL#` prompt
- Quick database operations: `db create/drop/use/list`
- Quick table operations: `table create/drop/truncate/rename`
- Direct SQL execution
- Command history navigation (Up/Down arrows)
- Comprehensive help system: `help <topic>`

### Database Support
- PostgreSQL
- MySQL
- SQLite
- SQL Server

### User Interface
- Clean TUI with keyboard navigation
- Connection management with presets
- Real-time query results
- Split-panel layout (database tree + results + sidebar)
- Syntax-aware SQL editor

## Quick Start

```bash
# Run the application
./sqlcmder

# Connect to database
# Use connection form or command line

# In command line (Ctrl+\)
SQL# db list                    # List databases
SQL# db use mydb                # Switch database
SQL# SELECT * FROM users;       # Execute SQL
SQL# help insert                # Get SQL syntax help
```

## Keyboard Shortcuts

- `Ctrl+\` - Open command line
- `Ctrl+P` - Command palette
- `Ctrl+F` - Search tree
- `Ctrl+Left/Right` - Switch panels
- `Up/Down` - Navigate history (in command line)
- `?` or `help` - Show help

## Configuration

Config file location: `./config.toml` (next to executable)

## Project Structure

```
sqlcmder/
├── app/                         # Application core
├── components/                  # TUI components
│   ├── command_line.go         # CMDER implementation
│   ├── home.go                 # Main window
│   ├── tree.go                 # Database tree
│   └── ...
├── drivers/                    # Database drivers
├── models/                     # Data models
├── internal/                   # Internal packages
│   ├── history/               # Query history
│   └── saved/                 # Saved queries
├── docs/                      # Documentation
│   └── CHANGELOG.md
├── config.toml                # Configuration
└── main.go                    # Entry point
```

## Credits

Based on [LazySQL](https://github.com/jorgerojas26/lazysql) by Jorge Rojas.

## License

MIT License
