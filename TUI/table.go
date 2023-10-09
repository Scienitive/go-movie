package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *TUI) fillTable(table *tview.Table) error {
	initialMovieCount := 10000000
	movies, err := t.getMovies(initialMovieCount, 0)
	if err != nil {
		return err
	}
	t.Movies = movies

	for row := 0; row < len(movies)+1; row++ {
		for col := 0; col < 6; col++ {
			color := tcell.ColorWhite
			if row == 0 {
				color = tcell.ColorYellow
			}
			align := tview.AlignLeft
			if row == 0 || col == 2 || col == 3 {
				align = tview.AlignCenter
			} else if col == 0 {
				align = tview.AlignRight
			}
			bgColor := tcell.NewRGBColor(80, 80, 80)
			if row == 0 {
				bgColor = tcell.NewRGBColor(40, 40, 40)
			} else if row%2 == 1 {
				bgColor = tcell.NewRGBColor(60, 60, 60)
			}
			table.SetCell(
				row,
				col,
				&tview.TableCell{
					Text:            t.textPlacer(movies, row, col),
					Color:           color,
					BackgroundColor: bgColor,
					Align:           align,
					Expansion:       5,
					NotSelectable:   row == 0,
				},
			)
		}
	}
	return nil
}

func (t *TUI) textPlacer(movies []Movie, row int, col int) string {
	if row == 0 {
		switch col {
		case 0:
			if t.MoviesOrder == 2 {
				return "[Q\r] Year ↓"
			} else if t.MoviesOrder == 3 {
				return "[Q\r] Year ↑"
			} else {
				return "[Q\r] Year"
			}
		case 1:
			if t.MoviesOrder == 0 {
				return "[W\r] Title ↓"
			} else if t.MoviesOrder == 1 {
				return "[W\r] Title ↑"
			} else {
				return "[W\r] Title"
			}
		case 2:
			if t.MoviesOrder == 4 {
				return "[E\r] Your Rating ↓"
			} else if t.MoviesOrder == 5 {
				return "[E\r] Your Rating ↑"
			} else {
				return "[E\r] Your Rating"
			}
		case 3:
			if t.MoviesOrder == 6 {
				return "[R\r] IMDB Rating ↓"
			} else if t.MoviesOrder == 7 {
				return "[R\r] IMDB Rating ↑"
			} else {
				return "[R\r] IMDB Rating"
			}
		case 4:
			if t.MoviesOrder == 12 {
				return "[T\r] Directors ↓"
			} else if t.MoviesOrder == 13 {
				return "[T\r] Directors ↑"
			} else {
				return "[T\r] Diretors"
			}
		case 5:
			if t.MoviesOrder == 10 {
				return "[Y\r] Genres ↓"
			} else if t.MoviesOrder == 11 {
				return "[Y\r] Genres ↑"
			} else {
				return "[Y\r] Genres"
			}
		}
	} else {
		switch col {
		case 0:
			return strconv.Itoa(movies[row-1].Year)
		case 1:
			return movies[row-1].Title
		case 2:
			if movies[row-1].Rating != nil {
				return strconv.Itoa(*movies[row-1].Rating)
			} else {
				return ""
			}
		case 3:
			if movies[row-1].ImdbRating != nil {
				return strconv.FormatFloat(float64(*movies[row-1].ImdbRating), 'f', -1, 32)
			} else {
				return ""
			}
		case 4:
			return strings.Join(movies[row-1].Directors, ", ")
		case 5:
			return strings.Join(movies[row-1].Genres, ", ")
		}
	}
	return ""
}

func (t *TUI) getMovies(limit int, skip int) ([]Movie, error) {
	theURL := fmt.Sprintf("http://localhost:%d/movies", t.Port)
	queryParams := url.Values{}

	queryParams.Add("limit", strconv.Itoa(limit))
	queryParams.Add("skip", strconv.Itoa(skip))
	queryParams.Add("order", strconv.Itoa(t.MoviesOrder))

	// Filter queryParams
	queryParams.Add("title", t.FilterTitle)
	queryParams.Add("yearMax", t.FilterYearMax)
	queryParams.Add("yearMin", t.FilterYearMin)
	queryParams.Add("ratingMax", t.FilterRatingMax)
	queryParams.Add("ratingMin", t.FilterRatingMin)
	queryParams.Add("imdbMax", t.FilterImdbRatingMax)
	queryParams.Add("imdbMin", t.FilterImdbRatingMin)
	queryParams.Add("directors", t.FilterDirectors)
	queryParams.Add("genres", t.FilterGenres)

	finalURL := theURL + "?" + queryParams.Encode()
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Cannot communicate with server.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	movies := []Movie{}
	err = json.Unmarshal(body, &movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}
