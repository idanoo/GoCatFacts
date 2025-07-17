package main

import (
	"bufio"
	"crypto/rand"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
)

var (
	port  = "8080"
	facts = map[int64]string{}
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
	i := 0
	for scanner.Scan() {
		facts[int64(i)] = scanner.Text()
		i++
	}

	log.Println("Loaded", len(facts), "cat facts")
	return scanner.Err()
}

// getRandomFact returns a random cat fact from the loaded facts
func getRandomFact() string {
	if len(facts) == 0 {
		return "No facts available."
	}

	counter := int64(len(facts) - 1)
	n, _ := rand.Int(rand.Reader, big.NewInt(counter))
	return facts[n.Int64()]
}
