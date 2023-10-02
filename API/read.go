package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

func (app *App) getMovieHandler(w http.ResponseWriter, r *http.Request) {
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

	orderSQL := ""
	switch order {
	case OrderTitleASC:
		orderSQL = "movies.title ASC NULLS LAST"
	case OrderYearDESC:
		orderSQL = "movies.year DESC NULLS LAST"
	case OrderYearASC:
		orderSQL = "movies.year ASC NULLS LAST"
	case OrderRatingDESC:
		orderSQL = "movies.rating DESC NULLS LAST"
	case OrderRatingASC:
		orderSQL = "movies.rating ASC NULLS LAST"
	case OrderImdbRatingDESC:
		orderSQL = "movies.imdbRating DESC NULLS LAST"
	case OrderImdbRatingASC:
		orderSQL = "movies.imdbRating ASC NULLS LAST"
	case OrderDateDESC:
		orderSQL = "movies.dateAdded DESC NULLS LAST"
	case OrderDateASC:
		orderSQL = "movies.dateAdded ASC NULLS LAST"
	case OrderGenresDESC:
		orderSQL = "genre_names DESC NULLS LAST"
	case OrderGenresASC:
		orderSQL = "genre_names ASC NULLS LAST"
	case OrderDirectorsDESC:
		orderSQL = "director_names DESC NULLS LAST"
	case OrderDirectorsASC:
		orderSQL = "director_names ASC NULLS LAST"
	default:
		orderSQL = "movies.title DESC NULLS LAST"
	}
	fmt.Println(orderSQL)
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
	rows, err := app.db.Query(SQL, limit, skip)
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
	movies := []MovieJSON{}
	for rows.Next() {
		err = rows.Scan(&id, &date, &title, &year, &rating, &imdbRating, &genres, &directors)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		var genresSplit []string
		var directorsSplit []string
		if genres != nil {
			genresSplit = strings.Split(*genres, ":")
		} else {
			genresSplit = []string{}
		}
		if directors != nil {
			directorsSplit = strings.Split(*directors, ":")
		} else {
			directorsSplit = []string{}
		}
		movies = append(movies, MovieJSON{
			ID:         id,
			Title:      title,
			Year:       year,
			Date:       date,
			Rating:     rating,
			ImdbRating: imdbRating,
			Genres:     genresSplit,
			Directors:  directorsSplit,
		})
	}

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
	// Parse the URL to get the ID
	idStr := r.URL.Path[len("/get/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie := MovieJSON{}
	var (
		genresStr    string
		directorsStr string
	)
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
	err = row.Scan(&movie.ID, &movie.Date, &movie.Title, &movie.Year, &movie.Rating, &movie.ImdbRating, &genresStr, &directorsStr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	movie.Genres = strings.Split(genresStr, ":")
	movie.Directors = strings.Split(directorsStr, ":")
	json, err := json.Marshal(&movie)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintln(w, string(json))
}
