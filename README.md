# moviterm

A movie rating app on terminal.

I made this project to learn more about Go programming language.

https://github.com/Scienitive/moviterm/assets/58879237/6098ddb4-e284-4b2d-b188-54e4105f1f02

## Usage

**Dependencies:**
- Go 1.21 or Later
- SQLite 3.30 or Later

```
To run API: 'make api' or if you want to specify a port 'make api PORT=xxxx'

To run TUI: 'make tui' or if you want to specify a port 'make tui PORT=xxxx'

To run CLI for IMDB: 'make cli-imdb FILE=(file path to ratings.csv)' or if you want to specify a port 'make cli-imdb PORT=xxxx FILE=(file path to ratings.csv)'

To run CLI for Letterboxd: 'make cli-letterboxd FILE=(file path to ratings.csv)' or if you want to specify a port 'make cli-letterboxd PORT=xxxx FILE=(file path to ratings.csv)'
```

**NOTE:** TUI and CLI won't work if the API is not running.

**NOTE:** Default PORT number is 8080.

## API

A CRUD API for movies. It uses SQLite as the database.

### Routes

#### `/movies` (POST)

Use this route to insert new movies.

```
JSON BODY:

{
  "title":       ... (Mandatory)
  "year":        ... (Mandatory)
  "rating":      ... (Between 1 and 10 [both including])
  "imdbRating":  ... (Between 1 and 10 [both including])
  "genres":      [... , ...]
  "directors":   [... , ...]
}
```

#### `/movies/force` (POST)

Use this route when you really want to insert another movie with same name and year information.

```
JSON BODY:

{
  "title":       ... (Mandatory)
  "year":        ... (Mandatory)
  "rating":      ... (Between 1 and 10 [both including])
  "imdbRating":  ... (Between 1 and 10 [both including])
  "genres":      [... , ...]
  "directors":   [... , ...]
}
```

#### `/movies` (GET)

Use this route to get all movies.

```
QUERY PARAMETERS:

limit: For limiting the number of movies that you get (Default is 10)
skip: For skipping the first x number of movies (Default is 0)
order: For changing the order of movies (Default is 0) (You need to provide a number between 0 and 13 [both including])

title: For filtering based on title
yearMax: For filtering based on movie release years
yearMin: For filtering based on movie release years
ratingMax: For filtering based on your ratings
ratingMin: For filtering based on your ratings
imdbMax: For filtering based on IMDB ratings
imdbMin: For filtering based on IMDB ratings
genres: For filtering based on genres (You can seperate multiple values with "," it acts like "OR")
directors: For filtering based on directors (You can seperate multiple values with "," it acts like "OR")
```

#### `/movies/{id}` (GET)

Use this route to only receive the information of one movie.

#### `/movies/{id}` (PUT)

Use this route to change the information of a movie.

```
JSON BODY:

{
  "title":       ...
  "year":        ...
  "rating":      ... (Between 1 and 10 [both including])
  "imdbRating":  ... (Between 1 and 10 [both including])
  "genres":      [... , ...]
  "directors":   [... , ...]
}
```

#### `/movies/{id}` (DELETE)

Use this route for deleting a movie.

## TUI

The application that consumes the API.

I could've made a web app but I wanted to make something different and made this terminal app.

![Visual](./Assets/tui.gif)

It's based on tview package.

## CLI

A CLI applicaton for feeding CSV data into the API.

You can export your ratings on IMDB or Letterboxd and import them into the database using this CLI tool.

![Visual](./Assets/cli.gif)

I suggest you to use IMDB because it's exported `ratings.csv` file has more information.

To export your ratings:

**FOR IMDB:** Go to `Your Ratings`, click the 3 dots, click `Export`.

**FOR LETTERBOXD:** Go to `Settings`, click `Data`, click `Export Your Data`. (We only need `ratings.csv` other ones are unnecessary)
