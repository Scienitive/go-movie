package main

import (
	"github.com/rivo/tview"
)

func (t *TUI) filterMovieButton() {
	t.FilterTitle = t.FilterForm.GetFormItem(0).(*tview.InputField).GetText()
	t.FilterYearMax = t.FilterForm.GetFormItem(1).(*tview.InputField).GetText()
	t.FilterYearMin = t.FilterForm.GetFormItem(2).(*tview.InputField).GetText()
	t.FilterRatingMax = t.FilterForm.GetFormItem(3).(*tview.InputField).GetText()
	t.FilterRatingMin = t.FilterForm.GetFormItem(4).(*tview.InputField).GetText()
	t.FilterImdbRatingMax = t.FilterForm.GetFormItem(5).(*tview.InputField).GetText()
	t.FilterImdbRatingMin = t.FilterForm.GetFormItem(6).(*tview.InputField).GetText()
	t.FilterDirectors = t.FilterForm.GetFormItem(7).(*tview.InputField).GetText()
	t.FilterGenres = t.FilterForm.GetFormItem(8).(*tview.InputField).GetText()

	t.Table.Clear()
	t.fillTable(t.Table)
	t.Pages.HidePage("filter")
}
