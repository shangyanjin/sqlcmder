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
├── cmd/
│   └── sqlcmder/                 # Entry point (main.go)
│
├── internal/
│   ├── tui/                      # Terminal UI components (bubbletea, tview, etc.)
│   ├── commands/                 # Command handler (query, export, import, etc.)
│   ├── db/                       # Database driver layer (pgsql/mysql/sqlite)
│   ├── config/                   # Config system (TOML/YAML/ENV)
│   ├── model/                    # Shared structs (DB config, query result, etc.)
│   ├── logger/                   # Logging system (zap/logrus/custom)
│   ├── utils/                    # Common helpers (string, file, env, time, etc.)
│   └── backup/                   # Database backup & restore logic
│
├── scripts/                      # Utility scripts (build, release, clean, test)
│   ├── build.sh                  # Cross-compile script for Linux/Mac/Win
│   ├── release.ps1               # Windows build + zip packaging
│   ├── backup_db.sh              # CLI database backup helper
│   └── init_config.ps1           # Generate default config for Windows
│
├── docs/                         # Documentation
│   ├── README.md                 # Main project doc
│   ├── CONFIG.md                 # Config format & examples
│   ├── COMMANDS.md               # CLI usage reference
│   ├── DB_SUPPORT.md             # Supported databases and drivers
│   └── DEV_GUIDE.md              # Developer contribution guide
│
├── tmp/                          # Temp files (session cache, query history)
│   ├── logs/                     # Runtime logs (if not system log)
│   └── query_cache/              # Cached query results (optional)
│
├── backup/                       # Auto or manual database backups
│   ├── pgsql/                    # PostgreSQL dumps
│   ├── mysql/                    # MySQL dumps
│   └── sqlite/                   # SQLite copies
│
├── examples/                     # Example config & query templates
│   ├── sample_config.toml
│   └── example_queries.sql
│
├── .env                          # Environment variables (optional)
├── .gitignore
├── go.mod
└── go.sum
```

## Credits

Based on [LazySQL](https://github.com/jorgerojas26/lazysql) by Jorge Rojas.

## License

MIT License
