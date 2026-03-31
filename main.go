package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

const (
	lastFactIDFile = "last_fact_id.json"
	factsFile      = "facts.txt"
)

var (
	port  = "8080"
	facts = map[int64]string{}
)

type CatFactor struct {
	LastFactId int64            `json:"last_fact_id"`
	Facts      map[int64]string `json:"-"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Set custom port if PORT environment variable is set
	portOverride := os.Getenv("PORT")
	if portOverride != "" {
		port = portOverride
	}

	// Load it up
	catFactor := CatFactor{}
	err := catFactor.loadFacts()
	if err != nil {
		slog.Error("Error loading facts", "error", err)
		os.Exit(1)
	}

	// Load last fact ID from JSON file
	err = catFactor.loadLastFactID()
	if err != nil {
		slog.Info("Error loading last fact ID, defaulting to 0", "error", err)
		catFactor.LastFactId = 0
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, catFactor.getRandomFact()+"\n")
	})

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error("Error listening and serving", "error", err)
		os.Exit(1)
	}
}

// loadFacts loads cat facts from a file named "facts.txt"
func (catFactor *CatFactor) loadFacts() error {
	file, err := os.Open(factsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	catFactor.Facts = make(map[int64]string)
	for scanner.Scan() {
		catFactor.Facts[int64(i)] = scanner.Text()
		i++
	}

	slog.Info("Loaded facts", "count", len(catFactor.Facts), "lastID", catFactor.LastFactId)

	return scanner.Err()
}

// getRandomFact returns a random cat fact from the loaded facts
func (catFactor *CatFactor) getRandomFact() string {
	if len(catFactor.Facts) == 0 {
		return "No facts available."
	}

	catFactor.LastFactId++

	// Loop over if needed
	if catFactor.LastFactId >= int64(len(catFactor.Facts)) {
		catFactor.LastFactId = 0
	}
	// Save the updated last fact ID to JSON file
	catFactor.saveLastFactID()

	return catFactor.Facts[catFactor.LastFactId]
}

// saveLastFactID saves the current last fact ID to a JSON file
func (catFactor *CatFactor) saveLastFactID() error {
	data, err := json.Marshal(catFactor)
	if err != nil {
		return err
	}

	return os.WriteFile(lastFactIDFile, data, 0644)
}

// loadLastFactID loads the last fact ID from a JSON file
func (catFactor *CatFactor) loadLastFactID() error {
	data, err := os.ReadFile(lastFactIDFile)
	if err != nil {
		// If file doesn't exist, start with 0
		catFactor.LastFactId = 0
		return nil
	}

	return json.Unmarshal(data, catFactor)
}
