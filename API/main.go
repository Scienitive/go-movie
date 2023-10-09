package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db *sql.DB
}

type Movie struct {
	Title      *string
	Year       *int
	Rating     *int
	ImdbRating *float32
	Genres     []string
	Directors  []string
}

func main() {
	fPort := flag.Int("port", 8080, "Port of the API")
	flag.Parse()
	app := initializeDatabase()

	http.HandleFunc("/movies", app.moviesRouter)
	http.HandleFunc("/movies/", app.moviesIdRouter)
	http.HandleFunc("/movies/force", app.insertMovieHandler(true))

	fmt.Printf("Listening on port %d...\n", *fPort)
	http.ListenAndServe(fmt.Sprintf(":%d", *fPort), nil)
}

func initializeDatabase() *App {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS movies
		(id INTEGER PRIMARY KEY, dateAdded INTEGER, title TEXT, year INTEGER, rating INTEGER, imdbRating REAL)`,
	)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS genres
		(id INTEGER PRIMARY KEY, name TEXT UNIQUE)`,
	)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS movies_genres
		(movieId INTEGER, genreId INTERGER,
		FOREIGN KEY (movieId) REFERENCES movies(id) ON DELETE CASCADE,
		FOREIGN KEY (genreId) REFERENCES genres(id) ON DELETE CASCADE
		PRIMARY KEY (movieId, genreId))`,
	)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS directors
		(id INTEGER PRIMARY KEY, name TEXT UNIQUE)`,
	)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS movies_directors
		(movieId INTEGER, directorId INTERGER,
		FOREIGN KEY (movieId) REFERENCES movies(id) ON DELETE CASCADE,
		FOREIGN KEY (directorId) REFERENCES directors(id) ON DELETE CASCADE,
		PRIMARY KEY (movieId, directorId))`,
	)
	if err != nil {
		panic(err.Error())
	}

	app := &App{db: db}
	return app
}

func (app *App) moviesRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.insertMovieHandler(false)(w, r)
	case http.MethodGet:
		app.getMovieHandler(w, r)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func (app *App) moviesIdRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		app.updateMovieHandler(w, r)
	case http.MethodGet:
		app.getMovieByIdHandler(w, r)
	case http.MethodDelete:
		app.deleteMovieHandler(w, r)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}
