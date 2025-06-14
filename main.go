package main

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var (
	port  = "8080"
	facts = []string{}
)

func main() {
	// Set custom port if PORT environment variable is set
	portOverride := os.Getenv("PORT")
	if portOverride != "" {
		port = portOverride
	}

	// Load it up
	err := loadFacts()
	if err != nil {
		log.Fatalf("Error loading facts: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, getRandomFact()+"\n")
	})
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error loading facts: %v", err)
	}
}

// loadFacts loads cat facts from a file named "facts.txt"
func loadFacts() error {
	file, err := os.Open("facts.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		facts = append(facts, scanner.Text())
	}

	return scanner.Err()
}

// getRandomFact returns a random cat fact from the loaded facts
func getRandomFact() string {
	if len(facts) == 0 {
		return "No facts available."
	}

	return facts[rand.Intn(len(facts))]
}
