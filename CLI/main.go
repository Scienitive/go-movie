package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type IMDB struct {
	Title      string
	Year       int
	Rating     int
	ImdbRating float32
	Genres     []string
	Directors  []string
}

type Letterboxd struct {
	Title  string
	Year   int
	Rating int
}

type DuplicateError struct {
	Message     string
	DuplicateID int
}

func main() {
	// Defining Flags
	var (
		fImdb  = flag.Bool("I", false, "IMDB Flag")
		fLboxd = flag.Bool("L", false, "Letterboxd Flag")
		fPort  = flag.Int("port", 8080, "Port of the API")
	)

	// Store the non-flag arguments
	flag.Parse()
	port := *fPort
	args := flag.Args()

	// Error Handling
	if *fImdb && *fLboxd {
		errorExit("Cannot specify both -I and -L flags simultaneously.")
	} else if !*fImdb && !*fLboxd {
		errorExit("You need to specify either -I (IMDB) or -L (Letterboxd) flag.")
	}
	if len(args) != 1 {
		errorExit("You need to provide '1' csv file.")
	}

	// Preparing file to reading
	file, err := os.Open(args[0])
	if err != nil {
		errorExit("Cannot open the file.")
	}

	// The Thing
	theURL := fmt.Sprintf("http://localhost:%d/movies", port)
	if *fImdb {
		readFileIMDB(file, theURL)
	} else {
		readFileLetterboxd(file, theURL)
	}
}

func readFileIMDB(file *os.File, url string) {
	fileScanner := bufio.NewScanner(file)
	i := 0
	for fileScanner.Scan() {
		if i == 0 {
			i++
			continue
		}
		line := fileScanner.Text()
		if len(line) == 0 {
			i++
			continue
		}
		split := split(line, ',')
		if len(split) != 13 {
			errorMessage := fmt.Sprintf("Line %d, Unexpected format.", i)
			errorExit(errorMessage)
		}
		if split[5] != "movie" {
			i++
			continue
		}

		entry := IMDB{
			Title:      strings.Trim(split[3], "\""),
			Year:       atoiExit(split[8], i),
			Rating:     atoiExit(split[1], i),
			ImdbRating: atofExit(split[6], i),
			Genres:     strings.Split(strings.Trim(split[9], "\""), ", "),
			Directors:  strings.Split(strings.Trim(split[12], "\""), ", "),
		}

		jsonData, err := json.Marshal(&entry)
		if err != nil {
			errorExit("Cannot turn the data into JSON")
		}

		client := http.Client{
			Timeout: time.Second,
		}
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			errorExit("Cannot send request to server.\nMake sure the API is on.")
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errorExit("Cannot read the response body.")
		}

		if resp.StatusCode == http.StatusCreated {
			fmt.Printf("%s (%d) added!\n", entry.Title, entry.Year)
		} else if resp.StatusCode == http.StatusConflict {
			contentType := resp.Header.Get("Content-Type")
			handleConflict(resp, body, entry.Title, entry.Year, jsonData, contentType == "application/json", url)
		} else {
			errorExit(string(body))
		}

		i++
	}
}

func readFileLetterboxd(file *os.File, url string) {
	fileScanner := bufio.NewScanner(file)
	i := 0
	for fileScanner.Scan() {
		if i == 0 {
			i++
			continue
		}
		line := fileScanner.Text()
		if len(line) == 0 {
			i++
			continue
		}
		split := split(line, ',')
		if len(split) != 5 {
			errorMessage := fmt.Sprintf("Line %d, Unexpected format.", i)
			errorExit(errorMessage)
		}

		entry := Letterboxd{
			Title:  strings.Trim(split[1], "\""),
			Year:   atoiExit(split[2], i),
			Rating: int(atofExit(split[4], i) * 2),
		}

		jsonData, err := json.Marshal(&entry)
		if err != nil {
			errorExit("Cannot turn the data into JSON")
		}

		client := http.Client{
			Timeout: time.Second,
		}
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			errorExit("Cannot send request to server.\nMake sure the API is on.")
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errorExit("Cannot read the response body.")
		}

		if resp.StatusCode == http.StatusCreated {
			fmt.Printf("%s (%d) added!\n", entry.Title, entry.Year)
		} else if resp.StatusCode == http.StatusConflict {
			contentType := resp.Header.Get("Content-Type")
			handleConflict(resp, body, entry.Title, entry.Year, jsonData, contentType == "application/json", url)
		} else {
			errorExit(string(body))
		}

		i++
	}
}

// Worst Code Ever
func handleConflict(resp *http.Response, body []byte, movieTitle string, movieYear int, jsonData []byte, isJSON bool, url string) {
	client := http.Client{
		Timeout: time.Second,
	}
ForeverLoop:
	for {
		// Printing the Information
		if isJSON {
			fmt.Printf("You already have %s (%d) in your database. What should we do?\n(1) Skip | (2) Create Another Entry | (3) Overwrite\n", movieTitle, movieYear)
		} else {
			fmt.Printf("You already have %s (%d) in your database. What should we do?\n(1) Skip | (2) Create Another Entry\n", movieTitle, movieYear)
		}

		// Waiting for User Input
		var value string
		_, err := fmt.Scan(&value)
		if err != nil {
			errorExit("User input error.")
		}

		// Conclusion
		switch value {
		case "1":
			break ForeverLoop
		case "2":
			_, err := client.Post(url+"/force", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				errorExit("Cannot force insert the movie.")
			}
			break ForeverLoop
		case "3":
			if !isJSON {
				break
			}
			var errInfo DuplicateError
			if err := json.Unmarshal(body, &errInfo); err != nil {
				fmt.Println(err.Error())
				errorExit("JSON Unmarshaling failed.")
			}
			movieID := errInfo.DuplicateID
			putURL := fmt.Sprintf("%s/%d", url, movieID)
			req, err := http.NewRequest("PUT", putURL, bytes.NewBuffer(jsonData))
			if err != nil {
				errorExit("Cannot update the movie.")
			}
			req.Header.Set("Content-Type", "application/json")
			client.Do(req)
			break ForeverLoop
		default:
			break
		}
	}
}

func split(input string, seperator rune) []string {
	var result []string
	current := ""
	inQuote := false

	for _, char := range input {
		if char == '"' {
			inQuote = !inQuote
		}

		if char == seperator && !inQuote {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}

	result = append(result, current)
	return result
}

func atoiExit(s string, i int) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		errorMessage := fmt.Sprintf("Line %d, Unexpected format.", i)
		errorExit(errorMessage)
	}
	return res
}

func atofExit(s string, i int) float32 {
	res, err := strconv.ParseFloat(s, 32)
	if err != nil {
		errorMessage := fmt.Sprintf("Line %d, Unexpected format.", i)
		errorExit(errorMessage)
	}
	return float32(res)
}

func errorExit(m string) {
	fmt.Fprintln(os.Stderr, "Error: ", m)
	os.Exit(1)
}
