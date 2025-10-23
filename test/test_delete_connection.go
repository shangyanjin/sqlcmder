package main

import (
	"testing"

	"sqlcmder/cmd/app"
	"sqlcmder/models"
)

// TestDeleteConnection tests the connection deletion functionality
func TestDeleteConnection(t *testing.T) {
	// Initialize app
	_ = app.App

	// Get initial connections count
	initialConnections := app.App.Connections()
	initialCount := len(initialConnections)

	// Add a test connection
	testConn := models.Connection{
		Name:     "Test Delete Connection",
		Driver:   "postgres",
		Hostname: "localhost",
		Port:     "5432",
		Username: "test",
		Password: "test",
		DBName:   "testdb",
	}

	newConnections := append(initialConnections, testConn)
	err := app.App.SaveConnections(newConnections)
	if err != nil {
		t.Fatalf("Failed to save test connection: %v", err)
	}

	// Verify connection was added
	afterAdd := app.App.Connections()
	if len(afterAdd) != initialCount+1 {
		t.Errorf("Expected %d connections after add, got %d", initialCount+1, len(afterAdd))
	}

	// Verify we can find the test connection
	found := false
	for _, conn := range afterAdd {
		if conn.Name == "Test Delete Connection" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Test connection should exist after adding")
	}

	// Delete the test connection
	afterDelete := make([]models.Connection, 0, len(afterAdd)-1)
	for _, conn := range afterAdd {
		if conn.Name != "Test Delete Connection" {
			afterDelete = append(afterDelete, conn)
		}
	}

	err = app.App.SaveConnections(afterDelete)
	if err != nil {
		t.Fatalf("Failed to save connections after delete: %v", err)
	}

	// Verify connection was deleted
	final := app.App.Connections()
	if len(final) != initialCount {
		t.Errorf("Expected %d connections after delete, got %d", initialCount, len(final))
	}

	// Verify the specific connection was removed
	for _, conn := range final {
		if conn.Name == "Test Delete Connection" {
			t.Errorf("Test connection should have been deleted but still exists")
		}
	}

	t.Logf("✓ Delete connection test passed: %d -> %d -> %d", initialCount, initialCount+1, len(final))
}
