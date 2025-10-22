# sqlcmder
A cross-platform Terminal CMD + TUI database management tool written in Go, personally customized and enhanced based on LazySQL, supporting PostgreSQL, MySQL, SQLite, and more.



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
