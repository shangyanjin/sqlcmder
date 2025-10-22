package ui

import (
	"fmt"
	"strings"

	"sqlcmder/internal/drivers"
	"sqlcmder/internal/helpers"
	"sqlcmder/models"
)

func InitFromArg(connectionString string) error {
	parsed, err := helpers.ParseConnectionString(connectionString)
	if err != nil {
		return fmt.Errorf("could not parse connection string: %s", err)
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
		return fmt.Errorf("could not handle database driver %s", connection.Driver)
	}

	err = newDBDriver.Connect(connection.DSN)
	if err != nil {
		return fmt.Errorf("could not connect to database %s: %s", connectionString, err)
	}
	mainPages.AddAndSwitchToPage(connection.DSN, NewHomePage(connection, newDBDriver).Flex, true)

	return nil
}
