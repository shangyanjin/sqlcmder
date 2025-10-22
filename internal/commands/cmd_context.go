package commands

import (
	"sqlcmder/internal/drivers"
	"sqlcmder/models"
)

// Context holds the current context for command execution
type Context struct {
	DB              drivers.Driver
	CurrentDatabase string
	CurrentTable    string
	Connection      string
	ConnectionModel *models.Connection // Full connection details for backup/import
}

