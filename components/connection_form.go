package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/internal/app"
	"sqlcmder/internal/drivers"
	"sqlcmder/internal/helpers"
	"sqlcmder/models"
)

type ConnectionForm struct {
	*tview.Flex
	*tview.Form
	StatusText *tview.TextView
	Action     string
	// Individual form fields for easy access
	DbTypeField *tview.InputField
	NameField   *tview.InputField
	HostField   *tview.InputField
	PortField   *tview.InputField
	UserField   *tview.InputField
	PassField   *tview.InputField
	DBNameField *tview.InputField
	SSLCheckbox *tview.Checkbox
	DSNField    *tview.InputField
}

func NewConnectionForm(connectionPages *models.ConnectionPages) *ConnectionForm {
	wrapper := tview.NewFlex()
	wrapper.SetDirection(tview.FlexColumnCSS)

	// Create individual form fields with defaults for PostgreSQL
	dbTypeField := tview.NewInputField().SetLabel("Database Type").SetText(drivers.DriverPostgres).SetFieldWidth(0)
	nameField := tview.NewInputField().SetLabel("Connection Name").SetFieldWidth(0)
	hostField := tview.NewInputField().SetLabel("Hostname").SetText("localhost").SetFieldWidth(0)
	portField := tview.NewInputField().SetLabel("Port").SetText("5432").SetFieldWidth(0)
	userField := tview.NewInputField().SetLabel("Username").SetText("postgres").SetFieldWidth(0)
	passField := tview.NewInputField().SetLabel("Password").SetText("postgres").SetFieldWidth(0)
	dbNameField := tview.NewInputField().SetLabel("DB Name").SetFieldWidth(0)
	sslCheckbox := tview.NewCheckbox().SetLabel("SSL Mode").SetChecked(false)
	dsnField := tview.NewInputField().SetLabel("DSN (Auto)").SetFieldWidth(0)

	// Helper function to auto-generate DSN
	updateDSN := func() {
		dbType := dbTypeField.GetText()
		hostname := hostField.GetText()
		port := portField.GetText()
		username := userField.GetText()
		password := passField.GetText()
		database := dbNameField.GetText()
		sslMode := sslCheckbox.IsChecked()

		// Build connection string
		var connectionString string
		switch dbType {
		case drivers.DriverPostgres:
			sslModeStr := "disable"
			if sslMode {
				sslModeStr = "require"
			}
			if username != "" && password != "" {
				connectionString = "postgres://" + username + ":" + password + "@" + hostname + ":" + port + "/" + database + "?sslmode=" + sslModeStr
			} else {
				connectionString = "postgres://" + hostname + ":" + port + "/" + database + "?sslmode=" + sslModeStr
			}
		case drivers.DriverMySQL:
			if username != "" && password != "" {
				connectionString = username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + database
			} else {
				connectionString = "@tcp(" + hostname + ":" + port + ")/" + database
			}
		case drivers.DriverSqlite:
			connectionString = database
		case drivers.DriverMSSQL:
			sslModeStr := "disable"
			if sslMode {
				sslModeStr = "true"
			}
			if username != "" && password != "" {
				connectionString = "sqlserver://" + username + ":" + password + "@" + hostname + ":" + port + "?database=" + database + "&encrypt=" + sslModeStr
			} else {
				connectionString = "sqlserver://" + hostname + ":" + port + "?database=" + database + "&encrypt=" + sslModeStr
			}
		}

		dsnField.SetText(connectionString)
	}

	// Set change handlers to auto-update DSN
	dbTypeField.SetChangedFunc(func(text string) { updateDSN() })
	hostField.SetChangedFunc(func(text string) { updateDSN() })
	portField.SetChangedFunc(func(text string) { updateDSN() })
	userField.SetChangedFunc(func(text string) { updateDSN() })
	passField.SetChangedFunc(func(text string) { updateDSN() })
	dbNameField.SetChangedFunc(func(text string) { updateDSN() })

	// Generate initial DSN
	updateDSN()

	// Set colors for all fields
	for _, field := range []*tview.InputField{dbTypeField, nameField, hostField, portField, userField, passField, dbNameField, dsnField} {
		field.SetFieldBackgroundColor(app.Styles.InverseTextColor)
		field.SetLabelColor(app.Styles.PrimaryTextColor)
		field.SetFieldTextColor(app.Styles.ContrastSecondaryTextColor)
	}
	sslCheckbox.SetFieldBackgroundColor(app.Styles.InverseTextColor)
	sslCheckbox.SetLabelColor(app.Styles.PrimaryTextColor)

	// Create left column form
	leftForm := tview.NewForm()
	leftForm.SetFieldBackgroundColor(app.Styles.InverseTextColor)
	leftForm.SetLabelColor(app.Styles.PrimaryTextColor)
	leftForm.AddFormItem(nameField)   // 1. Connection Name
	leftForm.AddFormItem(userField)   // 2. Username
	leftForm.AddFormItem(passField)   // 3. Password
	leftForm.AddFormItem(dbNameField) // 4. DB Name
	leftForm.SetBorder(false)

	// Create right column form
	rightForm := tview.NewForm()
	rightForm.SetFieldBackgroundColor(app.Styles.InverseTextColor)
	rightForm.SetLabelColor(app.Styles.PrimaryTextColor)
	rightForm.AddFormItem(dbTypeField) // 1. Database Type
	rightForm.AddFormItem(hostField)   // 2. Hostname
	rightForm.AddFormItem(portField)   // 3. Port
	rightForm.AddFormItem(dsnField)    // 4. DSN (Optional)
	rightForm.SetBorder(false)

	// Create two-column layout
	formLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	formLayout.AddItem(leftForm, 0, 1, true)
	formLayout.AddItem(rightForm, 0, 1, false)

	buttonsWrapper := tview.NewFlex().SetDirection(tview.FlexColumn)

	saveButton := tview.NewButton("[yellow]F1 [dark]Save")
	saveButton.SetStyle(tcell.StyleDefault.Background(app.Styles.PrimaryTextColor))
	saveButton.SetBorder(true)

	buttonsWrapper.AddItem(saveButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	testButton := tview.NewButton("[yellow]F2 [dark]Test")
	testButton.SetStyle(tcell.StyleDefault.Background(app.Styles.PrimaryTextColor))
	testButton.SetBorder(true)

	buttonsWrapper.AddItem(testButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	connectButton := tview.NewButton("[yellow]F3 [dark]Connect")
	connectButton.SetStyle(tcell.StyleDefault.Background(app.Styles.PrimaryTextColor))
	connectButton.SetBorder(true)

	buttonsWrapper.AddItem(connectButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	cancelButton := tview.NewButton("[yellow]Esc [dark]Cancel")
	cancelButton.SetStyle(tcell.StyleDefault.Background(tcell.Color(app.Styles.PrimaryTextColor)))
	cancelButton.SetBorder(true)

	buttonsWrapper.AddItem(cancelButton, 0, 1, false)

	statusText := tview.NewTextView()
	statusText.SetBorderPadding(0, 0, 0, 0)

	// Shortcuts hint
	shortcutsHint := tview.NewTextView()
	shortcutsHint.SetText("[yellow]Alt+P[white] PostgreSQL  [yellow]Alt+M[white] MySQL  [yellow]Alt+S[white] SQLite  [yellow]Alt+Q[white] SQL Server")
	shortcutsHint.SetTextAlign(tview.AlignCenter)
	shortcutsHint.SetDynamicColors(true)

	wrapper.AddItem(formLayout, 0, 1, true)
	wrapper.AddItem(statusText, 2, 0, false)
	wrapper.AddItem(buttonsWrapper, 3, 0, false)
	wrapper.AddItem(shortcutsHint, 1, 0, false)

	form := &ConnectionForm{
		Flex:        wrapper,
		Form:        leftForm, // Use left form for compatibility
		StatusText:  statusText,
		DbTypeField: dbTypeField,
		NameField:   nameField,
		HostField:   hostField,
		PortField:   portField,
		UserField:   userField,
		PassField:   passField,
		DBNameField: dbNameField,
		SSLCheckbox: sslCheckbox,
		DSNField:    dsnField,
	}

	// Define tab order: row by row (left to right, top to bottom)
	tabOrder := []tview.Primitive{
		nameField,   // Row 1 Left
		dbTypeField, // Row 1 Right
		userField,   // Row 2 Left
		hostField,   // Row 2 Right
		passField,   // Row 3 Left
		portField,   // Row 3 Right
		dbNameField, // Row 4 Left
		dsnField,    // Row 4 Right
	}

	// Setup custom tab navigation
	setupTabNavigation := func(field tview.Primitive, index int) {
		if f, ok := field.(*tview.InputField); ok {
			f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyTab {
					nextIndex := (index + 1) % len(tabOrder)
					App.SetFocus(tabOrder[nextIndex])
					return nil
				} else if event.Key() == tcell.KeyBacktab {
					prevIndex := (index - 1 + len(tabOrder)) % len(tabOrder)
					App.SetFocus(tabOrder[prevIndex])
					return nil
				}
				return event
			})
		}
	}

	for i, field := range tabOrder {
		setupTabNavigation(field, i)
	}

	// Set default focus to first input field (Connection Name)
	App.SetFocus(nameField)

	wrapper.SetInputCapture(form.inputCapture(connectionPages))

	return form
}

func (form *ConnectionForm) inputCapture(connectionPages *models.ConnectionPages) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		// Handle Alt+Key shortcuts for quick database type selection
		if event.Key() == tcell.KeyRune && event.Modifiers()&tcell.ModAlt != 0 {
			switch event.Rune() {
			case 'p', 'P': // Alt+P for PostgreSQL
				form.setDatabasePreset(drivers.DriverPostgres)
				return nil
			case 'm', 'M': // Alt+M for MySQL
				form.setDatabasePreset(drivers.DriverMySQL)
				return nil
			case 's', 'S': // Alt+S for SQLite
				form.setDatabasePreset(drivers.DriverSqlite)
				return nil
			case 'q', 'Q': // Alt+Q for SQL Server
				form.setDatabasePreset(drivers.DriverMSSQL)
				return nil
			}
		}

		if event.Key() == tcell.KeyEsc {
			connectionPages.SwitchToPage(pageNameConnectionSelection)
		} else if event.Key() == tcell.KeyF1 || event.Key() == tcell.KeyEnter {
			// Get form field values
			dbType := form.DbTypeField.GetText()
			connectionName := form.NameField.GetText()
			hostname := form.HostField.GetText()
			port := form.PortField.GetText()
			username := form.UserField.GetText()
			password := form.PassField.GetText()
			database := form.DBNameField.GetText()
			sslMode := form.SSLCheckbox.IsChecked()
			dsn := form.DSNField.GetText()

			if connectionName == "" {
				form.StatusText.SetText("Connection name is required").SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))
				return event
			}

			// Build connection string from form fields or use DSN directly
			var connectionString string
			if dsn != "" {
				connectionString = dsn
			} else {
				connectionString = form.buildConnectionString(dbType, hostname, port, username, password, database, sslMode)
			}

			// Validate connection string only if it's not empty
			if connectionString != "" {
				_, err := helpers.ParseConnectionString(connectionString)
				if err != nil {
					form.StatusText.SetText("Warning: " + err.Error() + " (saved anyway)").SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow))
				}
			}

			databases := app.App.Connections()
			newDatabases := make([]models.Connection, len(databases))

			parsedDatabaseData := models.Connection{
				Name:     connectionName,
				Provider: dbType,
				Hostname: hostname,
				Port:     port,
				Username: username,
				Password: password,
				DBName:   database,
				URL:      connectionString,
			}

			switch form.Action {
			case actionNewConnection:

				newDatabases = append(databases, parsedDatabaseData)
				err := app.App.SaveConnections(newDatabases)
				if err != nil {
					form.StatusText.SetText("Save failed: " + err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))
					return event
				}
				configPath := app.App.GetConfigFilePath()
				form.StatusText.SetText("Saved to: " + configPath).SetTextColor(app.Styles.TertiaryTextColor)

			case actionEditConnection:
				newDatabases = make([]models.Connection, len(databases))
				row, _ := connectionsTable.GetSelection()

				for i, database := range databases {
					if i == row {
						newDatabases[i] = parsedDatabaseData

						// newDatabases[i].Name = connectionName
						// newDatabases[i].Provider = database.Provider
						// newDatabases[i].User = parsed.User.Username()
						// newDatabases[i].Password, _ = parsed.User.Password()
						// newDatabases[i].Host = parsed.Hostname()
						// newDatabases[i].Port = parsed.Port()
						// newDatabases[i].Query = parsed.Query().Encode()
						// newDatabases[i].DBName = helpers.ParsedDBName(parsed.Path)
						// newDatabases[i].DSN = parsed.DSN
					} else {
						newDatabases[i] = database
					}
				}

				err := app.App.SaveConnections(newDatabases)
				if err != nil {
					form.StatusText.SetText("Save failed: " + err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))
					return event
				}
				configPath := app.App.GetConfigFilePath()
				form.StatusText.SetText("Updated: " + configPath).SetTextColor(app.Styles.TertiaryTextColor)
			}

			connectionsTable.SetConnections(newDatabases)
			connectionPages.SwitchToPage(pageNameConnectionSelection)

		} else if event.Key() == tcell.KeyF2 {
			// Get form field values for testing
			dbType := form.DbTypeField.GetText()
			hostname := form.HostField.GetText()
			port := form.PortField.GetText()
			username := form.UserField.GetText()
			password := form.PassField.GetText()
			database := form.DBNameField.GetText()
			sslMode := form.SSLCheckbox.IsChecked()
			dsn := form.DSNField.GetText()

			// Build connection string from form fields or use DSN directly
			var connectionString string
			if dsn != "" {
				connectionString = dsn
			} else {
				connectionString = form.buildConnectionString(dbType, hostname, port, username, password, database, sslMode)
			}

			go form.testConnection(connectionString)
		} else if event.Key() == tcell.KeyF3 {
			// F3 - Connect directly
			// TODO: Add direct connect logic if needed
		}
		return event
	}
}

