package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"sqlcmder/cmd/app"
	"sqlcmder/models"
)

type ResultsTableFilter struct {
	*tview.Flex
	Input         *tview.InputField
	Label         *tview.TextView
	currentFilter string
	subscribers   []chan models.StateChange
}

func NewResultsFilter() *ResultsTableFilter {
	recordsFilter := &ResultsTableFilter{
		Flex:  tview.NewFlex(),
		Input: tview.NewInputField(),
		Label: tview.NewTextView(),
	}
	// recordsFilter.SetBorder(true)  // Remove border to save space
	recordsFilter.SetDirection(tview.FlexRowCSS)
	recordsFilter.SetTitleAlign(tview.AlignCenter)
	recordsFilter.SetBorderPadding(0, 0, 1, 1)

	recordsFilter.Label.SetTextColor(app.Styles.TertiaryTextColor)
	recordsFilter.Label.SetText("WHERE")
	recordsFilter.Label.SetBorderPadding(0, 0, 0, 1)

	// Set input field styling - unified background color for input and autocomplete dropdown
	fieldBgColor := app.Styles.PrimitiveBackgroundColor

	recordsFilter.Input.SetPlaceholder("  Press / to search")
	recordsFilter.Input.SetPlaceholderStyle(tcell.StyleDefault.Foreground(app.Styles.PrimaryTextColor).Background(fieldBgColor))
	recordsFilter.Input.SetFieldBackgroundColor(fieldBgColor)
	recordsFilter.Input.SetFieldTextColor(app.Styles.PrimaryTextColor)
	recordsFilter.Input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			if recordsFilter.Input.GetText() != "" {
				recordsFilter.currentFilter = "WHERE " + recordsFilter.Input.GetText()
				recordsFilter.Publish("WHERE " + recordsFilter.Input.GetText())

			}
		case tcell.KeyEscape:
			if recordsFilter.Input.GetText() == "" {
				recordsFilter.currentFilter = ""
				recordsFilter.Input.SetText("")
			}

			recordsFilter.Publish("")

		}
	})

	// Autocomplete dropdown uses same background as input field
	recordsFilter.Input.SetAutocompleteStyles(
		fieldBgColor,
		tcell.StyleDefault.Foreground(app.Styles.PrimaryTextColor).Background(fieldBgColor),
		tcell.StyleDefault.Foreground(app.Styles.SecondaryTextColor).Background(fieldBgColor),
	)

	recordsFilter.AddItem(recordsFilter.Label, 6, 0, false)
	recordsFilter.AddItem(recordsFilter.Input, 0, 1, false)

	return recordsFilter
}

func (filter *ResultsTableFilter) Subscribe() chan models.StateChange {
	subscriber := make(chan models.StateChange)
	filter.subscribers = append(filter.subscribers, subscriber)
	return subscriber
}

func (filter *ResultsTableFilter) Publish(message string) {
	for _, sub := range filter.subscribers {
		sub <- models.StateChange{
			Key:   eventResultsTableFiltering,
			Value: message,
		}
	}
}

func (filter *ResultsTableFilter) GetCurrentFilter() string {
	return filter.currentFilter
}

// Function to blur
func (filter *ResultsTableFilter) RemoveHighlight() {
	filter.SetBorderColor(app.Styles.UnfocusedBorderColor)
	filter.Label.SetTextColor(app.Styles.UnfocusedAccentColor)
	filter.Input.SetPlaceholderTextColor(app.Styles.UnfocusedTextColor)
	filter.Input.SetFieldTextColor(app.Styles.UnfocusedTextColor)
}

func (filter *ResultsTableFilter) RemoveLocalHighlight() {
	filter.SetBorderColor(app.Styles.UnfocusedBorderColor)
	filter.Label.SetTextColor(app.Styles.UnfocusedAccentColor)
	filter.Input.SetPlaceholderTextColor(app.Styles.UnfocusedTextColor)
	filter.Input.SetFieldTextColor(app.Styles.UnfocusedTextColor)
}

func (filter *ResultsTableFilter) Highlight() {
	filter.SetBorderColor(app.Styles.PrimaryTextColor)
	filter.Label.SetTextColor(app.Styles.TertiaryTextColor)
	filter.Input.SetPlaceholderTextColor(app.Styles.PrimaryTextColor)
	filter.Input.SetFieldTextColor(app.Styles.PrimaryTextColor)
}

func (filter *ResultsTableFilter) HighlightLocal() {
	filter.SetBorderColor(app.Styles.PrimaryTextColor)
	filter.Label.SetTextColor(app.Styles.TertiaryTextColor)
	filter.Input.SetPlaceholderTextColor(app.Styles.PrimaryTextColor)
	filter.Input.SetFieldTextColor(app.Styles.PrimaryTextColor)
}
