package components

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/app"
	"sqlcmder/helpers/logger"
)

// CommandLine is a 2-row command line with message display and input
type CommandLine struct {
	*tview.Flex
	MessageView  *tview.TextView
	InputField   *tview.InputField
	OnCommand    func(cmd string)
	OnCancel     func()
	History      []string // Command history
	HistoryIndex int      // Current position in history (-1 means not browsing)
	TempInput    string   // Temporary storage for current input when browsing history
}

// NewCommandLine creates a new command line with 2 rows
func NewCommandLine() *CommandLine {
	// Message display (row 1)
	messageView := tview.NewTextView()
	messageView.SetDynamicColors(true)
	messageView.SetText("[grey]Ready")
	messageView.SetTextAlign(tview.AlignLeft)
	messageView.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)

	// Input field (row 2)
	inputField := tview.NewInputField()
	inputField.SetLabel("SQL# ")
	inputField.SetFieldWidth(0)
	inputField.SetFieldBackgroundColor(app.Styles.PrimitiveBackgroundColor)
	inputField.SetLabelColor(tcell.ColorYellow)
	inputField.SetFieldTextColor(app.Styles.PrimaryTextColor)

	// Container
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(messageView, 1, 0, false) // Row 1: message
	flex.AddItem(inputField, 1, 0, true)   // Row 2: input

	cl := &CommandLine{
		Flex:         flex,
		MessageView:  messageView,
		InputField:   inputField,
		History:      []string{},
		HistoryIndex: -1,
	}

	// Set reference for helper functions
	globalCommandLine = cl

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			cmd := strings.TrimSpace(inputField.GetText())
			if cmd != "" {
				// Add to history
				cl.AddToHistory(cmd)

				if cl.OnCommand != nil {
					cl.OnCommand(cmd)
				}
			}
			inputField.SetText("")
			cl.HistoryIndex = -1 // Reset history browsing
			// Don't call OnCancel here - let ShowError/ShowSuccess handle focus
			// Message will be set by command handlers
			return nil

		case tcell.KeyEsc:
			inputField.SetText("")
			cl.ShowMessage("[grey]Ready", false)
			cl.HistoryIndex = -1 // Reset history browsing
			// Esc should return focus to table
			if cl.OnCancel != nil {
				cl.OnCancel()
			}
			return nil

		case tcell.KeyUp:
			// Browse history backward (older commands)
			cl.NavigateHistory(-1, inputField)
			return nil

		case tcell.KeyDown:
			// Browse history forward (newer commands)
			cl.NavigateHistory(1, inputField)
			return nil
		}

		// Reset history index when user types
		if event.Key() == tcell.KeyRune {
			cl.HistoryIndex = -1
		}

		return event
	})

	return cl
}

// ShowMessage displays a message in the first row
func (cl *CommandLine) ShowMessage(message string, isError bool) {
	cl.MessageView.SetText(message)
}

// SetText sets the input text
func (cl *CommandLine) SetText(text string) {
	cl.InputField.SetText(text)
}

// GetText gets the input text
func (cl *CommandLine) GetText() string {
	return cl.InputField.GetText()
}

// SetLabel sets the input label
func (cl *CommandLine) SetLabel(label string) {
	cl.InputField.SetLabel(label)
}

// AddToHistory adds a command to history
func (cl *CommandLine) AddToHistory(cmd string) {
	// Don't add duplicate if it's the same as the last command
	if len(cl.History) > 0 && cl.History[len(cl.History)-1] == cmd {
		return
	}

	cl.History = append(cl.History, cmd)

	// Limit history to 100 commands
	if len(cl.History) > 100 {
		cl.History = cl.History[1:]
	}

	logger.Debug("Command added to history", map[string]any{
		"command":     cmd,
		"historySize": len(cl.History),
	})
}