func (form *ConnectionForm) testConnection(connectionString string) {
	parsed, err := helpers.ParseConnectionString(connectionString)
	if err != nil {
		form.StatusText.SetText(err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))
		return
	}

	form.StatusText.SetText("Connecting...").SetTextColor(app.Styles.TertiaryTextColor)

	var db drivers.Driver

	switch parsed.Driver {
	case drivers.DriverMySQL:
		db = &drivers.MySQL{}
	case drivers.DriverPostgres:
		db = &drivers.Postgres{}
	case drivers.DriverSqlite:
		db = &drivers.SQLite{}
	case drivers.DriverMSSQL:
		db = &drivers.MSSQL{}
	}

	err = db.TestConnection(connectionString)

	if err != nil {
		form.StatusText.SetText(err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))
	} else {
		form.StatusText.SetText("Connection success").SetTextColor(app.Styles.TertiaryTextColor)
	}
	App.ForceDraw()
}

func (form *ConnectionForm) SetAction(action string) {
	form.Action = action
}

// setDatabasePreset sets the database type and fills in default values
func (form *ConnectionForm) setDatabasePreset(dbType string) {
	// Set database type field
	form.DbTypeField.SetText(dbType)

	// Set default values based on database type
	switch dbType {
	case drivers.DriverPostgres:
		form.HostField.SetText("localhost")
		form.PortField.SetText("5432")
		form.UserField.SetText("postgres")
		form.PassField.SetText("postgres")
		form.DBNameField.SetText("")
		form.SSLCheckbox.SetChecked(false)
	case drivers.DriverMySQL:
		form.HostField.SetText("localhost")
		form.PortField.SetText("3306")
		form.UserField.SetText("root")
		form.PassField.SetText("root")
		form.DBNameField.SetText("")
		form.SSLCheckbox.SetChecked(false)
	case drivers.DriverSqlite:
		form.HostField.SetText("")
		form.PortField.SetText("")
		form.UserField.SetText("")
		form.PassField.SetText("")
		form.DBNameField.SetText("./sqlite.db")
		form.SSLCheckbox.SetChecked(false)
	case drivers.DriverMSSQL:
		form.HostField.SetText("localhost")
		form.PortField.SetText("1433")
		form.UserField.SetText("")
		form.PassField.SetText("")
		form.DBNameField.SetText("")
		form.SSLCheckbox.SetChecked(false)
	}

	form.StatusText.SetText("Preset: " + dbType + " | Use Tab to navigate between fields").SetTextColor(app.Styles.TertiaryTextColor)
}

