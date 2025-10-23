package ui

import (
	"net/url"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
	"sqlcmder/drivers"
	"sqlcmder/helpers"
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
	dbTypeField := tview.NewInputField().SetLabel("DB Type").SetText(drivers.DriverPostgres).SetFieldWidth(0)
	nameField := tview.NewInputField().SetLabel("Conn Name").SetFieldWidth(0)
	hostField := tview.NewInputField().SetLabel("Hostname").SetText("localhost").SetFieldWidth(0)
	portField := tview.NewInputField().SetLabel("Port").SetText("5432").SetFieldWidth(0)
	userField := tview.NewInputField().SetLabel("Username").SetText("postgres").SetFieldWidth(0)
	passField := tview.NewInputField().SetLabel("Password").SetText("postgres").SetFieldWidth(0)
	dbNameField := tview.NewInputField().SetLabel("DB Name").SetFieldWidth(0)
	sslCheckbox := tview.NewCheckbox().SetLabel("SSL Mode").SetChecked(false)
	dsnField := tview.NewInputField().SetLabel("DSN").SetFieldWidth(0)

	// Helper function to auto-generate DSN
	generateAutoDSN := func() string {
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
				connectionString = "mysql://" + username + ":" + password + "@" + hostname + ":" + port + "/" + database
			} else {
				connectionString = "mysql://" + hostname + ":" + port + "/" + database
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

		return connectionString
	}

	// Function to update DSN field with auto-generated value
	updateDSNField := func() {
		autoDSN := generateAutoDSN()
		dsnField.SetText(autoDSN)
	}

	// Generate initial DSN
	updateDSNField()

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

	saveButton := tview.NewButton("[yellow]F1 [dark]Save & Test")
	saveButton.SetStyle(tcell.StyleDefault.Background(app.Styles.ButtonBackgroundColor))
	saveButton.SetBorder(true)

	buttonsWrapper.AddItem(saveButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	autoButton := tview.NewButton("[yellow]F2 [dark]Auto DSN")
	autoButton.SetStyle(tcell.StyleDefault.Background(app.Styles.ButtonBackgroundColor))
	autoButton.SetBorder(true)

	buttonsWrapper.AddItem(autoButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	connectButton := tview.NewButton("[yellow]F3 [dark]Connect")
	connectButton.SetStyle(tcell.StyleDefault.Background(app.Styles.ButtonBackgroundColor))
	connectButton.SetBorder(true)

	buttonsWrapper.AddItem(connectButton, 0, 1, false)
	buttonsWrapper.AddItem(nil, 1, 0, false)

	cancelButton := tview.NewButton("[yellow]Esc [dark]Cancel")
	cancelButton.SetStyle(tcell.StyleDefault.Background(app.Styles.ButtonBackgroundColor))
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
			// F1 - Save & Test (test first, then save)
			go form.saveAndTestConnection()
		} else if event.Key() == tcell.KeyF2 {
			// F2 - Auto generate DSN
			form.autoGenerateDSN()
		} else if event.Key() == tcell.KeyF3 {
			// F3 - Save & Test + Connect
			go form.saveTestAndConnect()
		}
		return event
	}
}

