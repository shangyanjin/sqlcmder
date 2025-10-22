package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"sqlcmder/helpers"
	"sqlcmder/logger"
	"sqlcmder/models"
)

// BackupDatabase performs database backup
func BackupDatabase(filename string, ctx Context, onSuccess func(string), onError func(string)) {
	if ctx.ConnectionModel == nil {
		onError("Connection information not available")
		return
	}

	conn := ctx.ConnectionModel
	provider := strings.ToLower(conn.Driver)
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
		backupMySQL(filename, dbName, conn, onSuccess, onError)
	case "postgres", "postgresql":
		backupPostgreSQL(filename, dbName, conn, onSuccess, onError)
	case "sqlite":
		backupSQLite(filename, dbName, conn, onSuccess, onError)
	case "mssql", "sqlserver":
		backupMSSQL(filename, dbName, conn, onSuccess, onError)
	default:
		onError("Backup not supported for provider: " + provider)
	}
}

// ImportDatabase imports data from SQL file
func ImportDatabase(filename string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	if ctx.ConnectionModel == nil {
		onError("Connection information not available")
		return
	}

	conn := ctx.ConnectionModel
	provider := strings.ToLower(conn.Driver)
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
		importMySQL(filename, dbName, conn, onSuccess, onError, onRefresh)
	case "postgres", "postgresql":
		importPostgreSQL(filename, dbName, conn, onSuccess, onError, onRefresh)
	case "sqlite":
		importSQLite(filename, dbName, conn, onSuccess, onError, onRefresh)
	case "mssql", "sqlserver":
		importMSSQL(filename, dbName, conn, onSuccess, onError, onRefresh)
	default:
		onError("Import not supported for provider: " + provider)
	}
}

// MySQL backup using mysqldump
func backupMySQL(filename string, dbName string, conn *models.Connection, onSuccess func(string), onError func(string)) {
	args := []string{
		"-h", conn.Hostname,
		"-P", conn.Port,
		"-u", conn.Username,
	}

	if conn.Password != "" {
		args = append(args, fmt.Sprintf("-p%s", conn.Password))
	}

	args = append(args, dbName)

	backupDir := filepath.Join(".", "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		onError("Failed to create backup directory: " + err.Error())
		return
	}

	outputFile := filepath.Join(backupDir, filename)

	logger.Info("Executing mysqldump", map[string]any{
		"args":   args,
		"output": outputFile,
	})

	cmd := exec.Command("mysqldump", args...)
	outFile, err := os.Create(outputFile)
	if err != nil {
		onError("Failed to create backup file: " + err.Error())
		return
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		onError("Backup failed: " + err.Error() + " - " + stderr.String())
		return
	}

	onSuccess(fmt.Sprintf("Database backed up to: %s", outputFile))
}

// PostgreSQL backup using pg_dump
func backupPostgreSQL(filename string, dbName string, conn *models.Connection, onSuccess func(string), onError func(string)) {
	backupDir := filepath.Join(".", "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		onError("Failed to create backup directory: " + err.Error())
		return
	}

	outputFile := filepath.Join(backupDir, filename)

	args := []string{
		"-h", conn.Hostname,
		"-p", conn.Port,
		"-U", conn.Username,
		"-F", "p",
		"-f", outputFile,
		dbName,
	}

	logger.Info("Executing pg_dump", map[string]any{
		"args":   args,
		"output": outputFile,
	})

	cmd := exec.Command("pg_dump", args...)

	if conn.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", conn.Password))
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		onError("Backup failed: " + err.Error() + " - " + stderr.String())
		return
	}

	onSuccess(fmt.Sprintf("Database backed up to: %s", outputFile))
}

// SQLite backup by copying file
func backupSQLite(filename string, dbName string, conn *models.Connection, onSuccess func(string), onError func(string)) {
	sourceFile := conn.DBName
	if sourceFile == "" {
		sourceFile = dbName
	}

	backupDir := filepath.Join(".", "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		onError("Failed to create backup directory: " + err.Error())
		return
	}

	outputFile := filepath.Join(backupDir, filename)

	logger.Info("Copying SQLite database", map[string]any{
		"source": sourceFile,
		"dest":   outputFile,
	})

	source, err := os.Open(sourceFile)
	if err != nil {
		onError("Failed to open source file: " + err.Error())
		return
	}
	defer source.Close()

	dest, err := os.Create(outputFile)
	if err != nil {
		onError("Failed to create backup file: " + err.Error())
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, source); err != nil {
		onError("Failed to copy database: " + err.Error())
		return
	}

	onSuccess(fmt.Sprintf("Database backed up to: %s", outputFile))
}

