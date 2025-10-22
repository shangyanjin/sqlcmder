package commands

// ExecuteTableCommand handles table-related commands
func ExecuteTableCommand(args []string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	if len(args) == 0 {
		onError("Usage: table <create|drop|truncate|rename> <name>")
		return
	}

	action := args[0]
	switch action {
	case "create", "c":
		createTable(args, ctx, onSuccess, onError, onRefresh)
	case "drop", "d":
		dropTable(args, ctx, onSuccess, onError, onRefresh)
	case "truncate", "t":
		truncateTable(args, ctx, onSuccess, onError)
	case "rename", "r":
		renameTable(args, ctx, onSuccess, onError, onRefresh)
	default:
		onError("Unknown table command: " + action)
	}
}

func createTable(args []string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	if len(args) < 2 {
		onError("Usage: table create <name>")
		return
	}

	tableName := args[1]
	sql := "CREATE TABLE `" + tableName + "` (id INT AUTO_INCREMENT PRIMARY KEY)"
	_, err := ctx.DB.ExecuteDMLStatement(sql)
	if err != nil {
		onError("Failed to create table: " + err.Error())
	} else {
		onSuccess("Table '" + tableName + "' created")
		onRefresh()
	}
}

func dropTable(args []string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	if len(args) < 2 {
		onError("Usage: table drop <name>")
		return
	}

	tableName := args[1]
	_, err := ctx.DB.ExecuteDMLStatement("DROP TABLE `" + tableName + "`")
	if err != nil {
		onError("Failed to drop table: " + err.Error())
	} else {
		onSuccess("Table '" + tableName + "' dropped")
		onRefresh()
	}
}

func truncateTable(args []string, ctx Context, onSuccess func(string), onError func(string)) {
	if len(args) < 2 {
		onError("Usage: table truncate <name>")
		return
	}

	tableName := args[1]
	_, err := ctx.DB.ExecuteDMLStatement("TRUNCATE TABLE `" + tableName + "`")
	if err != nil {
		onError("Failed to truncate table: " + err.Error())
	} else {
		onSuccess("Table '" + tableName + "' truncated")
	}
}

func renameTable(args []string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	if len(args) < 3 {
		onError("Usage: table rename <old> <new>")
		return
	}

	oldName := args[1]
	newName := args[2]
	sql := "ALTER TABLE `" + oldName + "` RENAME TO `" + newName + "`"
	_, err := ctx.DB.ExecuteDMLStatement(sql)
	if err != nil {
		onError("Failed to rename table: " + err.Error())
	} else {
		onSuccess("Table renamed to '" + newName + "'")
		onRefresh()
	}
}

