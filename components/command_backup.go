package components

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"sqlcmder/helpers"
	"sqlcmder/helpers/logger"
	"sqlcmder/models"
)

// MySQL backup using mysqldump
func (cl *CommandLine) backupMySQL(filename string, dbName string, conn *models.Connection) {
	// Build mysqldump command
	args := []string{
		"-h", conn.Hostname,
		"-P", conn.Port,
		"-u", conn.Username,
	}

	if conn.Password != "" {
		args = append(args, fmt.Sprintf("-p%s", conn.Password))
	}

	args = append(args, dbName)

	// Create backup directory if needed
	backupDir := filepath.Join(".", "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		ShowError("Failed to create backup directory: " + err.Error())
		return
	}

	outputFile := filepath.Join(backupDir, filename)

	logger.Info("Executing mysqldump", map[string]any{
		"args":   args,
		"output": outputFile,
	})

	// Execute mysqldump
	cmd := exec.Command("mysqldump", args...)
	outFile, err := os.Create(outputFile)
	if err != nil {
		ShowError("Failed to create backup file: " + err.Error())
		return
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		ShowError("Backup failed: " + err.Error() + " - " + stderr.String())
		return
	}

	ShowSuccess(fmt.Sprintf("Database backed up to: %s", outputFile))
}

// PostgreSQL backup using pg_dump
func (cl *CommandLine) backupPostgreSQL(filename string, dbName string, conn *models.Connection) {
	// Build pg_dump command
	args := []string{
		"-h", conn.Hostname,
		"-p", conn.Port,
		"-U", conn.Username,
		"-F", "p", // Plain text format
		"-f", filename,
		dbName,
	}

	// Create backup directory if needed
	backupDir := filepath.Join(".", "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		ShowError("Failed to create backup directory: " + err.Error())
		return
	}

	outputFile := filepath.Join(backupDir, filename)
	args[5] = outputFile // Update -f argument with full path

	logger.Info("Executing pg_dump", map[string]any{
		"args":   args,
		"output": outputFile,
	})

	cmd := exec.Command("pg_dump", args...)

	// Set PGPASSWORD environment variable
	if conn.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", conn.Password))
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		ShowError("Backup failed: " + err.Error() + " - " + stderr.String())
		return
	}

	ShowSuccess(fmt.Sprintf("Database backed up to: %s", outputFile))
}

// SQLite backup by copying file
func (cl *CommandLine) backupSQLite(filename string, dbName string, conn *models.Connection) {
	// For SQLite, dbName is the file path
	sourceFile := conn.DBName
	if sourceFile == "" {
		sourceFile = dbName
	}

	// Create backup directory if needed
	backupDir := filepath.Join(".", "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		ShowError("Failed to create backup directory: " + err.Error())
		return
	}

	outputFile := filepath.Join(backupDir, filename)

	logger.Info("Copying SQLite database", map[string]any{
		"source": sourceFile,
		"dest":   outputFile,
	})

	// Copy file
	source, err := os.Open(sourceFile)
	if err != nil {
		ShowError("Failed to open source file: " + err.Error())
		return
	}
	defer source.Close()

	dest, err := os.Create(outputFile)
	if err != nil {
		ShowError("Failed to create backup file: " + err.Error())
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, source); err != nil {
		ShowError("Failed to copy database: " + err.Error())
		return
	}

	ShowSuccess(fmt.Sprintf("Database backed up to: %s", outputFile))
}

// MSSQL backup using T-SQL BACKUP DATABASE
func (cl *CommandLine) backupMSSQL(filename string, dbName string, conn *models.Connection) {
	// Create backup directory if needed
	backupDir := filepath.Join(".", "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		ShowError("Failed to create backup directory: " + err.Error())
		return
	}

	// For MSSQL, backup path must be accessible by SQL Server
	// Use absolute path
	absPath, err := filepath.Abs(filepath.Join(backupDir, filename))
	if err != nil {
		ShowError("Failed to resolve backup path: " + err.Error())
		return
	}

	// Build sqlcmd command
	backupSQL := fmt.Sprintf("BACKUP DATABASE [%s] TO DISK = '%s' WITH FORMAT, COMPRESSION", dbName, absPath)

	logger.Info("Executing MSSQL backup", map[string]any{
		"sql": backupSQL,
	})

	// Use sqlcmd command line tool
	args := []string{
		"-S", fmt.Sprintf("%s,%s", conn.Hostname, conn.Port),
		"-U", conn.Username,
		"-Q", backupSQL,
	}

	if conn.Password != "" {
		args = append(args, "-P", conn.Password)
	}

	cmd := exec.Command("sqlcmd", args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		ShowError("Backup failed: " + err.Error() + " - " + stderr.String())
		return
	}

	ShowSuccess(fmt.Sprintf("Database backed up to: %s", absPath))
}