// MSSQL backup using sqlcmd
func backupMSSQL(filename string, dbName string, conn *models.Connection, onSuccess func(string), onError func(string)) {
	backupDir := filepath.Join(".", "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		onError("Failed to create backup directory: " + err.Error())
		return
	}

	absPath, err := filepath.Abs(filepath.Join(backupDir, filename))
	if err != nil {
		onError("Failed to resolve backup path: " + err.Error())
		return
	}

	backupSQL := fmt.Sprintf("BACKUP DATABASE [%s] TO DISK = '%s' WITH FORMAT, COMPRESSION", dbName, absPath)

	logger.Info("Executing MSSQL backup", map[string]any{
		"sql": backupSQL,
	})

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
		onError("Backup failed: " + err.Error() + " - " + stderr.String())
		return
	}

	onSuccess(fmt.Sprintf("Database backed up to: %s", absPath))
}

// MySQL import using mysql client
func importMySQL(filename string, dbName string, conn *models.Connection, onSuccess func(string), onError func(string), onRefresh func()) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		backupFile := filepath.Join(".", "backup", filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			onError("Import file not found: " + filename)
			return
		}
		filename = backupFile
	}

	logger.Info("Importing MySQL database", map[string]any{
		"file":     filename,
		"database": dbName,
	})

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

	inFile, err := os.Open(filename)
	if err != nil {
		onError("Failed to open import file: " + err.Error())
		return
	}
	defer inFile.Close()

	cmd.Stdin = inFile
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		onError("Import failed: " + err.Error() + " - " + stderr.String())
		return
	}

	onSuccess("Database imported successfully")
	onRefresh()
}

// PostgreSQL import using psql
func importPostgreSQL(filename string, dbName string, conn *models.Connection, onSuccess func(string), onError func(string), onRefresh func()) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		backupFile := filepath.Join(".", "backup", filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			onError("Import file not found: " + filename)
			return
		}
		filename = backupFile
	}

	logger.Info("Importing PostgreSQL database", map[string]any{
		"file":     filename,
		"database": dbName,
	})

	args := []string{
		"-h", conn.Hostname,
		"-p", conn.Port,
		"-U", conn.Username,
		"-d", dbName,
		"-f", filename,
	}

	cmd := exec.Command("psql", args...)

	if conn.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", conn.Password))
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		onError("Import failed: " + err.Error() + " - " + stderr.String())
		return
	}

	onSuccess("Database imported successfully")
	onRefresh()
}

// SQLite import by reading SQL file
func importSQLite(filename string, dbName string, conn *models.Connection, onSuccess func(string), onError func(string), onRefresh func()) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		backupFile := filepath.Join(".", "backup", filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			onError("Import file not found: " + filename)
			return
		}
		filename = backupFile
	}

	logger.Info("Importing SQLite database", map[string]any{
		"file":     filename,
		"database": dbName,
	})

	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		onError("Failed to read import file: " + err.Error())
		return
	}

	sqlContent := string(sqlBytes)
	statements := strings.Split(sqlContent, ";")

	ctx := context.Background()
	successCount := 0

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		err := helpers.RunCommand(ctx, fmt.Sprintf("sqlite3 %s \"%s;\"", dbName, stmt), func(output string) {
			logger.Debug("SQL executed", map[string]any{"output": output})
		})

		if err != nil {
			logger.Error("Statement failed", map[string]any{"error": err.Error(), "sql": stmt})
		} else {
			successCount++
		}
	}

	onSuccess(fmt.Sprintf("Import completed: %d statements executed", successCount))
	onRefresh()
}

// MSSQL import using sqlcmd
func importMSSQL(filename string, dbName string, conn *models.Connection, onSuccess func(string), onError func(string), onRefresh func()) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		backupFile := filepath.Join(".", "backup", filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			onError("Import file not found: " + filename)
			return
		}
		filename = backupFile
	}

	logger.Info("Importing MSSQL database", map[string]any{
		"file":     filename,
		"database": dbName,
	})

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
		onError("Import failed: " + err.Error() + " - " + stderr.String())
		return
	}

	onSuccess("Database imported successfully")
	onRefresh()
}

