package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Schema of the return JSON
type MovieJSON struct {
	ID         int
	Date       string
	Title      string
	Year       int
	Rating     *int
	ImdbRating *float32
	Genres     []string
	Directors  []string
}

// Order Enum
const (
	OrderTitleDESC = iota
	OrderTitleASC
	OrderYearDESC
	OrderYearASC
	OrderRatingDESC
	OrderRatingASC
	OrderImdbRatingDESC
	OrderImdbRatingASC
	OrderDateDESC
	OrderDateASC
	OrderGenresDESC
	OrderGenresASC
	OrderDirectorsDESC
	OrderDirectorsASC
)

// Base GET Request
func (app *App) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Request Type Check
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the query parameters
	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 10
	}
	skip, err := strconv.Atoi(query.Get("skip"))
	if err != nil {
		skip = 0
	}
	order, err := strconv.Atoi(query.Get("order"))
	if err != nil {
		order = 0
	}

	// Prepeare the ORDER BY part of the SQL
	orderSQL := setOrder(order)

	// The SQL Query
	SQL := fmt.Sprintf(
		`SELECT
		movies.id, movies.dateAdded, movies.title, movies.year, movies.rating, movies.imdbRating,
		(
		SELECT GROUP_CONCAT(genres.name, ':')
		FROM movies_genres
		LEFT JOIN genres ON movies_genres.genreId = genres.id
		WHERE movies_genres.movieId = movies.id
		) AS genre_names,
		(
		SELECT GROUP_CONCAT(directors.name, ':')
		FROM movies_directors
		LEFT JOIN directors ON movies_directors.directorId = directors.id
		WHERE movies_directors.movieId = movies.id
		) AS director_names
		FROM movies
		ORDER BY %s
		LIMIT ?
		OFFSET ?`,
		orderSQL,
	)

	// Variables to hold the values from SQL Query
	var (
		id         int
		date       string
		title      string
		year       int
		rating     *int
		imdbRating *float32
		genres     *string
		directors  *string
	)

	// Executing the Query and stroring it in the movies slice
	rows, err := app.db.Query(SQL, limit, skip)
	movies := []MovieJSON{}
	for rows.Next() {
		err = rows.Scan(&id, &date, &title, &year, &rating, &imdbRating, &genres, &directors)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		movies = append(movies, MovieJSON{
			ID:         id,
			Title:      title,
			Year:       year,
			Date:       date,
			Rating:     rating,
			ImdbRating: imdbRating,
			Genres:     strPtrToSlice(genres),
			Directors:  strPtrToSlice(directors),
		})
	}

	// Turning slice into JSON and returning it
	json, err := json.Marshal(&movies)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintln(w, string(json))
}

func (app *App) getMovieByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Request Type Check
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the URL to get the ID
	idStr := r.URL.Path[len("/movies/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// The struct and variables for storing the SQL Query results
	movie := MovieJSON{}
	var (
		genres    *string
		directors *string
	)

	// The SQL Query and storing them
	row := app.db.QueryRow(
		`SELECT
		movies.id, movies.dateAdded, movies.title, movies.year, movies.rating, movies.imdbRating,
		(
		SELECT GROUP_CONCAT(genres.name, ':')
		FROM movies_genres
		LEFT JOIN genres ON movies_genres.genreId = genres.id
		WHERE movies_genres.movieId = movies.id
		) AS genre_names,
		(
		SELECT GROUP_CONCAT(directors.name, ':')
		FROM movies_directors
		LEFT JOIN directors ON movies_directors.directorId = directors.id
		WHERE movies_directors.movieId = movies.id
		) AS director_names
		FROM movies
		WHERE movies.id = ?`,
		id,
	)
	err = row.Scan(&movie.ID, &movie.Date, &movie.Title, &movie.Year, &movie.Rating, &movie.ImdbRating, &genres, &directors)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}
	movie.Genres = strPtrToSlice(genres)
	movie.Directors = strPtrToSlice(directors)

	// Turning slice into JSON and returning it
	json, err := json.Marshal(&movie)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintln(w, string(json))
}

func setOrder(o int) string {
	switch o {
	case OrderTitleASC:
		return "movies.title ASC NULLS LAST"
	case OrderYearDESC:
		return "movies.year DESC NULLS LAST"
	case OrderYearASC:
		return "movies.year ASC NULLS LAST"
	case OrderRatingDESC:
		return "movies.rating DESC NULLS LAST"
	case OrderRatingASC:
		return "movies.rating ASC NULLS LAST"
	case OrderImdbRatingDESC:
		return "movies.imdbRating DESC NULLS LAST"
	case OrderImdbRatingASC:
		return "movies.imdbRating ASC NULLS LAST"
	case OrderDateDESC:
		return "movies.dateAdded DESC NULLS LAST"
	case OrderDateASC:
		return "movies.dateAdded ASC NULLS LAST"
	case OrderGenresDESC:
		return "genre_names DESC NULLS LAST"
	case OrderGenresASC:
		return "genre_names ASC NULLS LAST"
	case OrderDirectorsDESC:
		return "director_names DESC NULLS LAST"
	case OrderDirectorsASC:
		return "director_names ASC NULLS LAST"
	default:
		return "movies.title DESC NULLS LAST"
	}
}

func strPtrToSlice(str *string) []string {
	var split []string
	if str != nil {
		split = strings.Split(*str, ":")
	} else {
		split = []string{}
	}
	return split
}
