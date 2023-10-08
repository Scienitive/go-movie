package main

import (
	"strconv"

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

func (t *TUI) setupFilterForm() {
	t.FilterButton.SetSelectedFunc(func() {
		t.FilterForm.GetFormItem(0).(*tview.InputField).SetText(t.FilterTitle)
		t.FilterForm.GetFormItem(1).(*tview.InputField).SetText(t.FilterYearMax)
		t.FilterForm.GetFormItem(2).(*tview.InputField).SetText(t.FilterYearMin)
		t.FilterForm.GetFormItem(3).(*tview.InputField).SetText(t.FilterRatingMax)
		t.FilterForm.GetFormItem(4).(*tview.InputField).SetText(t.FilterRatingMin)
		t.FilterForm.GetFormItem(5).(*tview.InputField).SetText(t.FilterImdbRatingMax)
		t.FilterForm.GetFormItem(6).(*tview.InputField).SetText(t.FilterImdbRatingMin)
		t.FilterForm.GetFormItem(7).(*tview.InputField).SetText(t.FilterDirectors)
		t.FilterForm.GetFormItem(8).(*tview.InputField).SetText(t.FilterGenres)
		t.Pages.ShowPage("filter")
		t.App.SetFocus(t.FilterForm.GetFormItem(0))
	})
	t.FilterForm.
		AddInputField("Title: ", "", 30, nil, nil).
		AddInputField("Year (Max): ", "", 10, func(textToCheck string, lastChar rune) bool {
			_, err := strconv.Atoi(textToCheck)
			if err != nil {
				return false
			}
			return true
		}, nil).
		AddInputField("Year (Min): ", "", 10, func(textToCheck string, lastChar rune) bool {
			_, err := strconv.Atoi(textToCheck)
			if err != nil {
				return false
			}
			return true
		}, nil).
		AddInputField("Your Rating (Max): ", "", 4, func(textToCheck string, lastChar rune) bool {
			val, err := strconv.Atoi(textToCheck)
			if err != nil {
				return false
			} else if val < 1 || val > 10 {
				return false
			}
			return true
		}, nil).
		AddInputField("Your Rating (Min): ", "", 4, func(textToCheck string, lastChar rune) bool {
			val, err := strconv.Atoi(textToCheck)
			if err != nil {
				return false
			} else if val < 1 || val > 10 {
				return false
			}
			return true
		}, nil).
		AddInputField("IMDB Rating (Max): ", "", 4, func(textToCheck string, lastChar rune) bool {
			afterDecimal := false
			afterDecimalCount := 0
			for _, c := range textToCheck {
				if c == '.' {
					afterDecimal = true
				} else if afterDecimal {
					afterDecimalCount++
				}
				if afterDecimalCount > 1 {
					return false
				}
			}
			val, err := strconv.ParseFloat(textToCheck, 32)
			if err != nil {
				return false
			} else if val < 1 || val > 10 {
				return false
			}
			return true
		}, nil).
		AddInputField("IMDB Rating (Min): ", "", 4, func(textToCheck string, lastChar rune) bool {
			afterDecimal := false
			afterDecimalCount := 0
			for _, c := range textToCheck {
				if c == '.' {
					afterDecimal = true
				} else if afterDecimal {
					afterDecimalCount++
				}
				if afterDecimalCount > 1 {
					return false
				}
			}
			val, err := strconv.ParseFloat(textToCheck, 32)
			if err != nil {
				return false
			} else if val < 1 || val > 10 {
				return false
			}
			return true
		}, nil).
		AddInputField("Directors: ", "", 30, nil, nil).
		AddInputField("Genres: ", "", 30, nil, nil).
		AddTextView("", "For adding multiple filters for genres or directors, seperate each value with a comma ','", 40, 3, false, false).
		AddButton("Filter", t.filterMovieButton)
}
