package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rivo/tview"
)

type MovieAdd struct {
	Title      string
	Year       int
	Rating     *int
	ImdbRating *float32
	Genres     []string
	Directors  []string
}

func (t *TUI) addMovieHandler() {
	t.AddForm.GetFormItem(0).(*tview.InputField).SetText("")
	t.AddForm.GetFormItem(1).(*tview.InputField).SetText("")
	t.AddForm.GetFormItem(2).(*tview.InputField).SetText("")
	t.AddForm.GetFormItem(3).(*tview.InputField).SetText("")
	t.AddForm.GetFormItem(4).(*tview.InputField).SetText("")
	t.AddForm.GetFormItem(5).(*tview.InputField).SetText("")
	t.AddForm.GetButton(0).SetLabel("Add")
	t.AddForm.GetButton(0).SetDisabled(true)
	t.AddForm.GetButton(0).SetSelectedFunc(t.addMovieButton)
	t.Pages.ShowPage("add")
	t.App.SetFocus(t.AddForm.GetFormItem(0))
}

func (t *TUI) addMovieButton() {
	movie := MovieAdd{}

	movie.Title = t.AddForm.GetFormItem(0).(*tview.InputField).GetText()
	year, err := strconv.Atoi(t.AddForm.GetFormItem(1).(*tview.InputField).GetText())
	if err != nil {
		panic(err)
	}
	movie.Year = year
	rating, err := strconv.Atoi(t.AddForm.GetFormItem(2).(*tview.InputField).GetText())
	if err == nil {
		movie.Rating = &rating
	}
	imdbRating, err := strconv.ParseFloat(t.AddForm.GetFormItem(3).(*tview.InputField).GetText(), 32)
	if err == nil {
		imdbRatingf32 := float32(imdbRating)
		movie.ImdbRating = &imdbRatingf32
	}
	genreSplit := strings.Split(t.AddForm.GetFormItem(4).(*tview.InputField).GetText(), ",")
	for _, str := range genreSplit {
		str = strings.TrimSpace(str)
		if str != "" {
			movie.Genres = append(movie.Genres, str)
		}
	}
	directorSplit := strings.Split(t.AddForm.GetFormItem(5).(*tview.InputField).GetText(), ",")
	for _, str := range directorSplit {
		str = strings.TrimSpace(str)
		if str != "" {
			movie.Directors = append(movie.Directors, str)
		}
	}

	jsonData, err := json.Marshal(&movie)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/movies", t.Port), bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusConflict {
		warningText := fmt.Sprintf("There already is a movie named %s (%d) in the database. Do you want to add another one?", movie.Title, movie.Year)
		t.WarningText.SetText(warningText)
		t.WarningOkButton.SetLabel("Yes")
		t.WarningNoButton.SetLabel("No")
		t.WarningOkButton.SetSelectedFunc(func() {
			req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/movies/force", t.Port), bytes.NewBuffer(jsonData))
			if err != nil {
				panic(err)
			}
			client.Do(req)
			t.Pages.HidePage("warning")
			t.Table.Clear()
			if err := t.fillTable(t.Table); err != nil {
				panic(err)
			}
		})
		t.WarningNoButton.SetSelectedFunc(func() {
			t.Pages.HidePage("warning")
			t.Pages.ShowPage("add")
		})
		t.Pages.ShowPage("warning")
		t.Pages.HidePage("add")
	} else if resp.StatusCode != http.StatusCreated {
		panic("Connection error with server")
	}

	t.Table.Clear()
	if err := t.fillTable(t.Table); err != nil {
		panic(err)
	}
	t.Pages.HidePage("add")
}

func (t *TUI) checkAddButton(text string) {
	title := t.AddForm.GetFormItem(0).(*tview.InputField).GetText()
	year := t.AddForm.GetFormItem(1).(*tview.InputField).GetText()
	button := t.AddForm.GetButton(0)
	if title != "" && year != "" {
		button.SetDisabled(false)
	} else {
		button.SetDisabled(true)
	}
}

func (t *TUI) setupAddForm() {
	t.AddButton.SetSelectedFunc(t.addMovieHandler)
	t.AddForm.
		AddInputField("Title: ", "", 30, nil, t.checkAddButton).
		AddInputField("Year: ", "", 10, func(textToCheck string, lastChar rune) bool {
			_, err := strconv.Atoi(textToCheck)
			if err != nil {
				return false
			}
			return true
		}, t.checkAddButton).
		AddInputField("Your Rating: ", "", 4, func(textToCheck string, lastChar rune) bool {
			val, err := strconv.Atoi(textToCheck)
			if err != nil {
				return false
			} else if val < 1 || val > 10 {
				return false
			}
			return true
		}, nil).
		AddInputField("IMDB Rating: ", "", 4, func(textToCheck string, lastChar rune) bool {
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
		AddInputField("Genres: ", "", 30, nil, nil).
		AddInputField("Directors: ", "", 30, nil, nil).
		AddTextView("", "For adding multiple genres or directors, seperate each value with a comma ','", 40, 3, false, false).
		AddButton("Add", t.addMovieButton)
	t.AddForm.GetButton(0).SetDisabled(true)
}
