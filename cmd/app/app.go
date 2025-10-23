package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/config"
	"sqlcmder/models"
)

var (
	App    *Application
	Styles *Theme
)

type Application struct {
	*tview.Application

	config    *config.Config
	context   context.Context
	cancelFn  context.CancelFunc
	waitGroup sync.WaitGroup
}

type Theme struct {
	tview.Theme

	SidebarTitleBorderColor string
	ButtonBackgroundColor   tcell.Color
	ButtonTextColor         tcell.Color
	ButtonUnselectedBgColor tcell.Color
	UnfocusedBorderColor    tcell.Color
	UnfocusedTextColor      tcell.Color
	UnfocusedAccentColor    tcell.Color
	SelectedTextColor       tcell.Color
}

func init() {
	ctx, cancel := context.WithCancel(context.Background())

	App = &Application{
		Application: tview.NewApplication(),
		config:      config.DefaultConfig(),
		context:     ctx,
		cancelFn:    cancel,
	}

	App.register()
	App.EnableMouse(true)
	App.EnablePaste(true)

	// Initialize with default theme (will be overridden after config loads)
	initializeTheme(models.ThemeDark)
}

// initializeTheme sets up the color scheme based on the theme name
func initializeTheme(themeName string) {
	// Set the active color scheme
	models.SetActiveColorScheme(themeName)
	scheme := models.ActiveColorScheme

	Styles = &Theme{
		Theme: tview.Theme{
			PrimitiveBackgroundColor:    scheme.PrimitiveBackgroundColor,
			ContrastBackgroundColor:     scheme.ContrastBackgroundColor,
			MoreContrastBackgroundColor: scheme.MoreContrastBgColor,
			BorderColor:                 scheme.Border,
			TitleColor:                  scheme.TextColor,
			GraphicsColor:               scheme.GraphicsColor,
			PrimaryTextColor:            scheme.TextColor,
			SecondaryTextColor:          scheme.AccentYellow,
			TertiaryTextColor:           scheme.AccentGreen,
			InverseTextColor:            scheme.InputBg,
			ContrastSecondaryTextColor:  scheme.SelectedTextColor,
		},
		SidebarTitleBorderColor: "#666A7E",
		ButtonBackgroundColor:   scheme.ButtonBg,
		ButtonTextColor:         scheme.ButtonTextColor,
		ButtonUnselectedBgColor: scheme.ButtonUnselectedBgColor,
		UnfocusedBorderColor:    scheme.UnfocusedBorder,
		UnfocusedTextColor:      scheme.UnfocusedText,
		UnfocusedAccentColor:    scheme.UnfocusedAccent,
		SelectedTextColor:       scheme.SelectedTextColor,
	}

	tview.Styles = Styles.Theme
}

// ApplyTheme applies the theme from the configuration
func (a *Application) ApplyTheme() {
	themeName := a.config.AppConfig.Theme
	if themeName == "" {
		themeName = models.ThemeDark
	}
	initializeTheme(themeName)
}

// Context returns the application context.
func (a *Application) Context() context.Context {
	return a.context
}

// Config returns the application configuration.
func (a *Application) Config() *models.AppConfig {
	return a.config.AppConfig
}

// GetConfig returns the full configuration object.
func (a *Application) GetConfig() *config.Config {
	return a.config
}

// Connections returns the database connections.
func (a *Application) Connections() []models.Connection {
	return a.config.Connections
}

// SaveConnections saves the database connections.
func (a *Application) SaveConnections(connections []models.Connection) error {
	return a.config.SaveConnections(connections)
}

// GetConfigFilePath returns the configuration file path.
func (a *Application) GetConfigFilePath() string {
	return a.config.ConfigFile
}

// Register adds a task to the wait group and returns a
// function that decrements the task count when called.
//
// The application will not stop until all registered tasks
// have finished by calling the returned function!
func (a *Application) Register() func() {
	a.waitGroup.Add(1)
	return a.waitGroup.Done
}

// Run starts and blocks until the application is stopped.
func (a *Application) Run(root *tview.Pages, configFile string) error {
	a.SetRoot(root, true)
	a.config.ConfigFile = configFile
	return a.Application.Run()
}

// Stop cancels the application context, waits for all
// tasks to finish, and then stops the application.
func (a *Application) Stop() {
	a.cancelFn()
	a.waitGroup.Wait()
	a.Application.Stop()
}

// register listens for interrupt and termination signals to
// gracefully handle shutdowns by calling the Stop method.
func (a *Application) register() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		a.Stop()
		<-c
		os.Exit(1)
	}()

	// Override the default input capture to listen for Ctrl+C
	// and make it send an interrupt signal to the channel to
	// trigger a graceful shutdown instead of closing the app
	// immediately without waiting for tasks to finish.
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			c <- os.Interrupt
			return nil
		}
		return event
	})
}
