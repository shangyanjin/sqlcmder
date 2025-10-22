package drivers

import (
	"sqlcmder/models"
)

type Driver interface {
	Connect(urlstr string) error
	TestConnection(urlstr string) error
	GetDatabases() ([]string, error)
	GetTables(database string) (map[string][]string, error)
	GetTableColumns(database, table string) ([][]string, error)
	GetConstraints(database, table string) ([][]string, error)
	GetForeignKeys(database, table string) ([][]string, error)
	GetIndexes(database, table string) ([][]string, error)
	GetRecords(database, table, where, sort string, offset, limit int) ([][]string, int, string, error)
	UpdateRecord(database, table, column, value, primaryKeyColumnName, primaryKeyValue string) error
	DeleteRecord(database, table string, primaryKeyColumnName, primaryKeyValue string) error
	ExecuteDMLStatement(query string) (string, error)
	ExecuteQuery(query string) ([][]string, int, error)
	ExecutePendingChanges(changes []models.DBDMLChange) error
	GetProvider() string
	GetPrimaryKeyColumnNames(database, table string) ([]string, error)

	FormatArg(arg any, colype models.CellValueType) any
	FormatArgForQueryString(arg any) string
	FormatReference(reference string) string
	FormatPlaceholder(index int) string

	// This converts a DML change to a query string with arg values
	DMLChangeToQueryString(change models.DBDMLChange) (string, error)

	// NOTE: This is used to get the primary key from the database table until I
	// find a better way to do it. See *ResultsTable.GetPrimaryKeyValue()
	SetProvider(provider string)
}