// saveAndTestConnection tests the connection first, then saves if test passes
func (form *ConnectionForm) saveAndTestConnection() {
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
		form.StatusText.SetText("Connection name is required").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
		return
	}

	// Build connection string with priority: custom DSN > auto-generated
	var connectionString string
	var dsnCustom string
	dsnAuto := form.buildConnectionString(dbType, hostname, port, username, password, database, sslMode)

	// Use custom DSN if provided, otherwise use auto-generated and show hint
	if dsn != "" {
		dsnCustom = dsn
		connectionString = dsn
	} else {
		connectionString = dsnAuto
		// Show hint that DSN was auto-generated
		form.StatusText.SetText("[green]DSN: " + dsnAuto).SetDynamicColors(true)
		App.Draw()
	}

	// Test connection first
	form.StatusText.SetText("Testing connection...").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.WarningColor))
	App.Draw()

	// Test the connection
	testResult := form.testConnectionSync(connectionString)
	if !testResult {
		form.StatusText.SetText("Connection test failed. Please check your settings.").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
		return
	}

	// If test passes, proceed with save
	form.StatusText.SetText("Connection test passed. Saving...").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.SuccessColor))
	App.Draw()

	// Validate connection string only if it's not empty
	if connectionString != "" {
		_, err := helpers.ParseConnectionString(connectionString)
		if err != nil {
			form.StatusText.SetText("Warning: " + err.Error() + " (saved anyway)").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.WarningColor))
		}
	}

	databases := app.App.Connections()
	newDatabases := make([]models.Connection, len(databases))

	parsedDatabaseData := models.Connection{
		Name:      connectionName,
		Driver:    dbType,
		Hostname:  hostname,
		Port:      port,
		Username:  username,
		Password:  password,
		DBName:    database,
		DSN:       connectionString, // Keep for backward compatibility
		DsnCustom: dsnCustom,
		DsnAuto:   dsnAuto,
		DsnValue:  connectionString,
	}

	switch form.Action {
	case actionNewConnection:
		newDatabases = append(databases, parsedDatabaseData)
		err := app.App.SaveConnections(newDatabases)
		if err != nil {
			form.StatusText.SetText("Save failed: " + err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
			return
		}
		configPath := app.App.GetConfigFilePath()
		form.StatusText.SetText("Saved to: " + configPath).SetTextColor(app.Styles.TertiaryTextColor)

	case actionEditConnection:
		newDatabases = make([]models.Connection, len(databases))
		row, _ := connectionsTable.GetSelection()

		for i, database := range databases {
			if i == row {
				newDatabases[i] = parsedDatabaseData
			} else {
				newDatabases[i] = database
			}
		}

		err := app.App.SaveConnections(newDatabases)
		if err != nil {
			form.StatusText.SetText("Save failed: " + err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
			return
		}
		configPath := app.App.GetConfigFilePath()
		form.StatusText.SetText("Saved to: " + configPath).SetTextColor(app.Styles.TertiaryTextColor)
	}

	connectionsTable.SetConnections(newDatabases)
	// Note: connectionPages is not available in this context,
	// the page switching will be handled by the calling function
}

// testConnectionSync tests connection synchronously and returns true if successful
func (form *ConnectionForm) testConnectionSync(connectionString string) bool {
	// Parse connection string to get driver type
	parsedURL, err := url.Parse(connectionString)
	if err != nil {
		return false
	}

	var driver drivers.Driver
	switch parsedURL.Scheme {
	case "postgres":
		driver = &drivers.Postgres{}
	case "mysql":
		driver = &drivers.MySQL{}
	case "sqlite":
		driver = &drivers.SQLite{}
	case "sqlserver":
		driver = &drivers.MSSQL{}
	default:
		return false
	}

	// Test connection
	err = driver.Connect(connectionString)
	if err != nil {
		return false
	}

	// Connection successful
	return true
}

func (form *ConnectionForm) testConnection(connectionString string) {
	parsed, err := helpers.ParseConnectionString(connectionString)
	if err != nil {
		form.StatusText.SetText(err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
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
		form.StatusText.SetText(err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
	} else {
		form.StatusText.SetText("Connection success").SetTextColor(app.Styles.TertiaryTextColor)
	}
	App.ForceDraw()
}

// SetAction sets the action for the connection form (new or edit)
func (form *ConnectionForm) SetAction(action string) {
	form.Action = action
}

// showDSNHint displays DSN information in StatusText
// For new connections: shows hint if DSN is empty
// For edit connections: shows current DSN value
func (form *ConnectionForm) showDSNHint() {
	dsn := form.DSNField.GetText()

	if dsn == "" {
		// For new connection or when DSN is empty
		form.StatusText.SetText("[green]DSN: (auto-generate if empty)").SetDynamicColors(true)
	} else {
		// For edit connection or when DSN has value
		form.StatusText.SetText("[green]DSN: " + dsn).SetDynamicColors(true)
	}
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
			connectionString = "mysql://" + username + ":" + password + "@" + hostname + ":" + port + "/" + database
		} else {
			connectionString = "mysql://" + hostname + ":" + port + "/" + database
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

// getOrAutoGenerateDSN returns the DSN from field or auto-generates it if empty
// Also updates StatusText with hint if DSN was empty
func (form *ConnectionForm) getOrAutoGenerateDSN() string {
	dsn := form.DSNField.GetText()

	// If DSN is empty, auto-generate it and show hint
	if dsn == "" {
		dbType := form.DbTypeField.GetText()
		hostname := form.HostField.GetText()
		port := form.PortField.GetText()
		username := form.UserField.GetText()
		password := form.PassField.GetText()
		database := form.DBNameField.GetText()
		sslMode := form.SSLCheckbox.IsChecked()

		dsn = form.buildConnectionString(dbType, hostname, port, username, password, database, sslMode)
		form.StatusText.SetText("[green]DSN: " + dsn).SetDynamicColors(true)
	}

	return dsn
}

// autoGenerateDSN generates DSN automatically and updates the DSN field
func (form *ConnectionForm) autoGenerateDSN() {
	dbType := form.DbTypeField.GetText()
	hostname := form.HostField.GetText()
	port := form.PortField.GetText()
	username := form.UserField.GetText()
	password := form.PassField.GetText()
	database := form.DBNameField.GetText()
	sslMode := form.SSLCheckbox.IsChecked()

	autoDSN := form.buildConnectionString(dbType, hostname, port, username, password, database, sslMode)
	form.DSNField.SetText(autoDSN)

	// Show status message with DSN value - consistent with getOrAutoGenerateDSN style
	form.StatusText.SetText("[green]DSN: " + autoDSN).SetDynamicColors(true)
}

// saveTestAndConnect tests the connection, saves it, and then connects to the database
func (form *ConnectionForm) saveTestAndConnect() {
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
		form.StatusText.SetText("Connection name is required").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
		return
	}

	// Build connection string with priority: custom DSN > auto-generated
	var connectionString string
	var dsnCustom string
	dsnAuto := form.buildConnectionString(dbType, hostname, port, username, password, database, sslMode)

	// Use custom DSN if provided, otherwise use auto-generated and show hint
	if dsn != "" {
		dsnCustom = dsn
		connectionString = dsn
	} else {
		connectionString = dsnAuto
		// Show hint that DSN was auto-generated
		form.StatusText.SetText("[green]DSN: " + dsnAuto).SetDynamicColors(true)
		App.Draw()
	}

	// Test connection first
	form.StatusText.SetText("Testing connection...").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.WarningColor))
	App.Draw()

	// Test the connection
	testResult := form.testConnectionSync(connectionString)
	if !testResult {
		form.StatusText.SetText("Connection test failed. Please check your settings.").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
		return
	}

	// If test passes, proceed with save
	form.StatusText.SetText("Connection test passed. Saving...").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.SuccessColor))
	App.Draw()

	// Validate connection string only if it's not empty
	if connectionString != "" {
		_, err := helpers.ParseConnectionString(connectionString)
		if err != nil {
			form.StatusText.SetText("Warning: " + err.Error() + " (saved anyway)").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.WarningColor))
		}
	}

	databases := app.App.Connections()
	newDatabases := make([]models.Connection, len(databases))

	parsedDatabaseData := models.Connection{
		Name:      connectionName,
		Driver:    dbType,
		Hostname:  hostname,
		Port:      port,
		Username:  username,
		Password:  password,
		DBName:    database,
		DSN:       connectionString, // Keep for backward compatibility
		DsnCustom: dsnCustom,
		DsnAuto:   dsnAuto,
		DsnValue:  connectionString,
	}

	switch form.Action {
	case actionNewConnection:
		newDatabases = append(databases, parsedDatabaseData)
		err := app.App.SaveConnections(newDatabases)
		if err != nil {
			form.StatusText.SetText("Save failed: " + err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
			return
		}
		configPath := app.App.GetConfigFilePath()
		form.StatusText.SetText("Saved to: " + configPath).SetTextColor(app.Styles.TertiaryTextColor)

	case actionEditConnection:
		newDatabases = make([]models.Connection, len(databases))
		row, _ := connectionsTable.GetSelection()

		for i, database := range databases {
			if i == row {
				newDatabases[i] = parsedDatabaseData
			} else {
				newDatabases[i] = database
			}
		}

		err := app.App.SaveConnections(newDatabases)
		if err != nil {
			form.StatusText.SetText("Save failed: " + err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
			return
		}
		configPath := app.App.GetConfigFilePath()
		form.StatusText.SetText("Saved to: " + configPath).SetTextColor(app.Styles.TertiaryTextColor)
	}

	connectionsTable.SetConnections(newDatabases)

	// Now connect to the database
	form.StatusText.SetText("Connecting to database...").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.WarningColor))
	App.Draw()

	// Create database driver
	var dbDriver drivers.Driver
	switch dbType {
	case drivers.DriverMySQL:
		dbDriver = &drivers.MySQL{}
	case drivers.DriverPostgres:
		dbDriver = &drivers.Postgres{}
	case drivers.DriverSqlite:
		dbDriver = &drivers.SQLite{}
	case drivers.DriverMSSQL:
		dbDriver = &drivers.MSSQL{}
	default:
		form.StatusText.SetText("Unsupported database type: " + dbType).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
		return
	}

	// Connect to database
	err := dbDriver.Connect(connectionString)
	if err != nil {
		form.StatusText.SetText("Connection failed: " + err.Error()).SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.ErrorColor))
		return
	}

	// Success - navigate to home page
	form.StatusText.SetText("Connected successfully!").SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.SuccessColor))
	App.Draw()

	// Create home page and navigate to it
	newHome := NewHomePage(parsedDatabaseData, dbDriver)
	newHome.Tree.SetCurrentNode(newHome.Tree.GetRoot())
	newHome.Tree.Wrapper.SetTitle(parsedDatabaseData.Name)

	// Add page to main pages and switch to it
	mainPages.AddAndSwitchToPage(parsedDatabaseData.Name, newHome, true)
	App.SetFocus(newHome.Tree)
	App.Draw()
}
