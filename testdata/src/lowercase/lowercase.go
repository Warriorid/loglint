package lowercase

import (
	"log"
	"log/slog"
)

func examplesLowercase() {
	log.Print("Starting server on port 8080")   // want `log message should start with a lowercase letter`
	log.Print("starting server on port 8080")
	slog.Error("Failed to connect to database") // want `log message should start with a lowercase letter`
	slog.Error("failed to connect to database")
	log.Print("Hello world")                    // want `log message should start with a lowercase letter`
	log.Print("hello world")
	log.Print("ERROR occurred")                 // want `log message should start with a lowercase letter`
	log.Print("error occurred")
}
