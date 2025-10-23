package ui

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	commands "sqlcmder/cli"
	"sqlcmder/cmd/app"
	"sqlcmder/drivers"
	"sqlcmder/helpers"
	"sqlcmder/keymap"
	"sqlcmder/logger"
	"sqlcmder/models"
)

type ConnectionSelection struct {
	*tview.Flex
	StatusText *tview.TextView
}

func NewConnectionSelection(connectionForm *ConnectionForm, connectionPages *models.ConnectionPages) *ConnectionSelection {
	wrapper := tview.NewFlex()

	wrapper.SetDirection(tview.FlexColumnCSS)

	buttonsWrapper := tview.NewFlex().SetDirection(tview.FlexRowCSS)

	newButton := tview.NewButton("[yellow]N[dark]ew")
	newButton.SetStyle(tcell.StyleDefault.Background(app.Styles.PrimitiveBackgroundColor))
	newButton.SetBorder(true)

	buttonsWrapper.AddItem(newButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	connectButton := tview.NewButton("[yellow]C[dark]onnect")
	connectButton.SetStyle(tcell.StyleDefault.Background(app.Styles.PrimitiveBackgroundColor))
	connectButton.SetBorder(true)

	buttonsWrapper.AddItem(connectButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	editButton := tview.NewButton("[yellow]E[dark]dit")
	editButton.SetStyle(tcell.StyleDefault.Background(app.Styles.PrimitiveBackgroundColor))
	editButton.SetBorder(true)

	buttonsWrapper.AddItem(editButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	deleteButton := tview.NewButton("[yellow]D[dark]elete")
	deleteButton.SetStyle(tcell.StyleDefault.Background(app.Styles.PrimitiveBackgroundColor))
	deleteButton.SetBorder(true)

	buttonsWrapper.AddItem(deleteButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	quitButton := tview.NewButton("[yellow]Q[dark]uit")
	quitButton.SetStyle(tcell.StyleDefault.Background(app.Styles.PrimitiveBackgroundColor))
	quitButton.SetBorder(true)

	buttonsWrapper.AddItem(quitButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	statusText := tview.NewTextView()
	statusText.SetBorderPadding(1, 1, 0, 0)

	// Hint text above buttons
	hintText := tview.NewTextView()
	hintText.SetText("  [yellow]Up/Down[white] Select  [yellow]Enter[white] Connect  [yellow]N[white]ew  [yellow]E[white]dit  [yellow]D[white]elete  [yellow]Q[white]uit")
	hintText.SetDynamicColors(true)
	hintText.SetTextAlign(tview.AlignCenter)
	hintText.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)

	wrapper.AddItem(NewConnectionsTable(), 0, 1, true)
	wrapper.AddItem(statusText, 4, 0, false)
	wrapper.AddItem(hintText, 1, 0, false) // Hint area (1 line)
	wrapper.AddItem(buttonsWrapper, 3, 0, false)

	cs := &ConnectionSelection{
		Flex:       wrapper,
		StatusText: statusText,
	}

	wrapper.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		connections := connectionsTable.GetConnections()

		command := keymap.Keymaps.Group(keymap.ConnectionGroup).Resolve(event)

		if len(connections) != 0 {
			row, _ := connectionsTable.GetSelection()
			selectedConnection := connections[row]

			switch command {
			case commands.Connect:
				go cs.Connect(selectedConnection)
			case commands.EditConnection:
				connectionPages.SwitchToPage(pageNameConnectionForm)
				connectionForm.NameField.SetText(selectedConnection.Name)
				connectionForm.DSNField.SetText(selectedConnection.GetDSN())
				connectionForm.StatusText.SetText("")
				// Show DSN hint/value for edit connection
				connectionForm.showDSNHint()

				connectionForm.SetAction(actionEditConnection)
				return nil
			case commands.DeleteConnection:
				// Capture the current row and connections in this scope to avoid closure issues
				currentRow := row
				currentConnections := connections
				selectedConnectionToDelete := selectedConnection

				logger.Info("Delete connection initiated", map[string]any{
					"connectionName":   selectedConnectionToDelete.Name,
					"row":              currentRow,
					"totalConnections": len(currentConnections),
				})

				confirmationModal := NewConfirmationModal("")

				confirmationModal.SetDoneFunc(func(_ int, buttonLabel string) {
					defer mainPages.RemovePage(pageNameConfirmation)

					logger.Info("Delete confirmation response", map[string]any{
						"buttonLabel":    buttonLabel,
						"connectionName": selectedConnectionToDelete.Name,
					})

					if buttonLabel == "Yes" {
						logger.Info("Deleting connection", map[string]any{
							"connectionName":   selectedConnectionToDelete.Name,
							"row":              currentRow,
							"totalConnections": len(currentConnections),
						})

						newConnections := append(currentConnections[:currentRow], currentConnections[currentRow+1:]...)

						logger.Info("New connections list after deletion", map[string]any{
							"totalConnectionsBefore": len(currentConnections),
							"totalConnectionsAfter":  len(newConnections),
							"deletedConnection":      selectedConnectionToDelete.Name,
						})

						err := app.App.SaveConnections(newConnections)
						if err != nil {
							logger.Error("Failed to save connections after delete", map[string]any{
								"error":          err.Error(),
								"connectionName": selectedConnectionToDelete.Name,
							})
							connectionsTable.SetError(err)
						} else {
							logger.Info("Successfully saved connections after deletion", map[string]any{
								"totalConnections":  len(newConnections),
								"deletedConnection": selectedConnectionToDelete.Name,
							})
							connectionsTable.SetConnections(newConnections)

							logger.Info("Updated connections table display", map[string]any{
								"totalConnections": len(newConnections),
							})
						}
					} else {
						logger.Info("Delete operation cancelled by user", map[string]any{
							"connectionName": selectedConnectionToDelete.Name,
						})
					}
				})

				mainPages.AddPage(pageNameConfirmation, confirmationModal, true, true)
				App.SetFocus(confirmationModal)

				return nil
			}
		}

		switch command {
		case commands.NewConnection:
			connectionForm.SetAction(actionNewConnection)
			// Reset to PostgreSQL defaults
			connectionForm.DbTypeField.SetText(drivers.DriverPostgres)
			connectionForm.NameField.SetText("")
			connectionForm.HostField.SetText("localhost")
			connectionForm.PortField.SetText("5432")
			connectionForm.UserField.SetText("postgres")
			connectionForm.PassField.SetText("postgres")
			connectionForm.DBNameField.SetText("")
			connectionForm.DSNField.SetText("")
			connectionForm.StatusText.SetText("")
			// Show DSN hint for new connection
			connectionForm.showDSNHint()
			connectionPages.SwitchToPage(pageNameConnectionForm)
		case commands.Quit:
			if wrapper.HasFocus() {
				app.App.Stop()
			}
		}

		return event
	})

	return cs
}

func (cs *ConnectionSelection) Connect(connection models.Connection) *tview.Application {
	if mainPages.HasPage(connection.Name) {
		mainPages.SwitchToPage(connection.Name)
		return App.Draw()
	}

	if len(connection.Commands) > 0 {

		// Contains variables -- both the generated port and user-defined.
		variables := map[string]string{}

		// Avoid getting the port when it's not requested.
		waitsForPort := strings.Contains(connection.GetDSN(), "${port}")
		waitsForPort = waitsForPort || slices.ContainsFunc(connection.Commands, func(command *models.Command) bool {
			return command.WaitForPort != ""
		})

		if waitsForPort {
			port, err := helpers.GetFreePort()
			if err != nil {
				cs.StatusText.SetText(err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
				return App.Draw()
			}
			// Add port variable for the auto-generated port.
			variables["port"] = port
		}

		for i, command := range connection.Commands {
			message := fmt.Sprintf("Running command %d/%d...", i+1, len(connection.Commands))
			cs.StatusText.SetText(message).SetTextColor(app.Styles.TertiaryTextColor)
			App.Draw()

			cmd := command.Command
			for variable, value := range variables {
				cmd = strings.ReplaceAll(cmd, "${"+variable+"}", value)
			}

			markCommandComplete := App.Register()
			onCommandDone, waitToCaptureVariable := setupOutputVariableCommand(variables, command, markCommandComplete)

			if err := helpers.RunCommand(App.Context(), cmd, onCommandDone); err != nil {
				cs.StatusText.SetText(err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
				return App.Draw()
			}

			waitToCaptureVariable()

			if command.WaitForPort != "" {
				interpolatedPort := command.WaitForPort
				for variable, value := range variables {
					interpolatedPort = strings.ReplaceAll(interpolatedPort, "${"+variable+"}", value)
				}

				if portInt, err := strconv.Atoi(interpolatedPort); err != nil || portInt < 0 || portInt >= 1<<16 {
					cs.StatusText.SetText("bad port: " + interpolatedPort).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
					return App.Draw()
				}

				message := fmt.Sprintf("Waiting for port %s...", interpolatedPort)
				cs.StatusText.SetText(message).SetTextColor(app.Styles.TertiaryTextColor)
				App.Draw()

				if err := helpers.WaitForPort(App.Context(), interpolatedPort); err != nil {
					cs.StatusText.SetText(err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
					return App.Draw()
				}
			}
		}

		// Replace variables in URL.
		for variable, value := range variables {
			if variable == "" || value == "" {
				continue
			}
			dsnValue := connection.GetDSN()
			dsnValue = strings.ReplaceAll(dsnValue, "${"+variable+"}", value)
			// Update the appropriate DSN field
			if connection.DsnCustom != "" {
				connection.DsnCustom = dsnValue
			} else if connection.DsnAuto != "" {
				connection.DsnAuto = dsnValue
			} else {
				connection.DSN = dsnValue
			}
			connection.SetDSNValue()
		}
	}

	cs.StatusText.SetText("Connecting...").SetTextColor(app.Styles.TertiaryTextColor)
	App.Draw()

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
	}

	err := newDBDriver.Connect(connection.GetDSN())
	if err != nil {
		cs.StatusText.SetText(err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
		return App.Draw()
	}

	selectedRow, selectedCol := connectionsTable.GetSelection()
	cell := connectionsTable.GetCell(selectedRow, selectedCol)
	cell.SetText(fmt.Sprintf("[green]* %s", cell.Text))
	cs.StatusText.SetText("")

	newHome := NewHomePage(connection, newDBDriver)
	newHome.Tree.SetCurrentNode(newHome.Tree.GetRoot())
	newHome.Tree.Wrapper.SetTitle(connection.Name)

	mainPages.AddAndSwitchToPage(connection.Name, newHome, true)
	App.SetFocus(newHome.Tree)

	return App.Draw()
}

// Produces two functions: [onCommandDone] should be passed to [helpers.RunCommand],
// and [captureVariable] should be called after. [captureVariable] will block until
// the output from the command is saved into [variables].
// If no [command.SaveOutputTo] is defined, [captureVariable] is a no-op.
func setupOutputVariableCommand(variables map[string]string, command *models.Command, markCommandComplete func()) (onCommandDone func(string), captureVariable func()) {
	if command.SaveOutputTo == "" {
		// No variable? Mark the command completed, but otherwise no-op.
		onCommandDone = func(_ string) { markCommandComplete() }
		return onCommandDone, func() {}
	}

	// When the command runs, the stdout will be passed through this channel.
	variableSaved := make(chan string)

	// To capture the variable, we receive from the channel; onCommandDone sends
	// on that channel, so we're just synchronizing with the completion of the command.
	captureVariable = func() {
		output := <-variableSaved
		variables[command.SaveOutputTo] = output
		logger.Debug("Saved command output to variable", map[string]any{"Variable": command.SaveOutputTo, "Output": output, "Command": command.Command})
	}

	onCommandDone = func(output string) {
		variableSaved <- output
		close(variableSaved)
		markCommandComplete()
	}

	return onCommandDone, captureVariable
}
