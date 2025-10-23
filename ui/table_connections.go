package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
	"sqlcmder/models"
)

type ConnectionsTable struct {
	*tview.Table
	Wrapper       *tview.Flex
	errorTextView *tview.TextView
	error         string
	connections   []models.Connection
}

var connectionsTable *ConnectionsTable

func NewConnectionsTable() *ConnectionsTable {
	wrapper := tview.NewFlex()

	errorTextView := tview.NewTextView()
	errorTextView.SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))

	table := &ConnectionsTable{
		Table:         tview.NewTable().SetSelectable(true, false),
		Wrapper:       wrapper,
		errorTextView: errorTextView,
	}

	table.SetOffset(5, 0)
	table.SetSelectedStyle(tcell.StyleDefault.Foreground(app.Styles.SecondaryTextColor).Background(app.Styles.PrimitiveBackgroundColor))

	// Set selection changed callback to update * marker
	table.SetSelectionChangedFunc(func(row, column int) {
		table.UpdateSelectionMarker(row)
	})

	wrapper.AddItem(table, 0, 1, true)
	table.SetConnections(app.App.Connections())

	connectionsTable = table

	return connectionsTable
}

func (ct *ConnectionsTable) AddConnection(connection models.Connection) {
	rowCount := ct.GetRowCount()
	ct.SetCellSimple(rowCount, 0, "  "+connection.Name) // Reserve space for * marker
	ct.connections = append(ct.connections, connection)
}

func (ct *ConnectionsTable) GetConnections() []models.Connection {
	return ct.connections
}

func (ct *ConnectionsTable) GetError() string {
	return ct.error
}

func (ct *ConnectionsTable) SetConnections(connections []models.Connection) {
	ct.connections = make([]models.Connection, 0)

	ct.Clear()

	for _, connection := range connections {
		ct.AddConnection(connection)
	}

	ct.Select(0, 0)

	// Initialize selection marker for first row
	if len(connections) > 0 {
		ct.UpdateSelectionMarker(0)
	}

	App.ForceDraw()
}

func (ct *ConnectionsTable) SetError(err error) {
	ct.error = err.Error()
	ct.errorTextView.SetText(ct.error)
}

// UpdateSelectionMarker updates the * marker for the selected row
func (ct *ConnectionsTable) UpdateSelectionMarker(selectedRow int) {
	// Update all rows
	for i := 0; i < ct.GetRowCount(); i++ {
		if i >= len(ct.connections) {
			break
		}

		cell := ct.GetCell(i, 0)
		if i == selectedRow {
			// Add * marker for selected row
			cell.SetText("[yellow]*[white] " + ct.connections[i].Name)
		} else {
			// Remove * marker for other rows
			cell.SetText("  " + ct.connections[i].Name)
		}
	}
}
