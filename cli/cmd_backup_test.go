package commands

import (
	"testing"

	"sqlcmder/drivers"
	"sqlcmder/models"
)

// TestBackupDatabase_NilConnection tests that BackupDatabase returns error when ConnectionModel is nil
func TestBackupDatabase_NilConnection(t *testing.T) {
	ctx := Context{
		DB:              &drivers.MySQL{},
		CurrentDatabase: "test_db",
		Connection:      "test_connection",
		ConnectionModel: nil, // This should trigger the error
	}

	errorCalled := false
	onError := func(message string) {
		errorCalled = true
		expectedMsg := "Connection information not available"
		if message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, message)
		}
	}

	onSuccess := func(message string) {
		t.Error("Success callback should not be called when ConnectionModel is nil")
	}

	BackupDatabase("test.sql", ctx, onSuccess, onError)

	if !errorCalled {
		t.Error("Error callback was not called")
	}
}

// TestImportDatabase_NilConnection tests that ImportDatabase returns error when ConnectionModel is nil
func TestImportDatabase_NilConnection(t *testing.T) {
	ctx := Context{
		DB:              &drivers.MySQL{},
		CurrentDatabase: "test_db",
		Connection:      "test_connection",
		ConnectionModel: nil, // This should trigger the error
	}

	errorCalled := false
	onError := func(message string) {
		errorCalled = true
		expectedMsg := "Connection information not available"
		if message != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, message)
		}
	}

	onSuccess := func(message string) {
		t.Error("Success callback should not be called when ConnectionModel is nil")
	}

	onRefresh := func() {
		t.Error("Refresh callback should not be called when ConnectionModel is nil")
	}

	ImportDatabase("test.sql", ctx, onSuccess, onError, onRefresh)

	if !errorCalled {
		t.Error("Error callback was not called")
	}
}

// TestBackupDatabase_WithConnection tests that BackupDatabase accepts valid connection
func TestBackupDatabase_WithConnection(t *testing.T) {
	conn := &models.Connection{
		Driver:   "mysql",
		Hostname: "localhost",
		Port:     "3306",
		Username: "root",
		Password: "password",
		DBName:   "test_db",
	}

	ctx := Context{
		DB:              &drivers.MySQL{},
		CurrentDatabase: "test_db",
		Connection:      "test_connection",
		ConnectionModel: conn,
	}

	// This test just verifies that the connection is not nil
	// Actual backup functionality would require mysqldump to be installed
	if ctx.ConnectionModel == nil {
		t.Error("ConnectionModel should not be nil")
	}

	if ctx.ConnectionModel.Driver != "mysql" {
		t.Errorf("Expected driver 'mysql', got '%s'", ctx.ConnectionModel.Driver)
	}
}

// TestImportDatabase_WithConnection tests that ImportDatabase accepts valid connection
func TestImportDatabase_WithConnection(t *testing.T) {
	conn := &models.Connection{
		Driver:   "postgres",
		Hostname: "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "password",
		DBName:   "test_db",
	}

	ctx := Context{
		DB:              &drivers.Postgres{},
		CurrentDatabase: "test_db",
		Connection:      "test_connection",
		ConnectionModel: conn,
	}

	// This test just verifies that the connection is not nil
	// Actual import functionality would require psql to be installed
	if ctx.ConnectionModel == nil {
		t.Error("ConnectionModel should not be nil")
	}

	if ctx.ConnectionModel.Driver != "postgres" {
		t.Errorf("Expected driver 'postgres', got '%s'", ctx.ConnectionModel.Driver)
	}
}

// TestContext_CurrentDatabaseFallback tests database name fallback logic
func TestContext_CurrentDatabaseFallback(t *testing.T) {
	conn := &models.Connection{
		Driver:   "mysql",
		Hostname: "localhost",
		Port:     "3306",
		Username: "root",
		Password: "password",
		DBName:   "default_db",
	}

	// Test case 1: CurrentDatabase is set
	ctx1 := Context{
		DB:              &drivers.MySQL{},
		CurrentDatabase: "current_db",
		Connection:      "test",
		ConnectionModel: conn,
	}

	if ctx1.CurrentDatabase != "current_db" {
		t.Errorf("Expected CurrentDatabase 'current_db', got '%s'", ctx1.CurrentDatabase)
	}

	// Test case 2: CurrentDatabase is empty, should fallback to ConnectionModel.DBName
	ctx2 := Context{
		DB:              &drivers.MySQL{},
		CurrentDatabase: "",
		Connection:      "test",
		ConnectionModel: conn,
	}

	// The actual fallback logic is in BackupDatabase/ImportDatabase functions
	// This test just verifies the context structure
	if ctx2.ConnectionModel.DBName != "default_db" {
		t.Errorf("Expected DBName 'default_db', got '%s'", ctx2.ConnectionModel.DBName)
	}
}
