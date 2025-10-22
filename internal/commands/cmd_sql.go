package commands

import (
	"sqlcmder/internal/helpers/logger"
)

// ExecuteSQL executes arbitrary SQL statement
func ExecuteSQL(sql string, ctx Context, onSuccess func(string), onError func(string), onRefresh func()) {
	logger.Info("Execute SQL", map[string]any{"sql": sql})

	_, err := ctx.DB.ExecuteDMLStatement(sql)
	if err != nil {
		onError("SQL Error: " + err.Error())
	} else {
		onSuccess("SQL executed successfully")
		onRefresh()
	}
}

