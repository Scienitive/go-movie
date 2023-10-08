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

func (t *TUI) editMovieHandler(row int, col int) {
	movie := t.Movies[row-1]

	t.EditMovieID = movie.ID
	t.AddForm.GetFormItem(0).(*tview.InputField).SetText(movie.Title)
	t.AddForm.GetFormItem(1).(*tview.InputField).SetText(strconv.Itoa(movie.Year))
	if movie.Rating != nil {
		t.AddForm.GetFormItem(2).(*tview.InputField).SetText(strconv.Itoa(*movie.Rating))
	}
	if movie.ImdbRating != nil {
		t.AddForm.GetFormItem(3).(*tview.InputField).SetText(strconv.FormatFloat(float64(*movie.ImdbRating), 'f', 1, 32))
	}
	t.AddForm.GetFormItem(4).(*tview.InputField).SetText(strings.Join(movie.Genres, ", "))
	t.AddForm.GetFormItem(5).(*tview.InputField).SetText(strings.Join(movie.Directors, ", "))
	t.AddForm.GetButton(0).SetLabel("Update")
	t.AddForm.GetButton(0).SetSelectedFunc(t.editMovieButton)
	t.AddForm.GetButton(0).SetDisabled(false)
	t.Pages.ShowPage("add")
	t.App.SetFocus(t.AddForm.GetFormItem(0))
}

func (t *TUI) editMovieButton() {
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

	theURL := fmt.Sprintf("http://localhost:8080/movies/%d", t.EditMovieID)
	req, err := http.NewRequest("PUT", theURL, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	client := http.Client{
		Timeout: time.Second / 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic("Connection error with server")
	}

	t.Table.Clear()
	t.fillTable(t.Table)
	t.Pages.HidePage("add")
}
