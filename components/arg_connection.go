package components

import (
	"fmt"
	"strings"

	"sqlcmder/drivers"
	"sqlcmder/helpers"
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
		Name:     "",
		Provider: parsed.Driver,
		DBName:   DBName,
		URL:      connectionString,
	}

	var newDBDriver drivers.Driver
	switch connection.Provider {
	case drivers.DriverMySQL:
		newDBDriver = &drivers.MySQL{}
	case drivers.DriverPostgres:
		newDBDriver = &drivers.Postgres{}
	case drivers.DriverSqlite:
		newDBDriver = &drivers.SQLite{}
	case drivers.DriverMSSQL:
		newDBDriver = &drivers.MSSQL{}
	default:
		return fmt.Errorf("could not handle database driver %s", connection.Provider)
	}

	err = newDBDriver.Connect(connection.URL)
	if err != nil {
		return fmt.Errorf("could not connect to database %s: %s", connectionString, err)
	}
	mainPages.AddAndSwitchToPage(connection.URL, NewHomePage(connection, newDBDriver).Flex, true)

	return nil
}
