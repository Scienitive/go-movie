package main

import (
	"flag"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	App             *tview.Application
	Pages           *tview.Pages
	MainGrid        *tview.Grid
	AddGrid         *tview.Grid
	FilterGrid      *tview.Grid
	TopGrid         *tview.Grid
	BottomGrid      *tview.Grid
	WarningGrid     *tview.Grid
	HeaderText      *tview.TextView
	Table           *tview.Table
	AddForm         *tview.Form
	FilterForm      *tview.Form
	AddButton       *tview.Button
	FilterButton    *tview.Button
	WarningText     *tview.TextView
	WarningOkButton *tview.Button
	WarningNoButton *tview.Button

	Port        int
	Movies      []Movie
	EditMovieID int
	MoviesOrder int

	FilterTitle         string
	FilterYearMax       string
	FilterYearMin       string
	FilterRatingMax     string
	FilterRatingMin     string
	FilterImdbRatingMax string
	FilterImdbRatingMin string
	FilterDirectors     string
	FilterGenres        string
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

func main() {
	fPort := flag.Int("port", 8080, "Port of the API")
	flag.Parse()
	t := initializeTUI(*fPort)

	// Setup elements
	t.Table.SetFixed(1, 0).SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(140, 140, 140))).
		SetEvaluateAllRows(true)
	t.Table.SetSelectedFunc(t.editMovieHandler)
	t.setupAddForm()
	t.setupFilterForm()

	// Layouts
	modalWidth := 60
	modalHeight := 30
	warningWidth := 40
	warningHeight := 4

	t.TopGrid.SetRows(0, 0, 0).SetColumns(0).
		AddItem(t.HeaderText, 1, 0, 1, 1, 0, 0, false)

	t.BottomGrid.SetRows(0, 0, 0).SetColumns(0, 0, 0, 0, 0).
		AddItem(t.AddButton, 1, 1, 1, 1, 0, 0, true).
		AddItem(t.FilterButton, 1, 3, 1, 1, 0, 0, false)

	t.MainGrid.SetRows(3, 0, 6).SetColumns(0).SetBorders(false).
		AddItem(t.TopGrid, 0, 0, 1, 1, 0, 0, false).
		AddItem(t.Table, 1, 0, 1, 1, 0, 0, true).
		AddItem(t.BottomGrid, 2, 0, 1, 1, 0, 0, false)

	t.AddGrid.SetColumns(0, modalWidth, 0).SetRows(0, modalHeight, 0).
		AddItem(t.AddForm, 1, 1, 1, 1, 0, 0, true)

	t.FilterGrid.SetColumns(0, modalWidth, 0).SetRows(0, modalHeight, 0).
		AddItem(t.FilterForm, 1, 1, 1, 1, 0, 0, true)

	t.WarningGrid.SetColumns(0, warningWidth, 0).SetRows(0, warningHeight, 0).
		AddItem(tview.NewGrid().SetRows(3, 1).SetColumns(0, 0).
			AddItem(t.WarningText, 0, 0, 1, 2, 0, 0, false).
			AddItem(t.WarningOkButton, 1, 0, 1, 1, 0, 0, true).
			AddItem(t.WarningNoButton, 1, 1, 1, 1, 0, 0, false),
			1, 1, 1, 1, 0, 0, true)

	// Configure apperances
	t.Table.SetTitle(" Table [Ctrl-K] ").SetBorder(true)
	t.BottomGrid.SetTitle(" Buttons [Ctrl-J] ").SetBorder(true)
	t.HeaderText.SetText("[Q, W, E, R, T, Y] For Sorting | [D] For Deleting | [Ctrl-D] and [Ctrl-U] For Fast Jumping | [ESC] For Exiting Popups | [Ctrl-C] For Exiting the App").SetTextAlign(tview.AlignCenter)
	t.HeaderText.SetTextAlign(tview.AlignCenter)
	t.AddForm.SetBorder(true)
	t.FilterForm.SetBorder(true)

	// Set Pages
	t.Pages.
		AddPage("main", t.MainGrid, true, true).
		AddPage("add", t.AddGrid, true, false).
		AddPage("filter", t.FilterGrid, true, false).
		AddPage("warning", t.WarningGrid, true, false)

	// Last Setups
	if err := t.fillTable(t.Table); err != nil {
		panic(err)
	}
	t.setKeyBindings()

	if err := t.App.SetRoot(t.Pages, true).SetFocus(t.Pages).Run(); err != nil {
		panic(err)
	}
}

func initializeTUI(port int) TUI {
	t := TUI{}
	t.App = tview.NewApplication()
	t.Pages = tview.NewPages()
	t.MainGrid = tview.NewGrid()
	t.AddGrid = tview.NewGrid()
	t.FilterGrid = tview.NewGrid()
	t.WarningGrid = tview.NewGrid()
	t.TopGrid = tview.NewGrid()
	t.BottomGrid = tview.NewGrid()
	t.HeaderText = tview.NewTextView()
	t.Table = tview.NewTable()
	t.AddForm = tview.NewForm()
	t.FilterForm = tview.NewForm()
	t.AddButton = tview.NewButton("Add Movie")
	t.FilterButton = tview.NewButton("Filter")
	t.WarningText = tview.NewTextView().SetTextAlign(tview.AlignCenter)
	t.WarningOkButton = tview.NewButton("")
	t.WarningNoButton = tview.NewButton("")
	t.Port = port

	return t
}