// MySQL import using mysql client
func (cl *CommandLine) importMySQL(filename string, dbName string, conn *models.Connection) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Try in backup directory
		backupFile := filepath.Join(".", "backup", filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			ShowError("Import file not found: " + filename)
			return
		}
		filename = backupFile
	}

	logger.Info("Importing MySQL database", map[string]any{
		"file":     filename,
		"database": dbName,
	})

	// Build mysql command
	args := []string{
		"-h", conn.Hostname,
		"-P", conn.Port,
		"-u", conn.Username,
	}

	if conn.Password != "" {
		args = append(args, fmt.Sprintf("-p%s", conn.Password))
	}

	args = append(args, dbName)

	cmd := exec.Command("mysql", args...)

	// Read SQL file
	inFile, err := os.Open(filename)
	if err != nil {
		ShowError("Failed to open import file: " + err.Error())
		return
	}
	defer inFile.Close()

	cmd.Stdin = inFile
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		ShowError("Import failed: " + err.Error() + " - " + stderr.String())
		return
	}

	ShowSuccess("Database imported successfully")
	RefreshTree()
}

// PostgreSQL import using psql
func (cl *CommandLine) importPostgreSQL(filename string, dbName string, conn *models.Connection) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Try in backup directory
		backupFile := filepath.Join(".", "backup", filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			ShowError("Import file not found: " + filename)
			return
		}
		filename = backupFile
	}

	logger.Info("Importing PostgreSQL database", map[string]any{
		"file":     filename,
		"database": dbName,
	})

	// Build psql command
	args := []string{
		"-h", conn.Hostname,
		"-p", conn.Port,
		"-U", conn.Username,
		"-d", dbName,
		"-f", filename,
	}

	cmd := exec.Command("psql", args...)

	// Set PGPASSWORD environment variable
	if conn.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", conn.Password))
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		ShowError("Import failed: " + err.Error() + " - " + stderr.String())
		return
	}

	ShowSuccess("Database imported successfully")
	RefreshTree()
}

// SQLite import by reading SQL file
func (cl *CommandLine) importSQLite(filename string, dbName string, conn *models.Connection) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Try in backup directory
		backupFile := filepath.Join(".", "backup", filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			ShowError("Import file not found: " + filename)
			return
		}
		filename = backupFile
	}

	logger.Info("Importing SQLite database", map[string]any{
		"file":     filename,
		"database": dbName,
	})

	// Read SQL file
	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		ShowError("Failed to read import file: " + err.Error())
		return
	}

	sqlContent := string(sqlBytes)

	// Split into individual statements (basic split by semicolon)
	statements := strings.Split(sqlContent, ";")

	ctx := context.Background()
	successCount := 0

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// Execute statement using helpers.RunCommand
		err := helpers.RunCommand(ctx, fmt.Sprintf("sqlite3 %s \"%s;\"", dbName, stmt), func(output string) {
			logger.Debug("SQL executed", map[string]any{"output": output})
		})

		if err != nil {
			logger.Error("Statement failed", map[string]any{"error": err.Error(), "sql": stmt})
		} else {
			successCount++
		}
	}

	ShowSuccess(fmt.Sprintf("Import completed: %d statements executed", successCount))
	RefreshTree()
}

// MSSQL import by reading SQL file
func (cl *CommandLine) importMSSQL(filename string, dbName string, conn *models.Connection) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Try in backup directory
		backupFile := filepath.Join(".", "backup", filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			ShowError("Import file not found: " + filename)
			return
		}
		filename = backupFile
	}

	logger.Info("Importing MSSQL database", map[string]any{
		"file":     filename,
		"database": dbName,
	})

	// Build sqlcmd command
	args := []string{
		"-S", fmt.Sprintf("%s,%s", conn.Hostname, conn.Port),
		"-U", conn.Username,
		"-d", dbName,
		"-i", filename,
	}

	if conn.Password != "" {
		args = append(args, "-P", conn.Password)
	}

	cmd := exec.Command("sqlcmd", args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		ShowError("Import failed: " + err.Error() + " - " + stderr.String())
		return
	}

	ShowSuccess("Database imported successfully")
	RefreshTree()
}
