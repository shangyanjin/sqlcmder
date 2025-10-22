package commands

import (
	"strings"
)

// ExecuteDatabaseCommand handles database-related commands
func ExecuteDatabaseCommand(args []string, ctx Context, onSuccess func(string), onError func(string), onInfo func(string), onRefresh func()) {
	if len(args) == 0 {
		onError("Usage: db <create|drop|use|list|backup|import> <name>")
		return
	}

	action := args[0]
	switch action {
	case "create", "c":
		createDatabase(args, ctx, onSuccess, onError, onRefresh)
	case "drop", "d":
		dropDatabase(args, ctx, onSuccess, onError, onRefresh)
	case "use", "u":
		useDatabase(args, ctx, onSuccess, onError, onRefresh)
	case "list", "ls", "l":
		listDatabases(ctx, onError, onInfo)
	case "backup", "b":
		if len(args) < 2 {
			onError("Usage: db backup <filename>")
			return
		}
		BackupDatabase(args[1], ctx, onSuccess, onError)
	case "import", "i":
		if len(args) < 2 {
			onError("Usage: db import <filename>")
			return
		}
		ImportDatabase(args[1], ctx, onSuccess, onError, onRefresh)
	default:
		onError("Unknown database command: " + action)
	}
}

func createDatabase(args []string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	if len(args) < 2 {
		onError("Usage: db create <name>")
		return
	}

	dbName := args[1]
	_, err := ctx.DB.ExecuteDMLStatement("CREATE DATABASE `" + dbName + "`")
	if err != nil {
		onError("Failed to create database: " + err.Error())
	} else {
		onSuccess("Database '" + dbName + "' created")
		onRefresh()
	}
}

func dropDatabase(args []string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	if len(args) < 2 {
		onError("Usage: db drop <name>")
		return
	}

	dbName := args[1]
	_, err := ctx.DB.ExecuteDMLStatement("DROP DATABASE `" + dbName + "`")
	if err != nil {
		onError("Failed to drop database: " + err.Error())
	} else {
		onSuccess("Database '" + dbName + "' dropped")
		onRefresh()
	}
}

func useDatabase(args []string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	if len(args) < 2 {
		onError("Usage: db use <name>")
		return
	}

	dbName := args[1]
	_, err := ctx.DB.ExecuteDMLStatement("USE `" + dbName + "`")
	if err != nil {
		onError("Failed to switch database: " + err.Error())
	} else {
		onSuccess("Switched to database '" + dbName + "'")
		onRefresh()
	}
}

func listDatabases(ctx Context, onError func(string), onInfo func(string)) {
	databases, err := ctx.DB.GetDatabases()
	if err != nil {
		onError("Failed to get databases: " + err.Error())
		return
	}
	onInfo("Databases:\n" + strings.Join(databases, "\n"))
}

