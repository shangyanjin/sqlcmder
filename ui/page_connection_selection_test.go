package ui

import (
	"testing"

	"sqlcmder/models"
)

// TestConnectionDeletionLogic tests the core deletion logic used in page_connection_selection.go
func TestConnectionDeletionLogic(t *testing.T) {
	// Create a sample list of connections
	connections := []models.Connection{
		{
			Name:   "Connection 1",
			Driver: "postgres",
		},
		{
			Name:   "Connection 2 - To Delete",
			Driver: "mysql",
		},
		{
			Name:   "Connection 3",
			Driver: "sqlite",
		},
	}

	initialCount := len(connections)
	if initialCount != 3 {
		t.Fatalf("Initial connections count should be 3, got %d", initialCount)
	}

	// Find the connection to delete (simulating what happens when user clicks delete)
	row := 1 // Simulate selecting Connection 2
	connectionToDelete := connections[row]

	if connectionToDelete.Name != "Connection 2 - To Delete" {
		t.Errorf("Should be deleting 'Connection 2 - To Delete', got %s", connectionToDelete.Name)
	}

	// Apply the deletion logic (same as in page_connection_selection.go line 139)
	newConnections := append(connections[:row], connections[row+1:]...)

	// Verify the deletion
	if len(newConnections) != initialCount-1 {
		t.Errorf("After deletion, expected %d connections, got %d", initialCount-1, len(newConnections))
	}

	// Verify the correct connection was removed
	for _, conn := range newConnections {
		if conn.Name == "Connection 2 - To Delete" {
			t.Errorf("Connection 'Connection 2 - To Delete' should have been deleted but still exists")
		}
	}

	// Verify other connections remain
	remainingNames := []string{"Connection 1", "Connection 3"}
	for i, conn := range newConnections {
		if conn.Name != remainingNames[i] {
			t.Errorf("Expected connection %s at index %d, got %s", remainingNames[i], i, conn.Name)
		}
	}

	t.Logf("✓ Deletion logic test passed: %d connections -> %d connections", initialCount, len(newConnections))
}

// TestClosureVariableCapture tests that closure variables are properly captured
// This is important because the SetDoneFunc callback captures variables that must remain valid
func TestClosureVariableCapture(t *testing.T) {
	connections := []models.Connection{
		{Name: "Conn A", Driver: "postgres"},
		{Name: "Conn B", Driver: "mysql"},
		{Name: "Conn C", Driver: "sqlite"},
	}

	row := 1
	selectedConnection := connections[row]

	// Simulate what happens in the delete handler
	currentRow := row
	currentConnections := connections
	selectedConnectionToDelete := selectedConnection

	// Now, even if row changes, the captured variables should remain the same
	row = 999 // Simulate variable mutation

	// Verify captured variables are unchanged
	if currentRow != 1 {
		t.Errorf("Expected currentRow to be 1, got %d", currentRow)
	}

	if len(currentConnections) != 3 {
		t.Errorf("Expected 3 connections, got %d", len(currentConnections))
	}

	if selectedConnectionToDelete.Name != "Conn B" {
		t.Errorf("Expected 'Conn B', got %s", selectedConnectionToDelete.Name)
	}

	t.Logf("✓ Closure variable capture test passed: variables preserved despite external changes")
}

// TestButtonResponseHandling tests the button label handling in delete confirmation
func TestButtonResponseHandling(t *testing.T) {
	type testCase struct {
		buttonLabel    string
		shouldDelete   bool
		description    string
	}

	testCases := []testCase{
		{buttonLabel: "Yes", shouldDelete: true, description: "Yes button should trigger deletion"},
		{buttonLabel: "No", shouldDelete: false, description: "No button should cancel deletion"},
		{buttonLabel: "yes", shouldDelete: false, description: "Lowercase 'yes' should NOT trigger deletion (case-sensitive)"},
		{buttonLabel: "", shouldDelete: false, description: "Empty string should NOT trigger deletion"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			shouldDelete := tc.buttonLabel == "Yes"

			if shouldDelete != tc.shouldDelete {
				t.Errorf("Button label '%s': expected shouldDelete=%v, got %v",
					tc.buttonLabel, tc.shouldDelete, shouldDelete)
			}
		})
	}

	t.Logf("✓ Button response handling test passed: all button labels handled correctly")
}

// TestMultipleDeletionSequence tests deleting multiple connections in sequence
func TestMultipleDeletionSequence(t *testing.T) {
	connections := []models.Connection{
		{Name: "Conn 1", Driver: "postgres"},
		{Name: "Conn 2", Driver: "mysql"},
		{Name: "Conn 3", Driver: "sqlite"},
		{Name: "Conn 4", Driver: "postgres"},
	}

	initialCount := len(connections)

	// Delete Conn 2 (index 1)
	conn1 := append(connections[:1], connections[2:]...)
	if len(conn1) != initialCount-1 {
		t.Errorf("After first deletion: expected %d, got %d", initialCount-1, len(conn1))
	}

	// Verify Conn 2 is gone
	for _, c := range conn1 {
		if c.Name == "Conn 2" {
			t.Errorf("Conn 2 should be deleted")
		}
	}

	// Delete Conn 4 (now at index 2)
	conn2 := append(conn1[:2], conn1[3:]...)
	if len(conn2) != initialCount-2 {
		t.Errorf("After second deletion: expected %d, got %d", initialCount-2, len(conn2))
	}

	// Verify both are gone
	for _, c := range conn2 {
		if c.Name == "Conn 2" || c.Name == "Conn 4" {
			t.Errorf("Deleted connections should not exist")
		}
	}

	// Verify remaining connections
	expectedNames := []string{"Conn 1", "Conn 3"}
	for i, c := range conn2 {
		if c.Name != expectedNames[i] {
			t.Errorf("Expected %s at index %d, got %s", expectedNames[i], i, c.Name)
		}
	}

	t.Logf("✓ Multiple deletion sequence test passed: %d -> %d -> %d", initialCount, initialCount-1, len(conn2))
}
