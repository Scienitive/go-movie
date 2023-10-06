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

type TUI struct {
	App          *tview.Application
	Pages        *tview.Pages
	MainGrid     *tview.Grid
	AddGrid      *tview.Grid
	FilterGrid   *tview.Grid
	BottomGrid   *tview.Grid
	HeaderText   *tview.TextView
	Table        *tview.Table
	AddForm      *tview.Form
	FilterForm   *tview.Form
	AddButton    *tview.Button
	FilterButton *tview.Button
}

type Movie struct {
	ID         int
	Date       int
	Title      string
	Year       int
	Rating     *int
	ImdbRating *float32
	Genres     []string
	Directors  []string
}

func initializeTUI() TUI {
	t := TUI{}
	t.App = tview.NewApplication()
	t.Pages = tview.NewPages()
	t.MainGrid = tview.NewGrid()
	t.AddGrid = tview.NewGrid()
	t.FilterGrid = tview.NewGrid()
	t.BottomGrid = tview.NewGrid()
	t.HeaderText = tview.NewTextView()
	t.Table = tview.NewTable()
	t.AddForm = tview.NewForm()
	t.FilterForm = tview.NewForm()
	t.AddButton = tview.NewButton("Add Movie")
	t.FilterButton = tview.NewButton("Filter")

	return t
}

func main() {
	t := initializeTUI()

	// Setup elements
	t.Table.SetFixed(1, 0).SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(140, 140, 140))).
		SetEvaluateAllRows(true)
	t.HeaderText.SetTextAlign(tview.AlignCenter)
	t.AddButton.SetSelectedFunc(func() {
		// fmt.Println("ASDADASDADADASDAS")
		t.Pages.ShowPage("add")
	})
	t.AddForm.
		AddInputField("Title: ", "", 20, nil, nil).
		AddInputField("Year: ", "", 10, func(textToCheck string, lastChar rune) bool {
			_, err := strconv.Atoi(textToCheck)
			if err != nil {
				return false
			}
			return true
		}, nil)

	// Layouts
	modalWidth := 40
	modalHeight := 40

	t.BottomGrid.SetRows(0, 0, 0).SetColumns(0, 0, 0, 0, 0).
		AddItem(t.AddButton, 1, 1, 1, 1, 0, 0, true).
		AddItem(t.FilterButton, 1, 3, 1, 1, 0, 0, false)

	t.MainGrid.SetRows(3, 0, 6).SetColumns(0).SetBorders(false).
		AddItem(t.HeaderText, 0, 0, 1, 1, 0, 0, false).
		AddItem(t.Table, 1, 0, 1, 1, 0, 0, true).
		AddItem(t.BottomGrid, 2, 0, 1, 1, 0, 0, false)

	t.AddGrid.SetColumns(0, modalWidth, 0).SetRows(0, modalHeight, 0).
		AddItem(t.AddForm, 1, 1, 1, 1, 0, 0, true)

	t.FilterGrid.SetColumns(0, modalWidth, 0).SetRows(0, modalHeight, 0).
		AddItem(t.FilterForm, 1, 1, 1, 1, 0, 0, true)

	// Configure apperances
	tableTitle := fmt.Sprintf(" Table [ Ctrl-K ] ")
	bottomTitle := fmt.Sprintf(" Buttons [ Ctrl-J ] ")
	t.Table.SetTitle(tableTitle).SetBorder(true)
	t.BottomGrid.SetTitle(bottomTitle).SetBorder(true)
	t.HeaderText.SetLabel("ZAZAZ").SetText("ASDASD")

	// Set Pages
	t.Pages.
		AddPage("main", t.MainGrid, true, true).
		AddPage("add", t.AddGrid, true, false).
		AddPage("filter", t.FilterGrid, true, false)

	fillTable(t.Table)
	t.setKeyBindings()

	if err := t.App.SetRoot(t.Pages, true).SetFocus(t.Pages).Run(); err != nil {
		panic(err)
	}
}

func (t *TUI) setKeyBindings() {
	// MainGrid Keybindings
	t.MainGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlJ:
			t.App.SetFocus(t.BottomGrid)
		case tcell.KeyCtrlK:
			t.App.SetFocus(t.Table)
		}

		return event
	})

	// BottomGrid Keybindings
	t.BottomGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			t.App.SetFocus(t.AddButton)
		case tcell.KeyRight:
			t.App.SetFocus(t.FilterButton)
		}

		switch event.Rune() {
		case 'h':
			t.App.SetFocus(t.AddButton)
		case 'l':
			t.App.SetFocus(t.FilterButton)
		}

		return event
	})

	// Modal Keybindings
	t.AddGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			t.AddForm.Clear(true)
			t.Pages.HidePage("add")
			t.App.SetFocus(t.Table)
		}

		return event
	})

	t.FilterGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			t.Pages.HidePage("filter")
		}

		return event
	})
}

func (t *TUI) addMovieHandler() {
	// app.Pages.ShowPage("modal")
	fmt.Println("SADASDSA")
}

func fillTable(table *tview.Table) error {
	initialMovieCount := 100000
	movies, err := getMovies(initialMovieCount, 0)
	if err != nil {
		return err
	}

	for row := 0; row < len(movies)+1; row++ {
		for col := 0; col < 7; col++ {
			color := tcell.ColorWhite
			if row == 0 {
				color = tcell.ColorYellow
			}
			align := tview.AlignLeft
			if row == 0 {
				align = tview.AlignCenter
			} else if col == 1 {
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
					Text:            textPlacer(movies, row, col),
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

func textPlacer(movies []Movie, row int, col int) string {
	if row == 0 {
		switch col {
		case 0:
			return "Date Added"
		case 1:
			return "Year"
		case 2:
			return "Title"
		case 3:
			return "Your Rating"
		case 4:
			return "IMDB Score"
		case 5:
			return "Directors"
		case 6:
			return "Genres"
		}
	} else {
		switch col {
		case 0:
			return time.Unix(int64(movies[row-1].Date), 0).Format(time.DateOnly)
		case 1:
			return strconv.Itoa(movies[row-1].Year)
		case 2:
			return movies[row-1].Title
		case 3:
			if movies[row-1].Rating != nil {
				return strconv.Itoa(*movies[row-1].Rating)
			} else {
				return ""
			}
		case 4:
			if movies[row-1].ImdbRating != nil {
				return strconv.FormatFloat(float64(*movies[row-1].ImdbRating), 'f', -1, 32)
			} else {
				return ""
			}
		case 5:
			return strings.Join(movies[row-1].Directors, ", ")
		case 6:
			return strings.Join(movies[row-1].Genres, ", ")
		}
	}
	return ""
}

func getMovies(limit int, skip int) ([]Movie, error) {
	theURL := "http://localhost:8080/movies"
	queryParams := url.Values{}

	queryParams.Add("limit", strconv.Itoa(limit))
	queryParams.Add("skip", strconv.Itoa(skip))

	finalURL := theURL + "?" + queryParams.Encode()
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 3 * time.Second,
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
