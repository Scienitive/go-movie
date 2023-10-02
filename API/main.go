package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	app := initializeDatabase()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/create", app.insertMovieHandler(false))
	r.Post("/force/create", app.insertMovieHandler(true))
	r.Get("/get/{id}", app.getMovieByIdHandler)
	r.Get("/get", app.getMovieHandler)
	r.Put("/update/{id}", app.updateMovieHandler)
	r.Delete("/delete/{id}", app.deleteMovieHandler)

	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

func initializeDatabase() *App {
	db, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS movies
		(id INTEGER PRIMARY KEY, dateAdded TEXT, title TEXT, year INTEGER, rating INTEGER, imdbRating REAL)`,
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
