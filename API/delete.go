package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (app *App) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Setting up a transaction for multiple statements
	tx, err := app.db.Begin()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	// Parse the URL to get the ID
	idStr := r.URL.Path[len("/delete/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		tx.Rollback()
		return
	}

	// Deleting the Movie
	_, err = tx.Exec(
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

	// Commit all changes
	if err := tx.Commit(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	fmt.Fprintln(w, "Movie successfully deleted.")
}
