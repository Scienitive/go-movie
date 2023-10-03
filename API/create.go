package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CreateConflictError struct {
	Message     string
	DuplicateID int
}

func (app *App) insertMovieHandler(forced bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request Type Check
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
			return
		}

		// Decoding JSON into Movie struct
		var movie Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil || movie.Title == nil || movie.Year == nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Setting up a transaction for multiple statements
		tx, err := app.db.Begin()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		// Multiple entry checking if it's not forced
		if !forced {
			rows, err := tx.Query(
				`SELECT id FROM movies WHERE title = ? AND year = ?`,
				movie.Title, movie.Year,
			)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			defer rows.Close()

			var idHolder int
			count := 0
			for rows.Next() {
				rows.Scan(&idHolder)
				count++
			}
			if count == 1 {
				e := CreateConflictError{
					Message:     "Duplicate entries\nUse '/movies/force' to create the entry",
					DuplicateID: idHolder,
				}
				json, err := json.Marshal(e)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					log.Println(err.Error())
					tx.Rollback()
					return
				}
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, string(json), http.StatusConflict)
				tx.Rollback()
				return
			} else if count > 1 {
				http.Error(w, "Duplicate entries\nUse '/movies/force' to create the entry", http.StatusConflict)
				tx.Rollback()
				return
			}
		}

		// Insert the movie
		res, err := tx.Exec(
			`INSERT INTO movies (dateAdded, title, year, rating, imdbRating)
			VALUES (?, ?, ?, ?, ?)`,
			time.Now().Unix(), movie.Title, movie.Year, movie.Rating, movie.ImdbRating,
		)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			tx.Rollback()
			return
		}
		movieId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			tx.Rollback()
			return
		}

		// Insert Genres
		for _, genre := range movie.Genres {
			_, err := tx.Exec(
				`INSERT OR IGNORE INTO genres (name)
				VALUES (?)`,
				genre,
			)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err.Error())
				tx.Rollback()
				return
			}
			var genreId int
			row := tx.QueryRow(
				`SELECT id FROM genres
				WHERE name = ?`,
				genre,
			)
			err = row.Scan(&genreId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err.Error())
				tx.Rollback()
				return
			}
			_, err = tx.Exec(
				`INSERT OR IGNORE INTO movies_genres (movieId, genreId)
				VALUES (?, ?)`,
				movieId, genreId,
			)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err.Error())
				tx.Rollback()
				return
			}
		}

		// Insert directors
		for _, director := range movie.Directors {
			_, err := tx.Exec(
				`INSERT OR IGNORE INTO directors (name)
				VALUES (?)`,
				director,
			)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err.Error())
				tx.Rollback()
				return
			}
			var directorId int
			row := tx.QueryRow(
				`SELECT id FROM directors
				WHERE name = ?`,
				director,
			)
			err = row.Scan(&directorId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err.Error())
				tx.Rollback()
				return
			}
			_, err = tx.Exec(
				`INSERT OR IGNORE INTO movies_directors (movieId, directorId)
				VALUES (?, ?)`,
				movieId, directorId,
			)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err.Error())
				tx.Rollback()
				return
			}
		}

		// Commit all changes
		if err := tx.Commit(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		fmt.Fprintln(w, "Movie successfully inserted.")
	}
}
