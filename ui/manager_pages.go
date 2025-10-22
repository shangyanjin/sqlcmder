package ui

import (
	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
)

var mainPages *tview.Pages

func MainPages() *tview.Pages {
	mainPages = tview.NewPages()
	mainPages.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)
	mainPages.AddPage(pageNameConnections, NewConnectionPages().Grid, true, true)
	return mainPages
}
