package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *TUI) setKeyBindings() {
	// Table Keybindings
	t.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q', 'w', 'e', 'r', 't', 'y':
			t.orderFilter(event.Rune())
		case 'd':
			row, _ := t.Table.GetSelection()
			if row <= len(t.Movies) {
				t.deleteMovie(t.Movies[row-1])
			}
		}

		switch event.Key() {
		case tcell.KeyCtrlD:
			if len(t.Movies) > 0 {
				row, _ := t.Table.GetSelection()
				if row+25 <= len(t.Movies) {
					t.Table.Select(row+25, 0)
				} else {
					t.Table.Select(len(t.Movies), 0)
				}
			}
		case tcell.KeyCtrlU:
			if len(t.Movies) > 0 {
				row, _ := t.Table.GetSelection()
				if row-25 >= 1 {
					t.Table.Select(row-25, 0)
				} else {
					t.Table.Select(1, 0)
				}
			}
		}

		return event
	})
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
			for i := 0; i < t.AddForm.GetFormItemCount()-2; i++ {
				t.AddForm.GetFormItem(i).(*tview.InputField).SetText("")
			}
			t.AddForm.GetButton(0).SetDisabled(true)
			t.Pages.HidePage("add")
			t.App.SetFocus(t.Table)
		case tcell.KeyCtrlJ, tcell.KeyDown:
			i, b := t.AddForm.GetFocusedItemIndex()
			switch i {
			case 0, 1, 2, 3, 4:
				t.App.SetFocus(t.AddForm.GetFormItem(i + 1))
			case 5:
				button := t.AddForm.GetButton(0)
				if button.IsDisabled() {
					t.App.SetFocus(t.AddForm.GetFormItem(0))
				} else {
					t.App.SetFocus(button)
				}
			}
			if b != -1 {
				t.App.SetFocus(t.AddForm.GetFormItem(0))
			}
		case tcell.KeyCtrlK, tcell.KeyUp:
			i, b := t.AddForm.GetFocusedItemIndex()
			switch i {
			case 1, 2, 3, 4, 5:
				t.App.SetFocus(t.AddForm.GetFormItem(i - 1))
			case 0:
				button := t.AddForm.GetButton(0)
				if button.IsDisabled() {
					t.App.SetFocus(t.AddForm.GetFormItem(5))
				} else {
					t.App.SetFocus(button)
				}
			}
			if b != -1 {
				t.App.SetFocus(t.AddForm.GetFormItem(5))
			}
		}

		return event
	})

	t.FilterGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			t.Pages.HidePage("filter")
		case tcell.KeyCtrlJ, tcell.KeyDown:
			i, b := t.FilterForm.GetFocusedItemIndex()
			switch i {
			case 0, 1, 2, 3, 4, 5, 6, 7:
				t.App.SetFocus(t.FilterForm.GetFormItem(i + 1))
			case 8:
				t.App.SetFocus(t.FilterForm.GetButton(0))
			}
			if b != -1 {
				t.App.SetFocus(t.FilterForm.GetFormItem(0))
			}
		case tcell.KeyCtrlK, tcell.KeyUp:
			i, b := t.FilterForm.GetFocusedItemIndex()
			switch i {
			case 1, 2, 3, 4, 5, 6, 7, 8:
				t.App.SetFocus(t.FilterForm.GetFormItem(i - 1))
			case 0:
				t.App.SetFocus(t.FilterForm.GetButton(0))
			}
			if b != -1 {
				t.App.SetFocus(t.FilterForm.GetFormItem(8))
			}
		}

		return event
	})

	t.WarningGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			t.App.SetFocus(t.WarningOkButton)
		case tcell.KeyRight:
			t.App.SetFocus(t.WarningNoButton)
		}

		switch event.Rune() {
		case 'h':
			t.App.SetFocus(t.WarningOkButton)
		case 'l':
			t.App.SetFocus(t.WarningNoButton)
		}

		return event
	})
}