// NavigateHistory navigates through command history
// direction: -1 for up (older), 1 for down (newer)
func (cl *CommandLine) NavigateHistory(direction int, inputField *tview.InputField) {
	if len(cl.History) == 0 {
		return
	}

	// First time browsing history - save current input
	if cl.HistoryIndex == -1 {
		cl.TempInput = inputField.GetText()
		if direction < 0 {
			// Going up - start from the end
			cl.HistoryIndex = len(cl.History)
		}
	}

	// Update index
	newIndex := cl.HistoryIndex + direction

	// Bounds checking
	if newIndex < 0 {
		newIndex = 0
	} else if newIndex > len(cl.History) {
		newIndex = len(cl.History)
	}

	cl.HistoryIndex = newIndex

	// Set input text
	if cl.HistoryIndex == len(cl.History) {
		// Reached the end - restore temp input
		inputField.SetText(cl.TempInput)
		logger.Debug("History navigation - restored temp input", nil)
	} else {
		// Show history command
		inputField.SetText(cl.History[cl.HistoryIndex])
		logger.Debug("History navigation", map[string]any{
			"index":   cl.HistoryIndex,
			"command": cl.History[cl.HistoryIndex],
		})
	}
}

// ExecuteCommand parses and executes a command
func (cl *CommandLine) ExecuteCommand(cmd string, ctx CommandContext) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	logger.Info("Command line execute", map[string]any{
		"command": command,
		"args":    args,
		"raw":     cmd,
	})

	switch command {
	case "db", "database":
		cl.handleDatabaseCommand(args, ctx)
	case "table", "tbl", "t":
		cl.handleTableCommand(args, ctx)
	case "help", "h", "?", "/?", "/help":
		cl.showCommandHelp(args)
	case "q", "quit", "exit":
		logger.Info("Quit command", nil)
		app.App.Stop()
	default:
		// Treat as SQL query
		logger.Info("Execute SQL", map[string]any{"sql": cmd})
		cl.executeSQL(cmd, ctx)
	}
}

func (cl *CommandLine) handleDatabaseCommand(args []string, ctx CommandContext) {
	if len(args) == 0 {
		ShowError("Usage: db <create|drop|use|list|backup|import> <name>")
		return
	}

	action := args[0]
	switch action {
	case "create", "c":
		if len(args) < 2 {
			ShowError("Usage: db create <name>")
			return
		}
		dbName := args[1]
		_, err := ctx.DB.ExecuteDMLStatement("CREATE DATABASE `" + dbName + "`")
		if err != nil {
			ShowError("Failed to create database: " + err.Error())
		} else {
			ShowSuccess("Database '" + dbName + "' created")
			RefreshTree()
		}
	case "drop", "d":
		if len(args) < 2 {
			ShowError("Usage: db drop <name>")
			return
		}
		dbName := args[1]
		_, err := ctx.DB.ExecuteDMLStatement("DROP DATABASE `" + dbName + "`")
		if err != nil {
			ShowError("Failed to drop database: " + err.Error())
		} else {
			ShowSuccess("Database '" + dbName + "' dropped")
			RefreshTree()
		}
	case "use", "u":
		if len(args) < 2 {
			ShowError("Usage: db use <name>")
			return
		}
		dbName := args[1]
		_, err := ctx.DB.ExecuteDMLStatement("USE `" + dbName + "`")
		if err != nil {
			ShowError("Failed to switch database: " + err.Error())
		} else {
			ShowSuccess("Switched to database '" + dbName + "'")
			RefreshTree()
		}
	case "list", "ls", "l":
		databases, err := ctx.DB.GetDatabases()
		if err != nil {
			ShowError("Failed to get databases: " + err.Error())
			return
		}
		ShowInfo("Databases:\n" + strings.Join(databases, "\n"))
	case "backup", "b":
		if len(args) < 2 {
			ShowError("Usage: db backup <filename>")
			return
		}
		filename := args[1]
		cl.handleBackup(filename, ctx)
	case "import", "i":
		if len(args) < 2 {
			ShowError("Usage: db import <filename>")
			return
		}
		filename := args[1]
		cl.handleImport(filename, ctx)
	default:
		ShowError("Unknown database command: " + action)
	}
}