// buildConnectionString constructs a database connection string from individual components
func (form *ConnectionForm) buildConnectionString(dbType, hostname, port, username, password, database string, sslMode bool) string {
	var connectionString string

	switch dbType {
	case drivers.DriverPostgres:
		sslModeStr := "disable"
		if sslMode {
			sslModeStr = "require"
		}
		if username != "" && password != "" {
			connectionString = "postgres://" + username + ":" + password + "@" + hostname + ":" + port + "/" + database + "?sslmode=" + sslModeStr
		} else {
			connectionString = "postgres://" + hostname + ":" + port + "/" + database + "?sslmode=" + sslModeStr
		}

	case drivers.DriverMySQL:
		if username != "" && password != "" {
			connectionString = username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + database
		} else {
			connectionString = "@tcp(" + hostname + ":" + port + ")/" + database
		}

	case drivers.DriverSqlite:
		connectionString = database // SQLite just needs the database file path

	case drivers.DriverMSSQL:
		sslModeStr := "disable"
		if sslMode {
			sslModeStr = "true"
		}
		if username != "" && password != "" {
			connectionString = "sqlserver://" + username + ":" + password + "@" + hostname + ":" + port + "?database=" + database + "&encrypt=" + sslModeStr
		} else {
			connectionString = "sqlserver://" + hostname + ":" + port + "?database=" + database + "&encrypt=" + sslModeStr
		}
	}

	return connectionString
}
