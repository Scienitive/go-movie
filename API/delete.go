package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (app *App) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Request Type Check
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	// Setting up a transaction for multiple statements
	tx, err := app.db.Begin()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	// Parse the URL to get the ID
	idStr := r.URL.Path[len("/movies/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusNotFound)
		tx.Rollback()
		return
	}

	// Deleting the Movie
	res, err := tx.Exec(
		`DELETE FROM movies
		WHERE id = ?`,
		id,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		tx.Rollback()
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		tx.Rollback()
		return
	}
	if count <= 0 {
		http.Error(w, "Invalid movie ID", http.StatusNotFound)
		tx.Rollback()
		return
	}

	// Commit all changes
	if err := tx.Commit(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "Movie successfully deleted.")
}
