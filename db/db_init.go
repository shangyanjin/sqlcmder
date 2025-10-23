package db

import (
	"fmt"
	"strings"

	"sqlcmder/drivers"
	"sqlcmder/helpers"
	"sqlcmder/models"
)

// InitFromArg initializes a database connection from a connection string argument
// Returns the connection, driver, and any error
func InitFromArg(connectionString string) (*models.Connection, drivers.Driver, error) {
	parsed, err := helpers.ParseConnectionString(connectionString)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse connection string: %s", err)
	}
	DBName := strings.Split(parsed.Normalize(",", "NULL", 0), ",")[3]

	if DBName == "NULL" {
		DBName = ""
	}

	connection := models.Connection{
		Name:   "",
		Driver: parsed.Driver,
		DBName: DBName,
		DSN:    connectionString,
	}

	var newDBDriver drivers.Driver
	switch connection.Driver {
	case drivers.DriverMySQL:
		newDBDriver = &drivers.MySQL{}
	case drivers.DriverPostgres:
		newDBDriver = &drivers.Postgres{}
	case drivers.DriverSqlite:
		newDBDriver = &drivers.SQLite{}
	case drivers.DriverMSSQL:
		newDBDriver = &drivers.MSSQL{}
	default:
		return nil, nil, fmt.Errorf("could not handle database driver %s", connection.Driver)
	}

	err = newDBDriver.Connect(connection.GetDSN())
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to database %s: %s", connectionString, err)
	}

	return &connection, newDBDriver, nil
}
