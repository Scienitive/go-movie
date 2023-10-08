package main

func (t *TUI) orderFilter(c rune) {
	switch c {
	case 'w':
		if t.MoviesOrder == 0 {
			t.MoviesOrder = 1
		} else {
			t.MoviesOrder = 0
		}
	case 'q':
		if t.MoviesOrder == 2 {
			t.MoviesOrder = 3
		} else {
			t.MoviesOrder = 2
		}
	case 'e':
		if t.MoviesOrder == 4 {
			t.MoviesOrder = 5
		} else {
			t.MoviesOrder = 4
		}
	case 'r':
		if t.MoviesOrder == 6 {
			t.MoviesOrder = 7
		} else {
			t.MoviesOrder = 6
		}
	case 'y':
		if t.MoviesOrder == 10 {
			t.MoviesOrder = 11
		} else {
			t.MoviesOrder = 10
		}
	case 't':
		if t.MoviesOrder == 12 {
			t.MoviesOrder = 13
		} else {
			t.MoviesOrder = 12
		}
	}

	t.Table.Clear()
	t.fillTable(t.Table)
}