func (cl *CommandLine) handleTableCommand(args []string, ctx CommandContext) {
	if len(args) == 0 {
		ShowError("Usage: table <create|drop|truncate|rename> <name>")
		return
	}

	action := args[0]
	switch action {
	case "create", "c":
		if len(args) < 2 {
			ShowError("Usage: table create <name>")
			return
		}
		tableName := args[1]
		sql := "CREATE TABLE `" + tableName + "` (id INT AUTO_INCREMENT PRIMARY KEY)"
		_, err := ctx.DB.ExecuteDMLStatement(sql)
		if err != nil {
			ShowError("Failed to create table: " + err.Error())
		} else {
			ShowSuccess("Table '" + tableName + "' created")
			RefreshTree()
		}
	case "drop", "d":
		if len(args) < 2 {
			ShowError("Usage: table drop <name>")
			return
		}
		tableName := args[1]
		_, err := ctx.DB.ExecuteDMLStatement("DROP TABLE `" + tableName + "`")
		if err != nil {
			ShowError("Failed to drop table: " + err.Error())
		} else {
			ShowSuccess("Table '" + tableName + "' dropped")
			RefreshTree()
		}
	case "truncate", "t":
		if len(args) < 2 {
			ShowError("Usage: table truncate <name>")
			return
		}
		tableName := args[1]
		_, err := ctx.DB.ExecuteDMLStatement("TRUNCATE TABLE `" + tableName + "`")
		if err != nil {
			ShowError("Failed to truncate table: " + err.Error())
		} else {
			ShowSuccess("Table '" + tableName + "' truncated")
		}
	case "rename", "r":
		if len(args) < 3 {
			ShowError("Usage: table rename <old> <new>")
			return
		}
		oldName := args[1]
		newName := args[2]
		sql := "ALTER TABLE `" + oldName + "` RENAME TO `" + newName + "`"
		_, err := ctx.DB.ExecuteDMLStatement(sql)
		if err != nil {
			ShowError("Failed to rename table: " + err.Error())
		} else {
			ShowSuccess("Table renamed to '" + newName + "'")
			RefreshTree()
		}
	default:
		ShowError("Unknown table command: " + action)
	}
}

func (cl *CommandLine) executeSQL(sql string, ctx CommandContext) {
	_, err := ctx.DB.ExecuteDMLStatement(sql)
	if err != nil {
		ShowError("SQL Error: " + err.Error())
	} else {
		ShowSuccess("SQL executed successfully")
		RefreshTree()
	}
}

func (cl *CommandLine) showCommandHelp(args []string) {
	// If specific command requested, show detailed help
	if len(args) > 0 {
		cl.showDetailedHelp(args[0])
		return
	}

	// Show general command list
	help := `[yellow]SQLCmder Command Line Help[white]

[yellow]Database Commands:[white]
  db create <name>     - Create new database
  db drop <name>       - Drop database
  db use <name>        - Switch to database
  db list              - List all databases
  db backup <file>     - Backup current database
  db import <file>     - Import SQL from file
  Aliases: database, db

[yellow]Table Commands:[white]
  table create <name>  - Create new table (interactive)
  table drop <name>    - Drop table
  table truncate <name>- Clear all table data
  table rename <old> <new> - Rename table
  Aliases: table, tbl, t

[yellow]SQL Commands:[white]
  SELECT ...           - Query data
  INSERT ...           - Insert data (type: help insert)
  UPDATE ...           - Update data (type: help update)
  DELETE ...           - Delete data (type: help delete)
  Any SQL statement    - Execute directly

[yellow]System Commands:[white]
  help [command]       - Show help (? /? /help)
  quit                 - Exit application (q, exit)

[yellow]Examples:[white]
  help insert          - Show INSERT syntax
  db create mydb       - Create database 'mydb'
  db backup mydb.sql   - Backup to ./backup/mydb.sql
  table drop users     - Drop table 'users'

[grey]Press Esc to close | Type command for details[white]`

	cl.showHelpModal(help)
}

