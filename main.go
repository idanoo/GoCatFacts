package gocatfacts

import "os"

var (
	port = "8080"
)

func main() {
	// Set custom port if PORT environment variable is set
	portOverride := os.Getenv("PORT")
	if portOverride != "" {
		port = portOverride
	}

}
