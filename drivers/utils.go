package drivers

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"sqlcmder/helpers/logger"
	"sqlcmder/models"
)

func queriesInTransaction(db *sql.DB, queries []models.Query) (err error) {
	trx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		rErr := trx.Rollback()
		// sql.ErrTxDone is returned when trx.Commit was already called
		if !errors.Is(rErr, sql.ErrTxDone) {
			err = errors.Join(err, rErr)
		}
	}()

	for _, query := range queries {
		if _, err := trx.Exec(query.Query, query.Args...); err != nil {
			return err
		}
	}
	if err := trx.Commit(); err != nil {
		return err
	}
	return nil
}

func buildInsertQueryString(formattedTableName string, columns []string, values []any, driver Driver) string {
	sanitizedValues := make([]string, len(values))

	for i, v := range values {
		sanitizedValues[i] = fmt.Sprintf("%v", driver.FormatArgForQueryString(v))
	}

	queryStr := "INSERT INTO " + formattedTableName
	queryStr += fmt.Sprintf(" (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(sanitizedValues, ", "))

	return queryStr
}

func buildInsertQuery(formattedTableName string, values []models.CellValue, driver Driver) models.Query {
	cols := make([]string, 0, len(values))
	args := make([]any, 0, len(values))
	placeholders := make([]string, 0, len(values))

	index := 1

	for _, value := range values {
		if value.Type != models.Default {
			cols = append(cols, driver.FormatReference(value.Column))
		}

		if value.Value != nil && value.Type != models.Default {
			placeholders = append(placeholders, driver.FormatPlaceholder(index))
			args = append(args, value.Value)
			index++
		}
	}

	queryStr := "INSERT INTO " + formattedTableName
	queryStr += fmt.Sprintf(" (%s) VALUES (%s)", strings.Join(cols, ", "), strings.Join(placeholders, ", "))

	newQuery := models.Query{
		Query: queryStr,
		Args:  args,
	}

	return newQuery
}

func buildUpdateQueryString(sanitizedTableName string, colNames []string, args []any, primaryKeyInfo []models.PrimaryKeyInfo, driver Driver) string {
	queryStr := "UPDATE " + sanitizedTableName

	sanitizedColNames := make([]string, len(colNames))
	for i, colName := range colNames {
		sanitizedColNames[i] = driver.FormatReference(colName)
	}

	sanitizedPrimaryKeyInfo := make([]models.PrimaryKeyInfo, len(primaryKeyInfo))
	for i, pki := range primaryKeyInfo {
		sanitizedPrimaryKeyInfo[i] = models.PrimaryKeyInfo{
			Name:  driver.FormatReference(pki.Name),
			Value: driver.FormatArgForQueryString(pki.Value),
		}
	}

	sanitizedArgs := make([]any, len(args))
	for i, arg := range args {
		sanitizedArgs[i] = driver.FormatArgForQueryString(arg)
	}

	for i, sanitizedColName := range sanitizedColNames {
		if i == 0 {
			queryStr += fmt.Sprintf(" SET %s = %s", sanitizedColName, sanitizedArgs[i])
		} else {
			queryStr += fmt.Sprintf(", %s = %s", sanitizedColName, sanitizedArgs[i])
		}
	}

	for i, sanitizedPki := range sanitizedPrimaryKeyInfo {
		if i == 0 {
			queryStr += fmt.Sprintf(" WHERE %s = %s", sanitizedPki.Name, sanitizedPki.Value)
		} else {
			queryStr += fmt.Sprintf(" AND %s = %s", sanitizedPki.Name, sanitizedPki.Value)
		}
	}

	return queryStr
}

func buildUpdateQuery(sanitizedTableName string, values []models.CellValue, primaryKeyInfo []models.PrimaryKeyInfo, driver Driver) models.Query {
	argsWithoutDefaults := []models.CellValue{}

	for _, arg := range values {
		if arg.Type != models.Default {
			argsWithoutDefaults = append(argsWithoutDefaults, arg)
		}
	}

	placeholders := buildPlaceholders(values, driver)

	sanitizedCols := []string{}
	for _, value := range values {
		sanitizedCols = append(sanitizedCols, driver.FormatReference(value.Column))
	}

	sanitizedArgs := make([]any, len(argsWithoutDefaults))
	for i, arg := range argsWithoutDefaults {
		if arg.Type != models.Default {
			sanitizedArgs[i] = driver.FormatArg(arg.Value, arg.Type)
		}
	}

	sanitizedPrimaryKeyInfo := make([]models.PrimaryKeyInfo, len(primaryKeyInfo))
	for i, primaryKey := range primaryKeyInfo {
		sanitizedPrimaryKeyInfo[i] = models.PrimaryKeyInfo{
			Name:  driver.FormatReference(primaryKey.Name),
			Value: primaryKey.Value,
		}
	}

	queryStr := "UPDATE " + sanitizedTableName

	for i, sanitizedCol := range sanitizedCols {
		placeholder := placeholders[i]
		reference := sanitizedCol
		if i == 0 {
			queryStr += fmt.Sprintf(" SET %s = %s", reference, placeholder)
		} else {
			queryStr += fmt.Sprintf(", %s = %s", reference, placeholder)
		}
	}

	for i, sanitizedPki := range sanitizedPrimaryKeyInfo {
		placeholder := driver.FormatPlaceholder(len(argsWithoutDefaults) + i + 1)
		reference := sanitizedPki.Name

		if i == 0 {
			queryStr += fmt.Sprintf(" WHERE %s = %s", reference, placeholder)
		} else {
			queryStr += fmt.Sprintf(" AND %s = %s", reference, placeholder)
		}
		sanitizedArgs = append(sanitizedArgs, sanitizedPki.Value)
	}

	logger.Info("buildUpdateQueryString", map[string]any{"queryStr": queryStr, "sanitizedArgs": sanitizedArgs})

	newQuery := models.Query{
		Query: queryStr,
		Args:  sanitizedArgs,
	}

	return newQuery
}

func buildDeleteQueryString(sanitizedTableName string, primaryKeyInfo []models.PrimaryKeyInfo, driver Driver) string {
	queryStr := "DELETE FROM " + sanitizedTableName

	sanitizedPrimaryKeyInfo := make([]models.PrimaryKeyInfo, len(primaryKeyInfo))
	for i, pki := range primaryKeyInfo {
		sanitizedPrimaryKeyInfo[i] = models.PrimaryKeyInfo{
			Name:  driver.FormatReference(pki.Name),
			Value: driver.FormatArgForQueryString(pki.Value),
		}
	}

	for i, sanitizedPki := range sanitizedPrimaryKeyInfo {
		if i == 0 {
			queryStr += fmt.Sprintf(" WHERE %s = %s", sanitizedPki.Name, sanitizedPki.Value)
		} else {
			queryStr += fmt.Sprintf(" AND %s = %s", sanitizedPki.Name, sanitizedPki.Value)
		}
	}

	return queryStr
}

func buildDeleteQuery(formattedTableName string, primaryKeyInfo []models.PrimaryKeyInfo, driver Driver) models.Query {
	queryStr := "DELETE FROM " + formattedTableName
	args := make([]any, len(primaryKeyInfo))

	sanitizedPrimaryKeyInfo := sanitizePrimaryKeyInfo(primaryKeyInfo, driver)

	for i, sanitizedPki := range sanitizedPrimaryKeyInfo {
		placeholder := driver.FormatPlaceholder(i + 1)
		reference := sanitizedPki.Name

		if i == 0 {
			queryStr += fmt.Sprintf(" WHERE %s = %s", reference, placeholder)
		} else {
			queryStr += fmt.Sprintf(" AND %s = %s", reference, placeholder)
		}
		args[i] = sanitizedPki.Value
	}

	return models.Query{
		Query: queryStr,
		Args:  args,
	}
}

func sanitizePrimaryKeyInfo(primaryKeyInfo []models.PrimaryKeyInfo, driver Driver) []models.PrimaryKeyInfo {
	sanitizedPrimaryKeyInfo := []models.PrimaryKeyInfo{}

	for _, pki := range primaryKeyInfo {
		sanitizedPrimaryKeyInfo = append(sanitizedPrimaryKeyInfo, models.PrimaryKeyInfo{
			Name:  driver.FormatReference(pki.Name),
			Value: pki.Value,
		})
	}

	return sanitizedPrimaryKeyInfo
}

func getColNamesAndArgsAsString(values []models.CellValue) ([]string, []any) {
	cols := []string{}
	v := []any{}

	for _, cell := range values {

		cols = append(cols, cell.Column)

		switch cell.Type {
		case models.Empty:
			v = append(v, "")
		case models.Null:
			v = append(v, "NULL")
		case models.Default:
			v = append(v, "DEFAULT")
		default:
			v = append(v, cell.Value)
		}
	}

	return cols, v
}

func buildPlaceholders(values []models.CellValue, driver Driver) []string {
	placeholders := []string{}

	index := 1

	for _, cell := range values {
		switch cell.Type {
		// case models.Empty:
		// placeholders = append(placeholders, "")
		// case models.Null:
		// 	placeholders = append(placeholders, "NULL")
		case models.Default:
			placeholders = append(placeholders, "DEFAULT")
			index--
		default:
			placeholders = append(placeholders, driver.FormatPlaceholder(index))
			index++
		}
	}
	return placeholders
}