func (cl *CommandLine) showDetailedHelp(topic string) {
	var help string

	topic = strings.ToLower(topic)
	switch topic {
	case "insert":
		help = `[yellow]INSERT Statement Syntax[white]

[yellow]Basic Syntax:[white]
  INSERT INTO table_name (column1, column2, ...)
  VALUES (value1, value2, ...);

[yellow]Insert Multiple Rows:[white]
  INSERT INTO table_name (column1, column2)
  VALUES 
    (value1a, value2a),
    (value1b, value2b),
    (value1c, value2c);

[yellow]Insert with SELECT:[white]
  INSERT INTO table_name (column1, column2)
  SELECT column1, column2 FROM other_table
  WHERE condition;

[yellow]Examples:[white]
  INSERT INTO users (name, email, age)
  VALUES ('John', 'john@example.com', 30);
  
  INSERT INTO users (name, email)
  VALUES ('Alice', 'alice@example.com'),
         ('Bob', 'bob@example.com');

[grey]Press Esc to close[white]`

	case "update":
		help = `[yellow]UPDATE Statement Syntax[white]

[yellow]Basic Syntax:[white]
  UPDATE table_name
  SET column1 = value1, column2 = value2, ...
  WHERE condition;

[yellow]Update with Calculation:[white]
  UPDATE products
  SET price = price * 1.1
  WHERE category = 'electronics';

[yellow]Update with JOIN:[white]
  UPDATE orders o
  JOIN customers c ON o.customer_id = c.id
  SET o.status = 'vip'
  WHERE c.level = 'premium';

[yellow]Examples:[white]
  UPDATE users SET age = 31 WHERE name = 'John';
  
  UPDATE products SET stock = stock - 1
  WHERE id = 100;

[grey]Press Esc to close[white]`

	case "delete":
		help = `[yellow]DELETE Statement Syntax[white]

[yellow]Basic Syntax:[white]
  DELETE FROM table_name
  WHERE condition;

[yellow]Delete All Rows:[white]
  DELETE FROM table_name;  -- Use with caution!
  TRUNCATE TABLE table_name;  -- Faster alternative

[yellow]Delete with JOIN:[white]
  DELETE o FROM orders o
  JOIN customers c ON o.customer_id = c.id
  WHERE c.status = 'inactive';

[yellow]Examples:[white]
  DELETE FROM users WHERE age < 18;
  
  DELETE FROM logs
  WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

[grey]Press Esc to close[white]`

	case "select":
		help = `[yellow]SELECT Statement Syntax[white]

[yellow]Basic Syntax:[white]
  SELECT column1, column2, ...
  FROM table_name
  WHERE condition
  ORDER BY column1 ASC/DESC
  LIMIT count OFFSET offset;

[yellow]Aggregation:[white]
  SELECT COUNT(*), AVG(price), MAX(price)
  FROM products
  GROUP BY category
  HAVING COUNT(*) > 5;

[yellow]JOIN:[white]
  SELECT o.id, c.name, o.total
  FROM orders o
  JOIN customers c ON o.customer_id = c.id
  WHERE o.status = 'completed';

[yellow]Examples:[white]
  SELECT * FROM users WHERE age > 25;
  
  SELECT name, email FROM users
  ORDER BY created_at DESC LIMIT 10;

[grey]Press Esc to close[white]`

	case "db", "database":
		help = `[yellow]Database Commands[white]

[yellow]Commands:[white]
  db create <name>     - Create new database
  db drop <name>       - Drop database (permanent!)
  db use <name>        - Switch to database
  db list              - List all databases
  db backup <file>     - Backup current database
  db import <file>     - Import SQL from file

[yellow]Backup/Import:[white]
  - MySQL: Uses mysqldump/mysql (requires tools)
  - PostgreSQL: Uses pg_dump/psql (requires tools)
  - SQLite: Direct file copy
  - MSSQL: Uses sqlcmd (requires tools)
  - Files saved to ./backup/ directory
  - Can use relative or absolute paths

[yellow]Examples:[white]
  db create myapp_dev
  db use myapp_dev
  db backup mydb_20231022.sql
  db import mydb_backup.sql
  db list

[grey]Press Esc to close[white]`

	case "table", "tbl":
		help = `[yellow]Table Commands[white]

[yellow]Commands:[white]
  table create <name>    - Create table (interactive)
  table drop <name>      - Drop table (permanent!)
  table truncate <name>  - Delete all data (fast)
  table rename <old> <new> - Rename table

[yellow]SQL Equivalent:[white]
  CREATE TABLE table_name (...);
  DROP TABLE table_name;
  TRUNCATE TABLE table_name;
  RENAME TABLE old_name TO new_name;

[yellow]Examples:[white]
  table create users
  table truncate logs
  table rename user users

[grey]Press Esc to close[white]`

	case "backup":
		help = `[yellow]Database Backup[white]

[yellow]Command:[white]
  db backup <filename>

[yellow]Description:[white]
  Backs up the current database to a file in the ./backup/ directory.
  The backup method depends on your database type:

[yellow]Database Types:[white]
  MySQL        - Uses mysqldump (must be installed)
  PostgreSQL   - Uses pg_dump (must be installed)
  SQLite       - Direct file copy
  MSSQL        - Uses sqlcmd (must be installed)

[yellow]Notes:[white]
  - Backups are saved to ./backup/ directory
  - Directory is created automatically if needed
  - For MySQL/PostgreSQL, ensure CLI tools are in PATH
  - SQLite backups are file copies (reliable)
  - MSSQL requires SQL Server access permissions

[yellow]Examples:[white]
  db backup mydb_20231022.sql
  db backup production_backup.sql
  db backup test.db.bak

[grey]Press Esc to close[white]`

	case "import":
		help = `[yellow]Database Import[white]

[yellow]Command:[white]
  db import <filename>

[yellow]Description:[white]
  Imports SQL statements from a file into the current database.
  The import method depends on your database type:

[yellow]Database Types:[white]
  MySQL        - Uses mysql client (must be installed)
  PostgreSQL   - Uses psql (must be installed)
  SQLite       - Executes SQL statements directly
  MSSQL        - Uses sqlcmd (must be installed)

[yellow]File Lookup:[white]
  1. Checks current directory first
  2. Then checks ./backup/ directory
  3. Can use relative or absolute paths

[yellow]Notes:[white]
  - Ensure database exists before importing
  - Large files may take time to process
  - SQLite imports are executed statement by statement
  - For MySQL/PostgreSQL, ensure CLI tools are in PATH
  - Import refreshes the database tree

[yellow]Examples:[white]
  db import mydb_backup.sql
  db import ./backup/production.sql
  db import C:\backups\data.sql

[grey]Press Esc to close[white]`

	default:
		help = fmt.Sprintf(`[yellow]Help Topic: %s[white]

[red]Unknown topic.[white] Available topics:
  insert, update, delete, select
  db, table, backup, import

Type [yellow]help[white] to see all commands.
Type [yellow]help <topic>[white] for specific help.

[grey]Press Esc to close[white]`, topic)
	}

	cl.showHelpModal(help)
}

