PORT=8080
FILE=

all:
	@echo "For more information: github.com/Scienitive/moviterm"

api:
	go run API/*.go -port $(PORT)

tui:
	go run TUI/*.go -port $(PORT)

cli-imdb:
	go run CLI/main.go -port $(PORT) -I $(FILE)

cli-letterboxd:
	go run CLI/main.go -port $(PORT) -L $(FILE)
