package main

import (
	"fmt"
	"net/http"
	"time"
)

func (t *TUI) deleteMovie(movie Movie) {
	warningText := fmt.Sprintf("You are about to delete %s (%d). Are you sure?", movie.Title, movie.Year)
	t.WarningText.SetText(warningText)
	t.WarningOkButton.SetLabel("Yes")
	t.WarningNoButton.SetLabel("No")
	t.WarningOkButton.SetSelectedFunc(func() {
		theURL := fmt.Sprintf("http://localhost:%d/movies/%d", t.Port, movie.ID)
		req, err := http.NewRequest("DELETE", theURL, nil)
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

		if resp.StatusCode != http.StatusNoContent {
			panic("Connection error with server")
		}

		t.Table.Clear()
		if err := t.fillTable(t.Table); err != nil {
			panic(err)
		}
		t.Pages.HidePage("warning")
	})
	t.WarningNoButton.SetSelectedFunc(func() {
		t.Pages.HidePage("warning")
	})
	t.Pages.ShowPage("warning")
}