func (cl *CommandLine) showHelpModal(content string) {
	// Create text view for help content
	textView := tview.NewTextView()
	textView.SetText(content)
	textView.SetDynamicColors(true)
	textView.SetWordWrap(true)
	textView.SetScrollable(true)
	textView.SetBorder(true)
	textView.SetTitle(" Help - Press Esc to close ")
	textView.SetTitleAlign(tview.AlignCenter)
	textView.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)
	textView.SetBorderColor(app.Styles.InverseTextColor)

	// Handle Esc to close
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			CloseModal()
			// Refocus to command line input
			if globalCommandLine != nil {
				app.App.SetFocus(globalCommandLine.InputField)
			}
			return nil
		}
		return event
	})

	// Show modal
	ShowModal(textView, 80, 25)
}

// handleBackup performs database backup
func (cl *CommandLine) handleBackup(filename string, ctx CommandContext) {
	if ctx.ConnectionModel == nil {
		ShowError("Connection information not available")
		return
	}

	conn := ctx.ConnectionModel
	provider := strings.ToLower(conn.Provider)
	dbName := ctx.CurrentDatabase
	if dbName == "" {
		dbName = conn.DBName
	}

	logger.Info("Database backup", map[string]any{
		"provider": provider,
		"database": dbName,
		"file":     filename,
	})

	switch provider {
	case "mysql":
		cl.backupMySQL(filename, dbName, conn)
	case "postgres", "postgresql":
		cl.backupPostgreSQL(filename, dbName, conn)
	case "sqlite":
		cl.backupSQLite(filename, dbName, conn)
	case "mssql", "sqlserver":
		cl.backupMSSQL(filename, dbName, conn)
	default:
		ShowError("Backup not supported for provider: " + provider)
	}
}

// handleImport imports data from SQL file
func (cl *CommandLine) handleImport(filename string, ctx CommandContext) {
	if ctx.ConnectionModel == nil {
		ShowError("Connection information not available")
		return
	}

	conn := ctx.ConnectionModel
	provider := strings.ToLower(conn.Provider)
	dbName := ctx.CurrentDatabase
	if dbName == "" {
		dbName = conn.DBName
	}

	logger.Info("Database import", map[string]any{
		"provider": provider,
		"database": dbName,
		"file":     filename,
	})

	switch provider {
	case "mysql":
		cl.importMySQL(filename, dbName, conn)
	case "postgres", "postgresql":
		cl.importPostgreSQL(filename, dbName, conn)
	case "sqlite":
		cl.importSQLite(filename, dbName, conn)
	case "mssql", "sqlserver":
		cl.importMSSQL(filename, dbName, conn)
	default:
		ShowError("Import not supported for provider: " + provider)
	}
}
